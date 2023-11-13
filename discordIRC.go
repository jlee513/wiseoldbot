package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"

	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
)

// Use of discordgo as an intro to discord IRC
func startDiscordIRC() {
	// Create a new discord session
	session, err := discordgo.New("Bot " + config.DiscBotToken)
	if err != nil {
		log.Fatal(err)
	}

	// Create handler for listening for submission messages
	session.AddHandler(listenForMessage)

	// Send intent
	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	_ = session.Open()
	defer func(discord *discordgo.Session) {
		err := discord.Close()
		if err != nil {

		}
	}(session)

	// Initialize the Hall Of fame
	//kickOffHallOfFameUpdate(session)
	fmt.Println("the bot is online!")

	// Block so that it continues to run the bot
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func listenForMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Don't handle message if it's created by the discord bot
	if message.Author.ID == session.State.User.ID {
		return
	}

	// Run certain tasks depending on the channel the message was posted in
	switch channel := message.ChannelID; channel {
	case config.DiscSubChan:
		listenForSubmission(session, message)
	case config.DiscSignUpChan:
		UpdateMemberList(session, message)
	default:
		// Return if the message was not posted in one of the channels we are handling
		return
	}
}

func listenForSubmission(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Split the names into an array by , then make an empty array with those names as keys for an easier lookup
	// instead of running a for loop inside a for loop when adding clan points
	whitespaceStrippedMessage := strings.Replace(message.Content, ", ", ",", -1)
	whitespaceStrippedMessage = strings.Replace(whitespaceStrippedMessage, " ,", ",", -1)

	names := strings.Split(whitespaceStrippedMessage, ",")

	// Before adding clanpoints, ensure that all the names used in the submission is valid and already created
	// in the #ponies-signup channel
	for _, name := range names {
		// Ensure that this person does not exist in the submissions map currently
		if _, ok := submissions[name]; !ok {
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
			err = session.ChannelMessageDelete(config.DiscSubChan, message.ID)
			if err != nil {
				return
			}

			return
		}
	}

	numberOfSubmissions := 0

	// Iterate through all the pictures and download them
	for _, submissionPicture := range message.Attachments {
		downloadSubmissionScreenshot(submissionPicture.ProxyURL)
		numberOfSubmissions++
	}

	// Iterate over the all the names in the submissions and add the number of submissions to their clan points
	for _, name := range names {
		submissions[name] = submissions[name] + numberOfSubmissions
	}

	// Update the #cp-leaderboard
	updateLeaderboard(session)

	// Once everything is finished, delete the message from the submission channel
	err := session.ChannelMessageDelete(config.DiscSubChan, message.ID)
	if err != nil {
		return
	}
}

// updateLeaderboard will update the cp-leaderboard channel in discord with a new ranking of everyone in the clan
func updateLeaderboard(session *discordgo.Session) {
	// Update the #cp-leaderboard
	keys := make([]string, 0, len(submissions))
	for key := range submissions {
		keys = append(keys, key)
	}

	// Sort the map based on the values
	sort.SliceStable(keys, func(i, j int) bool {
		return submissions[keys[i]] > submissions[keys[j]]
	})

	// Create the leaderboard message that will be sent
	leaderboard := ""
	for placement, k := range keys {
		leaderboard = leaderboard + strconv.Itoa(placement+1) + ") " + k + ": " + strconv.Itoa(submissions[k]) + "\n"
	}

	// Retrieve the one channel message and delete it in the leaderboard channel
	messages, err := session.ChannelMessages(config.DiscLeaderboardChan, 1, "", "", "")
	if err != nil {
		return
	}
	err = session.ChannelMessageDelete(config.DiscLeaderboardChan, messages[0].ID)
	if err != nil {
		return
	}

	_, err = session.ChannelMessageSendEmbed(config.DiscLeaderboardChan, embed.NewEmbed().
		SetTitle("Ponies Clan Points Leaderboard").
		SetDescription(fmt.Sprintf(leaderboard)).
		SetColor(0x1c1c1c).SetThumbnail("https://i.imgur.com/O4NzB95.png").MessageEmbed)
	if err != nil {
		return
	}
}

func downloadSubmissionScreenshot(submissionLink string) {
	// Build fileName from fullPath
	fileURL, err := url.Parse(submissionLink)
	if err != nil {
		log.Fatal(err)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	// Create blank file
	file, err := os.Create("submissions/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	resp, err := client.Get(submissionLink)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)

	defer file.Close()
}

func UpdateMemberList(session *discordgo.Session, message *discordgo.MessageCreate) {
	//TODO: SCRUB THE USERNAME SUBMITTED
	// Don't include the remove command in the RSN
	newMember := strings.Replace(message.Content, "!rm ", "", -1)

	// Create a private channel with the user submitting (will reuse if one exists)
	channel, err := session.UserChannelCreate(message.Author.ID)
	if err != nil {
		return
	}

	// Remove user from temple if the message prefix is "rm"
	re := regexp.MustCompile("(?i)^(!rm)\\s+.+$") // Case insensitive. Must start with "!rm". Must have atleast one space between "!rm" and the username. There must be text after "!rm". We use "!" at the beginning in case a user's name starts with "rm".
	if re.Match([]byte(message.Content)) {
		// Remove the user from the temple page
		removeNewMemberToTemple(newMember)

		if userExists(session, newMember, message.ChannelID) {
			submissions[newMember] = 0

			// Send a message on that channel
			_, err := session.ChannelMessageSend(channel.ID, "You have successfully removed a member: "+newMember)
			if err != nil {
				return
			}
		} else {
			_, err := session.ChannelMessageSend(channel.ID, "Member: "+newMember+" does not exist.")
			if err != nil {
				return
			}
		}

		// Once everything is finished, delete the message from the submission channel
		err = session.ChannelMessageDelete(config.DiscSignUpChan, message.ID)
		if err != nil {
			return
		}

		// Don't continue because the following code is to add a user
		return
	}

	// Ensure that this person does not exist in the submissions map currently
	if userExists(session, newMember, message.ChannelID) {
		_, err := session.ChannelMessageSend(channel.ID, "Member: "+newMember+" already exists.")
		if err != nil {
			return
		}
	} else {
		submissions[newMember] = 0

		// Send a message on that channel
		_, err := session.ChannelMessageSend(channel.ID, "You have successfully added new member: "+newMember)
		if err != nil {
			return
		}
	}

	// Add the user to the temple page
	addNewMemberToTemple(newMember)

	// Once everything is finished, delete the message from the submission channel
	err = session.ChannelMessageDelete(config.DiscSignUpChan, message.ID)
	if err != nil {
		return
	}
}

func addNewMemberToTemple(newMember string) {
	url := "https://templeosrs.com/api/add_group_member.php"
	method := "POST"

	payload := strings.NewReader("id=" + config.TempleGroupId + "&key=" + config.TempleGroupKey + "&players=" + newMember)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func kickOffHallOfFameUpdate(session *discordgo.Session) {
	slayerBosses := hallOfFameRequestInfo{Bosses: map[string]string{"sire": "https://i.imgur.com/GhbmqEB.png", "hydra": "https://i.imgur.com/25GU0Ph.png", "cerberus": "https://i.imgur.com/UoxGuQi.png", "grotesqueguardians": "https://i.imgur.com/M7ylVBZ.png", "kraken": "https://i.imgur.com/Q6EbJb1.png", "smokedevil": "https://i.imgur.com/2AYntQ5.png"}, DiscChan: config.DiscSlayerBossesChan}
	gwd := hallOfFameRequestInfo{Bosses: map[string]string{"commanderzilyana": "https://i.imgur.com/aNm4Ydd.png", "kreearra": "https://i.imgur.com/lX8SfgN.png", "kriltsutsaroth": "https://i.imgur.com/hh8cMvp.png", "nex": "https://i.imgur.com/pqiVQBC.png", "generalgraardor": "https://i.imgur.com/hljv9ZW.png"}, DiscChan: config.DiscGwdChan}
	wildy := hallOfFameRequestInfo{Bosses: map[string]string{"artio": "https://i.imgur.com/bw6zLpU.png", "callisto": "https://i.imgur.com/bw6zLpU.png", "calvarion": "https://i.imgur.com/v3KX75y.png", "vetion": "https://i.imgur.com/v3KX75y.png", "spindel": "https://i.imgur.com/4zknWSX.png", "venenatis": "https://i.imgur.com/4zknWSX.png", "chaoselemental": "https://i.imgur.com/YAvIpbm.png", "chaosfanatic": "https://i.imgur.com/azV2sD1.png", "crazyarchaeologist": "https://i.imgur.com/23LXv53.png", "scorpia": "https://i.imgur.com/9aaguxB.png"}, DiscChan: config.DiscWildyChan}
	other := hallOfFameRequestInfo{Bosses: map[string]string{"corporealbeast": "https://i.imgur.com/zEDN4Pf.png", "prime": "https://i.imgur.com/kJBtqHB.png", "rexbro": "https://i.imgur.com/PvlGWFZ.png", "supreme": "https://i.imgur.com/BOgkBuD.png", "gauntlet": "https://i.imgur.com/weiHWnz.png", "gauntlethard": "https://i.imgur.com/xzW4TGR.png", "giantmole": "https://i.imgur.com/coKk2pr.gif", "jad": "https://i.imgur.com/H9aO1Ot.png", "zuk": "https://i.imgur.com/mKstHza.png", "kq": "https://i.imgur.com/ZuaFoBR.png", "kbd": "https://i.imgur.com/r5vkw1s.png", "sarachnis": "https://i.imgur.com/98THH8O.png", "skotizo": "https://i.imgur.com/YUcQu4d.png", "muspah": "https://i.imgur.com/sW2cLQ2.png", "vorkath": "https://i.imgur.com/6biF3P2.png", "phosanis": "https://i.imgur.com/4aDkxms.png", "nightmare": "https://i.imgur.com/4aDkxms.png", "zulrah": "https://i.imgur.com/tPllWNF.png"}, DiscChan: config.DiscOtherChan}
	misc := hallOfFameRequestInfo{Bosses: map[string]string{"barrows": "https://i.imgur.com/ajoK20v.png", "hespori": "https://i.imgur.com/b0qYGHS.png", "mimic": "https://i.imgur.com/jC7yTC3.png", "obor": "https://i.imgur.com/dwLvSbR.png", "bryophyta": "https://i.imgur.com/3cdyp4X.png", "derangedarchaeologist": "https://i.imgur.com/cnHpevF.png", "wintertodt": "https://i.imgur.com/6oFef2Y.png", "zalcano": "https://i.imgur.com/edN11Nf.png", "tempoross": "https://i.imgur.com/fRj3JA4.png", "rift": "https://i.imgur.com/MOiyXeH.png"}, DiscChan: config.DiscMiscChan}
	dt2 := hallOfFameRequestInfo{Bosses: map[string]string{"duke": "https://i.imgur.com/RYPmrXy.png", "leviathan": "https://i.imgur.com/mEQRq5c.png", "whisperer": "https://i.imgur.com/cFGWb6Y.png", "vardorvis": "https://i.imgur.com/WMPuShZ.png"}, DiscChan: config.DiscDT2Chan}
	raids := hallOfFameRequestInfo{Bosses: map[string]string{"cox": "https://i.imgur.com/gxdWXtH.png", "coxcm": "https://i.imgur.com/gxdWXtH.png", "tob": "https://i.imgur.com/pW1sJAQ.png", "tobcm": "https://i.imgur.com/pW1sJAQ.png", "toa": "https://i.imgur.com/2GvzqGw.png", "toae": "https://i.imgur.com/2GvzqGw.png"}, DiscChan: config.DiscRaidsChan}
	pvp := hallOfFameRequestInfo{Bosses: map[string]string{"bhh": "https://i.imgur.com/zSQhlWk.png", "bhr": "https://i.imgur.com/Y3Sga7t.png", "lms": "https://i.imgur.com/rzW7ZXx.png", "arena": "https://i.imgur.com/uNP6Ggu.png", "zeal": "https://i.imgur.com/Ws7HvKL.png"}, DiscChan: config.DiscPVPChan}
	clues := hallOfFameRequestInfo{Bosses: map[string]string{"clueall": "https://i.imgur.com/wX3Ei7U.png", "cluebeginner": "https://i.imgur.com/fUmzJkW.png", "clueeasy": "https://i.imgur.com/phnSCHj.png", "cluemedium": "https://i.imgur.com/t5iH8Xa.png", "cluehard": "https://i.imgur.com/a0xwcGI.png", "clueelite": "https://i.imgur.com/ibNRk3G.png", "cluemaster": "https://i.imgur.com/12rCLVv.png"}, DiscChan: config.DiscCluesChan}

	updateHallOfFame(session, slayerBosses)
	updateHallOfFame(session, gwd)
	updateHallOfFame(session, wildy)
	updateHallOfFame(session, other)
	updateHallOfFame(session, misc)
	updateHallOfFame(session, dt2)
	updateHallOfFame(session, raids)
	updateHallOfFame(session, pvp)
	updateHallOfFame(session, clues)
	updateCollectionLog(session)
}

func removeNewMemberToTemple(newMember string) {
	url := "https://templeosrs.com/api/remove_group_member.php"
	method := "POST"

	payload := strings.NewReader("id=" + config.TempleGroupId + "&key=" + config.TempleGroupKey + "&players=" + newMember)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func userExists(session *discordgo.Session, member string, channelID string) (exists bool) {
	exists = false

	if _, ok := submissions[member]; ok {
		exists = true
	}

	return
}
