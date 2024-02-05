package service

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/gemalto/flume"
	"osrs-disc-bot/util"
	"strconv"
	"strings"
	"time"
)

func (s *Service) listenForLootLog(ctx context.Context, session *discordgo.Session, message *discordgo.MessageCreate) {
	logger := flume.FromContext(ctx)

	// Remove the bolding to get the string manipulation correct
	msg := strings.Replace(message.Content, "**", "", -1)
	logger.Debug("Message received from loot log channel: " + msg)
	player := ""

	// Will always add a clan point
	if strings.Contains(msg, " just received a new pet!") {
		player = strings.Split(msg, " just received a new pet!")[0]
		logger.Debug("Accepted: Loot Log Pet submission for " + player)
		if _, ok := s.cp[player]; ok {
			s.cp[player] += 1
		} else {
			// Check to see if there is a main player for this user
			logger.Debug("Player " + player + " is not a main, determining main account...")
			discordId := s.members[player].DiscordId

			foundMain := false

			for user, member := range s.members {
				if discordId == member.DiscordId && member.Main {
					logger.Debug("Player " + player + "'s main is: " + user)
					player = user
					s.cp[user] += 1
					s.updateCpLeaderboard(context.Background(), session)
					foundMain = true
				}
			}

			if !foundMain {
				logger.Error("Rejected: Player " + player + " does not exist and neither does a main.")
				return
			}
		}
		s.updateCpLeaderboard(context.Background(), session)
	} else if strings.Contains(msg, " just received a valuable drop:") {
		index1 := strings.Index(msg, "drop:")
		index2 := strings.Index(msg, "!")

		item := msg[index1+6 : index2]
		if _, ok := util.LootLogClanPoint[item]; ok {
			player = strings.Split(msg, " just received a valuable drop:")[0]
			logger.Debug("Accepted: Loot Log Valuable drop submission for " + player + " with item: " + item)
			if _, ok := s.cp[player]; ok {
				s.cp[player] += 1
			} else {
				// Check to see if there is a main player for this user
				logger.Debug("Player " + player + " is not a main, determining main account...")
				discordId := s.members[player].DiscordId

				foundMain := false

				for user, member := range s.members {
					if discordId == member.DiscordId && member.Main {
						logger.Debug("Player " + player + "'s main is: " + user)
						player = user
						s.cp[user] += 1
						s.updateCpLeaderboard(context.Background(), session)
						foundMain = true
					}
				}

				if !foundMain {
					logger.Error("Rejected: Player " + player + " does not exist and neither does a main.")
					return
				}
			}
			s.updateCpLeaderboard(context.Background(), session)
		} else {
			logger.Debug("Rejected: Item " + item + " is not on the list.")
			return
		}
	} else if strings.Contains(msg, " just received a new collection log item:") {
		index1 := strings.Index(msg, "item:")
		index2 := strings.Index(msg, "!")

		item := msg[index1+6 : index2]
		if _, ok := util.LootLogClanPoint[item]; ok {
			player = strings.Split(msg, " just received a new collection log item:")[0]
			logger.Debug("Accepted: Loot Log collection log submission for " + player + " with item: " + item)
			if _, ok := s.cp[player]; ok {
				s.cp[player] += 1
			} else {
				// Check to see if there is a main player for this user
				logger.Debug("Player " + player + " is not a main, determining main account...")
				discordId := s.members[player].DiscordId

				foundMain := false

				for user, member := range s.members {
					if discordId == member.DiscordId && member.Main {
						logger.Debug("Player " + player + "'s main is: " + user)
						player = user
						s.cp[user] += 1
						s.updateCpLeaderboard(context.Background(), session)
						foundMain = true
					}
				}

				if !foundMain {
					logger.Error("Rejected: Player " + player + " does not exist and neither does a main.")
					return
				}
			}
			s.updateCpLeaderboard(context.Background(), session)
		} else {
			logger.Debug("Rejected: Item " + item + " is not on the list.")
			return
		}
	} else {
		return
	}

	// Retrieve the bytes of the image
	resp, err := s.client.Get(message.Attachments[0].ProxyURL)
	if err != nil {
		logger.Error("Failed to download discord image: " + err.Error())
		return
	}
	defer resp.Body.Close()

	// Retrieve the access token
	accessToken, err := s.imageservice.GetNewAccessToken(context.Background(), s.config.ImgurRefreshToken, s.config.ImgurClientId, s.config.ImgurClientSecret)
	if err != nil {
		// We will retry 10 times to get a new access token
		counter := 0
		for err != nil {
			logger.Error("Failed to get access token for imgur, retrying... (count: " + strconv.Itoa(counter) + ")")
			if counter == 15 {
				logger.Error("Failed to get access token for imgur: " + err.Error())
				return
			}
			accessToken, err = s.imageservice.GetNewAccessToken(context.Background(), s.config.ImgurRefreshToken, s.config.ImgurClientId, s.config.ImgurClientSecret)
			if err != nil {
				counter++
			} else {
				break
			}
		}
	}
	submissionUrl := s.imageservice.Upload(context.Background(), accessToken, resp.Body)
	logger.Info("Successfully uploaded lootlog imgur url: " + submissionUrl)

	channel := s.checkOrCreateFeedbackChannel(ctx, session, "", 0, player)
	feedBackMsg := "<@" + strconv.Itoa(s.members[player].DiscordId) + ">\nYour automatic loot log clan point has been accepted\n\nAdded to main user: **" + player + "**\n\n" + message.Content + "\n" + submissionUrl
	_, err = session.ChannelMessageSend(channel, feedBackMsg)
	if err != nil {
		logger.Error("Failed to send message to cp information channel", err)
	}
	s.cpscreenshots[time.Now().Format("2006-01-02 15:04:05")] = util.CpScInfo{
		PlayersInvolved: player,
		URL:             submissionUrl,
	}
}
