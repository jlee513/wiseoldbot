package service

import (
	"context"
	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
	"osrs-disc-bot/util"
	"strconv"
)

/*
updateKcHOF will iterate over all the HallOfFameRequestInfos, grab the podium from temple for each
of the bosses, sort them, and make the discord call to create the emded with the boss name, image,
and podium finish with [kc]
*/
func (s *Service) updateKcHOF(ctx context.Context, session *discordgo.Session, allRequestInfo ...util.HofRequestInfo) {
	hofLeaderboard := make(map[string]int)

	for _, requestInfo := range allRequestInfo {
		s.log.Info("Running HOF update for Boss: " + requestInfo.Name)
		// First, delete all the messages within the channel
		messages, err := session.ChannelMessages(requestInfo.DiscChan, 50, "", requestInfo.AfterId, "")
		if err != nil {
			s.log.Error("Failed to get all messages for deletion from channel: " + requestInfo.Name)
			return
		}
		var messageIDs []string
		for _, message := range messages {
			messageIDs = append(messageIDs, message.ID)
		}
		err = session.ChannelMessagesBulkDelete(requestInfo.DiscChan, messageIDs)
		if err != nil {
			s.log.Error("Failed to delete all messages from channel: " + requestInfo.Name + ", will try one by one")
			for _, message := range messageIDs {
				err = session.ChannelMessageDelete(requestInfo.DiscChan, message)
				if err != nil {
					s.log.Error("Failed to delete messages one by one from channel: " + requestInfo.Name)
					return
				}
			}
		}

		// Now add all the bosses
		for _, bossInfo := range requestInfo.Bosses {
			podium, rankings := s.temple.GetPodiumFromTemple(ctx, bossInfo.BossName)
			s.log.Debug("Updating " + podium.Data.BossName)

			// Iterate over the players to get the different places for users to create the placements
			placements := ""
			for _, k := range rankings {
				switch k {
				case 1:
					placements = placements + ":first_place: "
					s.addToHOFLeaderboard(hofLeaderboard, podium.Data.Players[k].Username, 3)
					break
				case 2:
					placements = placements + ":second_place: "
					s.addToHOFLeaderboard(hofLeaderboard, podium.Data.Players[k].Username, 2)
					break
				case 3:
					placements = placements + ":third_place: "
					s.addToHOFLeaderboard(hofLeaderboard, podium.Data.Players[k].Username, 1)
					break
				}
				placements = placements + podium.Data.Players[k].Username + " [" + strconv.Itoa(podium.Data.Players[k].Kc) + "]\n"
			}

			// Send the Discord Embed message for the boss podium finish
			_, err = session.ChannelMessageSendEmbed(requestInfo.DiscChan, embed.NewEmbed().
				SetTitle(podium.Data.BossName).
				SetDescription(placements).
				SetColor(0x1c1c1c).SetThumbnail(bossInfo.ImageLink).MessageEmbed)
			if err != nil {
				s.log.Error("Failed to send message for boss: " + podium.Data.BossName)
				return
			}
		}
	}

	s.updateHOFLeaderboard(ctx, session, hofLeaderboard)
}

func (s *Service) updateSpeedHOF(ctx context.Context, session *discordgo.Session, allRequestInfo ...util.SpeedsRequestInfo) {
	for _, requestInfo := range allRequestInfo {
		s.log.Info("Running Speed HOF update for Boss: " + requestInfo.Name)
		// First, delete all the messages within the channel
		messages, err := session.ChannelMessages(requestInfo.DiscChan, 50, "", requestInfo.AfterId, "")
		if err != nil {
			s.log.Error("Failed to get all messages for deletion from channel: " + requestInfo.Name)
			return
		}
		var messageIDs []string
		for _, message := range messages {
			messageIDs = append(messageIDs, message.ID)
		}
		err = session.ChannelMessagesBulkDelete(requestInfo.DiscChan, messageIDs)
		if err != nil {
			s.log.Error("Failed to delete all messages from channel: " + requestInfo.Name + ", will try one by one")
			for _, message := range messageIDs {
				err = session.ChannelMessageDelete(requestInfo.DiscChan, message)
				if err != nil {
					s.log.Error("Failed to delete messages one by one from channel: " + requestInfo.Name)
					return
				}
			}
		}

		// Now add all the bosses
		for _, bossInfo := range requestInfo.Bosses {
			// Get the speed info
			speed := s.speed[bossInfo.BossName]

			// Send the Discord Embed message for the boss podium finish
			_, err = session.ChannelMessageSendEmbed(requestInfo.DiscChan, embed.NewEmbed().
				SetTitle(bossInfo.BossName).
				SetDescription("**Players:** "+speed.PlayersInvolved+"\n**Time:** "+speed.Time.Format("15:04:05.00")).
				SetColor(0x1c1c1c).SetThumbnail(speed.URL).MessageEmbed)
			if err != nil {
				s.log.Error("Failed to send message for boss: " + bossInfo.BossName)
				return
			}
		}
	}
}
