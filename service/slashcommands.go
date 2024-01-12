package service

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gemalto/flume"
	"log"
	"osrs-disc-bot/util"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (s *Service) initSlashCommands(ctx context.Context, session *discordgo.Session) {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "submission",
			Description: "Submit screenshots for events, clan points, and speed times",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "submission-type",
					Description: "Choose one of the following: Event, Clan Point, or Speed",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Event",
							Value: "Event",
						},
						{
							Name:  "Clan Point",
							Value: "Clan Point",
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
			Name:        "player-administration",
			Description: "Administration of players",
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
			Name:        "guide-administration",
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
	}

	// Iterate over all the commands and create the application command - we will save all the registered commands
	// into the service struct that will be used to delete all the commands on bot termination
	s.log.Info("Adding all commands...")
	s.registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, s.config.DiscGuildId, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		s.log.Debug("ADDING COMMAND: " + v.Name)
		s.registeredCommands[i] = cmd
	}
}

func (s *Service) removeSlashCommands(session *discordgo.Session) {
	s.log.Info("Removing all commands...")

	for _, v := range s.registeredCommands {
		s.log.Debug("REMOVING COMMAND: " + v.Name)
		err := session.ApplicationCommandDelete(session.State.User.ID, s.config.DiscGuildId, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}

func (s *Service) slashCommands(session *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid))
	defer func() { s.tid++ }()

	returnMessage := ""

	switch i.ApplicationCommandData().Name {
	case "submission":
		returnMessage = s.handleSlashSubmission(ctx, session, i)
		break
	case "player-administration":
		returnMessage = s.handlePlayerAdministrationSubmission(ctx, session, i)
		break
	case "guide-administration":
		s.handleGuideAdministrationSubmission(ctx, session, i)
		return
	default:
		s.log.Error("ERROR: UNKNOWN COMMAND USED")
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

	// Can only have either a screenshot or an imgur link
	url := ""
	if len(screenshot) == 0 && len(imgurUrl) == 0 {
		s.log.Error("No screenshot has been submitted")
		return "No image has been submitted - please provide either a screenshot or an imgur link in their respective sections."
	} else if len(screenshot) > 0 && len(imgurUrl) > 0 {
		s.log.Error("Two screenshots has been submitted")
		return "Two images has been submitted - please provide either a screenshot or an imgur link in their respective sections, not both."
	} else if len(imgurUrl) > 0 {
		if !strings.Contains(imgurUrl, "https://i.imgur.com") {
			s.log.Error("Incorrect link used in imgur url submission: " + imgurUrl)
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
	case "Clan Point":
		channelId = s.config.DiscCpApprovalChan
		break
	case "Speed":
		channelId = s.config.DiscSpeedApprovalChan
		break
	default:
		s.log.Error("Unknown submission type used: " + typeOfSubmission)
		return "Unknown submission type used. Please use the drop down options when selecting type."
	}

	// If we have speed but don't have the time, send back an error message
	if typeOfSubmission == "Speed" && len(options) < 5 {
		s.log.Error("Time or boss name not provided in submission")
		return "Time or boss name was not provided during submission, click the +2 more at the end of the submission and click speed-time and speed-boss in the popup to enter in the time and boss name."
	}

	// Ensure the player used is valid
	// Split the names into an array by , then make an empty array with those names as keys for an easier lookup
	// instead of running a for loop inside a for loop when adding clan points
	whitespaceStrippedMessage := strings.Replace(playersInvolved, ", ", ",", -1)
	whitespaceStrippedMessage = strings.Replace(whitespaceStrippedMessage, " ,", ",", -1)

	s.log.Debug("Submitted names: " + whitespaceStrippedMessage)
	names := strings.Split(whitespaceStrippedMessage, ",")
	for _, name := range names {
		if _, ok := s.cp[name]; !ok {
			// We have a submission for an unknown person, throw an error
			s.log.Error("Unknown player submitted: " + name)
			return "Unknown player submitted. Please ensure all the names are correct or sign-up the following person: " + name
		}
	}

	msgToBeApproved := ""
	if typeOfSubmission == "Speed" {
		if _, ok := util.SpeedBossNames[bossName]; !ok {
			s.log.Error("Incorrect boss name: ", bossName)
			return "Incorrect boss name. Please look https://discord.com/channels/1172535371905646612/1194975272487878707/1194975272487878707 to see which boss names work."
		}

		// Ensure the format is hh:mm:ss:mmm
		reg := regexp.MustCompile("^\\d\\d:\\d\\d:\\d\\d.\\d\\d$")
		if !reg.Match([]byte(speedTime)) {
			s.log.Error("Invalid time format: ", speedTime)
			return "Incorrect time format. Please use the following format: hh:mm:ss.mmm"
		}

		msgToBeApproved = fmt.Sprintf("<@&1194691758353821847>\nBoss Name: %s\nTime: %s\nPlayers Involved: %s\n%+v", bossName, speedTime, playersInvolved, url)
	} else {
		msgToBeApproved = fmt.Sprintf("<@&1194691758353821847>\nPlayers Involved: %s\n%+v", playersInvolved, url)
	}

	// If we have the submission is valid, send the submission information to the admin channel
	msg, err := session.ChannelMessageSend(channelId, msgToBeApproved)
	if err != nil {
		s.log.Error("Failed to send message to admin channel", err)
		return "Issue with submitting, please contact a dev to fix this issue."
	} else {
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

	ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid))
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

func (s *Service) handleCpApproval(ctx context.Context, session *discordgo.Session, r *discordgo.MessageReactionAdd) {
	switch r.Emoji.Name {
	case "✅":
		msg, _ := session.ChannelMessage(s.config.DiscCpApprovalChan, r.MessageID)

		index := strings.Index(msg.Content, "Involved:")
		index2 := strings.Index(msg.Content, "https://")
		playersInvolved := msg.Content[index+10 : index2-1]
		s.log.Debug("CP Approved for: " + playersInvolved)

		submissionUrl := ""

		// If the url is an imgur link, skip uploading to imgur
		if strings.Contains(msg.Content, "https://i.imgur.com") {
			submissionUrl = msg.Content[index2:]
		} else {
			// Retrieve the bytes of the image
			resp, err := s.client.Get(msg.Embeds[0].URL)
			if err != nil {
				s.log.Error("Failed to download discord image: " + err.Error())
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
						s.log.Error("Failed to get access token for imgur: " + err.Error())
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

		// Update the clan points
		// Split the names into an array by , then make an empty array with those names as keys for an easier lookup
		// instead of running a for loop inside a for loop when adding clan points
		whitespaceStrippedMessage := strings.Replace(playersInvolved, ", ", ",", -1)
		whitespaceStrippedMessage = strings.Replace(whitespaceStrippedMessage, " ,", ",", -1)

		names := strings.Split(whitespaceStrippedMessage, ",")
		for _, name := range names {
			s.log.Debug("Adding clan point to: " + name)
			s.cp[name] += 1
		}

		// Update the cp leaderboard
		s.updateCpLeaderboard(ctx, session)

		// Delete the screenshot in the page
		err := session.ChannelMessageDelete(s.config.DiscCpApprovalChan, r.MessageID)
		if err != nil {
			s.log.Error("Failed to delete cp approval message: " + err.Error())
		}
		s.log.Debug("Successfully added CPs for: " + playersInvolved)
	case "❌":
		// TODO: Find a way to let the user know that their submission has been rejected

		// Delete the screenshot in the page
		err := session.ChannelMessageDelete(s.config.DiscCpApprovalChan, r.MessageID)
		if err != nil {
			s.log.Error("Failed to delete cp approval message: " + err.Error())
		}
	}
}

func (s *Service) handleSpeedApproval(ctx context.Context, session *discordgo.Session, r *discordgo.MessageReactionAdd) {
	switch r.Emoji.Name {
	case "✅":
		msg, _ := session.ChannelMessage(s.config.DiscSpeedApprovalChan, r.MessageID)
		index := strings.Index(msg.Content, "Name:")
		index2 := strings.Index(msg.Content, "Time:")
		bossName := msg.Content[index+6 : index2-1]

		index = strings.Index(msg.Content, "Time:")
		index2 = strings.Index(msg.Content, "Players Involved:")
		speedTime := msg.Content[index+6 : index2-1]

		index = strings.Index(msg.Content, "Involved:")
		index2 = strings.Index(msg.Content, "https://")
		playersInvolved := msg.Content[index+10 : index2-1]

		s.log.Debug("Speed Approved for: " + playersInvolved + " with speedTime: " + speedTime + " at boss: " + bossName)

		submissionUrl := ""

		// If the url is an imgur link, skip uploading to imgur
		if strings.Contains(msg.Content, "https://i.imgur.com") {
			submissionUrl = msg.Content[index2:]
		} else {
			// Retrieve the bytes of the image
			resp, err := s.client.Get(msg.Embeds[0].URL)
			if err != nil {
				s.log.Error("Failed to download discord image: " + err.Error())
				return
			}
			defer resp.Body.Close()

			// Retrieve the access token
			accessToken, err := s.imgur.GetNewAccessToken(ctx, s.config.ImgurRefreshToken, s.config.ImgurClientId, s.config.ImgurClientSecret)
			if err != nil {
				s.log.Debug("Failed to get imgur access token, will retry...")
				// We will retry 10 times to get a new access token
				counter := 1
				for err != nil {
					s.log.Debug("Failed to get imgur access token, will retry (attempt " + strconv.Itoa(counter) + ")")
					if counter == 11 {
						s.log.Error("Failed to get access token for imgur: " + err.Error())
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
			s.log.Info("NEW TIME FOR BOSS: " + bossName)
			s.log.Info(fmt.Sprintf("Old time: %+v", s.speed[bossName].Time.Format("15:04:05.00")))
			s.log.Info(fmt.Sprintf("New Time: %+v", t.Format("15:04:05.00")))
			s.speed[bossName] = util.SpeedInfo{Time: t, PlayersInvolved: playersInvolved, URL: submissionUrl}
		} else {
			s.log.Info("KEEP TIME FOR BOSS: " + bossName)
			s.log.Info(fmt.Sprintf("Current time: %+v", s.speed[bossName].Time.Format("15:04:05.00")))
			s.log.Info(fmt.Sprintf("Submitted Time: %+v", t.Format("15:04:05.00")))
		}

		// Delete the screenshot in the page
		err := session.ChannelMessageDelete(s.config.DiscSpeedApprovalChan, r.MessageID)
		if err != nil {
			s.log.Error("Failed to delete cp approval message: " + err.Error())
		}

		s.log.Debug("Successfully handled Speed Time for: " + playersInvolved)
	case "❌":
		// TODO: Find a way to let the user know that their submission has been rejected

		// Delete the screenshot in the page
		err := session.ChannelMessageDelete(s.config.DiscCpApprovalChan, r.MessageID)
		if err != nil {
			s.log.Error("Failed to delete cp approval message: " + err.Error())
		}
	}
}

func (s *Service) handleEventApproval(ctx context.Context, session *discordgo.Session, r *discordgo.MessageReactionAdd) {
	switch r.Emoji.Name {
	case "✅":
		// TODO: Write when there's an event

		// Delete the screenshot in the page
		err := session.ChannelMessageDelete(s.config.DiscCpApprovalChan, r.MessageID)
		if err != nil {
			s.log.Error("Failed to delete cp approval message: " + err.Error())
		}
	case "❌":
		// TODO: Find a way to let the user know that their submission has been rejected

		// Delete the screenshot in the page
		err := session.ChannelMessageDelete(s.config.DiscCpApprovalChan, r.MessageID)
		if err != nil {
			s.log.Error("Failed to delete cp approval message: " + err.Error())
		}
	}
}

func (s *Service) handlePlayerAdministrationSubmission(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) string {
	options := i.ApplicationCommandData().Options

	option := options[0].Value.(string)
	player := options[1].Value.(string)

	switch option {
	case "Add":
		// Ensure that this person does not exist in the cp map currently
		if _, ok := s.cp[player]; ok {
			// Send the failed addition message in the previously created private channel
			s.log.Error("Member: " + player + " already exists.")
			msg := "Member: " + player + " already exists."
			return msg
		} else {
			s.cp[player] = 0
			s.temple.AddMemberToTemple(ctx, player, s.config.TempleGroupId, s.config.TempleGroupKey)

			s.log.Debug("You have successfully added a new member: " + player)
			msg := "You have successfully added a new member: " + player
			return msg
		}
	case "Remove":
		// Remove the user from the temple page
		s.temple.RemoveMemberFromTemple(ctx, player, s.config.TempleGroupId, s.config.TempleGroupKey)

		if _, ok := s.cp[player]; ok {
			delete(s.cp, player)

			s.log.Debug("You have successfully removed a member: " + player)
			msg := "You have successfully removed a member: " + player
			return msg

		} else {
			// Send the failed removal message in the previously created private channel
			s.log.Error("Member: " + player + " does not exist.")
			msg := "Member: " + player + " does not exist."
			return msg
		}
	default:
		return "Invalid player management option chosen."
	}
}

func (s *Service) handleGuideAdministrationSubmission(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	option := options[0].Value.(string)
	guide := options[1].Value.(string)

	switch option {
	case "Update":
		// Remove leading and trailing whitespaces
		msg := strings.TrimSpace(guide)
		s.log.Debug("Updating guide: " + msg)

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
			s.log.Info("Successfully updated the trio-cm guide!")
		case "tob":
			session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Received request - kicking off tob guide update. https://discord.com/channels/1172535371905646612/1184607458153484430",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			s.updateTobGuide(ctx, session)
			s.log.Info("Successfully updated the tob guide!")
		default:
			s.log.Error("Unknown guide chosen: " + guide)
		}
	default:
		s.log.Error("Invalid guide management option chosen.")
	}
}
