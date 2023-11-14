package service

import (
	"context"
	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
	"osrs-disc-bot/util"
	"strconv"
)

/*
updateHOF will iterate over all the HallOfFameRequestInfos, grab the podium from temple for each
of the bosses, sort them, and make the discord call to create the emded with the boss name, image,
and podium finish with [kc]
*/
func (s *Service) updateHOF(ctx context.Context, session *discordgo.Session, allRequestInfo ...util.HallOfFameRequestInfo) {
	for _, requestInfo := range allRequestInfo {
		// First, delete all the messages within the channel
		messages, err := session.ChannelMessages(requestInfo.DiscChan, 50, "", "", "")
		if err != nil {
			return
		}
		var messageIDs []string
		for _, message := range messages {
			messageIDs = append(messageIDs, message.ID)
		}
		err = session.ChannelMessagesBulkDelete(requestInfo.DiscChan, messageIDs)
		if err != nil {
			return
		}

		// Now add all the bosses
		for bossIdForTemple, imageURL := range requestInfo.Bosses {
			podium, rankings := s.temple.GetPodiumFromTemple(ctx, bossIdForTemple)

			// Iterate over the players to get the different places for users to create the placements
			placements := ""
			for _, k := range rankings {
				switch k {
				case 1:
					placements = placements + ":first_place: "
					break
				case 2:
					placements = placements + ":second_place: "
					break
				case 3:
					placements = placements + ":third_place: "
					break
				}
				placements = placements + podium.Data.Players[k].Username + " [" + strconv.Itoa(podium.Data.Players[k].Kc) + "]\n"
			}

			// Send the Discord Embed message for the boss podium finish
			_, err = session.ChannelMessageSendEmbed(requestInfo.DiscChan, embed.NewEmbed().
				SetTitle(podium.Data.BossName).
				SetDescription(placements).
				SetColor(0x1c1c1c).SetThumbnail(imageURL).MessageEmbed)
			if err != nil {
				return
			}
		}
		return
	}
}

/*
updateColLog will use all the users within the in memory submission map to create the podium
from collectionlog.net and their rankings. It will create an embed with the top 10 placements in
discord.
*/
func (s *Service) updateColLog(ctx context.Context, session *discordgo.Session) {
	podium, ranking := s.collectionLog.RetrieveCollectionLogAndOrder(ctx, s.submissions)

	// Create the leaderboard message that will be sent
	placements := ""
	for placement, k := range ranking {
		switch placement {
		case 0:
			placements = placements + ":one: "
			break
		case 1:
			placements = placements + ":two: "
			break
		case 2:
			placements = placements + ":three: "
			break
		case 3:
			placements = placements + ":four: "
			break
		case 4:
			placements = placements + ":five: "
			break
		case 5:
			placements = placements + ":six: "
			break
		case 6:
			placements = placements + ":seven: "
			break
		case 7:
			placements = placements + ":eight: "
			break
		case 8:
			placements = placements + ":nine: "
			break
		case 9:
			placements = placements + ":keycap_10: "
			break

		}

		placements = placements + k + " [" + strconv.Itoa(podium[k]) + "]\n"
	}

	// First, delete all the messages within the channel
	messages, err := session.ChannelMessages(s.config.DiscColChan, 10, "", "", "")
	if err != nil {
		return
	}
	var messageIDs []string
	for _, message := range messages {
		messageIDs = append(messageIDs, message.ID)
	}
	err = session.ChannelMessagesBulkDelete(s.config.DiscColChan, messageIDs)
	if err != nil {
		return
	}

	// Send the Discord Embed message for collection log
	_, err = session.ChannelMessageSendEmbed(s.config.DiscColChan, embed.NewEmbed().
		SetTitle("Collection Log Ranking").
		SetDescription(placements).
		SetColor(0x1c1c1c).SetThumbnail("https://i.imgur.com/otTd8Dg.png").MessageEmbed)
	if err != nil {
		return
	}

	// Send the Discord Embed message for instructions on how to get on the collection log hall of fame
	var msg string
	msg = msg + "1. Download the collection-log plugin\n"
	msg = msg + "2. Click the box to \"Allow collectionlog.net connections\"\n"
	msg = msg + "3. Click through the collection log (there will be a * next to the one you still need to click)\n"
	msg = msg + "4. Go to the collection log icon on the sidebar\n"
	msg = msg + "5. Click Account at the top and then upload collection log\n"
	_, err = session.ChannelMessageSendEmbed(s.config.DiscColChan, embed.NewEmbed().
		SetTitle("How To Get Onto The Collection Log HOF").
		SetDescription(msg).
		SetColor(0x1c1c1c).SetThumbnail("https://i.imgur.com/otTd8Dg.png").MessageEmbed)
	if err != nil {
		return
	}
}
