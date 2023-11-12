package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Use of discordgo as an intro to discord IRC
func startDiscordIRC() {
	// Create a new discord session
	session, err := discordgo.New("Bot " + config.DiscBotToken)
	if err != nil {
		log.Fatal(err)
	}

	// Create handler for listening for discord messages
	session.AddHandler(listenForMessage)

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

func listenForMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Don't handle message if it's created by the discord bot
	if message.Author.ID == session.State.User.ID || message.ChannelID != config.DiscSubChan {
		return
	}

	// If hello is read, world is responded
	if message.Content == "hello" {
		_, err := session.ChannelMessageSend(message.ChannelID, "hello")
		if err != nil {
			return
		}
	}
}
