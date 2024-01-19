package util

import (
	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
	"time"
)

func AppendToHofSpeedArr(hofName string) []*discordgo.ApplicationCommandOptionChoice {
	var choices []*discordgo.ApplicationCommandOptionChoice
	switch hofName {
	case "TzHaar":
		for _, option := range HofSpeedTzhaar {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  option.BossName,
				Value: option.BossName,
			})
		}
	case "Chambers Of Xeric":
		for _, option := range HofSpeedCox {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  option.BossName,
				Value: option.BossName,
			})
		}
	case "Chambers Of Xeric Challenge Mode":
		for _, option := range HofSpeedCoxCm {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  option.BossName,
				Value: option.BossName,
			})
		}
	case "Nightmare":
		for _, option := range HofSpeedNightmare {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  option.BossName,
				Value: option.BossName,
			})
		}
	case "Theatre Of Blood":
		for _, option := range HofSpeedTob {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  option.BossName,
				Value: option.BossName,
			})
		}
	case "Theatre Of Blood Hard Mode":
		for _, option := range HofSpeedTobHm {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  option.BossName,
				Value: option.BossName,
			})
		}
	case "Agility":
		for _, option := range HofSpeedAgility {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  option.BossName,
				Value: option.BossName,
			})
		}
	case "Tombs Of Amascut":
		for _, option := range HofSpeedToa {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  option.BossName,
				Value: option.BossName,
			})
		}
	case "Tombs Of Amascut Expert":
		for _, option := range HofSpeedToae {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  option.BossName,
				Value: option.BossName,
			})
		}
	case "Solo Bosses":
		for _, option := range HofSpeedSolo {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  option.BossName,
				Value: option.BossName,
			})
		}
	case "Nex":
		for _, option := range HofSpeedNex {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  option.BossName,
				Value: option.BossName,
			})
		}
	case "Slayer":
		for _, option := range HofSpeedSlayer {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  option.BossName,
				Value: option.BossName,
			})
		}
	case "Desert Treasure 2":
		for _, option := range HofSpeedDt2 {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  option.BossName,
				Value: option.BossName,
			})
		}
	}

	return choices
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

func DeleteBulkDiscordMessages(session *discordgo.Session, channel string, afterId string) error {
	messages, err := session.ChannelMessages(channel, 100, "", afterId, "")
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
	return nil
}

func SendDiscordEmbedMsg(session *discordgo.Session, channel, title, description, thumbnailUrl string) error {
	// Send the Discord Embed message for the boss podium finish
	_, err := session.ChannelMessageSendEmbed(channel, embed.NewEmbed().
		SetTitle(title).
		SetDescription(description).
		SetColor(0x1c1c1c).SetThumbnail(thumbnailUrl).MessageEmbed)
	return err
}
