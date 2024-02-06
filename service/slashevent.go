package service

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/gemalto/flume"
	"osrs-disc-bot/util"
)

func (s *Service) handleEventApproval(ctx context.Context, session *discordgo.Session, r *discordgo.MessageReactionAdd) {
	logger := flume.FromContext(ctx)
	switch r.Emoji.Name {
	case "✅":
		// TODO: Write when there's an event

		// Delete the screenshot in the page
		err := session.ChannelMessageDelete(s.config.DiscCpApprovalChan, r.MessageID)
		if err != nil {
			util.LogError(logger, s.config.DiscAuditChan, session, r.Member.User.Username, r.Member.User.AvatarURL(""), "Failed to delete cp approval message: "+err.Error())
		}
	case "❌":
		// TODO: Find a way to let the user know that their submission has been rejected

		// Delete the screenshot in the page
		err := session.ChannelMessageDelete(s.config.DiscCpApprovalChan, r.MessageID)
		if err != nil {
			util.LogError(logger, s.config.DiscAuditChan, session, r.Member.User.Username, r.Member.User.AvatarURL(""), "Failed to delete cp approval message: "+err.Error())
		}
	}
}
