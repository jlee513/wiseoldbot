package service

import (
	"context"
	"fmt"
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

	embed "github.com/Clinet/discordgo-embed"
	"github.com/gemalto/flume"
	"github.com/go-co-op/gocron"

	"github.com/bwmarrin/discordgo"
)

type Service struct {
	collectionLog collectionLog
	sheets        sheets
	imgur         imgur
	temple        temple
	runescape     runescape

	submissions   map[string]int
	cpscreenshots map[string]string
	log           flume.Logger
	tid           int

	config *util.Config
	client *http.Client

	scheduler *gocron.Scheduler
}

func NewService(config *util.Config, collectionLog collectionLog, sheets sheets, imgur imgur, temple temple, runescape runescape) *Service {
	logger := flume.New("service")
	if config.LogDebug {
		_ = flume.Configure(flume.Config{Development: true, Levels: "*"})
	} else {
		_ = flume.Configure(flume.Config{Development: true})
	}
	client := &http.Client{Timeout: 30 * time.Second}
	s := gocron.NewScheduler(time.UTC)
	s.SingletonModeAll()
	return &Service{
		collectionLog: collectionLog,
		sheets:        sheets,
		imgur:         imgur,
		temple:        temple,
		runescape:     runescape,

		submissions:   make(map[string]int),
		cpscreenshots: make(map[string]string),
		log:           logger,
		tid:           1,

		config: config,
		client: client,

		scheduler: s,
	}
}

// StartDiscordIRC uses discordgo as an intro to discord IRC
func (s *Service) StartDiscordIRC() {
	s.log.Info("Initializing OSRS Disc Bot...")
	ctx := context.Background()
	s.sheets.InitializeSubmissionsFromSheet(ctx, s.submissions)

	// Create a new discord session
	session, err := discordgo.New("Bot " + s.config.DiscBotToken)
	if err != nil {
		s.log.Error("Failed to start discord bot: " + err.Error())
		panic(err)
	}

	// Create handler for listening for submission messages
	session.AddHandler(s.listenForMessage)

	// Set the session's intent
	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// Open up the session and defer the closing of the session
	_ = session.Open()
	defer func(discord *discordgo.Session) {
		err := discord.Close()
		if err != nil {
			s.log.Error("Failed to stop discord bot: " + err.Error())
		}
	}(session)

	s.kickOffCron(ctx, session)
	s.log.Info("OSRS Disc Bot is now online!")

	// Block so that it continues to run the bot
	sigchan := make(chan os.Signal)
	signal.Notify(sigchan,
		// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGKILL, // "always fatal", "SIGKILL and SIGSTOP may not be caught by a program"
		syscall.SIGHUP,  // "terminal is disconnected"
		os.Interrupt,    // os interrupt
	)
	<-sigchan

	// Once the program is interrupted, update the Google Sheet clan points & screenshot sheets
	s.sheets.UpdateCpSheet(ctx, s.submissions)
	s.sheets.UpdateCpScreenshotsSheet(ctx, s.cpscreenshots)
	// Stop the cron scheduler
	s.scheduler.Stop()
	s.log.Info("OSRS Disc Bot is now offline!")
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

/*
listenForCPSubmission will listen on the submission channel looking for a comma separated name in the
text and an image in a single message. It will determine how many images and how many people and will
supply the correct number of clan points to everyone in the list
*/
func (s *Service) listenForCPSubmission(ctx context.Context, session *discordgo.Session, message *discordgo.MessageCreate) {
	logger := flume.FromContext(ctx)
	logger.Info("Clan Point submission triggered.")

	// Defer the deletion of the message
	defer func(messageId string) {
		// Once everything is finished, delete the message from the submission channel
		err := session.ChannelMessageDelete(s.config.DiscSubChan, messageId)
		if err != nil {
			logger.Error("Failed to delete channel message: " + err.Error())
			return
		}
	}(message.ID)

	// Remove leading and trailing whitespaces
	msg := strings.TrimSpace(message.Content)
	logger.Debug("SUBMISSION MESSAGE: " + msg)

	// First, check if an i.imgur.com URL is used as a submission
	startOfImgurUrl := strings.Index(msg, "https://i.imgur.com")
	imgurUrl := ""

	// Only i.imgur.com links will work - other links will throw an error
	otherUrl := strings.Index(msg, "https://")
	if startOfImgurUrl == -1 && otherUrl > -1 {
		logger.Error("Only https://i.imgur.com links are valid: " + msg)
		msg := "Only https://i.imgur.com links are valid. Either resubmit as an imgur or upload the photo to discord submission message."
		s.sendPrivateMessage(ctx, session, message.Author.ID, msg)
		return
	}

	// If we have an i.imgur link, take the link out
	if startOfImgurUrl > -1 {
		// We have an imgur link, determine if it's PNG or JPEG
		endOfUrlPNG := strings.Index(msg, ".png")
		endOfUrlJPEG := strings.Index(msg, ".jpeg")

		endOfUrl := -1
		if endOfUrlPNG > -1 {
			endOfUrl = endOfUrlPNG + 4
		} else if endOfUrlJPEG > -1 {
			endOfUrl = endOfUrlJPEG + 5
		} else {
			logger.Error("Another image type other than PNG or JPEG was provided.")
			msg := "Another image type other than PNG or JPEG was provided. Please resubmit with either PNG or JPEG."
			s.sendPrivateMessage(ctx, session, message.Author.ID, msg)
			return
		}

		imgurUrl = msg[startOfImgurUrl:endOfUrl]
		// If the start of the URL is at the beginning of the message...
		if startOfImgurUrl == 0 {
			// Set the rest of the message after the .png as the message
			msg = msg[endOfUrl+1:]
		} else {
			msg = msg[:startOfImgurUrl-1]
		}
	}

	// Split the names into an array by , then make an empty array with those names as keys for an easier lookup
	// instead of running a for loop inside a for loop when adding clan points
	whitespaceStrippedMessage := strings.Replace(msg, ", ", ",", -1)
	whitespaceStrippedMessage = strings.Replace(whitespaceStrippedMessage, " ,", ",", -1)

	logger.Debug("Submitted names: " + whitespaceStrippedMessage)

	names := strings.Split(whitespaceStrippedMessage, ",")

	// Before adding clanpoints, ensure that all the names used in the submission is valid and already created
	// in the #ponies-signup channel
	for _, name := range names {
		// Ensure that this person does not exist in the submissions map currently
		if _, ok := s.submissions[name]; !ok {
			logger.Error("Non clan member used in this submission: " + name)
			msg := "Non clan member used in this submission. Please add the user: \"" + name + "\" using the " +
				"https://discord.com/channels/1172535371905646612/1176891514325057566 channel and resubmit " +
				"the screenshot with the names."
			s.sendPrivateMessage(ctx, session, message.Author.ID, msg)
			return
		}
	}

	// Allow for more than 1 image per submission
	numberOfSubmissions := 0

	// If there is an imgur URL, there won't be an attachment to the submission
	if len(imgurUrl) > 0 {
		s.cpscreenshots[imgurUrl] = whitespaceStrippedMessage
		numberOfSubmissions++
	} else {
		// Iterate through all the pictures and download them
		for _, submissionPicture := range message.Attachments {
			if !strings.Contains(submissionPicture.ContentType, "image") {
				// Invalid submission
				logger.Error("Invalid submission content. Submitted type: " + submissionPicture.ContentType)
				msg := "Only image attachments are allowed. Either resubmit as an imgur or upload the attachment as a photo to the discord submission message."
				s.sendPrivateMessage(ctx, session, message.Author.ID, msg)
				return
			}
			logger.Info(submissionPicture.ContentType)
			// If it's an imgur link, save the link in the cpscreenshots map
			if strings.Contains(submissionPicture.ProxyURL, "media.discordapp.net") {
				// Retrieve the access token
				accessToken := s.imgur.GetNewAccessToken(ctx, s.config.ImgurRefreshToken, s.config.ImgurClientId, s.config.ImgurClientSecret)

				// Retrieve the bytes of the image
				resp, err := s.client.Get(submissionPicture.ProxyURL)
				if err != nil {
					logger.Error("Failed to download discord image: " + err.Error())
					msg := "Failed to download discord image - please wait before trying again."
					s.sendPrivateMessage(ctx, session, message.Author.ID, msg)
					return
				}
				defer resp.Body.Close()

				submissionUrl := s.imgur.Upload(ctx, accessToken, resp.Body)
				s.cpscreenshots[submissionUrl] = whitespaceStrippedMessage
			} else {
				// Invalid submission
				logger.Error("INVALID SUBMISSION: " + submissionPicture.ProxyURL)
				msg := "Invalid submission - please upload the picture to imgur before submitting again."
				s.sendPrivateMessage(ctx, session, message.Author.ID, msg)
				return
			}
			numberOfSubmissions++
		}
	}

	// Iterate over the all the names in the submissions and add the number of submissions to their clan points
	for _, name := range names {
		s.submissions[name] = s.submissions[name] + numberOfSubmissions
	}

	// Update the #cp-leaderboard
	s.updateLeaderboard(ctx, session)

	logger.Info("Clan Point submission successful.")
}

// updateLeaderboard will update the cp-leaderboard channel in discord with a new ranking of everyone in the clan
func (s *Service) updateLeaderboard(ctx context.Context, session *discordgo.Session) {
	logger := flume.FromContext(ctx)
	logger.Info("Running clan point leaderboard update...")

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
		logger.Error("ERROR RETRIEVING MESSAGES FROM DISCORD LEADERBOARD CHANNEL.")
		return
	}
	err = session.ChannelMessageDelete(s.config.DiscLeaderboardChan, messages[0].ID)
	if err != nil {
		logger.Error("ERROR DELETING MESSAGES FROM DISCORD LEADERBOARD CHANNEL.")
		return
	}

	// Send the Discord Embed message
	_, err = session.ChannelMessageSendEmbed(s.config.DiscLeaderboardChan, embed.NewEmbed().
		SetTitle("Ponies Clan Points Leaderboard").
		SetDescription(fmt.Sprintf(leaderboard)).
		SetColor(0x1c1c1c).SetThumbnail("https://i.imgur.com/O4NzB95.png").MessageEmbed)
	if err != nil {
		logger.Error("ERROR SENDING MESSAGES TO DISCORD LEADERBOARD CHANNEL.")
		return
	}

	logger.Info("Clan point leaderboard update successful.")
}

/*
updateMemberList will deal with the addition of a new member to the submission map and to temple as well
as with the removal of an existing member
*/
func (s *Service) updateMemberList(ctx context.Context, session *discordgo.Session, message *discordgo.MessageCreate) {
	//TODO: SCRUB THE USERNAME SUBMITTED
	// Don't include the remove command in the RSN
	logger := flume.FromContext(ctx)
	member := strings.Replace(message.Content, "!rm ", "", -1)

	// Defer the deletion of the message
	defer func(messageId string) {
		// Once everything is finished, delete the message from the submission channel
		err := session.ChannelMessageDelete(s.config.DiscSignUpChan, message.ID)
		if err != nil {
			logger.Error("Failed to delete channel message: " + err.Error())
			return
		}
	}(message.ID)

	// Remove user from temple if the message prefix is "rm"
	re := regexp.MustCompile("(?i)^(!rm)\\s+.+$") // Case-insensitive. Must start with "!rm". Must have atleast one space between "!rm" and the username. There must be text after "!rm". We use "!" at the beginning in case a user's name starts with "rm".
	if re.Match([]byte(message.Content)) {
		// Remove the user from the temple page

		if _, ok := s.submissions[member]; ok {
			delete(s.submissions, member)

			s.temple.RemoveMemberFromTemple(ctx, member, s.config.TempleGroupId, s.config.TempleGroupKey)
			logger.Info("Successfully removed user from temple group: " + member)
			msg := "You have successfully removed a member: " + member
			s.sendPrivateMessage(ctx, session, message.Author.ID, msg)

		} else {
			// Send the failed removal message in the previously created private channel
			logger.Error("Member: " + member + " does not exist.")
			msg := "Member: " + member + " does not exist."
			s.sendPrivateMessage(ctx, session, message.Author.ID, msg)
		}

		// Don't continue because the following code is to add a user
		return
	}

	// Ensure that this person does not exist in the submissions map currently
	if _, ok := s.submissions[member]; ok {
		// Send the failed addition message in the previously created private channel
		logger.Error("Member: " + member + " already exists.")
		msg := "Member: " + member + " already exists."
		s.sendPrivateMessage(ctx, session, message.Author.ID, msg)
	} else {
		s.submissions[member] = 0

		s.temple.AddMemberToTemple(ctx, member, s.config.TempleGroupId, s.config.TempleGroupKey)
		logger.Info("Successfully added new user to temple group: " + member)
		msg := "You have successfully added a new member: " + member
		s.sendPrivateMessage(ctx, session, message.Author.ID, msg)
	}
}

// kickOffCron will instantiate the HallOfFameRequestInfos and kick off the cron job
func (s *Service) kickOffCron(ctx context.Context, session *discordgo.Session) {
	s.log.Info("Initializing Hall Of Fame Cron Job...")
	slayerBosses := util.HallOfFameRequestInfo{Name: "Slayer Bosses", DiscChan: s.config.DiscSlayerBossesChan, Bosses: []util.BossInfo{
		{BossName: "sire", ImageLink: "https://i.imgur.com/GhbmqEB.png"},
		{BossName: "hydra", ImageLink: "https://i.imgur.com/25GU0Ph.png"},
		{BossName: "cerberus", ImageLink: "https://i.imgur.com/UoxGuQi.png"},
		{BossName: "grotesqueguardians", ImageLink: "https://i.imgur.com/M7ylVBZ.png"},
		{BossName: "kraken", ImageLink: "https://i.imgur.com/Q6EbJb1.png"},
		{BossName: "smokedevil", ImageLink: "https://i.imgur.com/2AYntQ5.png"},
	}}
	gwd := util.HallOfFameRequestInfo{Name: "GWD Bosses", DiscChan: s.config.DiscGwdChan, Bosses: []util.BossInfo{
		{BossName: "commanderzilyana", ImageLink: "https://i.imgur.com/aNm4Ydd.png"},
		{BossName: "kreearra", ImageLink: "https://i.imgur.com/lX8SfgN.png"},
		{BossName: "kriltsutsaroth", ImageLink: "https://i.imgur.com/hh8cMvp.png"},
		{BossName: "nex", ImageLink: "https://i.imgur.com/pqiVQBC.png"},
		{BossName: "generalgraardor", ImageLink: "https://i.imgur.com/hljv9ZW.png"},
	}}
	wildy := util.HallOfFameRequestInfo{Name: "Wildy Bosses", DiscChan: s.config.DiscWildyChan, Bosses: []util.BossInfo{
		{BossName: "artio", ImageLink: "https://i.imgur.com/bw6zLpU.png"},
		{BossName: "callisto", ImageLink: "https://i.imgur.com/bw6zLpU.png"},
		{BossName: "calvarion", ImageLink: "https://i.imgur.com/v3KX75y.png"},
		{BossName: "vetion", ImageLink: "https://i.imgur.com/v3KX75y.png"},
		{BossName: "spindel", ImageLink: "https://i.imgur.com/4zknWSX.png"},
		{BossName: "venenatis", ImageLink: "https://i.imgur.com/4zknWSX.png"},
		{BossName: "chaoselemental", ImageLink: "https://i.imgur.com/YAvIpbm.png"},
		{BossName: "chaosfanatic", ImageLink: "https://i.imgur.com/azV2sD1.png"},
		{BossName: "crazyarchaeologist", ImageLink: "https://i.imgur.com/23LXv53.png"},
		{BossName: "scorpia", ImageLink: "https://i.imgur.com/9aaguxB.png"},
	}}
	other := util.HallOfFameRequestInfo{Name: "Other Bosses", DiscChan: s.config.DiscOtherChan, Bosses: []util.BossInfo{
		{BossName: "corporealbeast", ImageLink: "https://i.imgur.com/zEDN4Pf.png"},
		{BossName: "prime", ImageLink: "https://i.imgur.com/kJBtqHB.png"},
		{BossName: "rexbro", ImageLink: "https://i.imgur.com/PvlGWFZ.png"},
		{BossName: "supreme", ImageLink: "https://i.imgur.com/BOgkBuD.png"},
		{BossName: "gauntlet", ImageLink: "https://i.imgur.com/weiHWnz.png"},
		{BossName: "gauntlethard", ImageLink: "https://i.imgur.com/xzW4TGR.png"},
		{BossName: "giantmole", ImageLink: "https://i.imgur.com/coKk2pr.png"},
		{BossName: "jad", ImageLink: "https://i.imgur.com/H9aO1Ot.png"},
		{BossName: "zuk", ImageLink: "https://i.imgur.com/mKstHza.png"},
		{BossName: "kq", ImageLink: "https://i.imgur.com/ZuaFoBR.png"},
		{BossName: "kbd", ImageLink: "https://i.imgur.com/r5vkw1s.png"},
		{BossName: "sarachnis", ImageLink: "https://i.imgur.com/98THH8O.png"},
		{BossName: "skotizo", ImageLink: "https://i.imgur.com/YUcQu4d.png"},
		{BossName: "muspah", ImageLink: "https://i.imgur.com/sW2cLQ2.png"},
		{BossName: "vorkath", ImageLink: "https://i.imgur.com/6biF3P2.png"},
		{BossName: "nightmare", ImageLink: "https://i.imgur.com/4aDkxms.png"},
		{BossName: "phosanis", ImageLink: "https://i.imgur.com/4aDkxms.png"},
		{BossName: "zulrah", ImageLink: "https://i.imgur.com/tPllWNF.png"},
	}}
	misc := util.HallOfFameRequestInfo{Name: "Miscellaneous Bosses", DiscChan: s.config.DiscMiscChan, Bosses: []util.BossInfo{
		{BossName: "barrows", ImageLink: "https://i.imgur.com/ajoK20v.png"},
		{BossName: "hespori", ImageLink: "https://i.imgur.com/b0qYGHS.png"},
		{BossName: "mimic", ImageLink: "https://i.imgur.com/jC7yTC3.png"},
		{BossName: "obor", ImageLink: "https://i.imgur.com/dwLvSbR.png"},
		{BossName: "bryophyta", ImageLink: "https://i.imgur.com/3cdyp4X.png"},
		{BossName: "derangedarchaeologist", ImageLink: "https://i.imgur.com/cnHpevF.png"},
		{BossName: "wintertodt", ImageLink: "https://i.imgur.com/6oFef2Y.png"},
		{BossName: "zalcano", ImageLink: "https://i.imgur.com/edN11Nf.png"},
		{BossName: "rift", ImageLink: "https://i.imgur.com/MOiyXeH.png"},
	}}
	dt2 := util.HallOfFameRequestInfo{Name: "Desert Treasure 2 Bosses", DiscChan: s.config.DiscDT2Chan, Bosses: []util.BossInfo{
		{BossName: "duke", ImageLink: "https://i.imgur.com/RYPmrXy.png"},
		{BossName: "leviathan", ImageLink: "https://i.imgur.com/mEQRq5c.png"},
		{BossName: "whisperer", ImageLink: "https://i.imgur.com/cFGWb6Y.png"},
		{BossName: "vardorvis", ImageLink: "https://i.imgur.com/WMPuShZ.png"},
	}}
	raids := util.HallOfFameRequestInfo{Name: "Raids Bosses", DiscChan: s.config.DiscRaidsChan, Bosses: []util.BossInfo{
		{BossName: "cox", ImageLink: "https://i.imgur.com/gxdWXtH.png"},
		{BossName: "coxcm", ImageLink: "https://i.imgur.com/gxdWXtH.png"},
		{BossName: "tob", ImageLink: "https://i.imgur.com/pW1sJAQ.png"},
		{BossName: "tobcm", ImageLink: "https://i.imgur.com/pW1sJAQ.png"},
		{BossName: "toa", ImageLink: "https://i.imgur.com/2GvzqGw.png"},
		{BossName: "toae", ImageLink: "https://i.imgur.com/2GvzqGw.png"},
	}}
	pvp := util.HallOfFameRequestInfo{Name: "PVP", DiscChan: s.config.DiscPVPChan, Bosses: []util.BossInfo{
		{BossName: "bhh", ImageLink: "https://i.imgur.com/zSQhlWk.png"},
		{BossName: "bhr", ImageLink: "https://i.imgur.com/Y3Sga7t.png"},
		{BossName: "lms", ImageLink: "https://i.imgur.com/rzW7ZXx.png"},
		{BossName: "arena", ImageLink: "https://i.imgur.com/uNP6Ggu.png"},
		{BossName: "zeal", ImageLink: "https://i.imgur.com/Ws7HvKL.png"},
	}}
	clues := util.HallOfFameRequestInfo{Name: "Clues", DiscChan: s.config.DiscCluesChan, Bosses: []util.BossInfo{
		{BossName: "cluebeginner", ImageLink: "https://i.imgur.com/fUmzJkW.png"},
		{BossName: "clueeasy", ImageLink: "https://i.imgur.com/phnSCHj.png"},
		{BossName: "cluemedium", ImageLink: "https://i.imgur.com/t5iH8Xa.png"},
		{BossName: "cluehard", ImageLink: "https://i.imgur.com/a0xwcGI.png"},
		{BossName: "clueelite", ImageLink: "https://i.imgur.com/ibNRk3G.png"},
		{BossName: "cluemaster", ImageLink: "https://i.imgur.com/12rCLVv.png"},
		{BossName: "clueall", ImageLink: "https://i.imgur.com/wX3Ei7U.png"},
	}}

	// Kick off a scheduled job at a configured time
	job, err := s.scheduler.Every(1).Day().At(s.config.CronKickoffTime).Do(func() {
		s.log.Debug("Running Cron Job to update the Hall Of Fame, Collection Log, and Leagues...")
		s.updateHOF(ctx, session, slayerBosses, gwd, wildy, other, misc, dt2, raids, pvp, clues)
		s.updateColLog(ctx, session)
		s.updateLeagues(ctx, session)
		s.log.Debug("Finished running Cron Job to update the Hall Of Fame, Collection Log, and Leagues!")
	})
	if err != nil {
		// handle the error related to setting up the job
		s.log.Error(fmt.Sprintf("Error creating cron job. Job: %#v, Error: %#v", job, err))
		return
	}
	job.SingletonMode()
	s.scheduler.StartAsync()

	return
}

func (s *Service) sendPrivateMessage(ctx context.Context, session *discordgo.Session, userId string, message string) {
	logger := flume.FromContext(ctx)
	// Create a private channel with the user submitting (will reuse if one exists)
	channel, err := session.UserChannelCreate(userId)
	if err != nil {
		logger.Error("Failed to create private message with the user: " + err.Error())
		return
	}

	// Send a message on that channel
	_, err = session.ChannelMessageSend(channel.ID, message)
	if err != nil {
		logger.Error("Failed to send private message to the user: " + err.Error())
		return
	}
}
