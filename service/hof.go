package service

import (
	"context"
	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
	"github.com/gemalto/flume"
	"osrs-disc-bot/util"
	"strconv"
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
		messages, err := session.ChannelMessages(requestInfo.DiscChan, 50, "", requestInfo.AfterId, "")
		if err != nil {
			logger.Error("Failed to get all messages for deletion from channel: " + requestInfo.Name)
			return
		}
		var messageIDs []string
		for _, message := range messages {
			messageIDs = append(messageIDs, message.ID)
		}

		if len(messageIDs) > 0 {
			err = session.ChannelMessagesBulkDelete(requestInfo.DiscChan, messageIDs)
			if err != nil {
				logger.Error("Failed to delete all messages from channel: " + requestInfo.Name + ", will try one by one")
				for _, message := range messageIDs {
					err = session.ChannelMessageDelete(requestInfo.DiscChan, message)
					if err != nil {
						logger.Error("Failed to delete messages one by one from channel: " + requestInfo.Name)
						return
					}
				}
			}
		}

		// Now add all the bosses
		for _, bossInfo := range requestInfo.Bosses {
			podium, rankings := s.temple.GetPodiumFromTemple(ctx, bossInfo.BossName)

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
				logger.Error("Failed to send message for boss: " + podium.Data.BossName)
				return
			}
		}
	}

	s.updateHOFLeaderboard(ctx, session, hofLeaderboard)
	logger.Info("Successfully finished updating KC Hall Of Fame")
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
	tob := util.SpeedsRequestInfo{Name: "Theatre Of Blood", DiscChan: s.config.DiscSpeedTOBChan, AfterId: "1195001367685779509", Bosses: util.HofSpeedTob}
	toa := util.SpeedsRequestInfo{Name: "Tombs Of Amascut", DiscChan: s.config.DiscSpeedTOAChan, AfterId: "1195001626604355656", Bosses: util.HofSpeedToa}
	agility := util.SpeedsRequestInfo{Name: "Agility", DiscChan: s.config.DiscSpeedAgilityChan, AfterId: "1195002755132174368", Bosses: util.HofSpeedAgility}

	var allRequestInfo []util.SpeedsRequestInfo
	for _, boss := range requestedBosses {
		switch boss {
		case "TzHaar":
			allRequestInfo = append(allRequestInfo, tzhaar)
			break
		case "Slayer":
			allRequestInfo = append(allRequestInfo, slayer)
			break
		case "Nightmare":
			allRequestInfo = append(allRequestInfo, nightmare)
			break
		case "Nex":
			allRequestInfo = append(allRequestInfo, nex)
			break
		case "Solo Bosses":
			allRequestInfo = append(allRequestInfo, solo)
			break
		case "Chambers Of Xeric":
			allRequestInfo = append(allRequestInfo, cox)
			break
		case "Theatre Of Blood":
			allRequestInfo = append(allRequestInfo, tob)
			break
		case "Tombs Of Amascut":
			allRequestInfo = append(allRequestInfo, toa)
			break
		case "Agility":
			allRequestInfo = append(allRequestInfo, agility)
			break
		}
	}

	for _, requestInfo := range allRequestInfo {
		logger.Info("Running Speed HOF update for Boss: " + requestInfo.Name)
		// First, delete all the messages within the channel
		messages, err := session.ChannelMessages(requestInfo.DiscChan, 50, "", requestInfo.AfterId, "")
		if err != nil {
			logger.Error("Failed to get all messages for deletion from channel: " + requestInfo.Name)
			return
		}
		var messageIDs []string
		for _, message := range messages {
			messageIDs = append(messageIDs, message.ID)
		}
		if len(messageIDs) > 0 {
			err = session.ChannelMessagesBulkDelete(requestInfo.DiscChan, messageIDs)
			if err != nil {
				logger.Error("Failed to delete all messages from channel: " + requestInfo.Name + ", will try one by one")
				for _, message := range messageIDs {
					err = session.ChannelMessageDelete(requestInfo.DiscChan, message)
					if err != nil {
						logger.Error("Failed to delete messages one by one from channel: " + requestInfo.Name)
						return
					}
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
				logger.Error("Failed to send message for boss: " + bossInfo.BossName)
				return
			}
		}
	}
	logger.Info("Successfully finished updating Speed Hall Of Fame")
}
