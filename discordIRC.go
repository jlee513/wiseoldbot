package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

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
		listenForCPSubmission(session, message)
	case config.DiscSignUpChan:
		UpdateMemberList(session, message)
	default:
		// Return if the message was not posted in one of the channels we are handling
		return
	}
}

func listenForCPSubmission(session *discordgo.Session, message *discordgo.MessageCreate) {
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
		// If it's an imgur link, save the link in the cpscreenshots map
		if strings.Contains(submissionPicture.ProxyURL, "imgur") {
			cpscreenshots[submissionPicture.ProxyURL] = whitespaceStrippedMessage
		} else if strings.Contains(submissionPicture.ProxyURL, "media.discordapp.net") {
			submissionUrl := uploadToImgur(submissionPicture.ProxyURL)
			cpscreenshots[submissionUrl] = whitespaceStrippedMessage
		} else {
		}
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

func userExists(session *discordgo.Session, member string, channelID string) (exists bool) {
	exists = false

	if _, ok := submissions[member]; ok {
		exists = true
	}

	return
}
