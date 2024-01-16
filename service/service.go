package service

import (
	"context"
	"fmt"
	"github.com/gemalto/flume"
	"github.com/go-co-op/gocron"
	"log"
	"net/http"
	"os"
	"os/signal"
	"osrs-disc-bot/util"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Service struct {
	collectionLog collectionLog
	sheets        sheets
	imgur         imgur
	temple        temple
	runescape     runescape

	cp               map[string]int
	cpscreenshots    map[string]string
	speed            map[string]util.SpeedInfo
	speedscreenshots map[string]util.SpeedScInfo
	log              flume.Logger
	tid              int

	feedback map[string]string

	registeredCommands []*discordgo.ApplicationCommand

	config *util.Config
	client *http.Client

	scheduler *gocron.Scheduler
}

func NewService(config *util.Config, collectionLog collectionLog, sheets sheets, imgur imgur, temple temple, runescape runescape) *Service {
	logger := flume.New("service")
	if config.LogDebug {
		_ = flume.Configure(flume.Config{Development: true, Levels: "*"})
	} else {
		_ = flume.Configure(flume.Config{Development: true})
	}
	client := &http.Client{Timeout: 30 * time.Second}
	s := gocron.NewScheduler(time.UTC)
	s.SingletonModeAll()
	return &Service{
		collectionLog: collectionLog,
		sheets:        sheets,
		imgur:         imgur,
		temple:        temple,
		runescape:     runescape,

		cp:               make(map[string]int),
		cpscreenshots:    make(map[string]string),
		speed:            make(map[string]util.SpeedInfo),
		speedscreenshots: make(map[string]util.SpeedScInfo),
		log:              logger,
		tid:              1,

		feedback: make(map[string]string),

		config: config,
		client: client,

		scheduler: s,
	}
}

// StartDiscordIRC uses discordgo as an intro to discord IRC
func (s *Service) StartDiscordIRC() {
	s.log.Info("Initializing OSRS Disc Bot...")
	ctx := context.Background()
	s.sheets.InitializeCpFromSheet(ctx, s.cp)
	s.sheets.InitializeSpeedsFromSheet(ctx, s.speed)
	s.sheets.InitializeFeedbackFromSheet(ctx, s.feedback)
	s.tid = s.sheets.InitializeTIDFromSheet(ctx)

	// Create a new discord session
	session, err := discordgo.New("Bot " + s.config.DiscBotToken)
	if err != nil {
		log.Fatal(err)
	}

	// Set the session's intent
	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// Open up the session and defer the closing of the session
	_ = session.Open()
	defer func(discord *discordgo.Session) {
		err := discord.Close()
		if err != nil {

		}
	}(session)

	// Create handler for listening for submission messages
	session.AddHandler(s.submissionApproval)
	session.AddHandler(s.slashCommands)
	session.AddHandler(s.listenForAllChannelMessages)

	// Kick off gocron for updating the Hall Of fame
	s.initCron(ctx, session)
	s.initSlashCommands(ctx, session)
	s.updatePpLeaderboard(ctx, session)

	s.log.Info("OSRS Disc Bot is now online!")
	s.blockUntilInterrupt(ctx, session)
	s.log.Info("OSRS Disc Bot is now offline!")
}

func (s *Service) blockUntilInterrupt(ctx context.Context, session *discordgo.Session) {
	// Block so that it continues to run the bot
	sigchan := make(chan os.Signal)
	signal.Notify(sigchan,
		// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGKILL, // "always fatal", "SIGKILL and SIGSTOP may not be caught by a program"
		syscall.SIGHUP,  // "terminal is disconnected"
		os.Interrupt,    // os interrupt
	)
	<-sigchan

	// Once the program is interrupted, update the Google Sheet clan points & screenshot sheets
	s.log.Debug("Running cp sheets updates...")
	s.sheets.UpdateCpSheet(ctx, s.cp)
	s.log.Debug("Running cp sc sheets updates...")
	s.sheets.UpdateCpScreenshotsSheet(ctx, s.cpscreenshots)
	s.log.Debug("Running speed updates...")
	s.sheets.UpdateSpeedSheet(ctx, s.speed)
	s.log.Debug("Running speed sc sheets updates...")
	s.sheets.UpdateSpeedScreenshotsSheet(ctx, s.speedscreenshots)
	s.log.Debug("Running feedback sheets updates...")
	s.sheets.UpdateFeedbackChannel(ctx, s.feedback)
	s.log.Debug("Running tid sheets updates...")
	s.sheets.UpdateTIDFromSheet(ctx, s.tid)
	s.log.Debug("Finished running sheets updates")

	// Delete the slash commands the bot creates
	//session.ApplicationCommandDelete(session.State.User.ID, s.config.DiscGuildId, "")
	//s.removeSlashCommands(session)

	// Stop the cron scheduler
	s.scheduler.Stop()
}

// initCron will instantiate the HallOfFameRequestInfos and kick off the cron job
func (s *Service) initCron(ctx context.Context, session *discordgo.Session) {
	s.log.Info("Initializing Hall Of Fame Cron Job...")

	// Kick off a scheduled job at a configured time
	job, err := s.scheduler.Every(1).Day().At(s.config.CronKickoffTime).Do(func() {
		s.log.Debug("Running Cron Job to update the Hall Of Fame, Collection Log, and Leagues...")
		s.updateKcHOF(ctx, session)
		s.updateSpeedHOF(ctx, session, "TzHaar", "Slayer", "Nightmare", "Nex", "Solo Bosses", "Chambers Of Xeric", "Theatre Of Blood", "Tombs Of Amascut", "Agility")
		s.updateColLog(ctx, session)
		s.log.Debug("Finished running Cron Job to update the Hall Of Fame, Collection Log, and Leagues!")
	})
	if err != nil {
		// handle the error related to setting up the job
		fmt.Printf("Job: %#v, Error: %#v", job, err)
	}
	job.SingletonMode()
	s.scheduler.StartAsync()
}

func (s *Service) listenForAllChannelMessages(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Don't handle message if it's created by the discord bot
	//if message.Author.ID != session.State.User.ID {
	//	return
	//}

	ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid).With("user", message.Author.Username))
	defer func() { s.tid++ }()

	// Run certain tasks depending on the channel the message was posted in
	switch channel := message.ChannelID; channel {
	case s.config.DiscLootLogChan:
		s.listenForLootLog(ctx, session, message)
	default:
		// Return if the message was not posted in one of the channels we are handling
		return
	}
}
