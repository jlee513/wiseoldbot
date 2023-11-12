package main

import (
	"fmt"
	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

// Use of discordgo as an intro to discord IRC
func startDiscordIRC() {
	// Create a new discord session
	session, err := discordgo.New("Bot " + config.DiscBotToken)
	if err != nil {
		log.Fatal(err)
	}

	slayerBosses := hallOfFameRequestInfo{
		Bosses: map[string]string{
			"sire":               "https://i.imgur.com/O4NzB95.png",
			"hydra":              "https://i.imgur.com/O4NzB95.png",
			"cerberus":           "https://i.imgur.com/O4NzB95.png",
			"grotesqueguardians": "https://i.imgur.com/O4NzB95.png",
			"kraken":             "https://i.imgur.com/O4NzB95.png",
			"smokedevil":         "https://i.imgur.com/O4NzB95.png",
		},
		DiscChan: config.DiscSlayerBossesChan,
	}

	gwd := hallOfFameRequestInfo{
		Bosses: map[string]string{
			"commanderzilyana": "https://i.imgur.com/O4NzB95.png",
			"kreearra":         "https://i.imgur.com/O4NzB95.png",
			"kriltsutsaroth":   "https://i.imgur.com/O4NzB95.png",
			"nex":              "https://i.imgur.com/O4NzB95.png",
			"generalgraardor":  "https://i.imgur.com/O4NzB95.png",
		},
		DiscChan: config.DiscGwdChan,
	}

	wildy := hallOfFameRequestInfo{
		Bosses: map[string]string{
			"artio":              "https://i.imgur.com/O4NzB95.png",
			"callisto":           "https://i.imgur.com/O4NzB95.png",
			"calvarion":          "https://i.imgur.com/O4NzB95.png",
			"vetion":             "https://i.imgur.com/O4NzB95.png",
			"spindel":            "https://i.imgur.com/O4NzB95.png",
			"venenatis":          "https://i.imgur.com/O4NzB95.png",
			"chaoselemental":     "https://i.imgur.com/O4NzB95.png",
			"chaosfanatic":       "https://i.imgur.com/O4NzB95.png",
			"crazyarchaeologist": "https://i.imgur.com/O4NzB95.png",
		},
		DiscChan: config.DiscWildyChan,
	}

	other := hallOfFameRequestInfo{
		Bosses: map[string]string{
			"corporealbeast": "https://i.imgur.com/O4NzB95.png",
			"prime":          "https://i.imgur.com/O4NzB95.png",
			"rexbro":         "https://i.imgur.com/O4NzB95.png",
			"supreme":        "https://i.imgur.com/O4NzB95.png",
			"gauntlet":       "https://i.imgur.com/O4NzB95.png",
			"gauntlethard":   "https://i.imgur.com/O4NzB95.png",
			"giantmole":      "https://i.imgur.com/EO9HfXe.gif",
			"jad":            "https://i.imgur.com/O4NzB95.png",
			"zuk":            "https://i.imgur.com/O4NzB95.png",
			"kq":             "https://i.imgur.com/O4NzB95.png",
			"kbd":            "https://i.imgur.com/O4NzB95.png",
			"sarachnis":      "https://i.imgur.com/O4NzB95.png",
			"skotizo":        "https://i.imgur.com/O4NzB95.png",
			"muspah":         "https://i.imgur.com/O4NzB95.png",
			"vorkath":        "https://i.imgur.com/O4NzB95.png",
			"phosanis":       "https://i.imgur.com/O4NzB95.png",
			"nightmare":      "https://i.imgur.com/O4NzB95.png",
			"zulrah":         "https://i.imgur.com/O4NzB95.png",
		},
		DiscChan: config.DiscOtherChan,
	}

	misc := hallOfFameRequestInfo{
		Bosses: map[string]string{
			"barrows":               "https://i.imgur.com/O4NzB95.png",
			"hespori":               "https://i.imgur.com/O4NzB95.png",
			"mimic":                 "https://i.imgur.com/O4NzB95.png",
			"obor":                  "https://i.imgur.com/O4NzB95.png",
			"bryophyta":             "https://i.imgur.com/O4NzB95.png",
			"derangedarchaeologist": "https://i.imgur.com/O4NzB95.png",
			"wintertodt":            "https://i.imgur.com/O4NzB95.png",
			"zalcano":               "https://i.imgur.com/O4NzB95.png",
			"tempoross":             "https://i.imgur.com/O4NzB95.png",
			"rift":                  "https://i.imgur.com/O4NzB95.png",
		},
		DiscChan: config.DiscMiscChan,
	}

	dt2 := hallOfFameRequestInfo{
		Bosses: map[string]string{
			"duke":      "https://i.imgur.com/O4NzB95.png",
			"leviathan": "https://i.imgur.com/O4NzB95.png",
			"whisperer": "https://i.imgur.com/O4NzB95.png",
			"vardorvis": "https://i.imgur.com/O4NzB95.png",
		},
		DiscChan: config.DiscDT2Chan,
	}

	raids := hallOfFameRequestInfo{
		Bosses: map[string]string{
			"cox":   "https://i.imgur.com/O4NzB95.png",
			"coxcm": "https://i.imgur.com/O4NzB95.png",
			"tob":   "https://i.imgur.com/O4NzB95.png",
			"tobcm": "https://i.imgur.com/O4NzB95.png",
			"toa":   "https://i.imgur.com/O4NzB95.png",
			"toae":  "https://i.imgur.com/O4NzB95.png",
		},
		DiscChan: config.DiscRaidsChan,
	}

	pvp := hallOfFameRequestInfo{
		Bosses: map[string]string{
			"bhh":   "https://i.imgur.com/O4NzB95.png",
			"bhr":   "https://i.imgur.com/O4NzB95.png",
			"lms":   "https://i.imgur.com/O4NzB95.png",
			"arena": "https://i.imgur.com/O4NzB95.png",
			"zeal":  "https://i.imgur.com/O4NzB95.png",
		},
		DiscChan: config.DiscPVPChan,
	}

	clues := hallOfFameRequestInfo{
		Bosses: map[string]string{
			"clueall":      "https://i.imgur.com/O4NzB95.png",
			"cluebeginner": "https://i.imgur.com/O4NzB95.png",
			"clueeasy":     "https://i.imgur.com/O4NzB95.png",
			"cluemedium":   "https://i.imgur.com/O4NzB95.png",
			"cluehard":     "https://i.imgur.com/O4NzB95.png",
			"clueelite":    "https://i.imgur.com/O4NzB95.png",
			"cluemaster":   "https://i.imgur.com/O4NzB95.png",
		},
		DiscChan: config.DiscCluesChan,
	}

	updateHallOfFame(session, slayerBosses)
	updateHallOfFame(session, gwd)
	updateHallOfFame(session, wildy)
	updateHallOfFame(session, other)
	updateHallOfFame(session, misc)
	updateHallOfFame(session, dt2)
	updateHallOfFame(session, raids)
	updateHallOfFame(session, pvp)
	updateHallOfFame(session, clues)

	// Create handler for listening for submission messages
	session.AddHandler(listenForSubmission)
	session.AddHandler(AddNewMember)

	// Send intent
	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	_ = session.Open()
	defer func(discord *discordgo.Session) {
		err := discord.Close()
		if err != nil {

		}
	}(session)

	fmt.Println("the bot is online!")

	// Block so that it continues to run the bot
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func listenForSubmission(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Don't handle message if it's created by the discord bot
	// Also, don't handle messages other than ones send in the submission channel
	if message.Author.ID == session.State.User.ID || message.ChannelID != config.DiscSubChan {
		return
	}

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

func AddNewMember(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Don't handle message if it's created by the discord bot
	// Also, don't handle messages other than ones send in the submission channel
	if message.Author.ID == session.State.User.ID || message.ChannelID != config.DiscSignUpChan {
		return
	}

	//TODO: SCRUB THE USERNAME SUBMITTED
	newMember := message.Content

	// Create a private channel with the user submitting (will reuse if one exists)
	channel, err := session.UserChannelCreate(message.Author.ID)
	if err != nil {
		return
	}

	// Ensure that this person does not exist in the submissions map currently
	if _, ok := submissions[newMember]; !ok {
		submissions[newMember] = 0
		// Send a message on that channel
		_, err = session.ChannelMessageSend(channel.ID, "You have successfully added new member: "+newMember)
		if err != nil {
			return
		}
	} else {
		_, err = session.ChannelMessageSend(channel.ID, "Member: "+newMember+" already exists.")
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
