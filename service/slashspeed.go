package service

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gemalto/flume"
	"osrs-disc-bot/util"
	"regexp"
	"strconv"
	"strings"
	"time"
)

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
		reg := regexp.MustCompile("^\\d\\d:\\d\\d:\\d\\d\\.\\d\\d$")
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

		// Delete the screenshot in the page
		err := session.ChannelMessageDelete(s.config.DiscSpeedApprovalChan, r.MessageID)
		if err != nil {
			logger.Error("Failed to delete cp approval message: " + err.Error())
		}

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
			newTimeMsg := fmt.Sprintf("New record for **%s** has been set!\n", bossName)
			newTimePlayers := strings.Split(playersInvolved, ",")
			oldTimePlayers := strings.Split(s.speed[bossName].PlayersInvolved, ",")

			for _, player := range newTimePlayers {
				newTimeMsg = newTimeMsg + "<@" + s.members[player].DiscordId + ">,"
			}
			newTimeMsg = newTimeMsg[:len(newTimeMsg)-1]
			newTimeMsg = newTimeMsg + " has beaten "

			for _, player := range oldTimePlayers {
				if strings.Compare(player, "null") == 0 {
					newTimeMsg = newTimeMsg + "null,"
					continue
				}
				newTimeMsg = newTimeMsg + "<@" + s.members[player].DiscordId + ">,"
			}
			newTimeMsg = newTimeMsg[:len(newTimeMsg)-1]
			newTimeMsg = newTimeMsg + "!\nOld time: " + s.speed[bossName].Time.Format("15:04:05.00") + "\nNew Time: " + t.Format("15:04:05.00") + fmt.Sprintf("\n%s", submissionUrl)

			_, err := session.ChannelMessageSend(s.config.DiscNewFastestTimeChan, newTimeMsg)
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

		// Send feedback to user
		channel := s.checkOrCreateFeedbackChannel(ctx, session, submitter, submitterId, "")
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
		channel := s.checkOrCreateFeedbackChannel(ctx, session, submitter, submitterId, "")
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
