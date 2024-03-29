package service

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gemalto/flume"
	"osrs-disc-bot/util"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

func (s *Service) handleSpeedSubmission(session *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid).With("user", i.Member.User.Username))
		defer func() { s.tid++ }()
		returnMessage := s.handleSpeedSubmissionCommand(ctx, session, i)
		err := util.InteractionRespond(session, i, returnMessage)
		if err != nil {
			s.log.Error("Failed to send interaction response: " + err.Error())
		}
	case discordgo.InteractionApplicationCommandAutocomplete:
		s.handleSpeedSubmissionAutocomplete(session, i)
	}
}

func (s *Service) handleSpeedSubmissionCommand(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) string {
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
		// If none of these are submitted, check to see if the image is dragged in as an attachment
		if i.Message != nil && len(i.Message.Attachments) == 0 {
			util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "No screenshot has been submitted")
			return "No image has been submitted - please provide either a screenshot or an imgur link in their respective sections."
		} else {
			screenshot = i.Message.Attachments[0].ProxyURL
		}
	} else if len(screenshot) > 0 && len(imgurUrl) > 0 {
		util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Two screenshots has been submitted")
		return "Two images has been submitted - please provide either a screenshot or an imgur link in their respective sections, not both."
	} else if len(imgurUrl) > 0 {
		if !strings.Contains(imgurUrl, "https://i.imgur.com") {
			util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Incorrect link used in imgur url submission: "+imgurUrl)
			return "Incorrect link used in imgur url submission, please use https://i.imgur.com when submitting using the imgur url option."
		} else {
			url = imgurUrl
		}
	} else {
		url = screenshot
	}

	// Ensure the player used is valid
	// Split the names into an array by , then make an empty array with those names as keys for an easier lookup
	// instead of running a for loop inside a for loop when adding Clan Points
	whitespaceStrippedMessage := util.WhiteStripCommas(playersInvolved)
	logger.Debug("Submitted names: " + whitespaceStrippedMessage)
	names := strings.Split(whitespaceStrippedMessage, ",")
	for _, name := range names {
		if _, ok := s.cp[name]; !ok {
			// Check to see if one of the names is not a main
			logger.Debug("Player " + name + " is not a main, determining main account...")
			discordId := s.members[name].DiscordId

			for user, member := range s.members {
				if discordId == member.DiscordId && member.Main {
					util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Player "+name+"'s main is: "+user+" - resubmit is required")
					return "Player " + name + "'s main is: " + user + ". Please resubmit using the main username."
				}
			}

			// We have a submission for an unknown person, throw an error
			util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Unknown player submitted: "+name)
			return "Please ensure all the names are correct or sign-up the following person: " + name
		}
	}

	// Ensure the boss name is okay
	if _, ok := s.speed[boss]; !ok {
		util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Incorrect boss name: "+boss)
		return "Incorrect boss name. Please look ensure to select one of the options for boss names."
	}

	// Ensure the format is hh:mm:ss:mmm
	reg := regexp.MustCompile("^\\d\\d:\\d\\d:\\d\\d\\.\\d\\d$")
	if !reg.Match([]byte(speedTime)) {
		util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Invalid time format: "+speedTime)
		return "Incorrect time format. Please use the following format: hh:mm:ss.ms"
	}

	msgToBeApproved := fmt.Sprintf("<@&1194691758353821847>\nSubmitter: %s\nUser Id: %s\nBoss Name: %s\nTime: %s\nPlayers Involved: %s\n%+v", i.Member.User.Username, i.Member.User.ID, boss, speedTime, playersInvolved, url)

	// If we have the submission is valid, send the submission information to the admin channel
	msg, err := session.ChannelMessageSend(s.config.DiscSpeedApprovalChan, msgToBeApproved)
	if err != nil {
		util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Failed to send message to admin channel: "+err.Error())
		return "Issue with submitting the speed submission, please contact a dev to fix this issue."
	} else {
		logger.Info("Submission sent to moderators for approval")
		// Add a check and x reaction to the message to accept or reject the submission
		session.MessageReactionAdd(s.config.DiscSpeedApprovalChan, msg.ID, "✅")
		session.MessageReactionAdd(s.config.DiscSpeedApprovalChan, msg.ID, "❌")
	}

	// If nothing wrong happened, send a happy message back to the submitter
	return "Speed submission successfully submitted! Awaiting approval from a moderator!"
}

func (s *Service) handleSpeedSubmissionAutocomplete(session *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	var choices []*discordgo.ApplicationCommandOptionChoice
	switch {
	// In this case there are multiple autocomplete options. The Focused field shows which option user is focused on.
	case data.Options[0].Focused:
		for category := range util.HofSpeedCategories {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  category,
				Value: category,
			})
		}
	case data.Options[1].Focused:
		for _, boss := range s.speedCategory[data.Options[0].Value.(string)] {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  boss,
				Value: boss,
			})
		}
	}

	err := util.InteractionRespondChoices(session, i, choices)
	if err != nil {
		s.log.Error("Failed to handle speed autocomplete: " + err.Error())
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
		submitterId, _ := strconv.Atoi(msg.Content[index+9 : index2-1])

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

		// Delete the screenshot in the page
		err := session.ChannelMessageDelete(s.config.DiscSpeedApprovalChan, r.MessageID)
		if err != nil {
			util.LogError(logger, s.config.DiscAuditChan, session, r.Member.User.Username, r.Member.User.AvatarURL(""), "Failed to delete cp approval message: "+err.Error())
		}

		submissionUrl := ""

		// If the url is an imgur link, skip uploading to imgur
		if strings.Contains(msg.Content, "https://i.imgur.com") {
			submissionUrl = msg.Content[index2:]
		} else {
			// Retrieve the bytes of the image
			resp, err := s.client.Get(msg.Embeds[0].URL)
			if err != nil {
				util.LogError(logger, s.config.DiscAuditChan, session, r.Member.User.Username, r.Member.User.AvatarURL(""), "Failed to download discord image: "+err.Error())
				return
			}
			defer resp.Body.Close()

			// Retrieve the access token
			accessToken, err := s.imageservice.GetNewAccessToken(ctx, s.config.ImgurRefreshToken, s.config.ImgurClientId, s.config.ImgurClientSecret)
			if err != nil {
				logger.Debug("Failed to get imgur access token, will retry...")
				// We will retry 10 times to get a new access token
				counter := 1
				for err != nil {
					logger.Debug("Failed to get imgur access token, will retry (attempt " + strconv.Itoa(counter) + ")")
					if counter == 11 {
						util.LogError(logger, s.config.DiscAuditChan, session, r.Member.User.Username, r.Member.User.AvatarURL(""), "Failed to get access token for imgur: "+err.Error())
						return
					}
					accessToken, err = s.imageservice.GetNewAccessToken(ctx, s.config.ImgurRefreshToken, s.config.ImgurClientId, s.config.ImgurClientSecret)
					if err != nil {
						counter++
						continue
					} else {
						break
					}
				}
			}
			submissionUrl = s.imageservice.Upload(ctx, accessToken, resp.Body)
		}
		submissionTime := time.Now().Format("2006-01-02 15:04:05")
		s.speedscreenshots[submissionTime] = util.SpeedScInfo{BossName: bossName, Time: speedTime, PlayersInvolved: playersInvolved, URL: submissionUrl}

		// Only change the current top speed if it's faster
		t := util.CalculateTime(speedTime)

		// If the submission time is faster than the current speed time for the boss, update it
		if t.Before(s.speed[bossName].Time) {
			logger.Info("NEW TIME FOR BOSS: " + bossName)
			// Add message into new fastest time channel
			newTimeMsg := fmt.Sprintf("New record for **%s** has been set!\n", bossName)
			newTimePlayers := strings.Split(util.WhiteStripCommas(playersInvolved), ",")
			oldTimePlayers := strings.Split(util.WhiteStripCommas(s.speed[bossName].PlayersInvolved), ",")

			for _, player := range newTimePlayers {
				newTimeMsg = newTimeMsg + "<@" + strconv.Itoa(s.members[player].DiscordId) + ">,"
			}
			newTimeMsg = newTimeMsg[:len(newTimeMsg)-1]
			newTimeMsg = newTimeMsg + " has beaten "

			for _, player := range oldTimePlayers {
				if strings.Compare(player, "null") == 0 {
					newTimeMsg = newTimeMsg + "null,"
					continue
				}
				newTimeMsg = newTimeMsg + "<@" + strconv.Itoa(s.members[player].DiscordId) + ">,"
			}
			newTimeMsg = newTimeMsg[:len(newTimeMsg)-1]
			newTimeMsg = newTimeMsg + "!\nOld time: " + s.speed[bossName].Time.Format("15:04:05.00") + "\nNew Time: " + t.Format("15:04:05.00") + fmt.Sprintf("\n%s", submissionUrl)

			_, err := session.ChannelMessageSend(s.config.DiscNewFastestTimeChan, newTimeMsg)
			if err != nil {
				util.LogError(logger, s.config.DiscAuditChan, session, r.Member.User.Username, r.Member.User.AvatarURL(""), "Failed to send message to cp information channel: "+err.Error())
			}

			logger.Info(fmt.Sprintf("Old time: %+v", s.speed[bossName].Time.Format("15:04:05.00")))
			logger.Info(fmt.Sprintf("New Time: %+v", t.Format("15:04:05.00")))

			// Determine category
			for category, bosses := range s.speedCategory {
				if slices.Contains(bosses, bossName) {
					s.speed[bossName] = util.SpeedInfo{Time: t, PlayersInvolved: playersInvolved, URL: submissionUrl, Category: category}
				}
			}

			// Update the Clan Points
			names := strings.Split(util.WhiteStripCommas(playersInvolved), ",")
			for _, name := range names {
				logger.Debug("Adding Clan Point to: " + name)
				s.cp[name] += 1
			}

			// Update the cp leaderboard
			s.updateCpLeaderboard(ctx, session, r.Member.User)

			// Update the boss leaderboard that was updated
			s.updateSpeedHOF(ctx, session, r.Member.User, s.speed[bossName].Category)

		} else {
			logger.Info("KEEP TIME FOR BOSS: " + bossName)
			logger.Info(fmt.Sprintf("Current time: %+v", s.speed[bossName].Time.Format("15:04:05.00")))
			logger.Info(fmt.Sprintf("Submitted Time: %+v", t.Format("15:04:05.00")))
		}

		// Send feedback to user
		channel := s.checkOrCreateFeedbackChannel(ctx, session, submitter, submitterId, "", r.Member.User)
		index = strings.Index(msg.Content, "Boss Name:")
		index2 = strings.Index(msg.Content, "https://cdn.discordapp.com/ephemeral-attachments")
		feedBackMsg := "<@" + strconv.Itoa(submitterId) + ">\nYour speed submission has been accepted\n\n" + msg.Content[index:index2] + "\n" + submissionUrl
		_, err = session.ChannelMessageSend(channel, feedBackMsg)
		if err != nil {
			util.LogError(logger, s.config.DiscAuditChan, session, r.Member.User.Username, r.Member.User.AvatarURL(""), "Failed to send message to cp information channel: "+err.Error())
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
		submitterId, _ := strconv.Atoi(msg.Content[index+9 : index2-1])

		submissionUrl := ""

		// If the url is an imgur link, skip uploading to imgur
		if strings.Contains(msg.Content, "https://i.imgur.com") {
			submissionUrl = msg.Content[index2:]
		} else {
			// Retrieve the bytes of the image
			resp, err := s.client.Get(msg.Embeds[0].URL)
			if err != nil {
				util.LogError(logger, s.config.DiscAuditChan, session, r.Member.User.Username, r.Member.User.AvatarURL(""), "Failed to download discord image: "+err.Error())
				return
			}
			defer resp.Body.Close()

			// Retrieve the access token
			accessToken, err := s.imageservice.GetNewAccessToken(ctx, s.config.ImgurRefreshToken, s.config.ImgurClientId, s.config.ImgurClientSecret)
			if err != nil {
				// We will retry 10 times to get a new access token
				counter := 0
				for err != nil {
					if counter == 10 {
						util.LogError(logger, s.config.DiscAuditChan, session, r.Member.User.Username, r.Member.User.AvatarURL(""), "Failed to get access token for imgur: "+err.Error())
						return
					}
					accessToken, err = s.imageservice.GetNewAccessToken(ctx, s.config.ImgurRefreshToken, s.config.ImgurClientId, s.config.ImgurClientSecret)
					if err != nil {
						counter++
					} else {
						break
					}
				}
			}
			submissionUrl = s.imageservice.Upload(ctx, accessToken, resp.Body)
		}

		// Send feedback to user
		channel := s.checkOrCreateFeedbackChannel(ctx, session, submitter, submitterId, "", r.Member.User)
		index = strings.Index(msg.Content, "Boss Name:")
		index2 = strings.Index(msg.Content, "https://cdn.discordapp.com/ephemeral-attachments")
		feedBackMsg := "<@" + strconv.Itoa(submitterId) + ">\nYour speed submission has been rejected\n\n" + msg.Content[index:index2] + "\n" + submissionUrl
		_, err := session.ChannelMessageSend(channel, feedBackMsg)
		if err != nil {
			util.LogError(logger, s.config.DiscAuditChan, session, r.Member.User.Username, r.Member.User.AvatarURL(""), "Failed to send message to cp information channel: "+err.Error())
		}

		// Delete the screenshot in the page
		err = session.ChannelMessageDelete(s.config.DiscSpeedApprovalChan, r.MessageID)
		if err != nil {
			util.LogError(logger, s.config.DiscAuditChan, session, r.Member.User.Username, r.Member.User.AvatarURL(""), "Failed to delete cp approval message: "+err.Error())
		}
	}
}
