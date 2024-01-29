package service

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/gemalto/flume"
	"osrs-disc-bot/util"
	"strings"
)

func (s *Service) handleGuide(session *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid).With("user", i.Member.User.Username))
		defer func() { s.tid++ }()
		s.guideAdminCommand(ctx, session, i)
	case discordgo.InteractionApplicationCommandAutocomplete:
		s.guideAdminAutocomplete(session, i)
	}
}

func (s *Service) guideAdminCommand(ctx context.Context, session *discordgo.Session, i *discordgo.InteractionCreate) {
	logger := flume.FromContext(ctx)
	options := i.ApplicationCommandData().Options

	option := options[0].Value.(string)
	guide := options[1].Value.(string)

	switch option {
	case "Update":
		// Remove leading and trailing whitespaces
		msg := strings.TrimSpace(guide)
		if _, ok := s.discGuides[guide]; ok {
			err := util.InteractionRespond(session, i, "Updating guide: "+guide)
			if err != nil {
				logger.Error("Failed to send guide interaction response: " + err.Error())
			}
			logger.Debug("Updating guide: " + msg)
			s.updateGuide(ctx, session, guide)
		} else {
			err := util.InteractionRespond(session, i, "Unknown guide chosen: "+guide)
			if err != nil {
				logger.Error("Failed to send interaction response: " + err.Error())
			}
			logger.Error("Unknown guide chosen: " + guide)
		}
	default:
		err := util.InteractionRespond(session, i, "Invalid guide management option chosen: "+option)
		if err != nil {
			logger.Error("Failed to send interaction response: " + err.Error())
		}
		logger.Error("Invalid guide management option chosen.")
	}
}

func (s *Service) guideAdminAutocomplete(session *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	var choices []*discordgo.ApplicationCommandOptionChoice
	switch {
	// In this case there are multiple autocomplete options. The Focused field shows which option user is focused on.
	case data.Options[1].Focused:
		for guideName, _ := range s.discGuides {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  guideName,
				Value: guideName,
			})
		}
	}
	err := util.InteractionRespondChoices(session, i, choices)
	if err != nil {
		s.log.Error("Failed to handle speed autocomplete: " + err.Error())
	}
}

func (s *Service) updateGuide(ctx context.Context, session *discordgo.Session, guideName string) {
	logger := flume.FromContext(ctx)
	guideInfos := s.discGuides[guideName]
	logger.Info("Updating guide: " + guideName + "...")

	for _, guideInfo := range guideInfos {
		logger.Info("Updating guide info: " + guideInfo.GuidePageName + "...")
		guide := s.pastebin.GetGuide(ctx, guideInfo.PastebinKey)
		guideArr := strings.Split(guide, "\n")

		err := util.DeleteBulkDiscordMessages(session, guideInfo.DiscChan)
		if err != nil {
			logger.Error("Failed to delete bulk discord messages: " + err.Error())
		}
		for _, line := range guideArr {
			// Remove leading and trailing whitespaces
			line = strings.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			_, err := session.ChannelMessageSend(guideInfo.DiscChan, line)
			if err != nil {
				logger.Error("ERROR SENDING MESSAGE: " + line + " TO: " + guideInfo.GuidePageName + " - " + err.Error())
				return
			}
		}
		logger.Info("Finished updating guide info: " + guideInfo.GuidePageName)
	}
	logger.Info("Finished updating guide: " + guideName)
}
