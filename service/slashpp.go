package service

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gemalto/flume"
	"osrs-disc-bot/util"
	"strconv"
	"strings"
	"time"
)

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
		// If none of these are submitted, check to see if the image is dragged in as an attachment
		if i.Message != nil && len(i.Message.Attachments) == 0 {
			logger.Error("No screenshot has been submitted")
			return "No image has been submitted - please provide either a screenshot or an imgur link in their respective sections."
		} else {
			screenshot = i.Message.Attachments[0].ProxyURL
		}
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
	whitespaceStrippedMessage := util.WhiteStripCommas(playersInvolved)
	logger.Debug("Submitted names: " + whitespaceStrippedMessage)
	names := strings.Split(whitespaceStrippedMessage, ",")
	var checkedNames []string
	for _, name := range names {
		if _, ok := s.cp[name]; !ok {
			// Check to see if one of the names is not a main
			logger.Debug("Player " + name + " is not a main, determining main account...")
			discordId := s.members[name].DiscordId

			for user, member := range s.members {
				if discordId == member.DiscordId && member.Main {
					logger.Error("Player " + name + "'s main is: " + user + " - resubmit is required")
					return "Player " + name + "'s main is: " + user + ". Please resubmit using the main username."
				}
			}

			// We have a submission for an unknown person, throw an error
			logger.Error("Unknown player submitted: " + name)
			return "Please ensure all the names are correct or sign-up the following person: " + name
		} else {
			checkedNames = append(checkedNames, name)
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

func (s *Service) handlePPApproval(ctx context.Context, session *discordgo.Session, r *discordgo.MessageReactionAdd) {
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
		submitterId, _ := strconv.Atoi(msg.Content[index+9 : index2-1])

		index = strings.Index(msg.Content, "Involved:")
		index2 = strings.Index(msg.Content, "https://")
		playersInvolved := msg.Content[index+10 : index2-1]
		logger.Info("CP Approved for: " + playersInvolved)

		submissionUrl := ""

		// Delete the screenshot in the page
		err := session.ChannelMessageDelete(s.config.DiscCpApprovalChan, r.MessageID)
		if err != nil {
			logger.Error("Failed to delete cp approval message: " + err.Error())
		}
		logger.Debug("Successfully added CPs for: " + playersInvolved)

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
			accessToken, err := s.imageservice.GetNewAccessToken(ctx, s.config.ImgurRefreshToken, s.config.ImgurClientId, s.config.ImgurClientSecret)
			if err != nil {
				// We will retry 10 times to get a new access token
				counter := 0
				for err != nil {
					if counter == 10 {
						logger.Error("Failed to get access token for imgur: " + err.Error())
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

		s.cpscreenshots[time.Now().Format("2006-01-02 15:04:05")] = util.CpScInfo{
			PlayersInvolved: playersInvolved,
			URL:             submissionUrl,
		}

		// Update the Ponies Points
		names := strings.Split(util.WhiteStripCommas(playersInvolved), ",")
		for _, name := range names {
			logger.Debug("Adding Ponies Point to: " + name)
			s.cp[name] += 1
		}

		// Update the cp leaderboard
		s.updatePpLeaderboard(ctx, session)

		// Send feedback to user
		channel := s.checkOrCreateFeedbackChannel(ctx, session, submitter, submitterId, "")
		index = strings.Index(msg.Content, "Players Involved:")
		feedBackMsg := "<@" + strconv.Itoa(submitterId) + ">\nYour ponies point submission has been accepted\n\n" + msg.Content[index:]
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
		submitterId, _ := strconv.Atoi(msg.Content[index+9 : index2-1])

		// Send feedback to user
		channel := s.checkOrCreateFeedbackChannel(ctx, session, submitter, submitterId, "")
		index = strings.Index(msg.Content, "Players Involved:")
		feedBackMsg := "<@" + strconv.Itoa(submitterId) + ">\nYour ponies point submission has been rejected\n\n" + msg.Content[index:]
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
