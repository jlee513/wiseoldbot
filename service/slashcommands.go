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
			Name:        "submission",
			Description: "Submit screenshots for events, Ponies Points, and speed times",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "submission-type",
					Description: "Choose one of the following: Event, Ponies Point, or Speed",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Event",
							Value: "Event",
						},
						{
							Name:  "Ponies Point",
							Value: "Ponies Point",
						},
						{
							Name:  "Speed",
							Value: "Speed",
						},
					},
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
					Name:        "imgur_link",
					Description: "Imgur link of the submission",
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "speed-time",
					Description: "Only use if making a speed submission in format: xx:xx:xx.xx",
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "speed-bossname",
					Description: "Only use if making a speed submission",
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
							},
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "player",
							Description: "Player name",
							Required:    true,
						},
					},
				},
				{
					Name:        "pp-instructions",
					Description: "Update Ponies Points instructions",
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
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "leaderboard",
							Description: "leaderboard name",
							Required:    true,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "Kc",
									Value: "Kc",
								},
								{
									Name:  "Speed",
									Value: "Speed",
								},
							},
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "thread",
							Description: "Name of the thread you want to update",
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
	case "submission":
		returnMessage = s.handleSlashSubmission(ctx, session, i)
		break
	case "guide-administration":
		s.handleGuideAdministrationSubmission(ctx, session, i)
		return
	case "admin":
		s.handleAdmin(ctx, session, i)
		return
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

/* All the slash commands handling functions */
func (s *Service) handleSlashSubmission(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) string {
	logger := flume.FromContext(ctx)
	options := i.ApplicationCommandData().Options
	channelId := ""

	typeOfSubmission := ""
	playersInvolved := ""
	screenshot := ""
	imgurUrl := ""
	speedTime := ""
	bossName := ""

	for _, option := range options {
		switch option.Name {
		case "submission-type":
			typeOfSubmission = option.Value.(string)
			break
		case "player-names":
			playersInvolved = option.Value.(string)
			break
		case "screenshot":
			screenshot = i.ApplicationCommandData().Resolved.Attachments[option.Value.(string)].URL
			break
		case "imgur_link":
			imgurUrl = option.Value.(string)
			break
		case "speed-time":
			speedTime = option.Value.(string)
			break
		case "speed-bossname":
			bossName = option.Value.(string)
			break
		}
	}

	logger.Info("Submission created by: " + i.Member.User.Username + " with type: " + typeOfSubmission)

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

	// Ensure the right submission type is used
	switch typeOfSubmission {
	case "Event":
		channelId = s.config.DiscEventApprovalChan
		break
	case "Ponies Point":
		channelId = s.config.DiscCpApprovalChan
		break
	case "Speed":
		channelId = s.config.DiscSpeedApprovalChan
		break
	default:
		logger.Error("Unknown submission type used: " + typeOfSubmission)
		return "Unknown submission type used. Please use the drop down options when selecting type."
	}

	// If we have speed but don't have the time, send back an error message
	if typeOfSubmission == "Speed" && len(options) < 5 {
		logger.Error("Time or boss name not provided in submission")
		return "Time or boss name was not provided during submission, click the +2 more at the end of the submission and click speed-time and speed-boss in the popup to enter in the time and boss name."
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

	msgToBeApproved := ""
	if typeOfSubmission == "Speed" {
		if _, ok := util.SpeedBossNameToCategory[bossName]; !ok {
			logger.Error("Incorrect boss name: ", bossName)
			return "Incorrect boss name. Please look https://discord.com/channels/1172535371905646612/1194975272487878707/1194975272487878707 to see which boss names work."
		}

		// Ensure the format is hh:mm:ss:mmm
		reg := regexp.MustCompile("^\\d\\d:\\d\\d:\\d\\d.\\d\\d$")
		if !reg.Match([]byte(speedTime)) {
			logger.Error("Invalid time format: ", speedTime)
			return "Incorrect time format. Please use the following format: hh:mm:ss.mmm"
		}

		msgToBeApproved = fmt.Sprintf("<@&1194691758353821847>\nSubmitter: %s\nUser Id: %s\nBoss Name: %s\nTime: %s\nPlayers Involved: %s\n%+v", i.Member.User.Username, i.Member.User.ID, bossName, speedTime, playersInvolved, url)
	} else {
		msgToBeApproved = fmt.Sprintf("<@&1194691758353821847>\nSubmitter: %s\nUser Id: %s\nPlayers Involved: %s\n%+v", i.Member.User.Username, i.Member.User.ID, playersInvolved, url)
	}

	// If we have the submission is valid, send the submission information to the admin channel
	msg, err := session.ChannelMessageSend(channelId, msgToBeApproved)
	if err != nil {
		logger.Error("Failed to send message to admin channel", err)
		return "Issue with submitting, please contact a dev to fix this issue."
	} else {
		logger.Info("Submission sent to moderators for approval")
		// Add a check and x reaction to the message to accept or reject the submission
		session.MessageReactionAdd(channelId, msg.ID, "✅")
		session.MessageReactionAdd(channelId, msg.ID, "❌")
	}

	// If nothing wrong happened, send a happy message back to the submitter
	return "Successfully submitted! Awaiting approval from a moderator!"
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
		break
	case s.config.DiscSpeedApprovalChan:
		s.handleSpeedApproval(ctx, session, r)
		break
	case s.config.DiscEventApprovalChan:
		s.handleEventApproval(ctx, session, r)
		break
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
		s.updateCpLeaderboard(ctx, session)

		// Delete the screenshot in the page
		err := session.ChannelMessageDelete(s.config.DiscCpApprovalChan, r.MessageID)
		if err != nil {
			logger.Error("Failed to delete cp approval message: " + err.Error())
		}
		logger.Debug("Successfully added CPs for: " + playersInvolved)

		// Send feedback to user
		channel := s.checkOrCreateFeedbackChannel(ctx, session, submitter, submitterId)
		index = strings.Index(msg.Content, ">")
		feedBackMsg := "<@" + submitterId + ">\nYour submission has been accepted\n" + msg.Content[index+1:]
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
		index = strings.Index(msg.Content, ">")
		feedBackMsg := "<@" + submitterId + ">\nYour submission has been rejected\n" + msg.Content[index+1:]
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
				break
			case 1:
				c, _ := strconv.Atoi(splitTime)
				t = t.Add(time.Duration(c) * time.Minute)
				break
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
				break
			}
		}

		// If the submission time is faster than the current speed time for the boss, update it
		if t.Before(s.speed[bossName].Time) {
			logger.Info("NEW TIME FOR BOSS: " + bossName)
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
			s.updateCpLeaderboard(ctx, session)

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
		index = strings.Index(msg.Content, ">")
		feedBackMsg := "<@" + submitterId + ">\nYour submission has been accepted\n" + msg.Content[index+1:]
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
		index = strings.Index(msg.Content, ">")
		feedBackMsg := "<@" + submitterId + ">\nYour submission has been rejected\n" + msg.Content[index+1:]
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
	options := i.ApplicationCommandData().Options[0].Options

	leaderboardName := ""
	threadName := ""
	for _, option := range options {
		switch option.Name {
		case "leaderboard":
			leaderboardName = option.Value.(string)
			break
		case "thread":
			threadName = option.Value.(string)
			break
		}
	}

	switch leaderboardName {
	case "Kc":
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

}

func (s *Service) handleAdmin(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) string {
	options := i.ApplicationCommandData().Options
	returnMessage := ""

	switch options[0].Name {
	case "player":
		returnMessage = s.handlePlayerAdministration(ctx, session, i)
		break
	case "pp-instructions":
		session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Updating Ponies Point Instructions",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		_ = s.updateCpInstructions(ctx, session)
	case "update-cp":
		returnMessage = s.updateCpPoints(ctx, session, i)
		break
	case "update-leaderboard":
		s.updateLeaderboard(ctx, session, i)
		return ""
	}

	return returnMessage
}

func (s *Service) updateCpInstructions(ctx context.Context, session *discordgo.Session) string {
	returnMessage := "Successfully updated CP Instructions!"
	logger := flume.FromContext(ctx)

	// First, delete all the messages within the channel
	messages, err := session.ChannelMessages(s.config.DiscCpInformationChan, 100, "", "", "")
	if err != nil {
		logger.Error("Failed to get all messages for deletion from channel: Ponies Points Information Channel")
		return "Failed to get all messages for deletion from channel: Ponies Points Information Channel"
	}
	var messageIDs []string
	for _, message := range messages {
		messageIDs = append(messageIDs, message.ID)
	}

	if len(messageIDs) > 0 {
		err = session.ChannelMessagesBulkDelete(s.config.DiscCpInformationChan, messageIDs)
		if err != nil {
			logger.Error("Failed to delete all messages from channel: Ponies Points Information Channel, will try one by one")
			for _, message := range messageIDs {
				err = session.ChannelMessageDelete(s.config.DiscCpInformationChan, message)
				if err != nil {
					logger.Error("Failed to delete messages one by one from channel: Ponies Points Information Channel")
					return "Failed to delete messages from channel: Ponies Points Information Channel"
				}
			}
		}
	} else {
		logger.Debug("No messages to delete - proceeding with posting")
	}

	cpSubmissionInstruction := []string{
		"# Instructions for ponies point and speed submission",
		"In order to manually submit for ponies points, use the /submissions command. There will be 2 mandatory fields which are automatically placed in your chat box and there are 4 optional fields which needs to be selected when pressing the +4 more at the end of the chat box",
		"## Mandatory Fields For All Submissions",
		"https://i.imgur.com/Gu2WKNC.png",
		"### submission-type\nThis has 3 options which needs to be selected from the options menu that pops up (Event, Ponies Point, Speed)",
		"https://i.imgur.com/dfv8MRb.png",
		"### player-names\nThis is a comma separated list of all participating ponies for the Ponies Point (i.e. H ana,Chapo,Calibre). Spaces are allowed",
		"https://i.imgur.com/dZ4auf1.png",
		"## Additional Fields",
		"https://i.imgur.com/CSs9vOW.png",
		"**NOTE: Only 1 or either the screenshot field or imgur_link field is acceptable. Using both will cause and error as well as using none!**",
		"### screenshot\nThis allows you to select an image from your computer to upload to the submission",
		"https://i.imgur.com/SGvWSt8.png",
		"### imgur_link\nThis allows you to put in an i.imgur.com url instead of an image upload",
		"https://i.imgur.com/TaoiTLG.png",
		"### speed-time:\nThis is required for speed submissions and must be in the format of hh:mm:ss.ms where hh = hours, mm = minutes, ss = seconds, and ms = milliseconds",
		"https://i.imgur.com/Lb7k6uP.png",
		"### speed-bossname\nThis is required for speed submissions. It must be one of the spelling and capitalization specific boss names found: https://discord.com/channels/1172535371905646612/1194975272487878707/1194975272487878707",
		"# Examples of submissions",
		"## Speed Submission using screenshot",
		"https://i.imgur.com/8IEdtLK.gif",
		"## CP Submission using imgur",
		"https://i.imgur.com/o1XTqZm.gif",
	}

	for _, msg := range cpSubmissionInstruction {
		_, err := session.ChannelMessageSend(s.config.DiscCpInformationChan, msg)
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
	currentString := "# All items that count towards Ponies Points"

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
		_, err := session.ChannelMessageSend(s.config.DiscCpInformationChan, msg)
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

	option := options[0].Value.(string)
	player := options[1].Value.(string)

	switch option {
	case "Add":
		// Ensure that this person does not exist in the cp map currently
		if _, ok := s.cp[player]; ok {
			// Send the failed addition message in the previously created private channel
			logger.Error("Member: " + player + " already exists.")
			msg := "Member: " + player + " already exists."
			return msg
		} else {
			s.cp[player] = 0
			s.temple.AddMemberToTemple(ctx, player, s.config.TempleGroupId, s.config.TempleGroupKey)

			logger.Debug("You have successfully added a new member: " + player)
			msg := "You have successfully added a new member: " + player
			return msg
		}
	case "Remove":
		// Remove the user from the temple page
		s.temple.RemoveMemberFromTemple(ctx, player, s.config.TempleGroupId, s.config.TempleGroupKey)

		if _, ok := s.cp[player]; ok {
			delete(s.cp, player)

			logger.Debug("You have successfully removed a member: " + player)
			msg := "You have successfully removed a member: " + player
			return msg

		} else {
			// Send the failed removal message in the previously created private channel
			logger.Error("Member: " + player + " does not exist.")
			msg := "Member: " + player + " does not exist."
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
