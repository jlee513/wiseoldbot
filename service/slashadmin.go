package service

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/gemalto/flume"
	"osrs-disc-bot/util"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (s *Service) handleAdmin(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	returnMessage := ""

	switch options[0].Name {
	case "player":
		returnMessage = s.handlePlayerAdministration(ctx, session, i)
		s.updatePpLeaderboard(ctx, session)
	case "submission-instructions":
		session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Updating Ponies Point Instructions",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		_ = s.updateSubmissionInstructions(ctx, session)
	case "update-cp":
		returnMessage = s.updateCpPoints(ctx, session, i)
	case "update-speed":
		returnMessage = s.updateSpeedAdmin(ctx, session, i)
	case "update-leaderboard":
		s.updateLeaderboard(ctx, session, i)
		return
	}

	session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: returnMessage,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	return
}

func (s *Service) updateSpeedAdmin(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) string {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		logger := flume.FromContext(ctx)
		options := i.ApplicationCommandData().Options[0].Options

		category := ""
		boss := ""

		for _, option := range options {
			switch option.Name {
			case "category":
				category = option.Value.(string)
			case "boss":
				boss = option.Value.(string)
			}
		}

		logger.Info("Resetting speed for: " + boss)

		// Ensure the boss name is okay
		if _, ok := util.SpeedBossNameToCategory[boss]; !ok {
			logger.Error("Incorrect boss name: ", boss)
			return "Incorrect boss name. Please look ensure to select one of the options for boss names."
		}

		// Convert the time string into time
		var t time.Time
		speedTimeSplit := []string{"22", "22", "22.60"}

		for index, splitTime := range speedTimeSplit {
			switch index {
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
					t = t.Add(time.Duration(c2) * time.Millisecond)
				} else {
					c, _ := strconv.Atoi(splitTime)
					t = t.Add(time.Duration(c) * time.Second)
				}
			}
		}

		s.speed[boss] = util.SpeedInfo{Time: t, PlayersInvolved: "null", URL: "https://i.imgur.com/34dg0da.png"}
		s.updateSpeedHOF(ctx, session, category)

		// If nothing wrong happened, send a happy message back to the submitter
		return "Speed submission successfully submitted! Awaiting approval from a moderator!"

		return ""
	case discordgo.InteractionApplicationCommandAutocomplete:
		logger := flume.FromContext(ctx)
		data := i.ApplicationCommandData().Options[0]
		var choices []*discordgo.ApplicationCommandOptionChoice
		switch {
		// In this case there are multiple autocomplete options. The Focused field shows which option user is focused on.
		case data.Options[0].Focused:
			choices = []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "TzHaar",
					Value: "TzHaar",
				},
				{
					Name:  "Chambers Of Xeric",
					Value: "Chambers Of Xeric",
				},
				{
					Name:  "Chambers Of Xeric Challenge Mode",
					Value: "Chambers Of Xeric Challenge Mode",
				},
				{
					Name:  "Nightmare",
					Value: "Nightmare",
				},
				{
					Name:  "Theatre Of Blood Hard Mode",
					Value: "Theatre Of Blood Hard Mode",
				},
				{
					Name:  "Agility",
					Value: "Agility",
				},
				{
					Name:  "Tombs Of Amascut Expert",
					Value: "Tombs Of Amascut Expert",
				},
				{
					Name:  "Solo Bosses",
					Value: "Solo Bosses",
				},
				{
					Name:  "Nex",
					Value: "Nex",
				},
				{
					Name:  "Slayer",
					Value: "Slayer",
				},
			}
		case data.Options[1].Focused:
			switch data.Options[0].Value.(string) {
			case "TzHaar":
				for _, option := range util.HofSpeedTzhaar {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  option.BossName,
						Value: option.BossName,
					})
				}
			case "Chambers Of Xeric":
				for _, option := range util.HofSpeedCox {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  option.BossName,
						Value: option.BossName,
					})
				}
			case "Chambers Of Xeric Challenge Mode":
				for _, option := range util.HofSpeedCoxCm {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  option.BossName,
						Value: option.BossName,
					})
				}
			case "Nightmare":
				for _, option := range util.HofSpeedNightmare {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  option.BossName,
						Value: option.BossName,
					})
				}
			case "Theatre Of Blood":
				for _, option := range util.HofSpeedTob {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  option.BossName,
						Value: option.BossName,
					})
				}
			case "Theatre Of Blood Hard Mode":
				for _, option := range util.HofSpeedTobHm {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  option.BossName,
						Value: option.BossName,
					})
				}
			case "Agility":
				for _, option := range util.HofSpeedAgility {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  option.BossName,
						Value: option.BossName,
					})
				}
			case "Tombs Of Amascut":
				for _, option := range util.HofSpeedToa {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  option.BossName,
						Value: option.BossName,
					})
				}
			case "Tombs Of Amascut Expert":
				for _, option := range util.HofSpeedToae {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  option.BossName,
						Value: option.BossName,
					})
				}
			case "Solo Bosses":
				for _, option := range util.HofSpeedSolo {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  option.BossName,
						Value: option.BossName,
					})
				}
			case "Nex":
				for _, option := range util.HofSpeedNex {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  option.BossName,
						Value: option.BossName,
					})
				}
			case "Slayer":
				for _, option := range util.HofSpeedSlayer {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  option.BossName,
						Value: option.BossName,
					})
				}
			}
		}

		err := session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		})
		if err != nil {
			logger.Error("Failed to handle speed autocomplete: " + err.Error())
		}
	}
	return ""
}

func (s *Service) updateSubmissionInstructions(ctx context.Context, session *discordgo.Session) string {
	returnMessage := "Successfully updated submission Instructions!"
	logger := flume.FromContext(ctx)

	// We will update the speed information first

	// First, delete all the messages within the channel
	messages, err := session.ChannelMessages(s.config.DiscSpeedSubInfoChan, 100, "", "", "")
	if err != nil {
		logger.Error("Failed to get all messages for deletion from channel: Speed Submission Info")
		return "Failed to get all messages for deletion from channel: Speed Submission Info"
	}
	var messageIDs []string
	for _, message := range messages {
		messageIDs = append(messageIDs, message.ID)
	}

	if len(messageIDs) > 0 {
		err = session.ChannelMessagesBulkDelete(s.config.DiscSpeedSubInfoChan, messageIDs)
		if err != nil {
			logger.Error("Failed to delete all messages from channel: Speed Submission Info, will try one by one")
			for _, message := range messageIDs {
				err = session.ChannelMessageDelete(s.config.DiscSpeedSubInfoChan, message)
				if err != nil {
					logger.Error("Failed to delete messages one by one from channel: Speed Submission Info")
					return "Failed to delete messages from channel: Speed Submission Info"
				}
			}
		}
	} else {
		logger.Debug("No messages to delete - proceeding with posting")
	}

	speedSubmissionInstruction := []string{
		"# Instructions for Speed Submissions",
		"In order to manually submit for speed times, use the /speed-submissions command. There will be **4 mandatory fields** which are automatically placed in your chat box and there are 2 optional fields which needs to be selected when pressing the +2 more at the end of the chat box",
		"## Mandatory Fields For Speed Submissions",
		"https://i.imgur.com/MK6BzCK.png",
		"### category\nThis has a list of all the speed categories you see in the hof-speeds forum. Select one of these in order to proceed in the submission",
		"https://i.imgur.com/uVDhk9U.png",
		"### boss\nThis has a list of all the bosses in the previously selected category. Select one of these options to make a speed submission for",
		"https://i.imgur.com/gXD9bHy.png",
		"### speed-time\nThe time must be in the format of hh:mm:ss.ms where hh = hours, mm = minutes, ss = seconds, and ms = milliseconds. The following example is 20 hours, 20 minutes, 20 seconds and 1 tick",
		"https://i.imgur.com/uzwDOL3.png",
		"### player-names\nThis is comma separated list of all the participating ponies members. Any non-members submitted will cause an error in the submission.",
		"https://i.imgur.com/ML14RzQ.png",
		"## Additional Fields",
		"https://i.imgur.com/dD4FKb9.png",
		"**NOTE: Only 1 or either the screenshot field or i-imgur-link field is acceptable. Using both will cause and error as well as using none!**",
		"### screenshot\nThis allows you to select an image from your computer to upload to the submission",
		"https://i.imgur.com/SGvWSt8.png",
		"### i-imgur-link\nThis allows you to put in an i.imgur.com url instead of an image upload",
		"https://i.imgur.com/TaoiTLG.png",
		"# Examples of submissions",
		"## Speed Submission using screenshot",
		"https://i.imgur.com/IlgOsfy.gif",
		"# What happens when your screenshot gets approved/denied",
		"Once a moderator approves/denies your submission, a message will popup in a private channel between you and the moderators with feedback on your submission",
		"https://i.imgur.com/V5AiTyZ.png",
		"If you have an issue, you can ask the moderators there about what was incorrect about your submission.",
	}

	for _, msg := range speedSubmissionInstruction {
		_, err := session.ChannelMessageSend(s.config.DiscSpeedSubInfoChan, msg)
		if err != nil {
			logger.Error("Failed to send message to cp information channel", err)
			return "Failed to send message to cp information channel"
		}
	}

	// Now we will update the clan points information

	// First, delete all the messages within the channel
	messages, err = session.ChannelMessages(s.config.DiscPPInfoChan, 100, "", "", "")
	if err != nil {
		logger.Error("Failed to get all messages for deletion from channel: Speed Submission Info")
		return "Failed to get all messages for deletion from channel: Speed Submission Info"
	}
	for _, message := range messages {
		messageIDs = append(messageIDs, message.ID)
	}

	if len(messageIDs) > 0 {
		err = session.ChannelMessagesBulkDelete(s.config.DiscPPInfoChan, messageIDs)
		if err != nil {
			logger.Error("Failed to delete all messages from channel: Speed Submission Info, will try one by one")
			for _, message := range messageIDs {
				err = session.ChannelMessageDelete(s.config.DiscPPInfoChan, message)
				if err != nil {
					logger.Error("Failed to delete messages one by one from channel: Speed Submission Info")
					return "Failed to delete messages from channel: Speed Submission Info"
				}
			}
		}
	} else {
		logger.Debug("No messages to delete - proceeding with posting")
	}

	ppSubmissionInstruction := []string{
		"# Instructions for Ponies Points Submissions",
		"In order to manually submit for ponies points, use the /pp-submissions command. There will be **1 mandatory field** which is automatically placed in your chat box and there are 2 optional fields which needs to be selected when pressing the +2 more at the end of the chat box",
		"## Mandatory Fields For Speed Submissions",
		"https://i.imgur.com/hi66ThP.png",
		"### player-names\nThis is comma separated list of all the participating ponies members. Any non-members submitted will cause an error in the submission.",
		"https://i.imgur.com/lzYUZUz.png",
		"## Additional Fields",
		"https://i.imgur.com/dD4FKb9.png",
		"**NOTE: Only 1 or either the screenshot field or i-imgur-link field is acceptable. Using both will cause and error as well as using none!**",
		"### screenshot\nThis allows you to select an image from your computer to upload to the submission",
		"https://i.imgur.com/SGvWSt8.png",
		"### i-imgur-link\nThis allows you to put in an i.imgur.com url instead of an image upload",
		"https://i.imgur.com/TaoiTLG.png",
		"# Examples of submissions",
		"## PP Submission using screenshot",
		"https://i.imgur.com/FAFCyim.gif",
		"# What happens when your screenshot gets approved/denied",
		"Once a moderator approves/denies your submission, a message will popup in a private channel between you and the moderators with feedback on your submission",
		"https://i.imgur.com/QUvB4oo.png",
		"If you have an issue, you can ask the moderators there about what was incorrect about your submission.",
	}

	for _, msg := range ppSubmissionInstruction {
		_, err := session.ChannelMessageSend(s.config.DiscPPInfoChan, msg)
		if err != nil {
			logger.Error("Failed to send message to cp information channel", err)
			return "Failed to send message to cp information channel"
		}
	}

	keys := make([]string, 0, len(util.LootLogClanPoint))

	for key := range util.LootLogClanPoint {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return util.LootLogClanPoint[keys[i]] < util.LootLogClanPoint[keys[j]]
	})

	var cpInstructions []string
	currentCategory := ""
	currentString := "# The following items will count for Ponies Points"

	for _, item := range keys {
		category := util.LootLogClanPoint[item]
		if strings.Compare(currentCategory, category) != 0 {
			cpInstructions = append(cpInstructions, currentString)
			currentCategory = category
			currentString = "## " + category + "\n"
		}
		currentString = currentString + "- " + item + "\n"
	}

	for _, msg := range cpInstructions {
		_, err := session.ChannelMessageSend(s.config.DiscPPInfoChan, msg)
		if err != nil {
			logger.Error("Failed to send message to cp information channel", err)
			return "Failed to send message to cp information channel"
		}
	}

	return returnMessage
}

func (s *Service) handlePlayerAdministration(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) string {
	logger := flume.FromContext(ctx)
	options := i.ApplicationCommandData().Options[0].Options

	option := ""
	name := ""
	newName := ""
	discordid := ""
	discordname := ""

	for _, iterOption := range options {
		switch iterOption.Name {
		case "option":
			option = iterOption.Value.(string)
		case "name":
			name = iterOption.Value.(string)
		case "new-name":
			newName = iterOption.Value.(string)
		case "discord-id":
			discordid = iterOption.Value.(string)
		case "discord-name":
			discordname = iterOption.Value.(string)
		}
	}

	switch option {
	case "Add":
		// If add, ensure that discord-id and discord-name are there
		if len(discordid) == 0 || len(discordname) == 0 {
			logger.Error("Discord ID and Discord Name are required for an addition to the clan")
			msg := "Discord ID and Discord Name are required for an addition to the clan"
			return msg
		}

		// Ensure that this person does not exist in the cp map currently
		if _, ok := s.cp[name]; ok {
			// Send the failed addition message in the previously created private channel
			logger.Error("Member: " + name + " already exists.")
			msg := "Member: " + name + " already exists."
			return msg
		} else {
			s.cp[name] = 0
			s.members[name] = util.MemberInfo{
				DiscordId:   discordid,
				DiscordName: discordname,
			}
			s.temple.AddMemberToTemple(ctx, name, s.config.TempleGroupId, s.config.TempleGroupKey)

			logger.Debug("You have successfully added a new member: " + name)
			msg := "You have successfully added a new member: " + name
			return msg
		}
	case "Remove":
		// Remove the user from the temple page
		s.temple.RemoveMemberFromTemple(ctx, name, s.config.TempleGroupId, s.config.TempleGroupKey)

		if _, ok := s.cp[name]; ok {
			delete(s.cp, name)
			delete(s.members, name)

			logger.Debug("You have successfully removed a member: " + name)
			msg := "You have successfully removed a member: " + name
			return msg

		} else {
			// Send the failed removal message in the previously created private channel
			logger.Error("Member: " + name + " does not exist.")
			msg := "Member: " + name + " does not exist."
			return msg
		}
	case "Name Change":
		if _, ok := s.cp[name]; ok {
			// Remove the user from the temple page and add new name
			s.temple.RemoveMemberFromTemple(ctx, name, s.config.TempleGroupId, s.config.TempleGroupKey)
			s.temple.AddMemberToTemple(ctx, newName, s.config.TempleGroupId, s.config.TempleGroupKey)

			// Update HOF Speed times from old name to new name
			updatedSpeedInfo := make(map[string]util.SpeedInfo)
			for boss, speedInfo := range s.speed {
				updatedSpeedInfo[boss] = util.SpeedInfo{
					PlayersInvolved: strings.Replace(speedInfo.PlayersInvolved, name, newName, -1),
					Time:            speedInfo.Time,
					URL:             speedInfo.URL,
				}
			}
			s.speed = updatedSpeedInfo
			s.members[newName] = s.members[name]
			s.cp[newName] = s.cp[name]
			delete(s.cp, name)
			delete(s.members, name)

			logger.Debug("You have successfully changed names from: " + name + " to: " + newName)
			msg := "You have successfully changed names from: " + name + " to: " + newName

			return msg

		} else {
			// Send the failed removal message in the previously created private channel
			logger.Error("Member: " + name + " does not exist.")
			msg := "Member: " + name + " does not exist."
			return msg
		}

	default:
		return "Invalid player management option chosen."
	}
}

func (s *Service) handleGuideAdministrationSubmission(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) {
	logger := flume.FromContext(ctx)
	options := i.ApplicationCommandData().Options

	option := options[0].Value.(string)
	guide := options[1].Value.(string)

	switch option {
	case "Update":
		// Remove leading and trailing whitespaces
		msg := strings.TrimSpace(guide)
		logger.Debug("Updating guide: " + msg)

		switch msg {
		case "trio-cm":
			session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Received request - kicking off trio-cm guide update. https://discord.com/channels/1172535371905646612/1183750806709735424",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			s.updateTrioCMGuide(ctx, session)
			logger.Info("Successfully updated the trio-cm guide!")
		case "tob":
			session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Received request - kicking off tob guide update. https://discord.com/channels/1172535371905646612/1184607458153484430",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			s.updateTobGuide(ctx, session)
			logger.Info("Successfully updated the tob guide!")
		default:
			logger.Error("Unknown guide chosen: " + guide)
		}
	default:
		logger.Error("Invalid guide management option chosen.")
	}
}

func (s *Service) updateCpPoints(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) string {
	// options := i.ApplicationCommandData().Options[0].Options
	return "Will do eventually"
}
