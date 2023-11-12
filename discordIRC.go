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

	for _, submissionPicture := range message.Attachments {
		_, err := session.ChannelMessageSend(message.ChannelID, submissionPicture.ProxyURL)
		if err != nil {
			return
		}

		downloadSubmissionScreenshot(submissionPicture.ProxyURL)
	}

	// If hello is read, world is responded
	if message.Content == "hello" {
		_, err := session.ChannelMessageSend(message.ChannelID, "hello")
		if err != nil {
			return
		}
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
