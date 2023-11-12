package main

import (
	"fmt"
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

	leaderboard := ""
	for _, k := range keys {
		leaderboard = leaderboard + k + ":" + strconv.Itoa(submissions[k]) + "\n"
	}
	_, err := session.ChannelMessageSend(config.DiscLeaderboardChan, leaderboard)
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

	// Ensure that this person does not exist in the submissions map currently
	if _, ok := submissions[newMember]; !ok {
		submissions[newMember] = 0
		_, err := session.ChannelMessageSend(message.ChannelID, "Successfully added new member: "+newMember)
		if err != nil {
			return
		}
	} else {
		_, err := session.ChannelMessageSend(message.ChannelID, "Member: "+newMember+" already exists.")
		if err != nil {
			return
		}
	}
}
