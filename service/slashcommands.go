package service

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/gemalto/flume"
	"log"
	"osrs-disc-bot/util"
	"strconv"
	"strings"
)

func (s *Service) initSlashCommands(ctx context.Context, session *discordgo.Session) {
	logger := flume.FromContext(ctx)
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "speed-submission",
			Description: "Speed submissions for ponies",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "category",
					Description:  "Category of speed",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
				{
					Name:         "boss",
					Description:  "Boss name submitting for",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "speed-time",
					Description: "Only use if making a speed submission in format: hh:mm:ss.ms",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "player-names",
					Description: "Comma separated list of players involved",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "screenshot",
					Description: "Screenshot of the submission",
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "i-imgur-link",
					Description: "Imgur link of the submission",
				},
			},
		},
		{
			Name:        "pp-submission",
			Description: "Ponies points submissions for ponies",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "player-names",
					Description: "Comma separated list of players involved",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "screenshot",
					Description: "Screenshot of the submission",
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "i-imgur-link",
					Description: "Imgur link of the submission",
				},
			},
		},
		{
			Name:        "guide",
			Description: "Administration of Guides",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option",
					Description: "Choose one of the following: Update",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Update",
							Value: "Update",
						},
					},
				},
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "guide",
					Description:  "Guide name",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		{
			Name:        "admin",
			Description: "Admin commands",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "player",
					Description: "Player administration commands",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "option",
							Description: "Choose one of the following: Add, Remove",
							Required:    true,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "Add",
									Value: "Add",
								},
								{
									Name:  "Remove",
									Value: "Remove",
								},
								{
									Name:  "Name Change",
									Value: "Name Change",
								},
								{
									Name:  "Update Main",
									Value: "Update Main",
								},
							},
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "name",
							Description: "Player name",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "discord-id",
							Description: "Player discord id (this is a number)",
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "discord-name",
							Description: "Player discord name",
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "new-name",
							Description: "New player name",
						},
						{
							Type:        discordgo.ApplicationCommandOptionBoolean,
							Name:        "main",
							Description: "Is this the main account?",
						},
					},
				},
				{
					Name:        "update-instructions",
					Description: "Update Submission instructions",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "update-points",
					Description: "Update Pp for player",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "player",
							Description: "Player name",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "option",
							Description: "Choose one of the following: Add, Remove",
							Required:    true,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "Add",
									Value: "Add",
								},
								{
									Name:  "Remove",
									Value: "Remove",
								},
							},
						},
						{
							Type:        discordgo.ApplicationCommandOptionInteger,
							Name:        "amount-of-pp",
							Description: "Amount of Pp to manage for player",
							Required:    true,
						},
					},
				},
				{
					Name:        "speed",
					Description: "Speed Administration",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "action",
							Description: "Action to perform on speed",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "Add",
									Value: "Add",
								},
								{
									Name:  "Remove",
									Value: "Remove",
								},
								{
									Name:  "Update",
									Value: "Update",
								},
								{
									Name:  "Reset",
									Value: "Reset",
								},
							},
						},
						{
							Name:         "category",
							Description:  "Category of speed",
							Type:         discordgo.ApplicationCommandOptionString,
							Required:     true,
							Autocomplete: true,
						},
						{
							Name:         "existing-boss",
							Description:  "Existing boss name submitting for",
							Type:         discordgo.ApplicationCommandOptionString,
							Required:     false,
							Autocomplete: true,
						},
						{
							Name:        "new-boss",
							Description: "New Boss name adding to category selected",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						},
						{
							Name:        "update-speed-time",
							Type:        discordgo.ApplicationCommandOptionString,
							Description: "Only use if making a speed submission in format: hh:mm:ss.ms",
							Required:    false,
						},
						{
							Name:        "update-player-names",
							Type:        discordgo.ApplicationCommandOptionString,
							Description: "Comma separated list of players involved",
							Required:    false,
						},
					},
				},
				{
					Name:        "update-leaderboard",
					Description: "Update Leaderboard for player",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "leaderboard",
							Description:  "leaderboard name",
							Required:     true,
							Autocomplete: true,
						},
						{
							Name:         "thread",
							Description:  "Name of the thread you want to update",
							Type:         discordgo.ApplicationCommandOptionString,
							Required:     true,
							Autocomplete: true,
						},
					},
				},
				{
					Name:        "update-sheets",
					Description: "Update google sheets",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "update-guides-map",
					Description: "Update Guides that are stored in Map",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
	}

	// Iterate over all the commands and create the application command - we will save all the registered commands
	// into the service struct that will be used to delete all the commands on bot termination
	logger.Info("Adding all commands...")
	s.registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, s.config.DiscGuildId, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		s.registeredCommands[i] = cmd
	}
}

//func (s *Service) removeSlashCommands(session *discordgo.Session) {
//	logger.Info("Removing all commands...")
//
//	for _, v := range s.registeredCommands {
//		logger.Debug("REMOVING COMMAND: " + v.Name)
//		err := session.ApplicationCommandDelete(session.State.User.ID, s.config.DiscGuildId, v.ID)
//		if err != nil {
//			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
//		}
//	}
//}

func (s *Service) slashCommands(session *discordgo.Session, i *discordgo.InteractionCreate) {

	switch i.ApplicationCommandData().Name {
	case "pp-submission":
		ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid).With("user", i.Member.User.Username))
		defer func() { s.tid++ }()
		returnMessage := s.handlePPSubmission(ctx, session, i)
		err := util.InteractionRespond(session, i, returnMessage)
		if err != nil {
			s.log.Error("Failed to send interaction response: " + err.Error())
		}
	case "guide":
		s.handleGuide(session, i)
		return
	case "admin":
		s.handleAdmin(session, i)
		return
	case "speed-submission":
		s.handleSpeedSubmission(session, i)
	default:
		s.log.Error("ERROR: UNKNOWN COMMAND USED: " + i.ApplicationCommandData().Name)
		err := util.InteractionRespond(session, i, "Error - Unknown command used: "+i.ApplicationCommandData().Name)
		if err != nil {
			s.log.Error("Failed to send interaction response: " + err.Error())
		}
	}
}

func (s *Service) submissionApproval(session *discordgo.Session, r *discordgo.MessageReactionAdd) {

	// Don't handle message if it's created by the discord bot
	if r.UserID == session.State.User.ID {
		return
	}

	switch r.ChannelID {
	case s.config.DiscCpApprovalChan:
		ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid).With("user", r.Member.User.Username))
		defer func() { s.tid++ }()
		s.handlePPApproval(ctx, session, r)
	case s.config.DiscSpeedApprovalChan:
		ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid).With("user", r.Member.User.Username))
		defer func() { s.tid++ }()
		s.handleSpeedApproval(ctx, session, r)
	case s.config.DiscEventApprovalChan:
		ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid).With("user", r.Member.User.Username))
		defer func() { s.tid++ }()
		s.handleEventApproval(ctx, session, r)
	}
}

func (s *Service) checkOrCreateFeedbackChannel(ctx context.Context, session *discordgo.Session, user string, userId int, name string) string {
	logger := flume.FromContext(ctx)

	if len(name) > 0 {
		// If we are provided the username, we can skip the iteration through s.members
		if len(s.members[name].Feedback) > 0 {
			logger.Debug("Feedback channel found for user: " + s.members[name].DiscordName)
			return s.members[name].Feedback
		} else {
			logger.Debug("Feedback channel not found for user: " + user + " - will proceed to create one")
		}
	} else {
		// Iterate through all the members to see that when the discord name shows up, feedback is also set
		for player, member := range s.members {
			if strings.Compare(member.DiscordName, user) == 0 {
				name = player
				if len(member.Feedback) > 0 {
					logger.Debug("Feedback channel found for user: " + user)
					return member.Feedback
				} else {
					logger.Debug("Feedback channel not found for user: " + user + " - will proceed to create one")
				}
			}
		}
	}

	// Ensure that the user & userid is filled in before proceeding to create the channel
	if len(user) == 0 || userId == 0 {
		user = s.members[name].DiscordName
		userId = s.members[name].DiscordId
	}

	// If feedback is not set, we will create a feedback channel and set it
	channel, err := session.GuildChannelCreateComplex(s.config.DiscGuildId, discordgo.GuildChannelCreateData{
		Name: "feedback-" + user,
		Type: discordgo.ChannelTypeGuildText,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:    strconv.Itoa(userId),
				Type:  discordgo.PermissionOverwriteTypeMember,
				Allow: discordgo.PermissionAllText,
			},
			{
				ID:   s.config.DiscGuildId,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: discordgo.PermissionViewChannel,
			},
			{
				// Moderator Rank
				ID:    "1194691758353821847",
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionAllText,
			},
		},
	})
	if err != nil {
		logger.Error("Failed to create private text channel: " + err.Error())
	}
	s.members[name] = util.MemberInfo{
		DiscordId:   s.members[name].DiscordId,
		DiscordName: s.members[name].DiscordName,
		Feedback:    channel.ID,
	}
	logger.Debug("Successfully created feedback channel for user: " + user)
	return channel.ID

}
