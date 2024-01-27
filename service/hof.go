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

/*
updateKcHOF will iterate over all the HallOfFameRequestInfos, grab the podium from temple for each
of the bosses, sort them, and make the discord call to create the emded with the boss name, image,
and podium finish with [kc]
*/
func (s *Service) updateKcHOF(ctx context.Context, session *discordgo.Session) {
	logger := flume.FromContext(ctx)
	hofLeaderboard := make(map[string]int)

	// HOF KC
	slayerBosses := util.HofRequestInfo{Name: "Slayer", DiscChan: s.config.DiscSlayerBossesChan, AfterId: "1194801106291785778", Bosses: util.HofSlayerBosses}
	gwd := util.HofRequestInfo{Name: "Godwars Dungeon", DiscChan: s.config.DiscGwdChan, AfterId: "1194801166429724884", Bosses: util.HofGWDBosses}
	wildy := util.HofRequestInfo{Name: "Wilderness", DiscChan: s.config.DiscWildyChan, AfterId: "1194801335376285726", Bosses: util.HofWildyBosses}
	other := util.HofRequestInfo{Name: "Other Bosses", DiscChan: s.config.DiscOtherChan, AfterId: "1194801512870846535", Bosses: util.HofOtherBosses}
	misc := util.HofRequestInfo{Name: "Miscellaneous", DiscChan: s.config.DiscMiscChan, AfterId: "1194804397507620935", Bosses: util.HofMiscBosses}
	dt2 := util.HofRequestInfo{Name: "Desert Treasure 2", DiscChan: s.config.DiscDT2Chan, AfterId: "1194802032855498832", Bosses: util.HofDT2Bosses}
	raids := util.HofRequestInfo{Name: "Raids", DiscChan: s.config.DiscRaidsChan, AfterId: "1194802206487089182", Bosses: util.HofRaidsBosses}
	pvp := util.HofRequestInfo{Name: "PVP", DiscChan: s.config.DiscPVPChan, AfterId: "1194802450209718272", Bosses: util.HofPVPBosses}
	clues := util.HofRequestInfo{Name: "Clues", DiscChan: s.config.DiscCluesChan, AfterId: "1194802590270103582", Bosses: util.HofCluesBosses}

	allRequestInfo := []util.HofRequestInfo{
		slayerBosses, gwd, wildy, other, misc, dt2, raids, pvp, clues,
	}

	for _, requestInfo := range allRequestInfo {
		logger.Info("Running HOF update for Boss: " + requestInfo.Name)
		// First, delete all the messages within the channel
		err := util.DeleteBulkDiscordMessages(session, requestInfo.DiscChan, requestInfo.AfterId)
		if err != nil {
			logger.Error("Failed to delete discord messages: " + err.Error())
		}

		// Now add all the bosses
		for _, bossInfo := range requestInfo.Bosses {
			kcs, rankings := s.temple.GetKCsFromTemple(ctx, bossInfo.BossName)

			// Group all the kcs under main players
			mainKcs := make(map[int]util.Player)
			altKcs := make(map[int]util.Player)

			// Go through the rankings and determine whether they are a main or not
			for ranking := range rankings {
				member := s.members[s.templeUsernames[strings.ToLower(kcs.Data.Players[ranking].Username)]]
				// If they are a main, check to see if there are any alts there already calculated
				if member.Main {
					// If it exists within alt kcs, add it to the current member and delete it from alt kcs
					if _, ok := altKcs[member.DiscordId]; ok {
						updatedMemberKcs := util.Player{
							Username: kcs.Data.Players[ranking].Username,
							Kc:       kcs.Data.Players[ranking].Kc + altKcs[member.DiscordId].Kc,
						}
						mainKcs[member.DiscordId] = updatedMemberKcs
						delete(altKcs, member.DiscordId)
					} else {
						// If it doesn't exist, just set the mainKcs object with the member object
						mainKcs[member.DiscordId] = kcs.Data.Players[ranking]
					}
				} else {
					// If it exists within the main kcs, add to it
					if _, ok := mainKcs[member.DiscordId]; ok {
						updatedMemberKcs := util.Player{
							Username: mainKcs[member.DiscordId].Username,
							Kc:       kcs.Data.Players[ranking].Kc + mainKcs[member.DiscordId].Kc,
						}
						mainKcs[member.DiscordId] = updatedMemberKcs
					} else if _, ok := altKcs[member.DiscordId]; ok {
						// Check to see if an altKcs exists - if it does, add it to that
						updatedMemberKcs := util.Player{
							Username: altKcs[member.DiscordId].Username,
							Kc:       kcs.Data.Players[ranking].Kc + altKcs[member.DiscordId].Kc,
						}
						altKcs[member.DiscordId] = updatedMemberKcs
					} else {
						// Otherwise, just add it to altKcs
						altKcs[member.DiscordId] = kcs.Data.Players[ranking]
					}
				}
			}

			// Re-assign rankings for mainKcs
			updatedRankings := make([]int, 0, len(mainKcs))

			for key := range mainKcs {
				updatedRankings = append(updatedRankings, key)
			}

			sort.SliceStable(updatedRankings, func(i, j int) bool {
				return mainKcs[updatedRankings[i]].Kc > mainKcs[updatedRankings[j]].Kc
			})

			// Iterate over the players to get the different places for users to create the placements
			placements := ""
			for placement, updatedRank := range updatedRankings {
				switch placement {
				case 0:
					placements = placements + ":first_place: "
					s.addToHOFLeaderboard(hofLeaderboard, mainKcs[updatedRank].Username, 3)
					placements = placements + s.templeUsernames[strings.ToLower(mainKcs[updatedRank].Username)] + " [" + strconv.Itoa(mainKcs[updatedRank].Kc) + "]\n"
				case 1:
					placements = placements + ":second_place: "
					s.addToHOFLeaderboard(hofLeaderboard, mainKcs[updatedRank].Username, 2)
					placements = placements + s.templeUsernames[strings.ToLower(mainKcs[updatedRank].Username)] + " [" + strconv.Itoa(mainKcs[updatedRank].Kc) + "]\n"
				case 2:
					placements = placements + ":third_place: "
					s.addToHOFLeaderboard(hofLeaderboard, mainKcs[updatedRank].Username, 1)
					placements = placements + s.templeUsernames[strings.ToLower(mainKcs[updatedRank].Username)] + " [" + strconv.Itoa(mainKcs[updatedRank].Kc) + "]\n"
				}
			}

			// Beautify some names
			bossName := kcs.Data.BossName
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
			err = util.SendDiscordEmbedMsg(session, requestInfo.DiscChan, bossName, placements, bossInfo.ImageLink)
			if err != nil {
				logger.Error("Failed to send message for boss: " + kcs.Data.BossName)
				return
			}
		}
	}

	s.updateHOFLeaderboard(ctx, session, hofLeaderboard)
	logger.Info("Successfully updated KC Hall Of Fame")
}

func (s *Service) updateSpeedHOF(ctx context.Context, session *discordgo.Session, requestedBosses ...string) {
	logger := flume.FromContext(ctx)

	// HOF Speed
	tzhaar := util.SpeedsRequestInfo{Name: "TzHaar", DiscChan: s.config.DiscSpeedTzhaarChan, AfterId: "1194999599652425778", Bosses: util.HofSpeedTzhaar}
	slayer := util.SpeedsRequestInfo{Name: "Slayer", DiscChan: s.config.DiscSpeedSlayerChan, AfterId: "1194999714710573078", Bosses: util.HofSpeedSlayer}
	nightmare := util.SpeedsRequestInfo{Name: "Nightmare", DiscChan: s.config.DiscSpeedNightmareChan, AfterId: "1195000377288958023", Bosses: util.HofSpeedNightmare}
	nex := util.SpeedsRequestInfo{Name: "Nex", DiscChan: s.config.DiscSpeedNexChan, AfterId: "1195000695594684416", Bosses: util.HofSpeedNex}
	solo := util.SpeedsRequestInfo{Name: "Solo Bosses", DiscChan: s.config.DiscSpeedSoloChan, AfterId: "1195000959911350294", Bosses: util.HofSpeedSolo}
	cox := util.SpeedsRequestInfo{Name: "Chambers Of Xeric", DiscChan: s.config.DiscSpeedCOXChan, AfterId: "1195001187276161155", Bosses: util.HofSpeedCox}
	coxcm := util.SpeedsRequestInfo{Name: "Chambers Of Xeric Challenge Mode", DiscChan: s.config.DiscSpeedCOXCMChan, AfterId: "1196437479713165322", Bosses: util.HofSpeedCoxCm}
	tob := util.SpeedsRequestInfo{Name: "Theatre Of Blood", DiscChan: s.config.DiscSpeedTOBChan, AfterId: "1195001367685779509", Bosses: util.HofSpeedTob}
	tobhm := util.SpeedsRequestInfo{Name: "Theatre Of Blood Hard Mode", DiscChan: s.config.DiscSpeedTOBHMChan, AfterId: "1196437597053001758", Bosses: util.HofSpeedTobHm}
	toa := util.SpeedsRequestInfo{Name: "Tombs Of Amascut", DiscChan: s.config.DiscSpeedTOAChan, AfterId: "1195001626604355656", Bosses: util.HofSpeedToa}
	toae := util.SpeedsRequestInfo{Name: "Tombs Of Amascut Expert", DiscChan: s.config.DiscSpeedTOAEChan, AfterId: "1196437695522672650", Bosses: util.HofSpeedToae}
	agility := util.SpeedsRequestInfo{Name: "Agility", DiscChan: s.config.DiscSpeedAgilityChan, AfterId: "1195002755132174368", Bosses: util.HofSpeedAgility}
	dt2 := util.SpeedsRequestInfo{Name: "Desert Treasure 2", Bosses: util.HofSpeedDt2, DiscChan: s.config.DiscSpeedDt2Chan, AfterId: "1197895349935812639"}

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

	for _, requestInfo := range allRequestInfo {
		logger.Info("Running Speed HOF update for Boss: " + requestInfo.Name)
		// First, delete all the messages within the channel
		err := util.DeleteBulkDiscordMessages(session, requestInfo.DiscChan, requestInfo.AfterId)
		if err != nil {
			logger.Error("Failed to bulk delete discord messages: " + err.Error())
		}

		// Now add all the bosses
		for _, bossInfo := range requestInfo.Bosses {
			// Get the speed info
			speed := s.speed[bossInfo.BossName]
			err = util.SendDiscordEmbedMsg(session, requestInfo.DiscChan, bossInfo.BossName, "**Players:** "+speed.PlayersInvolved+"\n**Time:** "+speed.Time.Format("15:04:05.00"), speed.URL)
			if err != nil {
				logger.Error("Failed to send message for boss: " + bossInfo.BossName)
				return
			}
		}
	}
	logger.Info("Successfully updated Speed Hall Of Fame")
}
