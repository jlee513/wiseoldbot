package service

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/gemalto/flume"
	"osrs-disc-bot/util"
	"sort"
	"strconv"
	"strings"
)

func (s *Service) updateLeaderboard(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		logger := flume.FromContext(ctx)
		options := i.ApplicationCommandData().Options[0].Options

		leaderboardName := ""
		threadName := ""
		for _, option := range options {
			switch option.Name {
			case "leaderboard":
				leaderboardName = option.Value.(string)
			case "thread":
				threadName = option.Value.(string)
			}
		}

		switch leaderboardName {
		case "Kc":
			logger.Info("Admin invoked Kc Hall Of Fame Update: ", i.Member.User.Username)
			err := util.InteractionRespond(session, i, "Updating Leaderboard: "+leaderboardName)
			if err != nil {
				logger.Error("Failed to send interaction response: " + err.Error())
			}
			// If kc is updating, always update all of them
			s.updateKcHOF(ctx, session)
		case "Speed":
			logger.Info("Admin invoked Speed Hall Of Fame Update: ", i.Member.User.Username)
			if _, ok := util.HofSpeedCategories[threadName]; ok {
				err := util.InteractionRespond(session, i, "Updating Leaderboard: "+leaderboardName+" thread: "+threadName)
				if err != nil {
					logger.Error("Failed to send interaction response: " + err.Error())
				}
				s.updateSpeedHOF(ctx, session, threadName)
			} else if strings.Compare(threadName, "All") == 0 {
				err := util.InteractionRespond(session, i, "Updating All Speed Leaderboards")
				if err != nil {
					logger.Error("Failed to send interaction response: " + err.Error())
				}
				s.updateSpeedHOF(ctx, session, "TzHaar", "Slayer", "Nightmare", "Nex", "Solo Bosses", "Chambers Of Xeric", "Chambers Of Xeric Challenge Mode", "Theatre Of Blood", "Theatre Of Blood Hard Mode", "Tombs Of Amascut", "Tombs Of Amascut Expert", "Agility", "Desert Treasure 2")
			}
		default:
			err := util.InteractionRespond(session, i, "Unknown leaderboard submitted - please submit a proper leaderboard name")
			if err != nil {
				logger.Error("Failed to send interaction response: " + err.Error())
			}
			logger.Error("Unknown leaderboard submitted - please submit a proper leaderboard name")
		}
	case discordgo.InteractionApplicationCommandAutocomplete:
		logger := flume.FromContext(ctx)
		data := i.ApplicationCommandData()
		var choices []*discordgo.ApplicationCommandOptionChoice
		switch {
		case data.Options[0].Options[0].Focused:
			choices = []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "Kc",
					Value: "Kc",
				},
				{
					Name:  "Speed",
					Value: "Speed",
				},
			}
		// In this case there are multiple autocomplete options. The Focused field shows which option user is focused on.
		case data.Options[0].Options[1].Focused:
			switch data.Options[0].Options[0].Value.(string) {
			case "Kc":
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  "All",
					Value: "All",
				})
			case "Speed":
				for thread, _ := range util.HofSpeedCategories {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  thread,
						Value: thread,
					})
				}
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  "All",
					Value: "All",
				})
			}
		}

		err := util.InteractionRespondChoices(session, i, choices)
		if err != nil {
			logger.Error("Failed to handle admin autocomplete options: " + err.Error())
		}
	}
}

func (s *Service) addToHOFLeaderboard(hofLeaderboard map[string]int, player string, points int) {
	if _, ok := hofLeaderboard[player]; ok {
		hofLeaderboard[player] = hofLeaderboard[player] + points
	} else {
		hofLeaderboard[player] = points
	}
}

func (s *Service) updateHOFLeaderboard(ctx context.Context, session *discordgo.Session, hofLeaderboard map[string]int) {
	logger := flume.FromContext(ctx)
	keys := make([]string, 0, len(hofLeaderboard))
	for key := range hofLeaderboard {
		keys = append(keys, key)
	}

	// Sort the map based on the values
	sort.SliceStable(keys, func(i, j int) bool {
		return hofLeaderboard[keys[i]] > hofLeaderboard[keys[j]]
	})

	// First, delete all the messages within the channel
	err := util.DeleteBulkDiscordMessages(session, s.config.DiscHOFLeaderboardChan, "1196540272486649988")
	if err != nil {
		logger.Error("Failed to bulk delete discord messages: " + err.Error())
	}

	// Iterate over the players to get the different places for users to create the placements
	// Create the leaderboard message that will be sent
	placements := ""
	for placement, player := range keys {
		placements = placements + strconv.Itoa(placement+1) + ". " + player + " [" + strconv.Itoa(hofLeaderboard[player]) + "]\n"
	}

	// Send the Discord Embed message for the leaderboard
	err = util.SendDiscordEmbedMsg(session, s.config.DiscHOFLeaderboardChan, "Ponies Hall Of Fame Leaderboard", placements, "https://i.imgur.com/wbxOjrR.jpeg")
	if err != nil {
		logger.Error("Failed to send message for hof leaderboard: " + err.Error())
		return
	}

	// Send the Discord Embed message for instructions on how the rankings work
	var msg string
	msg = msg + "In order to get onto this leaderboard, you must have a podium finish of one of the HOF Bosses. Also, " +
		"the HOF uses TempleOSRS to get kcs - please turn on the XP Tracker plugin on Runelite and check the TempleOSRS option.\n\n"
	msg = msg + "3 points for :first_place:\n2 points for :second_place:\n1 points for :third_place:"
	err = util.SendDiscordEmbedMsg(session, s.config.DiscHOFLeaderboardChan, "How To Get Onto The Ponies HOF Leaderboard", msg, "https://i.imgur.com/wbxOjrR.jpeg")
	if err != nil {
		logger.Error("Failed to send message for hof leaderboard instructions: " + err.Error())
		return
	}
}

// updatePpLeaderboard will update the cp-leaderboard channel in discord with a new ranking of everyone in the clan
func (s *Service) updatePpLeaderboard(ctx context.Context, session *discordgo.Session) {
	logger := flume.FromContext(ctx)
	// Update the #cp-leaderboard
	keys := make([]string, 0, len(s.cp))
	for key := range s.cp {
		keys = append(keys, key)
	}

	// Sort the map based on the values
	sort.SliceStable(keys, func(i, j int) bool {
		return s.cp[keys[i]] > s.cp[keys[j]]
	})

	// Create the leaderboard message that will be sent
	leaderboard := ""
	for placement, k := range keys {
		if s.cp[k] > 0 {
			leaderboard = leaderboard + strconv.Itoa(placement+1) + ". " + k + ": " + strconv.Itoa(s.cp[k]) + "\n"
		}
	}

	// Retrieve the one channel message and delete it in the leaderboard channel
	messages, err := session.ChannelMessages(s.config.DiscPPLeaderboardChan, 1, "", "1196540063564177561", "")
	if err != nil {
		logger.Error("ERROR RETRIEVING MESSAGES FROM DISCORD LEADERBOARD CHANNEL")
		return
	}
	err = session.ChannelMessageDelete(s.config.DiscPPLeaderboardChan, messages[0].ID)
	if err != nil {
		logger.Error("ERROR DELETING MESSAGES FROM DISCORD LEADERBOARD CHANNEL")
		return
	}

	// Send the Discord Embed message
	err = util.SendDiscordEmbedMsg(session, s.config.DiscPPLeaderboardChan, "Ponies Points Leaderboard", leaderboard, "https://i.imgur.com/wbxOjrR.jpeg")
	if err != nil {
		logger.Error("ERROR SENDING MESSAGES TO DISCORD LEADERBOARD CHANNEL: " + err.Error())
		return
	}
}

/*
updateColLog will use all the users within the in memory submission map to create the podium
from collectionlog.net and their rankings. It will create an embed with the top 10 placements in
discord.
*/
func (s *Service) updateColLog(ctx context.Context, session *discordgo.Session) error {
	logger := flume.FromContext(ctx)
	logger.Info("Running collection log hiscores update...")

	podium, ranking := s.collectionLog.RetrieveCollectionLogAndOrder(ctx, s.cp)

	// Create the leaderboard message that will be sent
	placements := ""
	for placement, k := range ranking {
		placements = placements + strconv.Itoa(placement+1) + ". " + k + " [" + strconv.Itoa(podium[k]) + "]\n"
	}

	// First, delete all the messages within the channel
	err := util.DeleteBulkDiscordMessages(session, s.config.DiscColChan, "1196541219581460530")
	if err != nil {
		logger.Error("Failed to bulk delete discord messages: " + err.Error())
	}

	// Send the Discord Embed message for collection log
	err = util.SendDiscordEmbedMsg(session, s.config.DiscColChan, "Collection Log Leaderboard", placements, "https://i.imgur.com/otTd8Dg.png")
	if err != nil {
		logger.Error("Failed to send discord emded message" + err.Error())
		return err
	}

	// Send the Discord Embed message for instructions on how to get on the collection log hall of fame
	var msg string
	msg = msg + "1. Download the Collection Log plugin\n"
	msg = msg + "2. Click the box to \"Allow collectionlog.net connections\"\n"
	msg = msg + "3. Click through your collection log in game (there will be a * next to the one you still need to click)\n"
	msg = msg + "4. Go to the collection log icon on the sidebar\n"
	msg = msg + "5. Click Account at the top and then upload collection log\n"
	err = util.SendDiscordEmbedMsg(session, s.config.DiscColChan, "How To Get Onto The Collection Log HOF", msg, "https://i.imgur.com/otTd8Dg.png")
	if err != nil {
		logger.Error("Failed to send discord emded message" + err.Error())
		return err
	}

	logger.Info("Collection log hiscores update successful.")
	return nil
}

func (s *Service) updateLeagues(ctx context.Context, session *discordgo.Session) {
	logger := flume.FromContext(ctx)
	logger.Info("Running leagues hiscores update.")

	// First, delete all the messages within the channel
	err := util.DeleteBulkDiscordMessages(session, s.config.DiscLeaguesChan, "1196540715983974531")
	if err != nil {
		logger.Error("Failed to delete bulk discord messages: " + err.Error())
	}

	leaguesPodium, ranking := s.runescape.GetLeaguesPodiumFromRS(ctx, s.cp)
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

		placements = placements + strconv.Itoa(placement+1) + ". " + player + " [" + strconv.Itoa(leaguesPodium[player]) + "] "

		var bronze, iron, steel, mithril, adamant, runeTier, dragon = 2500, 5000, 10000, 18000, 28000, 42000, 56000
		points := leaguesPodium[player]
		switch true {
		case points >= bronze && points < iron:
			placements = placements + "<:bronze_trophy:1178124557933101056>\n"
		case points >= iron && points < steel:
			placements = placements + "<:iron_trophy:1178124556846780458>\n"
		case points >= steel && points < mithril:
			placements = placements + "<:steel_trophy:1178124555718508617>\n"
		case points >= mithril && points < adamant:
			placements = placements + "<:mithril_trophy:1178124554820931755>\n"
		case points >= adamant && points < runeTier:
			placements = placements + "<:adamant_trophy:1178124552971231365>\n"
		case points >= runeTier && points < dragon:
			placements = placements + "<:rune_trophy:1178124551188664350>\n"
		case points >= dragon:
			placements = placements + "<:dragon_trophy:1178124549141839893>\n"
		default:
			placements = placements + "\n"
		}
	}

	// Send the Discord Embed message for the boss podium finish
	err = util.SendDiscordEmbedMsg(session, s.config.DiscLeaguesChan, "Ponies Trailblazer Reloaded League Standings", placements, "https://i.imgur.com/wbxOjrR.jpeg")
	if err != nil {
		logger.Error("Failed to send message for leagues podium.")
		return
	}

	logger.Info("Leagues hiscores update successful.")
}
