package service

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"osrs-disc-bot/util"
	"strconv"
	"strings"
)

func (s *Service) listenForLootLog(session *discordgo.Session, message *discordgo.MessageCreate) {
	player := ""

	// Will always add a clan point
	if strings.Contains(message.Content, " just received a new pet!") {
		player = strings.Split(message.Content, " just received a new pet!")[0]
		s.log.Info("Loot Log Pet submission for " + player)
		if _, ok := s.cp[player]; ok {
			s.cp[player] += 1
		} else {
			return
		}
		s.updateCpLeaderboard(context.Background(), session)
	} else if strings.Contains(message.Content, " just received a valuable drop:") {
		index1 := strings.Index(message.Content, "drop:")
		index2 := strings.Index(message.Content, "!")

		item := message.Content[index1+6 : index2]
		if _, ok := util.LootLogClanPoint[item]; ok {
			player = strings.Split(message.Content, " just received a valuable drop:")[0]
			s.log.Info("Loot Log Valuable drop submission for " + player + " with item: " + item)
			if _, ok := s.cp[player]; ok {
				s.cp[player] += 1
			} else {
				return
			}
			s.updateCpLeaderboard(context.Background(), session)
		} else {
			return
		}
	} else if strings.Contains(message.Content, " just received a new collection log item:") {
		index1 := strings.Index(message.Content, "item:")
		index2 := strings.Index(message.Content, "!")

		item := message.Content[index1+6 : index2]
		if _, ok := util.LootLogClanPoint[item]; ok {
			player = strings.Split(message.Content, " just received a new collection log item:")[0]
			s.log.Info("Loot Log collection log submission for " + player + " with item: " + item)
			if _, ok := s.cp[player]; ok {
				s.cp[player] += 1
			} else {
				return
			}
			s.updateCpLeaderboard(context.Background(), session)
		} else {
			return
		}
	} else {
		return
	}

	// Retrieve the bytes of the image
	resp, err := s.client.Get(message.Attachments[0].ProxyURL)
	if err != nil {
		s.log.Error("Failed to download discord image: " + err.Error())
		return
	}
	defer resp.Body.Close()

	// Retrieve the access token
	accessToken, err := s.imgur.GetNewAccessToken(context.Background(), s.config.ImgurRefreshToken, s.config.ImgurClientId, s.config.ImgurClientSecret)
	if err != nil {
		// We will retry 10 times to get a new access token
		counter := 0
		for err != nil {
			s.log.Error("Failed to get access token for imgur, retrying... (count: " + strconv.Itoa(counter) + ")")
			if counter == 15 {
				s.log.Error("Failed to get access token for imgur: " + err.Error())
				return
			}
			accessToken, err = s.imgur.GetNewAccessToken(context.Background(), s.config.ImgurRefreshToken, s.config.ImgurClientId, s.config.ImgurClientSecret)
			if err != nil {
				counter++
			} else {
				break
			}
		}
	}
	submissionUrl := s.imgur.Upload(context.Background(), accessToken, resp.Body)
	s.log.Info("Successfully uploaded imgur url: " + submissionUrl)
	s.cpscreenshots[submissionUrl] = player
}
