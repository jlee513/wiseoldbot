package service

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/gemalto/flume"
	"osrs-disc-bot/util"
	"sort"
	"strconv"
	"strings"
	"sync"
)

/*
updateKcHOF will iterate over all the HallOfFameRequestInfos, grab the podium from temple for each
of the bosses, sort them, and make the discord call to create the emded with the boss name, image,
and podium finish with [kc]
*/
func (s *Service) updateKcHOF(ctx context.Context, session *discordgo.Session) {
	logger := flume.FromContext(ctx)
	hofLeaderboard := make(map[string]int)

	logger.Info("Running HOF update for KCs from Temple")

	/* Format:
	{
		"data": {
			"memberlist": {
				"Mager": {
					"player": "Mager",
					...
					"bosses": {
						"Clue_all": 1638,
						"Clue_beginner": 264,
						"Clue_easy": 501,
						...
				}
			}
		}
	}

	We need to parse through all the members and combine main/alt kcs. First, we need to retrieve the kcs from temple
	*/
	kcs := s.temple.GetKCsFromTemple(ctx)
	updatedMemberList := make(map[string]util.HallOfFameBossInfo)

	// Once we get the kcs, we need to iterate over all the keys and change the names to be their actual names stored
	for name, kc := range kcs.Data.Memberlist {
		updatedMemberList[s.templeUsernames[strings.ToLower(name)]] = kc
	}

	kcs.Data.Memberlist = updatedMemberList

	// HOF KC
	slayerBosses := util.HofRequestInfo{Name: "Slayer", DiscChan: s.config.DiscSlayerBossesChan, Bosses: util.HofSlayerBosses}
	gwd := util.HofRequestInfo{Name: "Godwars Dungeon", DiscChan: s.config.DiscGwdChan, Bosses: util.HofGWDBosses}
	wildy := util.HofRequestInfo{Name: "Wilderness", DiscChan: s.config.DiscWildyChan, Bosses: util.HofWildyBosses}
	other := util.HofRequestInfo{Name: "Other Bosses", DiscChan: s.config.DiscOtherChan, Bosses: util.HofOtherBosses}
	misc := util.HofRequestInfo{Name: "Miscellaneous", DiscChan: s.config.DiscMiscChan, Bosses: util.HofMiscBosses}
	dt2 := util.HofRequestInfo{Name: "Desert Treasure 2", DiscChan: s.config.DiscDT2Chan, Bosses: util.HofDT2Bosses}
	raids := util.HofRequestInfo{Name: "Raids", DiscChan: s.config.DiscRaidsChan, Bosses: util.HofRaidsBosses}
	pvp := util.HofRequestInfo{Name: "PVP", DiscChan: s.config.DiscPVPChan, Bosses: util.HofPVPBosses}
	clues := util.HofRequestInfo{Name: "Clues", DiscChan: s.config.DiscCluesChan, Bosses: util.HofCluesBosses}

	kcList := []util.HofRequestInfo{
		slayerBosses, gwd, wildy, other, misc, dt2, raids, pvp, clues,
	}

	var wg sync.WaitGroup

	// We will iterate through all the kcList and parse out all the important information
	for _, kcItem := range kcList {
		wg.Add(1)

		// Kick off a goroutine for each of the updates
		go func(kcItem util.HofRequestInfo) {
			logger.Debug("Running HOF Section: " + kcItem.Name)
			defer wg.Done()
			// First, delete all the messages within the channel
			err := util.DeleteBulkDiscordMessages(session, kcItem.DiscChan)
			if err != nil {
				logger.Error("Failed to delete discord messages: " + err.Error())
			}

			// For each of the bosses, we need to iterate over the list of mainAndAlts to add them up
			// kcs.Data.Memberlist["Mager"].Bosses["Callisto"] <- this is how you access the kcs
			for _, boss := range kcItem.Bosses {
				addedUpKcs := make(map[string]int)
				for main, alts := range s.mainAndAlts {
					if kcs.Data.Memberlist[main].Bosses[boss.BossName] == nil {
						continue
					}
					totalKc := int(kcs.Data.Memberlist[main].Bosses[boss.BossName].(float64))
					for _, alt := range alts {
						if kcs.Data.Memberlist[alt].Bosses[boss.BossName] == nil {
							continue
						}
						totalKc += int(kcs.Data.Memberlist[alt].Bosses[boss.BossName].(float64))
					}
					addedUpKcs[main] = totalKc
				}
				// Sort the addedUpKcs based on the value (which is the addedUpKc)
				updatedRankings := make([]string, 0, len(addedUpKcs))

				for key := range addedUpKcs {
					updatedRankings = append(updatedRankings, key)
				}

				sort.SliceStable(updatedRankings, func(i, j int) bool {
					return addedUpKcs[updatedRankings[i]] > addedUpKcs[updatedRankings[j]]
				})

				//Iterate over the players to get the different places for users to create the placements
				placements := ""
				for placement, player := range updatedRankings {
					switch placement {
					case 0:
						placements = placements + ":first_place: "
						s.addToHOFLeaderboard(hofLeaderboard, player, 3)
						placements = placements + s.templeUsernames[strings.ToLower(player)] + " [" + strconv.Itoa(addedUpKcs[player]) + "]\n"
					case 1:
						placements = placements + ":second_place: "
						s.addToHOFLeaderboard(hofLeaderboard, player, 2)
						placements = placements + s.templeUsernames[strings.ToLower(player)] + " [" + strconv.Itoa(addedUpKcs[player]) + "]\n"
					case 2:
						placements = placements + ":third_place: "
						s.addToHOFLeaderboard(hofLeaderboard, player, 1)
						placements = placements + s.templeUsernames[strings.ToLower(player)] + " [" + strconv.Itoa(addedUpKcs[player]) + "]\n"
					}
				}

				// Beautify some names
				bossName := boss.BossName
				switch bossName {
				case "Clue_beginner":
					bossName = "Beginner Clue"
				case "Clue_easy":
					bossName = "Easy Clue"
				case "Clue_medium":
					bossName = "Medium Clue"
				case "Clue_hard":
					bossName = "Hard Clue"
				case "Clue_elite":
					bossName = "Elite Clue"
				case "Clue_master":
					bossName = "Master Clue"
				case "Clue_all":
					bossName = "All Clues"
				case "Bounty Hunter Hunter":
					bossName = "Bounty Hunter - Hunter"
				case "Bounty Hunter Rogue":
					bossName = "Bounty Hunter - Rogue"
				case "Theatre of Blood Challenge Mode":
					bossName = "Theatre of Blood Hard Mode"
				}

				// Send the Discord Embed message for the boss podium finish
				err = util.SendDiscordEmbedMsg(session, kcItem.DiscChan, bossName, placements, boss.ImageLink)
				if err != nil {
					logger.Error("Failed to send message for boss: " + bossName)
					return
				}
			}
		}(kcItem)
	}

	wg.Wait()
	s.updateHOFLeaderboard(ctx, session, hofLeaderboard)
	logger.Info("Successfully updated KC Hall Of Fame")
}

func (s *Service) updateSpeedHOF(ctx context.Context, session *discordgo.Session, requestedBosses ...string) {
	logger := flume.FromContext(ctx)

	// HOF Speed
	tzhaar := util.SpeedsRequestInfo{Name: "TzHaar", DiscChan: s.config.DiscSpeedTzhaarChan}
	slayer := util.SpeedsRequestInfo{Name: "Slayer", DiscChan: s.config.DiscSpeedSlayerChan}
	nightmare := util.SpeedsRequestInfo{Name: "Nightmare", DiscChan: s.config.DiscSpeedNightmareChan}
	nex := util.SpeedsRequestInfo{Name: "Nex", DiscChan: s.config.DiscSpeedNexChan}
	solo := util.SpeedsRequestInfo{Name: "Solo Bosses", DiscChan: s.config.DiscSpeedSoloChan}
	cox := util.SpeedsRequestInfo{Name: "Chambers Of Xeric", DiscChan: s.config.DiscSpeedCOXChan}
	coxcm := util.SpeedsRequestInfo{Name: "Chambers Of Xeric Challenge Mode", DiscChan: s.config.DiscSpeedCOXCMChan}
	tob := util.SpeedsRequestInfo{Name: "Theatre Of Blood", DiscChan: s.config.DiscSpeedTOBChan}
	tobhm := util.SpeedsRequestInfo{Name: "Theatre Of Blood Hard Mode", DiscChan: s.config.DiscSpeedTOBHMChan}
	toa := util.SpeedsRequestInfo{Name: "Tombs Of Amascut", DiscChan: s.config.DiscSpeedTOAChan}
	toae := util.SpeedsRequestInfo{Name: "Tombs Of Amascut Expert", DiscChan: s.config.DiscSpeedTOAEChan}
	agility := util.SpeedsRequestInfo{Name: "Agility", DiscChan: s.config.DiscSpeedAgilityChan}
	dt2 := util.SpeedsRequestInfo{Name: "Desert Treasure 2", DiscChan: s.config.DiscSpeedDt2Chan}

	var allRequestInfo []util.SpeedsRequestInfo
	for _, boss := range requestedBosses {
		switch boss {
		case "TzHaar":
			allRequestInfo = append(allRequestInfo, tzhaar)
		case "Slayer":
			allRequestInfo = append(allRequestInfo, slayer)
		case "Nightmare":
			allRequestInfo = append(allRequestInfo, nightmare)
		case "Nex":
			allRequestInfo = append(allRequestInfo, nex)
		case "Solo Bosses":
			allRequestInfo = append(allRequestInfo, solo)
		case "Chambers Of Xeric":
			allRequestInfo = append(allRequestInfo, cox)
		case "Chambers Of Xeric Challenge Mode":
			allRequestInfo = append(allRequestInfo, coxcm)
		case "Theatre Of Blood":
			allRequestInfo = append(allRequestInfo, tob)
		case "Theatre Of Blood Hard Mode":
			allRequestInfo = append(allRequestInfo, tobhm)
		case "Tombs Of Amascut":
			allRequestInfo = append(allRequestInfo, toa)
		case "Tombs Of Amascut Expert":
			allRequestInfo = append(allRequestInfo, toae)
		case "Agility":
			allRequestInfo = append(allRequestInfo, agility)
		case "Desert Treasure 2":
			allRequestInfo = append(allRequestInfo, dt2)
		}
	}

	var wg sync.WaitGroup

	for _, requestInfo := range allRequestInfo {
		wg.Add(1)

		go func(requestInfo util.SpeedsRequestInfo) {
			defer wg.Done()
			logger.Debug("Running Speed HOF update for Boss: " + requestInfo.Name)
			// First, delete all the messages within the channel
			err := util.DeleteBulkDiscordMessages(session, requestInfo.DiscChan)
			if err != nil {
				logger.Error("Failed to bulk delete discord messages: " + err.Error())
			}

			// Now add all the bosses
			for _, bossName := range s.speedCategory[requestInfo.Name] {
				// Get the speed info
				speed := s.speed[bossName]
				err = util.SendDiscordEmbedMsg(session, requestInfo.DiscChan, bossName, "**Players:** "+speed.PlayersInvolved+"\n**Time:** "+speed.Time.Format("15:04:05.00"), speed.URL)
				if err != nil {
					logger.Error("Failed to send message for boss: " + bossName)
					return
				}
			}
		}(requestInfo)

	}

	wg.Wait()
	logger.Info("Successfully updated Speed Hall Of Fame")
}
