package service

import (
	"fmt"
	embed "github.com/Clinet/discordgo-embed"
	"github.com/gemalto/flume"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"osrs-disc-bot/util"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Service struct {
	collectionLog collectionLog
	sheets        sheets
	imgur         imgur
	temple        temple

	submissions   map[string]int
	cpscreenshots map[string]string
	log           flume.Logger
	tid           int
	config        *util.Config
	client        *http.Client
}

func NewService(config *util.Config, collectionLog collectionLog, sheets sheets, imgur imgur, temple temple) *Service {
	logger := flume.New("service")
	client := &http.Client{Timeout: 30 * time.Second}
	return &Service{
		collectionLog: collectionLog,
		sheets:        sheets,
		imgur:         imgur,
		temple:        temple,

		submissions:   make(map[string]int),
		cpscreenshots: make(map[string]string),
		log:           logger,
		tid:           0,

		config: config,
		client: client,
	}
}

// Use of discordgo as an intro to discord IRC
func (s *Service) StartDiscordIRC() {
	s.sheets.InitializeSubmissionsFromSheet(s.submissions)

	// Create a new discord session
	session, err := discordgo.New("Bot " + s.config.DiscBotToken)
	if err != nil {
		log.Fatal(err)
	}

	// Create handler for listening for submission messages
	session.AddHandler(s.listenForMessage)

	// Send intent
	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	_ = session.Open()
	defer func(discord *discordgo.Session) {
		err := discord.Close()
		if err != nil {

		}
	}(session)

	// Initialize the Hall Of fame
	s.kickOffHallOfFameUpdate(session)
	fmt.Println("the bot is online!")

	// Throw a go func that will capture signal interrupts and will populate the submissions file
	go func() {
		sigchan := make(chan os.Signal)
		signal.Notify(sigchan,
			// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
			syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
			syscall.SIGINT,  // Ctrl+C
			syscall.SIGQUIT, // Ctrl-\
			syscall.SIGKILL, // "always fatal", "SIGKILL and SIGSTOP may not be caught by a program"
			syscall.SIGHUP,  // "terminal is disconnected"
		)
		<-sigchan

		// Once the program is interrupted, update the google clan points sheet
		s.sheets.UpdateCpSheet(s.submissions)
		s.sheets.UpdateCpScreenshotsSheet(s.cpscreenshots)
	}()

	// Block so that it continues to run the bot
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func (s *Service) listenForMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Don't handle message if it's created by the discord bot
	if message.Author.ID == session.State.User.ID {
		return
	}

	ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid))
	defer func() { s.tid++ }()

	// Run certain tasks depending on the channel the message was posted in
	switch channel := message.ChannelID; channel {
	case s.config.DiscSubChan:
		s.listenForCPSubmission(ctx, session, message)
	case s.config.DiscSignUpChan:
		s.updateMemberList(ctx, session, message)
	default:
		// Return if the message was not posted in one of the channels we are handling
		return
	}
}

func (s *Service) listenForCPSubmission(ctx context.Context, session *discordgo.Session, message *discordgo.MessageCreate) {
	// Split the names into an array by , then make an empty array with those names as keys for an easier lookup
	// instead of running a for loop inside a for loop when adding clan points
	whitespaceStrippedMessage := strings.Replace(message.Content, ", ", ",", -1)
	whitespaceStrippedMessage = strings.Replace(whitespaceStrippedMessage, " ,", ",", -1)

	names := strings.Split(whitespaceStrippedMessage, ",")

	// Before adding clanpoints, ensure that all the names used in the submission is valid and already created
	// in the #ponies-signup channel
	for _, name := range names {
		// Ensure that this person does not exist in the submissions map currently
		if _, ok := s.submissions[name]; !ok {
			// Create a private channel with the user submitting (will reuse if one exists)
			channel, err := session.UserChannelCreate(message.Author.ID)
			if err != nil {
				return
			}

			// Send a message on that channel
			_, err = session.ChannelMessageSend(channel.ID, "Non clan member used in this submission. "+
				"Please add the user: \""+name+"\" using the https://discord.com/channels/1172535371905646612/1173253913303056524 channel and resubmit the screenshot with the names.")
			if err != nil {
				return
			}

			// Once everything is finished, delete the message from the submission channel
			err = session.ChannelMessageDelete(s.config.DiscSubChan, message.ID)
			if err != nil {
				return
			}

			return
		}
	}

	numberOfSubmissions := 0

	// Iterate through all the pictures and download them
	for _, submissionPicture := range message.Attachments {
		// If it's an imgur link, save the link in the cpscreenshots map
		if strings.Contains(submissionPicture.ProxyURL, "imgur") {
			s.cpscreenshots[submissionPicture.ProxyURL] = whitespaceStrippedMessage
		} else if strings.Contains(submissionPicture.ProxyURL, "media.discordapp.net") {
			// Retrieve the access token
			accessToken := s.imgur.GetNewAccessToken(ctx, s.config.ImgurRefreshToken, s.config.ImgurClientId, s.config.ImgurClientSecret)

			// Retrieve the bytes of the image
			resp, err := s.client.Get(submissionPicture.ProxyURL)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			submissionUrl := s.imgur.Upload(ctx, accessToken, resp.Body)
			s.cpscreenshots[submissionUrl] = whitespaceStrippedMessage
		} else {
		}
		numberOfSubmissions++
	}

	// Iterate over the all the names in the submissions and add the number of submissions to their clan points
	for _, name := range names {
		s.submissions[name] = s.submissions[name] + numberOfSubmissions
	}

	// Update the #cp-leaderboard
	s.updateLeaderboard(session)

	// Once everything is finished, delete the message from the submission channel
	err := session.ChannelMessageDelete(s.config.DiscSubChan, message.ID)
	if err != nil {
		return
	}
}

// updateLeaderboard will update the cp-leaderboard channel in discord with a new ranking of everyone in the clan
func (s *Service) updateLeaderboard(session *discordgo.Session) {
	// Update the #cp-leaderboard
	keys := make([]string, 0, len(s.submissions))
	for key := range s.submissions {
		keys = append(keys, key)
	}

	// Sort the map based on the values
	sort.SliceStable(keys, func(i, j int) bool {
		return s.submissions[keys[i]] > s.submissions[keys[j]]
	})

	// Create the leaderboard message that will be sent
	leaderboard := ""
	for placement, k := range keys {
		leaderboard = leaderboard + strconv.Itoa(placement+1) + ") " + k + ": " + strconv.Itoa(s.submissions[k]) + "\n"
	}

	// Retrieve the one channel message and delete it in the leaderboard channel
	messages, err := session.ChannelMessages(s.config.DiscLeaderboardChan, 1, "", "", "")
	if err != nil {
		return
	}
	err = session.ChannelMessageDelete(s.config.DiscLeaderboardChan, messages[0].ID)
	if err != nil {
		return
	}

	_, err = session.ChannelMessageSendEmbed(s.config.DiscLeaderboardChan, embed.NewEmbed().
		SetTitle("Ponies Clan Points Leaderboard").
		SetDescription(fmt.Sprintf(leaderboard)).
		SetColor(0x1c1c1c).SetThumbnail("https://i.imgur.com/O4NzB95.png").MessageEmbed)
	if err != nil {
		return
	}
}

func (s *Service) updateMemberList(ctx context.Context, session *discordgo.Session, message *discordgo.MessageCreate) {
	//TODO: SCRUB THE USERNAME SUBMITTED
	// Don't include the remove command in the RSN
	member := strings.Replace(message.Content, "!rm ", "", -1)

	// Create a private channel with the user submitting (will reuse if one exists)
	channel, err := session.UserChannelCreate(message.Author.ID)
	if err != nil {
		return
	}

	// Remove user from temple if the message prefix is "rm"
	re := regexp.MustCompile("(?i)^(!rm)\\s+.+$") // Case insensitive. Must start with "!rm". Must have atleast one space between "!rm" and the username. There must be text after "!rm". We use "!" at the beginning in case a user's name starts with "rm".
	if re.Match([]byte(message.Content)) {
		// Remove the user from the temple page
		s.temple.RemoveMemberFromTemple(member, s.config.TempleGroupId, s.config.TempleGroupKey)

		if s.userExists(member) {
			s.submissions[member] = 0

			// Send a message on that channel
			_, err := session.ChannelMessageSend(channel.ID, "You have successfully removed a member: "+member)
			if err != nil {
				return
			}
		} else {
			_, err := session.ChannelMessageSend(channel.ID, "Member: "+member+" does not exist.")
			if err != nil {
				return
			}
		}

		// Once everything is finished, delete the message from the submission channel
		err = session.ChannelMessageDelete(s.config.DiscSignUpChan, message.ID)
		if err != nil {
			return
		}

		// Don't continue because the following code is to add a user
		return
	}

	// Ensure that this person does not exist in the submissions map currently
	if s.userExists(member) {
		_, err := session.ChannelMessageSend(channel.ID, "Member: "+member+" already exists.")
		if err != nil {
			return
		}
	} else {
		s.submissions[member] = 0

		// Send a message on that channel
		_, err := session.ChannelMessageSend(channel.ID, "You have successfully added new member: "+member)
		if err != nil {
			return
		}
	}

	// Add the user to the temple page
	s.temple.AddMemberToTemple(member, s.config.TempleGroupId, s.config.TempleGroupKey)

	// Once everything is finished, delete the message from the submission channel
	err = session.ChannelMessageDelete(s.config.DiscSignUpChan, message.ID)
	if err != nil {
		return
	}
}

func (s *Service) userExists(member string) (exists bool) {
	exists = false

	if _, ok := s.submissions[member]; ok {
		exists = true
	}

	return
}

func (s *Service) kickOffHallOfFameUpdate(session *discordgo.Session) {
	//slayerBosses := util.HallOfFameRequestInfo{Bosses: map[string]string{"sire": "https://i.imgur.com/GhbmqEB.png", "hydra": "https://i.imgur.com/25GU0Ph.png", "cerberus": "https://i.imgur.com/UoxGuQi.png", "grotesqueguardians": "https://i.imgur.com/M7ylVBZ.png", "kraken": "https://i.imgur.com/Q6EbJb1.png", "smokedevil": "https://i.imgur.com/2AYntQ5.png"}, DiscChan: s.config.DiscSlayerBossesChan}
	//gwd := util.HallOfFameRequestInfo{Bosses: map[string]string{"commanderzilyana": "https://i.imgur.com/aNm4Ydd.png", "kreearra": "https://i.imgur.com/lX8SfgN.png", "kriltsutsaroth": "https://i.imgur.com/hh8cMvp.png", "nex": "https://i.imgur.com/pqiVQBC.png", "generalgraardor": "https://i.imgur.com/hljv9ZW.png"}, DiscChan: s.config.DiscGwdChan}
	//wildy := util.HallOfFameRequestInfo{Bosses: map[string]string{"artio": "https://i.imgur.com/bw6zLpU.png", "callisto": "https://i.imgur.com/bw6zLpU.png", "calvarion": "https://i.imgur.com/v3KX75y.png", "vetion": "https://i.imgur.com/v3KX75y.png", "spindel": "https://i.imgur.com/4zknWSX.png", "venenatis": "https://i.imgur.com/4zknWSX.png", "chaoselemental": "https://i.imgur.com/YAvIpbm.png", "chaosfanatic": "https://i.imgur.com/azV2sD1.png", "crazyarchaeologist": "https://i.imgur.com/23LXv53.png", "scorpia": "https://i.imgur.com/9aaguxB.png"}, DiscChan: s.config.DiscWildyChan}
	//other := util.HallOfFameRequestInfo{Bosses: map[string]string{"corporealbeast": "https://i.imgur.com/zEDN4Pf.png", "prime": "https://i.imgur.com/kJBtqHB.png", "rexbro": "https://i.imgur.com/PvlGWFZ.png", "supreme": "https://i.imgur.com/BOgkBuD.png", "gauntlet": "https://i.imgur.com/weiHWnz.png", "gauntlethard": "https://i.imgur.com/xzW4TGR.png", "giantmole": "https://i.imgur.com/coKk2pr.gif", "jad": "https://i.imgur.com/H9aO1Ot.png", "zuk": "https://i.imgur.com/mKstHza.png", "kq": "https://i.imgur.com/ZuaFoBR.png", "kbd": "https://i.imgur.com/r5vkw1s.png", "sarachnis": "https://i.imgur.com/98THH8O.png", "skotizo": "https://i.imgur.com/YUcQu4d.png", "muspah": "https://i.imgur.com/sW2cLQ2.png", "vorkath": "https://i.imgur.com/6biF3P2.png", "phosanis": "https://i.imgur.com/4aDkxms.png", "nightmare": "https://i.imgur.com/4aDkxms.png", "zulrah": "https://i.imgur.com/tPllWNF.png"}, DiscChan: s.config.DiscOtherChan}
	//misc := util.HallOfFameRequestInfo{Bosses: map[string]string{"barrows": "https://i.imgur.com/ajoK20v.png", "hespori": "https://i.imgur.com/b0qYGHS.png", "mimic": "https://i.imgur.com/jC7yTC3.png", "obor": "https://i.imgur.com/dwLvSbR.png", "bryophyta": "https://i.imgur.com/3cdyp4X.png", "derangedarchaeologist": "https://i.imgur.com/cnHpevF.png", "wintertodt": "https://i.imgur.com/6oFef2Y.png", "zalcano": "https://i.imgur.com/edN11Nf.png", "tempoross": "https://i.imgur.com/fRj3JA4.png", "rift": "https://i.imgur.com/MOiyXeH.png"}, DiscChan: s.config.DiscMiscChan}
	//dt2 := util.HallOfFameRequestInfo{Bosses: map[string]string{"duke": "https://i.imgur.com/RYPmrXy.png", "leviathan": "https://i.imgur.com/mEQRq5c.png", "whisperer": "https://i.imgur.com/cFGWb6Y.png", "vardorvis": "https://i.imgur.com/WMPuShZ.png"}, DiscChan: s.config.DiscDT2Chan}
	//raids := util.HallOfFameRequestInfo{Bosses: map[string]string{"cox": "https://i.imgur.com/gxdWXtH.png", "coxcm": "https://i.imgur.com/gxdWXtH.png", "tob": "https://i.imgur.com/pW1sJAQ.png", "tobcm": "https://i.imgur.com/pW1sJAQ.png", "toa": "https://i.imgur.com/2GvzqGw.png", "toae": "https://i.imgur.com/2GvzqGw.png"}, DiscChan: s.config.DiscRaidsChan}
	//pvp := util.HallOfFameRequestInfo{Bosses: map[string]string{"bhh": "https://i.imgur.com/zSQhlWk.png", "bhr": "https://i.imgur.com/Y3Sga7t.png", "lms": "https://i.imgur.com/rzW7ZXx.png", "arena": "https://i.imgur.com/uNP6Ggu.png", "zeal": "https://i.imgur.com/Ws7HvKL.png"}, DiscChan: s.config.DiscPVPChan}
	//clues := util.HallOfFameRequestInfo{Bosses: map[string]string{"clueall": "https://i.imgur.com/wX3Ei7U.png", "cluebeginner": "https://i.imgur.com/fUmzJkW.png", "clueeasy": "https://i.imgur.com/phnSCHj.png", "cluemedium": "https://i.imgur.com/t5iH8Xa.png", "cluehard": "https://i.imgur.com/a0xwcGI.png", "clueelite": "https://i.imgur.com/ibNRk3G.png", "cluemaster": "https://i.imgur.com/12rCLVv.png"}, DiscChan: s.config.DiscCluesChan}
	//
	//s.updateHallOfFame(session, slayerBosses)
	//s.updateHallOfFame(session, gwd)
	//s.updateHallOfFame(session, wildy)
	//s.updateHallOfFame(session, other)
	//s.updateHallOfFame(session, misc)
	//s.updateHallOfFame(session, dt2)
	//s.updateHallOfFame(session, raids)
	//s.updateHallOfFame(session, pvp)
	//s.updateHallOfFame(session, clues)
	s.updateCollectionLog(session)
}

func (s *Service) updateHallOfFame(session *discordgo.Session, requestInfo util.HallOfFameRequestInfo) {
	// First, delete all the messages within the channel
	messages, err := session.ChannelMessages(requestInfo.DiscChan, 50, "", "", "")
	if err != nil {
		return
	}

	var messageIDs []string
	for _, message := range messages {
		messageIDs = append(messageIDs, message.ID)
	}

	err = session.ChannelMessagesBulkDelete(requestInfo.DiscChan, messageIDs)
	if err != nil {
		return
	}

	// Now add all the bosses
	for bossIdForTemple, imageURL := range requestInfo.Bosses {
		podium, rankings := s.temple.GetPodiumFromTemple(bossIdForTemple)

		// Iterate over the players to get the different places for users to create the placements
		placements := ""
		for _, k := range rankings {
			switch k {
			case 1:
				placements = placements + ":first_place: "
				break
			case 2:
				placements = placements + ":second_place: "
				break
			case 3:
				placements = placements + ":third_place: "
				break
			}
			placements = placements + podium.Data.Players[k].Username + " [" + strconv.Itoa(podium.Data.Players[k].Kc) + "]\n"
		}

		_, err = session.ChannelMessageSendEmbed(requestInfo.DiscChan, embed.NewEmbed().
			SetTitle(podium.Data.BossName).
			SetDescription(placements).
			SetColor(0x1c1c1c).SetThumbnail(imageURL).MessageEmbed)
		if err != nil {
			return
		}
	}
	return
}

func (s *Service) updateCollectionLog(session *discordgo.Session) {
	podium, ranking := s.collectionLog.RetrieveCollectionLogAndOrder(s.submissions)

	// Create the leaderboard message that will be sent
	placements := ""
	for placement, k := range ranking {
		switch placement {
		case 0:
			placements = placements + ":one: "
			break
		case 1:
			placements = placements + ":two: "
			break
		case 2:
			placements = placements + ":three: "
			break
		case 3:
			placements = placements + ":four: "
			break
		case 4:
			placements = placements + ":five: "
			break
		case 5:
			placements = placements + ":six: "
			break
		case 6:
			placements = placements + ":seven: "
			break
		case 7:
			placements = placements + ":eight: "
			break
		case 8:
			placements = placements + ":nine: "
			break
		case 9:
			placements = placements + ":keycap_10: "
			break

		}

		placements = placements + k + " [" + strconv.Itoa(podium[k]) + "]\n"
	}

	// First, delete all the messages within the channel
	messages, err := session.ChannelMessages(s.config.DiscColChan, 10, "", "", "")
	if err != nil {
		return
	}

	var messageIDs []string
	for _, message := range messages {
		messageIDs = append(messageIDs, message.ID)
	}

	err = session.ChannelMessagesBulkDelete(s.config.DiscColChan, messageIDs)
	if err != nil {
		return
	}

	// Send the collection log message
	_, err = session.ChannelMessageSendEmbed(s.config.DiscColChan, embed.NewEmbed().
		SetTitle("Collection Log Ranking").
		SetDescription(placements).
		SetColor(0x1c1c1c).SetThumbnail("https://i.imgur.com/otTd8Dg.png").MessageEmbed)
	if err != nil {
		return
	}

	// Send the instructions on how to get on the collection log hall of fame
	var msg string
	msg = msg + "1. Download the collection-log plugin\n"
	msg = msg + "2. Click the box to \"Allow collectionlog.net connections\"\n"
	msg = msg + "3. Click through the collection log (there will be a * next to the one you still need to click)\n"
	msg = msg + "4. Go to the collection log icon on the sidebar\n"
	msg = msg + "5. Click Account at the top and then upload collection log\n"
	_, err = session.ChannelMessageSendEmbed(s.config.DiscColChan, embed.NewEmbed().
		SetTitle("How To Get Onto The Collection Log HOF").
		SetDescription(msg).
		SetColor(0x1c1c1c).SetThumbnail("https://i.imgur.com/otTd8Dg.png").MessageEmbed)
	if err != nil {
		return
	}
}
