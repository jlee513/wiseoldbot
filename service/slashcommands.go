package service

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gemalto/flume"
	"log"
	"osrs-disc-bot/util"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (s *Service) initSlashCommands(ctx context.Context, session *discordgo.Session) {
	logger := flume.FromContext(ctx)
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "speed-submission",
			Description: "Speed submissions for ponies",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "category",
					Description:  "Category of speed",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
				{
					Name:         "boss",
					Description:  "Boss name submitting for",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "speed-time",
					Description: "Only use if making a speed submission in format: hh:mm:ss.ms",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "player-names",
					Description: "Comma separated list of players involved",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "screenshot",
					Description: "Screenshot of the submission",
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "i-imgur-link",
					Description: "Imgur link of the submission",
				},
			},
		},
		{
			Name:        "pp-submission",
			Description: "Ponies points submissions for ponies",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "player-names",
					Description: "Comma separated list of players involved",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "screenshot",
					Description: "Screenshot of the submission",
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "i-imgur-link",
					Description: "Imgur link of the submission",
				},
			},
		},
		{
			Name:        "guide",
			Description: "Administration of Guides",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option",
					Description: "Choose one of the following: Update",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Update",
							Value: "Update",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "guide",
					Description: "Guide name",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "trio-cm",
							Value: "trio-cm",
						},
						{
							Name:  "tob",
							Value: "tob",
						},
					},
				},
			},
		},
		{
			Name:        "admin",
			Description: "Admin commands",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "player",
					Description: "Player administration commands",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "option",
							Description: "Choose one of the following: Add, Remove",
							Required:    true,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "Add",
									Value: "Add",
								},
								{
									Name:  "Remove",
									Value: "Remove",
								},
								{
									Name:  "Name Change",
									Value: "Name Change",
								},
							},
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "name",
							Description: "Player name",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "new-name",
							Description: "New player name",
						},
					},
				},
				{
					Name:        "submission-instructions",
					Description: "Update Submission instructions",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "update-points",
					Description: "Update Pp for player",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "player",
							Description: "Player name",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionInteger,
							Name:        "amount-of-pp",
							Description: "Amount of Pp to manage for player",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "option",
							Description: "Choose one of the following: Add, Remove",
							Required:    true,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "Add",
									Value: "Add",
								},
								{
									Name:  "Remove",
									Value: "Remove",
								},
							},
						},
					},
				},
				{
					Name:        "update-leaderboard",
					Description: "Update Leaderboard for player",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "leaderboard",
							Description:  "leaderboard name",
							Required:     true,
							Autocomplete: true,
						},
						{
							Name:         "thread",
							Description:  "Name of the thread you want to update",
							Type:         discordgo.ApplicationCommandOptionString,
							Required:     true,
							Autocomplete: true,
						},
					},
				},
			},
		},
	}

	// Iterate over all the commands and create the application command - we will save all the registered commands
	// into the service struct that will be used to delete all the commands on bot termination
	logger.Info("Adding all commands...")
	s.registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, s.config.DiscGuildId, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		logger.Debug("ADDING COMMAND: " + v.Name)
		s.registeredCommands[i] = cmd
	}
}

//func (s *Service) removeSlashCommands(session *discordgo.Session) {
//	logger.Info("Removing all commands...")
//
//	for _, v := range s.registeredCommands {
//		logger.Debug("REMOVING COMMAND: " + v.Name)
//		err := session.ApplicationCommandDelete(session.State.User.ID, s.config.DiscGuildId, v.ID)
//		if err != nil {
//			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
//		}
//	}
//}

func (s *Service) slashCommands(session *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid).With("user", i.Member.User.Username))
	logger := flume.FromContext(ctx)
	defer func() { s.tid++ }()

	returnMessage := ""

	switch i.ApplicationCommandData().Name {
	case "pp-submission":
		returnMessage = s.handlePPSubmission(ctx, session, i)
	case "guide":
		s.handleGuideAdministrationSubmission(ctx, session, i)
		return
	case "admin":
		s.handleAdmin(ctx, session, i)
		return
	case "speed-submission":
		returnMessage = s.handleSpeedSubmission(ctx, session, i)
	default:
		logger.Error("ERROR: UNKNOWN COMMAND USED: " + i.ApplicationCommandData().Name)
		returnMessage = "Error: Unknown Command Used"
	}

	session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: returnMessage,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func (s *Service) handleSpeedSubmission(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) string {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		logger := flume.FromContext(ctx)
		options := i.ApplicationCommandData().Options

		boss := ""
		speedTime := ""
		screenshot := ""
		imgurUrl := ""
		playersInvolved := ""

		for _, option := range options {
			switch option.Name {
			case "boss":
				boss = option.Value.(string)
			case "speed-time":
				speedTime = option.Value.(string)
			case "player-names":
				playersInvolved = option.Value.(string)
			case "screenshot":
				screenshot = i.ApplicationCommandData().Resolved.Attachments[option.Value.(string)].URL
			case "i-imgur-link":
				imgurUrl = option.Value.(string)
			}
		}

		logger.Info("Speed submission created by: " + i.Member.User.Username)

		// Can only have either a screenshot or an imgur link
		url := ""
		if len(screenshot) == 0 && len(imgurUrl) == 0 {
			logger.Error("No screenshot has been submitted")
			return "No image has been submitted - please provide either a screenshot or an imgur link in their respective sections."
		} else if len(screenshot) > 0 && len(imgurUrl) > 0 {
			logger.Error("Two screenshots has been submitted")
			return "Two images has been submitted - please provide either a screenshot or an imgur link in their respective sections, not both."
		} else if len(imgurUrl) > 0 {
			if !strings.Contains(imgurUrl, "https://i.imgur.com") {
				logger.Error("Incorrect link used in imgur url submission: " + imgurUrl)
				return "Incorrect link used in imgur url submission, please use https://i.imgur.com when submitting using the imgur url option."
			} else {
				url = imgurUrl
			}
		} else {
			url = screenshot
		}

		// Ensure the player used is valid
		// Split the names into an array by , then make an empty array with those names as keys for an easier lookup
		// instead of running a for loop inside a for loop when adding Ponies Points
		whitespaceStrippedMessage := strings.Replace(playersInvolved, ", ", ",", -1)
		whitespaceStrippedMessage = strings.Replace(whitespaceStrippedMessage, " ,", ",", -1)

		logger.Debug("Submitted names: " + whitespaceStrippedMessage)
		names := strings.Split(whitespaceStrippedMessage, ",")
		for _, name := range names {
			if _, ok := s.cp[name]; !ok {
				// We have a submission for an unknown person, throw an error
				logger.Error("Unknown player submitted: " + name)
				return "Unknown player submitted. Please ensure all the names are correct or sign-up the following person: " + name
			}
		}

		// Ensure the boss name is okay
		if _, ok := util.SpeedBossNameToCategory[boss]; !ok {
			logger.Error("Incorrect boss name: ", boss)
			return "Incorrect boss name. Please look ensure to select one of the options for boss names."
		}

		// Ensure the format is hh:mm:ss:mmm
		reg := regexp.MustCompile("^\\d\\d:\\d\\d:\\d\\d.\\d\\d$")
		if !reg.Match([]byte(speedTime)) {
			logger.Error("Invalid time format: ", speedTime)
			return "Incorrect time format. Please use the following format: hh:mm:ss.ms"
		}

		msgToBeApproved := fmt.Sprintf("<@&1194691758353821847>\nSubmitter: %s\nUser Id: %s\nBoss Name: %s\nTime: %s\nPlayers Involved: %s\n%+v", i.Member.User.Username, i.Member.User.ID, boss, speedTime, playersInvolved, url)

		// If we have the submission is valid, send the submission information to the admin channel
		msg, err := session.ChannelMessageSend(s.config.DiscSpeedApprovalChan, msgToBeApproved)
		if err != nil {
			logger.Error("Failed to send message to admin channel", err)
			return "Issue with submitting the speed submission, please contact a dev to fix this issue."
		} else {
			logger.Info("Submission sent to moderators for approval")
			// Add a check and x reaction to the message to accept or reject the submission
			session.MessageReactionAdd(s.config.DiscSpeedApprovalChan, msg.ID, "✅")
			session.MessageReactionAdd(s.config.DiscSpeedApprovalChan, msg.ID, "❌")
		}

		// If nothing wrong happened, send a happy message back to the submitter
		return "Speed submission successfully submitted! Awaiting approval from a moderator!"

		return ""
	case discordgo.InteractionApplicationCommandAutocomplete:
		logger := flume.FromContext(ctx)
		data := i.ApplicationCommandData()
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

/* All the slash commands handling functions */
func (s *Service) handlePPSubmission(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) string {
	logger := flume.FromContext(ctx)
	options := i.ApplicationCommandData().Options

	playersInvolved := ""
	screenshot := ""
	imgurUrl := ""

	for _, option := range options {
		switch option.Name {
		case "player-names":
			playersInvolved = option.Value.(string)
		case "screenshot":
			screenshot = i.ApplicationCommandData().Resolved.Attachments[option.Value.(string)].URL
		case "i-imgur-link":
			imgurUrl = option.Value.(string)
		}
	}

	logger.Info("PP submission created by: " + i.Member.User.Username)

	// Can only have either a screenshot or an imgur link
	url := ""
	if len(screenshot) == 0 && len(imgurUrl) == 0 {
		logger.Error("No screenshot has been submitted")
		return "No image has been submitted - please provide either a screenshot or an imgur link in their respective sections."
	} else if len(screenshot) > 0 && len(imgurUrl) > 0 {
		logger.Error("Two screenshots has been submitted")
		return "Two images has been submitted - please provide either a screenshot or an imgur link in their respective sections, not both."
	} else if len(imgurUrl) > 0 {
		if !strings.Contains(imgurUrl, "https://i.imgur.com") {
			logger.Error("Incorrect link used in imgur url submission: " + imgurUrl)
			return "Incorrect link used in imgur url submission, please use https://i.imgur.com when submitting using the imgur url option."
		} else {
			url = imgurUrl
		}
	} else {
		url = screenshot
	}

	// Ensure the player used is valid
	// Split the names into an array by , then make an empty array with those names as keys for an easier lookup
	// instead of running a for loop inside a for loop when adding Ponies Points
	whitespaceStrippedMessage := strings.Replace(playersInvolved, ", ", ",", -1)
	whitespaceStrippedMessage = strings.Replace(whitespaceStrippedMessage, " ,", ",", -1)

	logger.Debug("Submitted names: " + whitespaceStrippedMessage)
	names := strings.Split(whitespaceStrippedMessage, ",")
	for _, name := range names {
		if _, ok := s.cp[name]; !ok {
			// We have a submission for an unknown person, throw an error
			logger.Error("Unknown player submitted: " + name)
			return "Unknown player submitted. Please ensure all the names are correct or sign-up the following person: " + name
		}
	}

	msgToBeApproved := fmt.Sprintf("<@&1194691758353821847>\nSubmitter: %s\nUser Id: %s\nPlayers Involved: %s\n%+v", i.Member.User.Username, i.Member.User.ID, playersInvolved, url)

	// If we have the submission is valid, send the submission information to the admin channel
	msg, err := session.ChannelMessageSend(s.config.DiscCpApprovalChan, msgToBeApproved)
	if err != nil {
		logger.Error("Failed to send message to admin channel", err)
		return "Issue with submitting, please contact a dev to fix this issue."
	} else {
		logger.Info("Submission sent to moderators for approval")
		// Add a check and x reaction to the message to accept or reject the submission
		session.MessageReactionAdd(s.config.DiscCpApprovalChan, msg.ID, "✅")
		session.MessageReactionAdd(s.config.DiscCpApprovalChan, msg.ID, "❌")
	}

	// If nothing wrong happened, send a happy message back to the submitter
	return "Ponies points successfully submitted! Awaiting approval from a moderator!"
}

func (s *Service) submissionApproval(session *discordgo.Session, r *discordgo.MessageReactionAdd) {

	// Don't handle message if it's created by the discord bot
	if r.UserID == session.State.User.ID {
		return
	}

	ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid).With("user", r.Member.User.Username))
	defer func() { s.tid++ }()

	switch r.ChannelID {
	case s.config.DiscCpApprovalChan:
		s.handleCpApproval(ctx, session, r)
	case s.config.DiscSpeedApprovalChan:
		s.handleSpeedApproval(ctx, session, r)
	case s.config.DiscEventApprovalChan:
		s.handleEventApproval(ctx, session, r)
	}
}

func (s *Service) checkOrCreateFeedbackChannel(ctx context.Context, session *discordgo.Session, user, userId string) string {
	logger := flume.FromContext(ctx)
	if feedbackChannelId, ok := s.feedback[user]; ok {
		return feedbackChannelId
	} else {
		channel, err := session.GuildChannelCreateComplex(s.config.DiscGuildId, discordgo.GuildChannelCreateData{
			Name: "feedback-" + user,
			Type: discordgo.ChannelTypeGuildText,
			PermissionOverwrites: []*discordgo.PermissionOverwrite{
				{
					ID:    userId,
					Type:  discordgo.PermissionOverwriteTypeMember,
					Allow: discordgo.PermissionAllText,
				},
				{
					ID:   s.config.DiscGuildId,
					Type: discordgo.PermissionOverwriteTypeRole,
					Deny: discordgo.PermissionViewChannel,
				},
				{
					// Moderator Rank
					ID:    "1194691758353821847",
					Type:  discordgo.PermissionOverwriteTypeRole,
					Allow: discordgo.PermissionAllText,
				},
			},
		})
		if err != nil {
			logger.Error("Failed to create private text channel: " + err.Error())
		}
		s.feedback[user] = "ponies" + channel.ID
		return channel.ID
	}
}

func (s *Service) handleCpApproval(ctx context.Context, session *discordgo.Session, r *discordgo.MessageReactionAdd) {
	logger := flume.FromContext(ctx)
	switch r.Emoji.Name {
	case "✅":
		logger.Info("Ponies Point submission approved by: " + r.Member.User.Username)
		msg, _ := session.ChannelMessage(s.config.DiscCpApprovalChan, r.MessageID)

		index := strings.Index(msg.Content, "Submitter:")
		index2 := strings.Index(msg.Content, "User Id:")
		submitter := msg.Content[index+11 : index2-1]

		index = strings.Index(msg.Content, "User Id:")
		index2 = strings.Index(msg.Content, "Players Involved:")
		submitterId := msg.Content[index+9 : index2-1]

		index = strings.Index(msg.Content, "Involved:")
		index2 = strings.Index(msg.Content, "https://")
		playersInvolved := msg.Content[index+10 : index2-1]
		logger.Info("CP Approved for: " + playersInvolved)

		submissionUrl := ""

		// If the url is an imgur link, skip uploading to imgur
		if strings.Contains(msg.Content, "https://i.imgur.com") {
			submissionUrl = msg.Content[index2:]
		} else {
			// Retrieve the bytes of the image
			resp, err := s.client.Get(msg.Embeds[0].URL)
			if err != nil {
				logger.Error("Failed to download discord image: " + err.Error())
				return
			}
			defer resp.Body.Close()

			// Retrieve the access token
			accessToken, err := s.imgur.GetNewAccessToken(ctx, s.config.ImgurRefreshToken, s.config.ImgurClientId, s.config.ImgurClientSecret)
			if err != nil {
				// We will retry 10 times to get a new access token
				counter := 0
				for err != nil {
					if counter == 10 {
						logger.Error("Failed to get access token for imgur: " + err.Error())
						return
					}
					accessToken, err = s.imgur.GetNewAccessToken(ctx, s.config.ImgurRefreshToken, s.config.ImgurClientId, s.config.ImgurClientSecret)
					if err != nil {
						counter++
					} else {
						break
					}
				}
			}
			submissionUrl = s.imgur.Upload(ctx, accessToken, resp.Body)
		}

		s.cpscreenshots[submissionUrl] = playersInvolved

		// Update the Ponies Points
		// Split the names into an array by , then make an empty array with those names as keys for an easier lookup
		// instead of running a for loop inside a for loop when adding Ponies Points
		whitespaceStrippedMessage := strings.Replace(playersInvolved, ", ", ",", -1)
		whitespaceStrippedMessage = strings.Replace(whitespaceStrippedMessage, " ,", ",", -1)

		names := strings.Split(whitespaceStrippedMessage, ",")
		for _, name := range names {
			logger.Debug("Adding Ponies Point to: " + name)
			s.cp[name] += 1
		}

		// Update the cp leaderboard
		s.updatePpLeaderboard(ctx, session)

		// Delete the screenshot in the page
		err := session.ChannelMessageDelete(s.config.DiscCpApprovalChan, r.MessageID)
		if err != nil {
			logger.Error("Failed to delete cp approval message: " + err.Error())
		}
		logger.Debug("Successfully added CPs for: " + playersInvolved)

		// Send feedback to user
		channel := s.checkOrCreateFeedbackChannel(ctx, session, submitter, submitterId)
		index = strings.Index(msg.Content, "Players Involved:")
		feedBackMsg := "<@" + submitterId + ">\nYour ponies point submission has been accepted\n\n" + msg.Content[index+1:]
		_, err = session.ChannelMessageSend(channel, feedBackMsg)
		if err != nil {
			logger.Error("Failed to send message to cp information channel", err)
		}
	case "❌":
		logger.Info("Ponies Point submission denied by: " + r.Member.User.Username)

		msg, _ := session.ChannelMessage(s.config.DiscCpApprovalChan, r.MessageID)

		index := strings.Index(msg.Content, "Submitter:")
		index2 := strings.Index(msg.Content, "User Id:")
		submitter := msg.Content[index+11 : index2-1]

		index = strings.Index(msg.Content, "User Id:")
		index2 = strings.Index(msg.Content, "Players Involved:")
		submitterId := msg.Content[index+9 : index2-1]

		// Send feedback to user
		channel := s.checkOrCreateFeedbackChannel(ctx, session, submitter, submitterId)
		index = strings.Index(msg.Content, "Players Involved:")
		feedBackMsg := "<@" + submitterId + ">\nYour ponies point submission has been rejected\n\n" + msg.Content[index:]
		_, err := session.ChannelMessageSend(channel, feedBackMsg)
		if err != nil {
			logger.Error("Failed to send message to cp information channel", err)
		}

		// Delete the screenshot in the page
		err = session.ChannelMessageDelete(s.config.DiscCpApprovalChan, r.MessageID)
		if err != nil {
			logger.Error("Failed to delete cp approval message: " + err.Error())
		}
	}
}

func (s *Service) handleSpeedApproval(ctx context.Context, session *discordgo.Session, r *discordgo.MessageReactionAdd) {
	logger := flume.FromContext(ctx)
	switch r.Emoji.Name {
	case "✅":
		logger.Info("Speed submission approved by: " + r.Member.User.Username)
		msg, _ := session.ChannelMessage(s.config.DiscSpeedApprovalChan, r.MessageID)
		index := strings.Index(msg.Content, "Submitter:")
		index2 := strings.Index(msg.Content, "User Id:")
		submitter := msg.Content[index+11 : index2-1]

		index = strings.Index(msg.Content, "User Id:")
		index2 = strings.Index(msg.Content, "Boss Name:")
		submitterId := msg.Content[index+9 : index2-1]

		index = strings.Index(msg.Content, "Name:")
		index2 = strings.Index(msg.Content, "Time:")
		bossName := msg.Content[index+6 : index2-1]

		index = strings.Index(msg.Content, "Time:")
		index2 = strings.Index(msg.Content, "Players Involved:")
		speedTime := msg.Content[index+6 : index2-1]

		index = strings.Index(msg.Content, "Involved:")
		index2 = strings.Index(msg.Content, "https://")
		playersInvolved := msg.Content[index+10 : index2-1]

		logger.Debug("Speed Approved for: " + playersInvolved + " with speedTime: " + speedTime + " at boss: " + bossName)

		submissionUrl := ""

		// If the url is an imgur link, skip uploading to imgur
		if strings.Contains(msg.Content, "https://i.imgur.com") {
			submissionUrl = msg.Content[index2:]
		} else {
			// Retrieve the bytes of the image
			resp, err := s.client.Get(msg.Embeds[0].URL)
			if err != nil {
				logger.Error("Failed to download discord image: " + err.Error())
				return
			}
			defer resp.Body.Close()

			// Retrieve the access token
			accessToken, err := s.imgur.GetNewAccessToken(ctx, s.config.ImgurRefreshToken, s.config.ImgurClientId, s.config.ImgurClientSecret)
			if err != nil {
				logger.Debug("Failed to get imgur access token, will retry...")
				// We will retry 10 times to get a new access token
				counter := 1
				for err != nil {
					logger.Debug("Failed to get imgur access token, will retry (attempt " + strconv.Itoa(counter) + ")")
					if counter == 11 {
						logger.Error("Failed to get access token for imgur: " + err.Error())
						return
					}
					accessToken, err = s.imgur.GetNewAccessToken(ctx, s.config.ImgurRefreshToken, s.config.ImgurClientId, s.config.ImgurClientSecret)
					if err != nil {
						counter++
						continue
					} else {
						break
					}
				}
			}
			submissionUrl = s.imgur.Upload(ctx, accessToken, resp.Body)
		}
		submissionTime := time.Now().Format("2006-01-02 15:04:05")
		s.speedscreenshots[submissionTime] = util.SpeedScInfo{BossName: bossName, Time: speedTime, PlayersInvolved: playersInvolved, URL: submissionUrl}

		// Only change the current top speed if it's faster
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

		// If the submission time is faster than the current speed time for the boss, update it
		if t.Before(s.speed[bossName].Time) {
			logger.Info("NEW TIME FOR BOSS: " + bossName)
			// Add message into new fastest time channel
			_, err := session.ChannelMessageSend(s.config.DiscNewFastestTimeChan, fmt.Sprintf("New record for %s has been set!\n%s has beaten %s!\nOld time: %+v\nNew Time: %+v\n%s", bossName, playersInvolved, s.speed[bossName].PlayersInvolved, s.speed[bossName].Time.Format("15:04:05.00"), t.Format("15:04:05.00"), submissionUrl))
			if err != nil {
				logger.Error("Failed to send message to cp information channel", err)
			}

			logger.Info(fmt.Sprintf("Old time: %+v", s.speed[bossName].Time.Format("15:04:05.00")))
			logger.Info(fmt.Sprintf("New Time: %+v", t.Format("15:04:05.00")))
			s.speed[bossName] = util.SpeedInfo{Time: t, PlayersInvolved: playersInvolved, URL: submissionUrl}

			// Update the Ponies Points
			// Split the names into an array by , then make an empty array with those names as keys for an easier lookup
			// instead of running a for loop inside a for loop when adding Ponies Points
			whitespaceStrippedMessage := strings.Replace(playersInvolved, ", ", ",", -1)
			whitespaceStrippedMessage = strings.Replace(whitespaceStrippedMessage, " ,", ",", -1)

			names := strings.Split(whitespaceStrippedMessage, ",")
			for _, name := range names {
				logger.Debug("Adding Ponies Point to: " + name)
				s.cp[name] += 1
			}

			// Update the cp leaderboard
			s.updatePpLeaderboard(ctx, session)

			// Update the boss leaderboard that was updated
			s.updateSpeedHOF(ctx, session, util.SpeedBossNameToCategory[bossName])

		} else {
			logger.Info("KEEP TIME FOR BOSS: " + bossName)
			logger.Info(fmt.Sprintf("Current time: %+v", s.speed[bossName].Time.Format("15:04:05.00")))
			logger.Info(fmt.Sprintf("Submitted Time: %+v", t.Format("15:04:05.00")))
		}

		// Delete the screenshot in the page
		err := session.ChannelMessageDelete(s.config.DiscSpeedApprovalChan, r.MessageID)
		if err != nil {
			logger.Error("Failed to delete cp approval message: " + err.Error())
		}

		// Send feedback to user
		channel := s.checkOrCreateFeedbackChannel(ctx, session, submitter, submitterId)
		index = strings.Index(msg.Content, "Boss Name:")
		feedBackMsg := "<@" + submitterId + ">\nYour speed submission has been accepted\n\n" + msg.Content[index:]
		_, err = session.ChannelMessageSend(channel, feedBackMsg)
		if err != nil {
			logger.Error("Failed to send message to cp information channel", err)
		}

		logger.Debug("Successfully handled Speed Time for: " + playersInvolved)
	case "❌":
		logger.Info("Speed submission denied by: " + r.Member.User.Username)
		msg, _ := session.ChannelMessage(s.config.DiscSpeedApprovalChan, r.MessageID)
		index := strings.Index(msg.Content, "Submitter:")
		index2 := strings.Index(msg.Content, "User Id:")
		submitter := msg.Content[index+11 : index2-1]

		index = strings.Index(msg.Content, "User Id:")
		index2 = strings.Index(msg.Content, "Boss Name:")
		submitterId := msg.Content[index+9 : index2-1]

		// Send feedback to user
		channel := s.checkOrCreateFeedbackChannel(ctx, session, submitter, submitterId)
		index = strings.Index(msg.Content, "Boss Name:")
		feedBackMsg := "<@" + submitterId + ">\nYour speed submission has been rejected\n\n" + msg.Content[index:]
		_, err := session.ChannelMessageSend(channel, feedBackMsg)
		if err != nil {
			logger.Error("Failed to send message to cp information channel", err)
		}

		// Delete the screenshot in the page
		err = session.ChannelMessageDelete(s.config.DiscSpeedApprovalChan, r.MessageID)
		if err != nil {
			logger.Error("Failed to delete cp approval message: " + err.Error())
		}
	}
}

func (s *Service) handleEventApproval(ctx context.Context, session *discordgo.Session, r *discordgo.MessageReactionAdd) {
	logger := flume.FromContext(ctx)
	switch r.Emoji.Name {
	case "✅":
		// TODO: Write when there's an event

		// Delete the screenshot in the page
		err := session.ChannelMessageDelete(s.config.DiscCpApprovalChan, r.MessageID)
		if err != nil {
			logger.Error("Failed to delete cp approval message: " + err.Error())
		}
	case "❌":
		// TODO: Find a way to let the user know that their submission has been rejected

		// Delete the screenshot in the page
		err := session.ChannelMessageDelete(s.config.DiscCpApprovalChan, r.MessageID)
		if err != nil {
			logger.Error("Failed to delete cp approval message: " + err.Error())
		}
	}
}

func (s *Service) updateLeaderboard(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		logger := flume.FromContext(ctx)
		options := i.ApplicationCommandData().Options[0].Options

		leaderboardName := ""
		threadName := ""
		for _, option := range options {
			switch option.Name {
			case "leaderboard":
				leaderboardName = option.Value.(string)
			case "thread":
				threadName = option.Value.(string)
			}
		}

		switch leaderboardName {
		case "Kc":
			logger.Info("Admin invoked Kc Hall Of Fame Update: ", i.Member.User.Username)
			session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Updating Leaderboard: " + leaderboardName,
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			// If kc is updating, always update all of them
			s.updateKcHOF(ctx, session)
		case "Speed":
			logger.Info("Admin invoked Speed Hall Of Fame Update: ", i.Member.User.Username)
			if _, ok := util.HofSpeedCategories[threadName]; ok {
				session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Updating Leaderboard: " + leaderboardName + " thread: " + threadName,
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				s.updateSpeedHOF(ctx, session, threadName)
			}
		default:
			session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Updating Leaderboard: " + leaderboardName + " thread: " + threadName,
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
	case discordgo.InteractionApplicationCommandAutocomplete:
		logger := flume.FromContext(ctx)
		data := i.ApplicationCommandData()
		var choices []*discordgo.ApplicationCommandOptionChoice
		switch {
		case data.Options[0].Options[0].Focused:
			choices = []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "Kc",
					Value: "Kc",
				},
				{
					Name:  "Speed",
					Value: "Speed",
				},
			}
		// In this case there are multiple autocomplete options. The Focused field shows which option user is focused on.
		case data.Options[0].Options[1].Focused:
			switch data.Options[0].Options[0].Value.(string) {
			case "Kc":
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  "All",
					Value: "All",
				})
			case "Speed":
				for thread, _ := range util.HofSpeedCategories {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  thread,
						Value: thread,
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
			logger.Error("Failed to handle admin autocomplete options: " + err.Error())
		}
	}
}

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

func (s *Service) updateCpPoints(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) string {
	// options := i.ApplicationCommandData().Options[0].Options
	return "Will do eventually"
}

func (s *Service) handlePlayerAdministration(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) string {
	logger := flume.FromContext(ctx)
	options := i.ApplicationCommandData().Options[0].Options

	option := ""
	name := ""
	newName := ""

	for _, iterOption := range options {
		switch iterOption.Name {
		case "option":
			option = iterOption.Value.(string)
		case "name":
			name = iterOption.Value.(string)
		case "new-name":
			newName = iterOption.Value.(string)
		}
	}

	switch option {
	case "Add":
		// Ensure that this person does not exist in the cp map currently
		if _, ok := s.cp[name]; ok {
			// Send the failed addition message in the previously created private channel
			logger.Error("Member: " + name + " already exists.")
			msg := "Member: " + name + " already exists."
			return msg
		} else {
			s.cp[name] = 0
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

		if _, ok := s.cp[name]; ok {
			s.cp[newName] = s.cp[name]
			delete(s.cp, name)

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
