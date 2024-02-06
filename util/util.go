package util

import (
	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
	"github.com/gemalto/flume"
	"strconv"
	"strings"
	"time"
)

func LogError(logger flume.Logger, auditChannel string, session *discordgo.Session, user string, userUrl string, errorMessage string) {
	err := SendDiscordEmbedMsg(session, auditChannel, "Error: "+user, errorMessage, userUrl)
	if err != nil {
		logger.Error("ERROR SENDING MESSAGE: " + errorMessage + " TO AUDIT CHANNEL - " + err.Error())
		return
	}
	logger.Error(errorMessage)
}

func LogAdminAction(logger flume.Logger, auditChannel string, admin string, adminUrl string, session *discordgo.Session, adminMessage string) {
	err := SendDiscordEmbedMsg(session, auditChannel, "Admin Action: "+admin, adminMessage, adminUrl)
	if err != nil {
		logger.Error("ERROR SENDING MESSAGE: " + adminMessage + " TO AUDIT CHANNEL - " + err.Error())
		return
	}
	logger.Info(adminMessage)
}

func CalculateTime(speedTime string) time.Time {
	var t time.Time
	speedTimeSplit := strings.Split(speedTime, ":")

	for i, splitTime := range speedTimeSplit {
		switch i {
		case 0:
			c, _ := strconv.Atoi(splitTime)
			t = t.Add(time.Duration(c) * time.Hour)
		case 1:
			c, _ := strconv.Atoi(splitTime)
			t = t.Add(time.Duration(c) * time.Minute)
		case 2:
			if strings.Contains(splitTime, ".") {
				milliAndSeconds := strings.Split(splitTime, ".")
				c, _ := strconv.Atoi(milliAndSeconds[0])
				c2, _ := strconv.Atoi(milliAndSeconds[1])
				t = t.Add(time.Duration(c) * time.Second)
				t = t.Add(time.Duration(c2) * time.Millisecond * 10)
			} else {
				c, _ := strconv.Atoi(splitTime)
				t = t.Add(time.Duration(c) * time.Second)
			}
		}
	}

	return t
}

func WhiteStripCommas(msg string) string {
	whitespaceStrippedMessage := strings.Replace(msg, ", ", ",", -1)
	whitespaceStrippedMessage = strings.Replace(whitespaceStrippedMessage, " ,", ",", -1)

	return whitespaceStrippedMessage
}

func InteractionRespond(session *discordgo.Session, i *discordgo.InteractionCreate, content string) error {
	if len(content) == 0 {
		content = "Generic response message"
	}
	return session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func InteractionRespondChoices(session *discordgo.Session, i *discordgo.InteractionCreate, choices []*discordgo.ApplicationCommandOptionChoice) error {
	return session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

func DeleteBulkDiscordMessages(session *discordgo.Session, channel string) error {
	messages, err := session.ChannelMessages(channel, 100, "", channel, "")
	if err != nil {
		return err
	}
	var messageIDs []string
	for _, message := range messages {
		messageIDs = append(messageIDs, message.ID)
	}

	if len(messageIDs) > 0 {
		err = session.ChannelMessagesBulkDelete(channel, messageIDs)
		if err != nil {
			for _, message := range messageIDs {
				err = session.ChannelMessageDelete(channel, message)
				if err != nil {
					return err
				}
			}
		}
	}

	// There are some guides that are larger than 100 messages - run it again just in case

	return nil
}

func SendDiscordEmbedMsg(session *discordgo.Session, channel, title, description, thumbnailUrl string) error {
	// Send the Discord Embed message for the boss podium finish
	if len(thumbnailUrl) > 0 {
		_, err := session.ChannelMessageSendEmbed(channel, embed.NewEmbed().
			SetTitle(title).
			SetDescription(description).
			SetColor(0x1c1c1c).
			SetThumbnail(thumbnailUrl).MessageEmbed)
		return err
	} else {
		_, err := session.ChannelMessageSendEmbed(channel, embed.NewEmbed().
			SetTitle(title).
			SetDescription(description).
			SetColor(0x1c1c1c).MessageEmbed)
		return err
	}
}
