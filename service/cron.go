package service

import (
	"context"
	"osrs-disc-bot/util"
	"strconv"

	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
	"osrs-disc-bot/util"
	"sort"
	"strconv"
)

/*
updateHOF will iterate over all the HallOfFameRequestInfos, grab the podium from temple for each
of the bosses, sort them, and make the discord call to create the emded with the boss name, image,
and podium finish with [kc]
*/
func (s *Service) updateHOF(ctx context.Context, session *discordgo.Session, allRequestInfo ...util.HallOfFameRequestInfo) {
	hofLeaderboard := make(map[string]int)

	for _, requestInfo := range allRequestInfo {
		s.log.Info("Running HOF update for Boss: " + requestInfo.Name)
		// First, delete all the messages within the channel
		messages, err := session.ChannelMessages(requestInfo.DiscChan, 50, "", "", "")
		if err != nil {
			s.log.Error("Failed to get all messages for deletion from channel: " + requestInfo.Name)
			return err
		}
		var messageIDs []string
		for _, message := range messages {
			messageIDs = append(messageIDs, message.ID)
		}
		err = session.ChannelMessagesBulkDelete(requestInfo.DiscChan, messageIDs)
		if err != nil {
			s.log.Error("Failed to delete all messages from channel: " + requestInfo.Name)
			return err
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
				return err
			}
		}
	}

	s.updateHOFLeaderboard(ctx, session, hofLeaderboard)
}

func (s *Service) addToHOFLeaderboard(hofLeaderboard map[string]int, player string, points int) {
	if _, ok := hofLeaderboard[player]; ok {
		hofLeaderboard[player] = hofLeaderboard[player] + points
	} else {
		hofLeaderboard[player] = points
	}
}

func (s *Service) updateHOFLeaderboard(ctx context.Context, session *discordgo.Session, hofLeaderboard map[string]int) {
	keys := make([]string, 0, len(hofLeaderboard))
	for key := range hofLeaderboard {
		keys = append(keys, key)
	}

	// Sort the map based on the values
	sort.SliceStable(keys, func(i, j int) bool {
		return hofLeaderboard[keys[i]] > hofLeaderboard[keys[j]]
	})

	// First, delete all the messages within the channel
	messages, err := session.ChannelMessages(s.config.DiscHOFLeaderboardChan, 50, "", "", "")
	if err != nil {
		s.log.Error("Failed to get all messages for deletion from the leagues podium channel")
		return
	}
	var messageIDs []string
	for _, message := range messages {
		messageIDs = append(messageIDs, message.ID)
	}
	err = session.ChannelMessagesBulkDelete(s.config.DiscHOFLeaderboardChan, messageIDs)
	if err != nil {
		s.log.Error("Failed to delete all messages from the leagues podium channel")
		return
	}

	// Iterate over the players to get the different places for users to create the placements
	// Create the leaderboard message that will be sent
	placements := ""
	for placement, player := range keys {
		placements = placements + strconv.Itoa(placement+1) + ") " + player + " [" + strconv.Itoa(hofLeaderboard[player]) + "]\n"
	}

	// Send the Discord Embed message for the leaderboard
	_, err = session.ChannelMessageSendEmbed(s.config.DiscHOFLeaderboardChan, embed.NewEmbed().
		SetTitle("Ponies HOF Leaderboard").
		SetDescription(placements).
		SetColor(0x1c1c1c).SetThumbnail("https://i.imgur.com/O4NzB95.png").MessageEmbed)
	if err != nil {
		s.log.Error("Failed to send message for leagues podium")
		return
	}

	// Send the Discord Embed message for instructions on how the rankings work
	var msg string
	msg = msg + "In order to get onto this leaderboard, you must have a podium finish of one of the HOF Bosses.\n\n"
	msg = msg + "3 points for :first_place:\n2 points for :second_place:\n1 points for :third_place:"
	_, err = session.ChannelMessageSendEmbed(s.config.DiscHOFLeaderboardChan, embed.NewEmbed().
		SetTitle("How To Get Onto The Collection Log HOF").
		SetDescription(msg).
		SetColor(0x1c1c1c).SetThumbnail("https://i.imgur.com/O4NzB95.png").MessageEmbed)
	if err != nil {
		return
	}
}

/*
updateColLog will use all the users within the in memory submission map to create the podium
from collectionlog.net and their rankings. It will create an embed with the top 10 placements in
discord.
*/
func (s *Service) updateColLog(ctx context.Context, session *discordgo.Session) error {
	s.log.Info("Running collection log hiscores update...")

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
		s.log.Error("Failed to retrieve last 10 messages" + err.Error())
		return err
	}
	var messageIDs []string
	for _, message := range messages {
		messageIDs = append(messageIDs, message.ID)
	}
	err = session.ChannelMessagesBulkDelete(s.config.DiscColChan, messageIDs)
	if err != nil {
		s.log.Error("Failed to bulk delete the collection log" + err.Error())
		return err
	}

	// Send the Discord Embed message for collection log
	_, err = session.ChannelMessageSendEmbed(s.config.DiscColChan, embed.NewEmbed().
		SetTitle("Collection Log Ranking").
		SetDescription(placements).
		SetColor(0x1c1c1c).SetThumbnail("https://i.imgur.com/otTd8Dg.png").MessageEmbed)
	if err != nil {
		s.log.Error("Failed to send discord emded message" + err.Error())
		return err
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
		s.log.Error("Failed to send discord emded message" + err.Error())
		return err
	}

	s.log.Info("Collection log hiscores update successful.")
	return nil
}

func (s *Service) updateLeagues(ctx context.Context, session *discordgo.Session) {
	s.log.Info("Running leagues hiscores update.")

	// First, delete all the messages within the channel
	messages, err := session.ChannelMessages(s.config.DiscLeaguesChan, 50, "", "", "")
	if err != nil {
		s.log.Error("Failed to get all messages for deletion from the leagues podium channel.")
		return
	}
	var messageIDs []string
	for _, message := range messages {
		messageIDs = append(messageIDs, message.ID)
	}
	err = session.ChannelMessagesBulkDelete(s.config.DiscLeaguesChan, messageIDs)
	if err != nil {
		s.log.Error("Failed to delete all messages from the leagues podium channel.")
		return
	}

	leaguesPodium, ranking := s.runescape.GetLeaguesPodiumFromRS(ctx, s.submissions)
	// Iterate over the players to get the different places for users to create the placements
	// Create the leaderboard message that will be sent
	placements := "<:Executioner:1176594739366219806> __**TIER 8**__ <:Executioner:1176594739366219806>\n"
	tier := 8
	var t8, t7, t6, t5, t4, t3, t2 = 24000, 15000, 7500, 4000, 2000, 1200, 500

	for placement, player := range ranking {
		if tier == 8 && leaguesPodium[player] < t8 {
			tier = 7
			placements = placements + "\n<:Weapon_Master:1176595216338260060>__**TIER 7**__<:Weapon_Master:1176595216338260060>\n"
		} else if tier == 7 && leaguesPodium[player] < t7 {
			tier = 6
			placements = placements + "\n<:Ruinous_Powers:1176595214828318741>__**TIER 6**__<:Ruinous_Powers:1176595214828318741>\n"
		} else if tier == 6 && leaguesPodium[player] < t6 {
			tier = 5
			placements = placements + "\n<:Bloodthirsty:1176595214014623834>__**TIER 5**__<:Bloodthirsty:1176595214014623834>\n"
		} else if tier == 5 && leaguesPodium[player] < t5 {
			tier = 4
			placements = placements + "\n<:Brawlers_Resolve:1176595213280628856>__**TIER 4**__<:Brawlers_Resolve:1176595213280628856>\n"
		} else if tier == 4 && leaguesPodium[player] < t4 {
			tier = 3
			placements = placements + "\n<:Fire_Sale:1176595212219453441>__**TIER 3**__<:Fire_Sale:1176595212219453441>\n"
		} else if tier == 3 && leaguesPodium[player] < t3 {
			tier = 2
			placements = placements + "\n<:Globetrotter:1176595211833577482>__**TIER 2**__<:Globetrotter:1176595211833577482>\n"
		} else if tier == 2 && leaguesPodium[player] < t2 {
			tier = 1
			placements = placements + "\n__**TIER 1**__\n"
		}

		placements = placements + strconv.Itoa(placement+1) + ") " + player + " [" + strconv.Itoa(leaguesPodium[player]) + "] "

		var bronze, iron, steel, mithril, adamant, runeTier, dragon = 2500, 5000, 10000, 18000, 28000, 42000, 56000
		points := leaguesPodium[player]
		switch true {
		case points >= bronze && points < iron:
			placements = placements + "<:bronze_trophy:1178124557933101056>\n"
			break
		case points >= iron && points < steel:
			placements = placements + "<:iron_trophy:1178124556846780458>\n"
			break
		case points >= steel && points < mithril:
			placements = placements + "<:steel_trophy:1178124555718508617>\n"
			break
		case points >= mithril && points < adamant:
			placements = placements + "<:mithril_trophy:1178124554820931755>\n"
			break
		case points >= adamant && points < runeTier:
			placements = placements + "<:adamant_trophy:1178124552971231365>\n"
			break
		case points >= runeTier && points < dragon:
			placements = placements + "<:rune_trophy:1178124551188664350>\n"
			break
		case points >= dragon:
			placements = placements + "<:dragon_trophy:1178124549141839893>\n"
			break
		default:
			placements = placements + "\n"
		}
	}

	// Send the Discord Embed message for the boss podium finish
	_, err = session.ChannelMessageSendEmbed(s.config.DiscLeaguesChan, embed.NewEmbed().
		SetTitle("Ponies League Standings").
		SetDescription(placements).
		SetColor(0x1c1c1c).SetThumbnail("https://i.imgur.com/O4NzB95.png").MessageEmbed)
	if err != nil {
		s.log.Error("Failed to send message for leagues podium.")
		return
	}

	s.log.Info("Leagues hiscores update successful.")
}
