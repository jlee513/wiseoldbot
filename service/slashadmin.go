package service

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gemalto/flume"
	"osrs-disc-bot/util"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (s *Service) handleAdmin(session *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	switch options[0].Name {
	case "player":
		ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid).With("user", i.Member.User.Username))
		logger := flume.FromContext(ctx)
		defer func() { s.tid++ }()
		returnMessage := s.handlePlayerAdministration(ctx, session, i)
		err := util.InteractionRespond(session, i, returnMessage)
		if err != nil {
			util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Failed to send admin interaction response: "+err.Error())
		}
		s.updateCpLeaderboard(ctx, session, i.Member.User)
	case "instructions":
		ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid).With("user", i.Member.User.Username))
		logger := flume.FromContext(ctx)
		defer func() { s.tid++ }()
		err := util.InteractionRespond(session, i, "Updating Clan Point Instructions")
		if err != nil {
			util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Failed to send interaction response: "+err.Error())
		}
		util.LogAdminAction(logger, s.config.DiscAuditChan, i.Member.User.Username, i.Member.User.AvatarURL(""), session, "Admin invoked instructions command")
		_ = s.updateSubmissionInstructions(ctx, session, i.Member.User)
		return
	case "points":
		ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid).With("user", i.Member.User.Username))
		logger := flume.FromContext(ctx)
		defer func() { s.tid++ }()
		returnMessage := s.updateCpPoints(ctx, session, i)
		err := util.InteractionRespond(session, i, returnMessage)
		if err != nil {
			util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Failed to send admin interaction response: "+err.Error())
		}
	case "speed":
		s.speedAdmin(session, i)
	case "leaderboard":
		s.updateLeaderboard(session, i)
	case "sheets":
		ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid).With("user", i.Member.User.Username))
		logger := flume.FromContext(ctx)
		defer func() { s.tid++ }()
		err := util.InteractionRespond(session, i, "Updating Google Sheets")
		if err != nil {
			util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Failed to send interaction response: "+err.Error())
		}
		util.LogAdminAction(logger, s.config.DiscAuditChan, i.Member.User.Username, i.Member.User.AvatarURL(""), session, fmt.Sprintf("Admin invoked sheets update command"))
		s.updateAllGoogleSheets(ctx)
		return
	case "guides-map":
		ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid).With("user", i.Member.User.Username))
		logger := flume.FromContext(ctx)
		util.LogAdminAction(logger, s.config.DiscAuditChan, i.Member.User.Username, i.Member.User.AvatarURL(""), session, fmt.Sprintf("Admin invoked guides-map update command"))
		defer func() { s.tid++ }()
		logger.Info("Updating Guides stored in Map...")
		err := util.InteractionRespond(session, i, "Updating Guides stored in Map")
		if err != nil {
			util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Failed to send interaction response: "+err.Error())
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

func (s *Service) speedAdmin(session *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid).With("user", i.Member.User.Username))
		defer func() { s.tid++ }()
		s.speedAdminCommand(ctx, session, i)
	case discordgo.InteractionApplicationCommandAutocomplete:
		s.speedAdminAutocomplete(session, i)
	}
}

func (s *Service) speedAdminCommand(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options[0].Options
	logger := flume.FromContext(ctx)

	action := ""
	category := ""
	existingBoss := ""
	newBoss := ""
	updateSpeedTime := ""
	updatePlayerNames := ""
	imgurUrl := ""

	for _, option := range options {
		switch option.Name {
		case "option":
			action = option.Value.(string)
		case "category":
			category = option.Value.(string)
		case "existing-boss":
			existingBoss = option.Value.(string)
		case "new-boss":
			newBoss = option.Value.(string)
		case "update-speed-time":
			updateSpeedTime = option.Value.(string)
		case "update-player-names":
			updatePlayerNames = option.Value.(string)
		case "imgur-url":
			imgurUrl = option.Value.(string)
		}
	}

	switch action {
	case "Reset":
		util.LogAdminAction(logger, s.config.DiscAuditChan, i.Member.User.Username, i.Member.User.AvatarURL(""), session, fmt.Sprintf("Admin invoked speed reset command with options: option=%s, category=%s, existing-boss=%s, new-boss=%s, update-speed-time=%s, update-player-names=%s, imgur-url=%s", action, category, existingBoss, newBoss, updateSpeedTime, updatePlayerNames, imgurUrl))
		s.resetSpeed(ctx, session, i, existingBoss, category)
	case "Add":
		util.LogAdminAction(logger, s.config.DiscAuditChan, i.Member.User.Username, i.Member.User.AvatarURL(""), session, fmt.Sprintf("Admin invoked speed add command with options: option=%s, category=%s, existing-boss=%s, new-boss=%s, update-speed-time=%s, update-player-names=%s, imgur-url=%s", action, category, existingBoss, newBoss, updateSpeedTime, updatePlayerNames, imgurUrl))
		s.addNewSpeed(ctx, session, i, newBoss, category)
	case "Update":
		util.LogAdminAction(logger, s.config.DiscAuditChan, i.Member.User.Username, i.Member.User.AvatarURL(""), session, fmt.Sprintf("Admin invoked speed update command with options: option=%s, category=%s, existing-boss=%s, new-boss=%s, update-speed-time=%s, update-player-names=%s, imgur-url=%s", action, category, existingBoss, newBoss, updateSpeedTime, updatePlayerNames, imgurUrl))
		s.updateSpeed(ctx, session, i, existingBoss, category, updateSpeedTime, updatePlayerNames, imgurUrl)
	case "Remove":
		util.LogAdminAction(logger, s.config.DiscAuditChan, i.Member.User.Username, i.Member.User.AvatarURL(""), session, fmt.Sprintf("Admin invoked speed remove command with options: option=%s, category=%s, existing-boss=%s, new-boss=%s, update-speed-time=%s, update-player-names=%s, imgur-url=%s", action, category, existingBoss, newBoss, updateSpeedTime, updatePlayerNames, imgurUrl))
		s.removeSpeed(ctx, session, i, existingBoss, category)
	}

}

func (s *Service) updateSpeed(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate, existingBoss, category, updateSpeedTime, updatePlayerNames, imgurUrl string) {
	logger := flume.FromContext(ctx)
	logger.Info("Updating speed for: " + existingBoss)

	err := util.InteractionRespond(session, i, "Updating speed for: "+existingBoss)
	if err != nil {
		util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Failed to send admin interaction response: "+err.Error())
	}

	// Ensure the boss name exists
	if _, ok := s.speed[existingBoss]; !ok {
		util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Incorrect boss name: "+existingBoss)
		return
	}

	speed := time.Time{}
	playersInvolved := updatePlayerNames
	url := imgurUrl

	if len(updateSpeedTime) > 0 {
		logger.Debug("Updating time for: " + existingBoss + " to: " + updateSpeedTime)
		// Ensure the format is hh:mm:ss:mmm
		reg := regexp.MustCompile("^\\d\\d:\\d\\d:\\d\\d\\.\\d\\d$")
		if !reg.Match([]byte(updateSpeedTime)) {
			util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Invalid time format: "+updateSpeedTime)
			return
		}
		speed = util.CalculateTime(updateSpeedTime)
	} else {
		speed = s.speed[existingBoss].Time
	}

	if len(updatePlayerNames) == 0 {
		playersInvolved = s.speed[existingBoss].PlayersInvolved
	} else {
		logger.Debug("Updating player names for: " + existingBoss + " to: " + updatePlayerNames)
	}

	if len(url) == 0 {
		url = s.speed[existingBoss].URL
	} else {
		logger.Debug("Updating URL for: " + existingBoss + " to: " + imgurUrl)
	}

	s.speed[existingBoss] = util.SpeedInfo{
		PlayersInvolved: playersInvolved,
		Time:            speed,
		URL:             url,
		Category:        category,
	}

	s.updateSpeedHOF(ctx, session, i.Member.User, category)
	logger.Info("Successfully updated speed for: " + existingBoss)
}

func (s *Service) resetSpeed(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate, existingBoss, category string) {
	logger := flume.FromContext(ctx)
	logger.Info("Resetting speed for: " + existingBoss)

	err := util.InteractionRespond(session, i, "Resetting speed for: "+existingBoss)
	if err != nil {
		util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Failed to send admin interaction response: "+err.Error())
	}

	// Ensure the boss name is okay
	if _, ok := s.speed[existingBoss]; !ok {
		util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Incorrect boss name: "+existingBoss)
		return
	}

	// Convert the time string into time
	t := util.CalculateTime("22:22:22.60")
	s.speed[existingBoss] = util.SpeedInfo{Time: t, PlayersInvolved: "null", URL: "https://i.imgur.com/34dg0da.png", Category: category}
	s.updateSpeedHOF(ctx, session, i.Member.User, category)
	logger.Info("Successfully reset speed for: " + existingBoss)
}

func (s *Service) addNewSpeed(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate, newBoss, category string) {
	logger := flume.FromContext(ctx)
	logger.Info("Adding new speed for: " + newBoss)

	err := util.InteractionRespond(session, i, "Adding new speed for: "+newBoss)
	if err != nil {
		util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Failed to send admin interaction response: "+err.Error())
	}

	// Ensure the boss name is okay
	if _, ok := s.speed[newBoss]; ok {
		util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Inputted in existing boss: "+newBoss)
		return
	}

	// Need to append the new boss to the category
	s.speedCategory[category] = append(s.speedCategory[category], newBoss)

	// Convert the time string into time
	t := util.CalculateTime("22:22:22.60")
	s.speed[newBoss] = util.SpeedInfo{Time: t, PlayersInvolved: "null", URL: "https://i.imgur.com/34dg0da.png", Category: category}
	s.updateSpeedHOF(ctx, session, i.Member.User, category)
	logger.Info("Successfully added speed for: " + newBoss)
}

func (s *Service) removeSpeed(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate, existingBoss, category string) {
	logger := flume.FromContext(ctx)
	logger.Info("Removing speed for: " + existingBoss)

	err := util.InteractionRespond(session, i, "Removing speed for: "+existingBoss)
	if err != nil {
		util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Failed to send admin interaction response: "+err.Error())
	}

	// Ensure the boss name is okay
	if _, ok := s.speed[existingBoss]; !ok {
		util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Inputted in non-existant boss: "+existingBoss)
		return
	}

	// Need to remove the new boss to the category
	for place, boss := range s.speedCategory[category] {
		if strings.Compare(boss, existingBoss) == 0 {
			s.speedCategory[category] = append(s.speedCategory[category][:place], s.speedCategory[category][place+1:]...)
			break
		}
	}

	delete(s.speed, existingBoss)
	s.updateSpeedHOF(ctx, session, i.Member.User, category)
	logger.Info("Successfully removed speed for: " + existingBoss)
}

func (s *Service) speedAdminAutocomplete(session *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData().Options[0]
	var choices []*discordgo.ApplicationCommandOptionChoice
	switch {
	// In this case there are multiple autocomplete options. The Focused field shows which option user is focused on.
	case data.Options[1].Focused:
		for category := range util.HofSpeedCategories {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  category,
				Value: category,
			})
		}
	case data.Options[2].Focused:
		for _, boss := range s.speedCategory[data.Options[1].Value.(string)] {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  boss,
				Value: boss,
			})
		}
	}

	err := util.InteractionRespondChoices(session, i, choices)
	if err != nil {
		s.log.Error("Failed to handle speed autocomplete: " + err.Error())
	}
}

func (s *Service) updateSubmissionInstructions(ctx context.Context, session *discordgo.Session, invokedUser *discordgo.User) string {
	returnMessage := "Successfully updated submission Instructions!"
	logger := flume.FromContext(ctx)

	// First, delete all the messages within the channel
	err := util.DeleteBulkDiscordMessages(session, s.config.DiscSpeedSubInfoChan)

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
			util.LogError(logger, s.config.DiscAuditChan, session, invokedUser.Username, invokedUser.AvatarURL(""), "Failed to send message to cp information channel: "+err.Error())
			return "Failed to send message to cp information channel"
		}
	}

	err = util.DeleteBulkDiscordMessages(session, s.config.DiscCpInfoChan)
	if err != nil {
		util.LogError(logger, s.config.DiscAuditChan, session, invokedUser.Username, invokedUser.AvatarURL(""), "Failed to delete bulk discord messages: "+err.Error())
	}

	ppSubmissionInstruction := []string{
		"# Instructions for Clan Points Submissions",
		"In order to manually submit for ponies points, use the /cp-submissions command. There will be **1 mandatory field** which is automatically placed in your chat box and there are 2 optional fields which needs to be selected when pressing the +2 more at the end of the chat box",
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
		"# ******Cp Submission using screenshot******",
		"https://i.imgur.com/FAFCyim.gif",
		"\n✎﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏﹏\n",
		"# What happens when your screenshot gets approved/denied",
		"Once a moderator approves/denies your submission, a message will popup in a private channel between you and the moderators with feedback on your submission",
		"https://i.imgur.com/QUvB4oo.png",
		"If you have an issue, you can ask the moderators there about what was incorrect about your submission.",
	}

	logger.Debug("Running speed submission instruction update")
	for _, msg := range ppSubmissionInstruction {
		_, err := session.ChannelMessageSend(s.config.DiscCpInfoChan, msg)
		if err != nil {
			util.LogError(logger, s.config.DiscAuditChan, session, invokedUser.Username, invokedUser.AvatarURL(""), "Failed to send message to cp information channel:"+err.Error())
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
	currentString := "# The following items will count for Clan Points"

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
		_, err := session.ChannelMessageSend(s.config.DiscCpInfoChan, msg)
		if err != nil {
			util.LogError(logger, s.config.DiscAuditChan, session, invokedUser.Username, invokedUser.AvatarURL(""), "Failed to send message to cp information channel: "+err.Error())
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
		util.LogAdminAction(logger, s.config.DiscAuditChan, i.Member.User.Username, i.Member.User.AvatarURL(""), session, fmt.Sprintf("Admin invoked player add command with options: option=%s, name=%s, newName=%s, discordid=%d, discordname=%s, main=%t", option, name, newName, discordid, discordname, main))
		// If add, ensure that discord-id and discord-name are there
		if discordid == 0 || len(discordname) == 0 {
			util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Discord ID and Discord Name are required for an addition to the clan")
			msg := "Discord ID and Discord Name are required for an addition to the clan"
			return msg
		}

		// Ensure that this person does not exist in the members map
		if _, ok := s.members[name]; ok {
			// Send the failed addition message in the previously created private channel
			util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Member: "+name+" already exists.")
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
				util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "First member for the discord user: "+discordname+". This is required to be a main.")
				msg := "First member for the discord user: " + discordname + ". This is required to be a main."
				return msg
			}

			// Ensure there can only be 1 main at a time
			if existingMember && main {
				// Send the failed addition message in the previously created private channel
				util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "There can only be 1 main at a time per user: "+discordname)
				msg := "There can only be 1 main at a time per user: " + discordname
				return msg
			}

			// Only add a clan point entry for a main
			if main {
				s.cp[name] = 0
				s.mainAndAlts[name] = []string{}
			} else {
				mainName := ""
				for name, member := range s.members {
					if member.DiscordId == discordid && member.Main {
						mainName = name
					}
				}
				s.mainAndAlts[mainName] = append(s.mainAndAlts[mainName], name)
			}

			s.members[name] = util.MemberInfo{
				DiscordId:   discordid,
				DiscordName: discordname,
				Main:        main,
			}

			s.temple.AddMemberToTemple(ctx, name)
			s.templeUsernames[strings.ToLower(name)] = name
			logger.Debug("You have successfully added a new member: " + name)
			msg := "You have successfully added a new member: " + name
			return msg
		}
	case "Remove":
		util.LogAdminAction(logger, s.config.DiscAuditChan, i.Member.User.Username, i.Member.User.AvatarURL(""), session, fmt.Sprintf("Admin invoked player remove command with options: option=%s, name=%s, newName=%s, discordid=%d, discordname=%s, main=%t", option, name, newName, discordid, discordname, main))
		// Remove the user from the temple page
		s.temple.RemoveMemberFromTemple(ctx, name)

		if _, ok := s.members[name]; ok {
			// If the account we're deleting is a main, check to see if there are any other accounts for this discord user
			// If there is, just assign the first instance as main - if not, just delete
			if s.members[name].Main {
				logger.Debug("Deleting player is a main, searching for another account to become main...")
				if len(s.mainAndAlts[name]) > 0 {
					newMain := s.mainAndAlts[name][0]
					logger.Debug("Found new main: " + newMain + ". Using as new main for user: " + s.members[newMain].DiscordName)
					s.members[newMain] = util.MemberInfo{
						DiscordId:   s.members[newMain].DiscordId,
						DiscordName: s.members[newMain].DiscordName,
						Feedback:    s.members[newMain].Feedback,
						Main:        true,
					}
					s.cp[newMain] = s.cp[name]

					// Update HOF Speed times from deleted user to new main user
					updatedSpeedInfo := make(map[string]util.SpeedInfo)
					for boss, speedInfo := range s.speed {
						updatedSpeedInfo[boss] = util.SpeedInfo{
							PlayersInvolved: strings.Replace(speedInfo.PlayersInvolved, name, newMain, -1),
							Time:            speedInfo.Time,
							URL:             speedInfo.URL,
						}
					}
					s.speed = updatedSpeedInfo

					// Promote newMain as the main of the group of alts as well
					s.mainAndAlts[newMain] = s.mainAndAlts[name][1:]
					delete(s.mainAndAlts, name)
					delete(s.members, name)
					delete(s.cp, name)
				} else {
					// NOTE: We will keep the speed times for the deleted user, if you want to remove this speed
					// time, it must be invoked using the admin command
					delete(s.mainAndAlts, name)
					delete(s.members, name)
					delete(s.cp, name)
				}
			} else {
				// Determine the main account and remove this name from the list of alts
				for main, alts := range s.mainAndAlts {
					for place, alt := range alts {
						if strings.Compare(alt, name) == 0 {
							s.mainAndAlts[main] = append(s.mainAndAlts[main][:place], s.mainAndAlts[main][place+1:]...)
						}
					}
				}

				delete(s.members, name)
			}

			logger.Debug("You have successfully removed a member: " + name)
			msg := "You have successfully removed a member: " + name
			return msg

		} else {
			// Send the failed removal message in the previously created private channel
			util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Member: "+name+" does not exist.")
			msg := "Member: " + name + " does not exist."
			return msg
		}
	case "Name Change":
		util.LogAdminAction(logger, s.config.DiscAuditChan, i.Member.User.Username, i.Member.User.AvatarURL(""), session, fmt.Sprintf("Admin invoked player name change command with options: option=%s, name=%s, newName=%s, discordid=%d, discordname=%s, main=%t", option, name, newName, discordid, discordname, main))
		if _, ok := s.members[name]; ok {
			// Remove the user from the temple page and add new name
			s.temple.RemoveMemberFromTemple(ctx, name)
			s.temple.AddMemberToTemple(ctx, newName)

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

			// If the name change is for a main, change the key to the mainAndAlts
			if _, ok := s.mainAndAlts[name]; ok {
				s.mainAndAlts[newName] = s.mainAndAlts[name]
				delete(s.mainAndAlts, name)
			} else {
				// If the name change if for an alt, determine the main and change the value of the alt
				mainName := ""
				for main, alts := range s.mainAndAlts {
					for place, alt := range alts {
						if strings.Compare(alt, name) == 0 {
							mainName = main
							s.mainAndAlts[main] = append(s.mainAndAlts[main][:place], s.mainAndAlts[main][place+1:]...)
						}
					}
				}
				s.mainAndAlts[mainName] = append(s.mainAndAlts[mainName], newName)
			}

			s.members[newName] = s.members[name]
			s.cp[newName] = s.cp[name]
			s.templeUsernames[strings.ToLower(newName)] = newName
			delete(s.cp, name)
			delete(s.members, name)

			logger.Debug("You have successfully changed names from: " + name + " to: " + newName)
			msg := "You have successfully changed names from: " + name + " to: " + newName

			return msg

		} else {
			// Send the failed removal message in the previously created private channel
			util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Member: "+name+" does not exist.")
			msg := "Member: " + name + " does not exist."
			return msg
		}
	case "Update Main":
		util.LogAdminAction(logger, s.config.DiscAuditChan, i.Member.User.Username, i.Member.User.AvatarURL(""), session, fmt.Sprintf("Admin invoked player update main command with options: option=%s, name=%s, newName=%s, discordid=%d, discordname=%s, main=%t", option, name, newName, discordid, discordname, main))
		// Ensure it is a main first
		if _, ok := s.cp[name]; !ok {
			util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Name: "+name+" is not a main")
			msg := "Name: " + name + " is not a main. Ensure a main is used in the name section and the transferring name (not a main) is in the new-name section"
			return msg
		} else if _, ok := s.cp[newName]; ok {
			// Ensure the new-name is not a main
			util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "New Name: "+newName+" is a main")
			msg := "New Name: " + newName + " is a main. Ensure a main is used in the name section and the transferring name (not a main) is in the new-name section"
			return msg
		}

		// Ensure the name and newName belong to the same discord user
		if s.members[name].DiscordId != s.members[newName].DiscordId {
			util.LogError(logger, s.config.DiscAuditChan, session, i.Member.User.Username, i.Member.User.AvatarURL(""), "Main: "+name+" and New Main: "+newName+" do not belong to the same discord id.")
			msg := "Main: " + name + " and New Main: " + newName + " do not belong to the same discord id."
			return msg
		}

		// Use the new name as the clan points owner
		s.cp[newName] = s.cp[name]
		delete(s.cp, name)

		// Use the new name as the key in the mainAndAlts
		for place, alt := range s.mainAndAlts[name] {
			if strings.Compare(alt, newName) == 0 {
				s.mainAndAlts[newName] = append(s.mainAndAlts[name][:place], s.mainAndAlts[name][place+1:]...)
				s.mainAndAlts[newName] = append(s.mainAndAlts[newName], name)
				delete(s.mainAndAlts, name)
			}
		}

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

func (s *Service) updateCpPoints(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) string {
	options := i.ApplicationCommandData().Options[0].Options
	logger := flume.FromContext(ctx)

	player := ""
	cp := 0
	addOrRemove := ""

	for _, option := range options {
		switch option.Name {
		case "player":
			player = option.Value.(string)
		case "amount-of-cp":
			cp = int(option.Value.(float64))
		case "option":
			addOrRemove = option.Value.(string)
		}
	}

	// Check to see if there are multiple people
	players := util.WhiteStripCommas(player)
	listOfPlayers := strings.Split(players, ",")
	util.LogAdminAction(logger, s.config.DiscAuditChan, i.Member.User.Username, i.Member.User.AvatarURL(""), session, fmt.Sprintf("Admin invoked points add command with options: player=%s, cp=%d, addOrRemove=%s", player, cp, addOrRemove))
	for _, playerName := range listOfPlayers {
		switch addOrRemove {
		case "Add":
			logger.Info("Adding " + strconv.Itoa(cp) + " clan point(s) to " + playerName)
			s.cp[playerName] += cp
		case "Remove":
			logger.Info("Removing " + strconv.Itoa(cp) + " clan point(s) to " + playerName)
			if s.cp[playerName]-cp < 0 {
				s.cp[playerName] = 0
			} else {
				s.cp[playerName] -= cp
			}
		}
	}

	s.updateCpLeaderboard(ctx, session, i.Member.User)

	logger.Info("Successfully managed cp for " + player)
	return "Successfully managed cp for " + player
}
