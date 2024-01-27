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

func (s *Service) handleAdmin(session *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	switch options[0].Name {
	case "player":
		ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid).With("user", i.Member.User.Username))
		defer func() { s.tid++ }()
		returnMessage := s.handlePlayerAdministration(ctx, session, i)
		err := util.InteractionRespond(session, i, returnMessage)
		if err != nil {
			s.log.Error("Failed to send admin interaction response: " + err.Error())
		}
		s.updatePpLeaderboard(ctx, session)
	case "update-instructions":
		ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid).With("user", i.Member.User.Username))
		logger := flume.FromContext(ctx)
		defer func() { s.tid++ }()
		err := util.InteractionRespond(session, i, "Updating Ponies Point Instructions")
		if err != nil {
			logger.Error("Failed to send interaction response: " + err.Error())
		}
		_ = s.updateSubmissionInstructions(ctx, session)
		return
	case "update-points":
		ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid).With("user", i.Member.User.Username))
		logger := flume.FromContext(ctx)
		defer func() { s.tid++ }()
		returnMessage := s.updatePPPoints(ctx, session, i)
		err := util.InteractionRespond(session, i, returnMessage)
		if err != nil {
			logger.Error("Failed to send admin interaction response: " + err.Error())
		}
	case "reset-speed":
		s.resetSpeedAdmin(session, i)
	case "update-leaderboard":
		s.updateLeaderboard(session, i)
	case "update-sheets":
		ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid).With("user", i.Member.User.Username))
		logger := flume.FromContext(ctx)
		defer func() { s.tid++ }()
		err := util.InteractionRespond(session, i, "Updating Google Sheets")
		if err != nil {
			logger.Error("Failed to send interaction response: " + err.Error())
		}
		s.updateAllGoogleSheets(ctx)
		return
	case "update-guides-map":
		ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid).With("user", i.Member.User.Username))
		logger := flume.FromContext(ctx)
		defer func() { s.tid++ }()
		logger.Info("Updating Guides stored in Map...")
		err := util.InteractionRespond(session, i, "Updating Guides stored in Map")
		if err != nil {
			logger.Error("Failed to send interaction response: " + err.Error())
		}
		s.pastebin.UpdateGuideList(ctx, s.discGuides)
		for _, discGuidesInfos := range s.discGuides {
			for _, discGuidesInfo := range discGuidesInfos {
				logger.Debug("Updating: " + discGuidesInfo.GuidePageName)
			}
		}
		logger.Info("Finished updating Guides stored in Map")
		return
	}

	return
}

func (s *Service) resetSpeedAdmin(session *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid).With("user", i.Member.User.Username))
		defer func() { s.tid++ }()
		s.resetSpeedAdminCommand(ctx, session, i)
	case discordgo.InteractionApplicationCommandAutocomplete:
		s.resetSpeedAdminAutocomplete(session, i)
	}
}

func (s *Service) resetSpeedAdminCommand(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options[0].Options
	logger := flume.FromContext(ctx)

	category := ""
	boss := ""

	for _, option := range options {
		switch option.Name {
		case "category":
			category = option.Value.(string)
		case "boss":
			boss = option.Value.(string)
		}
	}

	logger.Info("Resetting speed for: " + boss)

	err := util.InteractionRespond(session, i, "Resetting speed for: "+boss)
	if err != nil {
		logger.Error("Failed to send admin interaction response: " + err.Error())
	}

	// Ensure the boss name is okay
	if _, ok := util.SpeedBossNameToCategory[boss]; !ok {
		logger.Error("Incorrect boss name: ", boss)
	}

	// Convert the time string into time
	t := util.CalculateTime("22:22:22.60")
	s.speed[boss] = util.SpeedInfo{Time: t, PlayersInvolved: "null", URL: "https://i.imgur.com/34dg0da.png"}
	s.updateSpeedHOF(ctx, session, category)
}

func (s *Service) resetSpeedAdminAutocomplete(session *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData().Options[0]
	var choices []*discordgo.ApplicationCommandOptionChoice
	switch {
	// In this case there are multiple autocomplete options. The Focused field shows which option user is focused on.
	case data.Options[0].Focused:
		choices = util.SpeedAutocompleteCategories
	case data.Options[1].Focused:
		choices = util.AppendToHofSpeedArr(data.Options[0].Value.(string))
	}

	err := util.InteractionRespondChoices(session, i, choices)
	if err != nil {
		s.log.Error("Failed to handle speed autocomplete: " + err.Error())
	}
}

func (s *Service) updateSubmissionInstructions(ctx context.Context, session *discordgo.Session) string {
	returnMessage := "Successfully updated submission Instructions!"
	logger := flume.FromContext(ctx)

	// First, delete all the messages within the channel
	err := util.DeleteBulkDiscordMessages(session, s.config.DiscSpeedSubInfoChan, "")

	speedSubmissionInstruction := []string{
		"# Instructions for Speed Submissions",
		"In order to manually submit for speed times, use the /speed-submissions command. There will be **4 mandatory fields** which are automatically placed in your chat box and there are 2 optional fields which needs to be selected when pressing the +2 more at the end of the chat box",
		"\n✎﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏\n",
		"## ******__Mandatory Fields For Speed Submissions__******",
		"https://i.imgur.com/MK6BzCK.png",
		"### <:ponies:1197979241959145636> category\nThis has a list of all the speed categories you see in the hof-speeds forum. Select one of these in order to proceed in the submission",
		"https://i.imgur.com/uVDhk9U.png",
		"### <:ponies:1197979241959145636> boss\nThis has a list of all the bosses in the previously selected category. Select one of these options to make a speed submission for",
		"https://i.imgur.com/gXD9bHy.png",
		"### <:ponies:1197979241959145636> speed-time\nThe time must be in the format of hh:mm:ss.ms where hh = hours, mm = minutes, ss = seconds, and ms = milliseconds. The following example is 20 hours, 20 minutes, 20 seconds and 1 tick",
		"https://i.imgur.com/uzwDOL3.png",
		"### <:ponies:1197979241959145636> player-names\nThis is comma separated list of all the participating ponies members. Any non-members submitted will cause an error in the submission.",
		"https://i.imgur.com/ML14RzQ.png",
		"## ******__Additional Fields__******",
		"https://i.imgur.com/dD4FKb9.png",
		"**NOTE: Only 1 or either the screenshot field or i-imgur-link field is acceptable. Using both will cause and error as well as using none!**",
		"### <:ponies:1197979241959145636> screenshot\nThis allows you to select an image from your computer to upload to the submission",
		"https://i.imgur.com/SGvWSt8.png",
		"### <:ponies:1197979241959145636> i-imgur-link\nThis allows you to put in an i.imgur.com url instead of an image upload",
		"https://i.imgur.com/TaoiTLG.png",
		"\n✎﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏\n",
		"# ******Speed Submission using screenshot******",
		"https://i.imgur.com/IlgOsfy.gif",
		"\n✎﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏\n",
		"# What happens when your screenshot gets approved/denied",
		"Once a moderator approves/denies your submission, a message will popup in a private channel between you and the moderators with feedback on your submission",
		"https://i.imgur.com/V5AiTyZ.png",
		"If you have an issue, you can ask the moderators there about what was incorrect about your submission.",
	}

	logger.Debug("Running speed submission instruction update")
	for _, msg := range speedSubmissionInstruction {
		_, err := session.ChannelMessageSend(s.config.DiscSpeedSubInfoChan, msg)
		if err != nil {
			logger.Error("Failed to send message to cp information channel", err)
			return "Failed to send message to cp information channel"
		}
	}

	err = util.DeleteBulkDiscordMessages(session, s.config.DiscPPInfoChan, "")
	if err != nil {
		logger.Error("Failed to delete bulk discord messages: " + err.Error())
	}

	ppSubmissionInstruction := []string{
		"# Instructions for Ponies Points Submissions",
		"In order to manually submit for ponies points, use the /pp-submissions command. There will be **1 mandatory field** which is automatically placed in your chat box and there are 2 optional fields which needs to be selected when pressing the +2 more at the end of the chat box",
		"\n✎﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏\n",
		"## ******Mandatory Fields For Speed Submissions******",
		"https://i.imgur.com/hi66ThP.png",
		"### <:ponies:1197979241959145636> player-names\nThis is comma separated list of all the participating ponies members. Any non-members submitted will cause an error in the submission.",
		"https://i.imgur.com/lzYUZUz.png",
		"## ******Additional Fields******",
		"https://i.imgur.com/dD4FKb9.png",
		"**NOTE: Only 1 or either the screenshot field or i-imgur-link field is acceptable. Using both will cause and error as well as using none!**",
		"### <:ponies:1197979241959145636> screenshot\nThis allows you to select an image from your computer to upload to the submission",
		"https://i.imgur.com/SGvWSt8.png",
		"### <:ponies:1197979241959145636> i-imgur-link\nThis allows you to put in an i.imgur.com url instead of an image upload",
		"https://i.imgur.com/TaoiTLG.png",
		"\n✎﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏\n",
		"# ******PP Submission using screenshot******",
		"https://i.imgur.com/FAFCyim.gif",
		"\n✎﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏\n",
		"# What happens when your screenshot gets approved/denied",
		"Once a moderator approves/denies your submission, a message will popup in a private channel between you and the moderators with feedback on your submission",
		"https://i.imgur.com/QUvB4oo.png",
		"If you have an issue, you can ask the moderators there about what was incorrect about your submission.",
	}

	logger.Debug("Running speed submission instruction update")
	for _, msg := range ppSubmissionInstruction {
		_, err := session.ChannelMessageSend(s.config.DiscPPInfoChan, msg)
		if err != nil {
			logger.Error("Failed to send message to cp information channel", err)
			return "Failed to send message to cp information channel"
		}
	}

	keys := make([]string, 0, len(util.LootLogClanPoint))

	for key := range util.LootLogClanPoint {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return util.LootLogClanPoint[keys[i]] < util.LootLogClanPoint[keys[j]]
	})

	var cpInstructions []string
	currentCategory := ""
	currentString := "# The following items will count for Ponies Points"

	for _, item := range keys {
		category := util.LootLogClanPoint[item]
		if strings.Compare(currentCategory, category) != 0 {
			cpInstructions = append(cpInstructions, currentString)
			currentCategory = category
			currentString = "## " + category + "\n"
		}
		currentString = currentString + "- " + item + "\n"
	}

	for _, msg := range cpInstructions {
		_, err := session.ChannelMessageSend(s.config.DiscPPInfoChan, msg)
		if err != nil {
			logger.Error("Failed to send message to cp information channel", err)
			return "Failed to send message to cp information channel"
		}
	}

	logger.Debug("Successfully updated speed submission instruction!")
	return returnMessage
}

func (s *Service) handlePlayerAdministration(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) string {
	logger := flume.FromContext(ctx)
	options := i.ApplicationCommandData().Options[0].Options

	option := ""
	name := ""
	newName := ""
	discordid := 0
	discordname := ""
	main := false

	for _, iterOption := range options {
		switch iterOption.Name {
		case "option":
			option = iterOption.Value.(string)
		case "name":
			name = iterOption.Value.(string)
		case "new-name":
			newName = iterOption.Value.(string)
		case "discord-id":
			discordidStr := iterOption.Value.(string)
			var err error
			discordid, err = strconv.Atoi(discordidStr)
			if err != nil {
				msg := "Discord ID needs to be a number!"
				return msg
			}
		case "discord-name":
			discordname = iterOption.Value.(string)
		case "main":
			main = iterOption.Value.(bool)
		}
	}

	switch option {
	case "Add":
		// If add, ensure that discord-id and discord-name are there
		if discordid == 0 || len(discordname) == 0 {
			logger.Error("Discord ID and Discord Name are required for an addition to the clan")
			msg := "Discord ID and Discord Name are required for an addition to the clan"
			return msg
		}

		// Ensure that this person does not exist in the members map
		if _, ok := s.members[name]; ok {
			// Send the failed addition message in the previously created private channel
			logger.Error("Member: " + name + " already exists.")
			msg := "Member: " + name + " already exists."
			return msg
		} else {
			existingMember := false
			for _, member := range s.members {
				if member.DiscordId == discordid {
					existingMember = true
					break
				}
			}

			// Ensure that if it's the first account being added for a particular discord user, it has to be a main
			if !existingMember && !main {
				// Send the failed addition message in the previously created private channel
				logger.Error("First member for the discord user: " + discordname + ". This is required to be a main.")
				msg := "First member for the discord user: " + discordname + ". This is required to be a main."
				return msg
			}

			// Ensure there can only be 1 main at a time
			if existingMember && main {
				// Send the failed addition message in the previously created private channel
				logger.Error("There can only be 1 main at a time per user: " + discordname)
				msg := "There can only be 1 main at a time per user: " + discordname
				return msg
			}

			// Only add a clan point entry for a main
			if main {
				s.cp[name] = 0
			}
			s.members[name] = util.MemberInfo{
				DiscordId:   discordid,
				DiscordName: discordname,
				Main:        main,
			}
			s.temple.AddMemberToTemple(ctx, name, s.config.TempleGroupId, s.config.TempleGroupKey)

			logger.Debug("You have successfully added a new member: " + name)
			msg := "You have successfully added a new member: " + name
			return msg
		}
	case "Remove":
		// Remove the user from the temple page
		s.temple.RemoveMemberFromTemple(ctx, name, s.config.TempleGroupId, s.config.TempleGroupKey)

		if _, ok := s.members[name]; ok {
			// If the account we're deleting is a main, check to see if there are any other accounts for this discord user
			// If there is, just assign the first instance as main - if not, just delete
			if s.members[name].Main {
				logger.Debug("Deleting player is a main, searching for another account to become main...")
				discordId := s.members[name].DiscordId
				delete(s.members, name)
				for user, member := range s.members {
					if discordId == member.DiscordId {
						logger.Debug("Found new main: " + user + ". Using as new main for user: " + s.members[user].DiscordName)
						s.members[user] = util.MemberInfo{
							DiscordId:   s.members[user].DiscordId,
							DiscordName: s.members[user].DiscordName,
							Feedback:    s.members[user].Feedback,
							Main:        true,
						}
						s.cp[user] = s.cp[name]

						// Update HOF Speed times from deleted user to new main user
						updatedSpeedInfo := make(map[string]util.SpeedInfo)
						for boss, speedInfo := range s.speed {
							updatedSpeedInfo[boss] = util.SpeedInfo{
								PlayersInvolved: strings.Replace(speedInfo.PlayersInvolved, name, newName, -1),
								Time:            speedInfo.Time,
								URL:             speedInfo.URL,
							}
						}
						s.speed = updatedSpeedInfo
					}
				}
				delete(s.cp, name)
			} else {
				delete(s.members, name)
				delete(s.cp, name)
			}

			logger.Debug("You have successfully removed a member: " + name)
			msg := "You have successfully removed a member: " + name
			return msg

		} else {
			// Send the failed removal message in the previously created private channel
			logger.Error("Member: " + name + " does not exist.")
			msg := "Member: " + name + " does not exist."
			return msg
		}
	case "Name Change":
		if _, ok := s.members[name]; ok {
			// Remove the user from the temple page and add new name
			s.temple.RemoveMemberFromTemple(ctx, name, s.config.TempleGroupId, s.config.TempleGroupKey)
			s.temple.AddMemberToTemple(ctx, newName, s.config.TempleGroupId, s.config.TempleGroupKey)

			// Update HOF Speed times from old name to new name
			updatedSpeedInfo := make(map[string]util.SpeedInfo)
			for boss, speedInfo := range s.speed {
				updatedSpeedInfo[boss] = util.SpeedInfo{
					PlayersInvolved: strings.Replace(speedInfo.PlayersInvolved, name, newName, -1),
					Time:            speedInfo.Time,
					URL:             speedInfo.URL,
				}
			}
			s.speed = updatedSpeedInfo
			s.members[newName] = s.members[name]
			s.cp[newName] = s.cp[name]
			delete(s.cp, name)
			delete(s.members, name)

			logger.Debug("You have successfully changed names from: " + name + " to: " + newName)
			msg := "You have successfully changed names from: " + name + " to: " + newName

			return msg

		} else {
			// Send the failed removal message in the previously created private channel
			logger.Error("Member: " + name + " does not exist.")
			msg := "Member: " + name + " does not exist."
			return msg
		}
	case "Update Main":
		// Ensure it is a main first
		if _, ok := s.cp[name]; !ok {
			logger.Error("Name: " + name + " is not a main")
			msg := "Name: " + name + " is not a main. Ensure a main is used in the name section and the transferring name (not a main) is in the new-name section"
			return msg
		} else if _, ok := s.cp[newName]; ok {
			// Ensure the new-name is not a main
			logger.Error("New Name: " + newName + " is a main")
			msg := "New Name: " + newName + " is a main. Ensure a main is used in the name section and the transferring name (not a main) is in the new-name section"
			return msg
		}

		// Ensure the name and newName belong to the same discord user
		if s.members[name].DiscordId != s.members[newName].DiscordId {
			logger.Error("Main: " + name + " and New Main: " + newName + " do not belong to the same discord id.")
			msg := "Main: " + name + " and New Main: " + newName + " do not belong to the same discord id."
			return msg
		}

		// Use the new name as the clan points owner
		s.cp[newName] = s.cp[name]
		delete(s.cp, name)

		// Set the main to true for the newName and set main to false for name
		s.members[name] = util.MemberInfo{
			DiscordId:   s.members[name].DiscordId,
			DiscordName: s.members[name].DiscordName,
			Feedback:    s.members[name].Feedback,
			Main:        false,
		}
		s.members[newName] = util.MemberInfo{
			DiscordId:   s.members[newName].DiscordId,
			DiscordName: s.members[newName].DiscordName,
			Feedback:    s.members[newName].Feedback,
			Main:        true,
		}

		// Update HOF Speed times from name to new name
		updatedSpeedInfo := make(map[string]util.SpeedInfo)
		for boss, speedInfo := range s.speed {
			updatedSpeedInfo[boss] = util.SpeedInfo{
				PlayersInvolved: strings.Replace(speedInfo.PlayersInvolved, name, newName, -1),
				Time:            speedInfo.Time,
				URL:             speedInfo.URL,
			}
		}
		s.speed = updatedSpeedInfo
		logger.Debug("You have successfully updated main from: " + name + " to: " + newName)
		msg := "You have successfully updated main from: " + name + " to: " + newName
		return msg
	}

	return "Invalid player management option chosen."
}

func (s *Service) updatePPPoints(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) string {
	options := i.ApplicationCommandData().Options[0].Options
	logger := flume.FromContext(ctx)

	player := ""
	pp := 0
	addOrRemove := ""

	for _, option := range options {
		switch option.Name {
		case "player":
			player = option.Value.(string)
		case "amount-of-pp":
			pp = int(option.Value.(float64))
		case "option":
			addOrRemove = option.Value.(string)
		}
	}

	switch addOrRemove {
	case "Add":
		logger.Info("Adding " + strconv.Itoa(pp) + " ponies point(s) to " + player)
		s.cp[player] += pp
	case "Remove":
		logger.Info("Removing " + strconv.Itoa(pp) + " ponies point(s) to " + player)
		if s.cp[player]-pp < 0 {
			s.cp[player] = 0
		} else {
			s.cp[player] -= pp
		}
	}

	s.updatePpLeaderboard(ctx, session)

	return "Successfully managed pp for " + player
}
