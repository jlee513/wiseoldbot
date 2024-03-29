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
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Service struct {
	collectionLog collectionLog
	sheets        sheets
	imageservice  imageservice
	temple        temple
	runescape     runescape
	pastebin      pastebin

	cp               map[string]int
	cpscreenshots    map[string]util.CpScInfo
	speed            map[string]util.SpeedInfo
	speedscreenshots map[string]util.SpeedScInfo
	members          map[string]util.MemberInfo
	mainAndAlts      map[string][]string
	templeUsernames  map[string]string
	discGuides       map[string][]util.GuideInfo
	speedCategory    map[string][]string

	log                flume.Logger
	tid                int
	registeredCommands []*discordgo.ApplicationCommand

	config *util.Config
	client *http.Client

	scheduler *gocron.Scheduler
}

func NewService(config *util.Config, collectionLog collectionLog, sheets sheets, imageservice imageservice, temple temple, runescape runescape, pastebin pastebin) *Service {
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
		imageservice:  imageservice,
		temple:        temple,
		runescape:     runescape,
		pastebin:      pastebin,

		cp:               make(map[string]int),
		cpscreenshots:    make(map[string]util.CpScInfo),
		speed:            make(map[string]util.SpeedInfo),
		speedscreenshots: make(map[string]util.SpeedScInfo),
		members:          make(map[string]util.MemberInfo),
		mainAndAlts:      make(map[string][]string),
		templeUsernames:  make(map[string]string),
		discGuides:       make(map[string][]util.GuideInfo),
		speedCategory:    make(map[string][]string),

		log:    logger,
		tid:    1,
		config: config,
		client: client,

		scheduler: s,
	}
}

// StartDiscordIRC uses discordgo as an intro to discord IRC
func (s *Service) StartDiscordIRC() {
	s.log.Info("Initializing OSRS Disc Bot...")
	ctx := flume.WithLogger(context.Background(), s.log)

	// Initialize from Sheets
	s.sheets.InitializeCpFromSheet(ctx, s.cp)
	s.sheets.InitializeSpeedsFromSheet(ctx, s.speed)
	s.sheets.InitializeMembersFromSheet(ctx, s.members)
	s.tid = s.sheets.InitializeTIDFromSheet(ctx)

	// Initialize from Pastebin
	s.pastebin.UpdateGuideList(ctx, s.discGuides)

	// Determine extra information
	s.determineMainAndAlts(ctx)
	s.determineSpeedCategory(ctx)

	// Set templeUsernames to all lowercase to avoid issues with capitalization's and preserve original username
	// capitalization for HOF and other leaderboards
	for name := range s.members {
		s.templeUsernames[strings.ToLower(name)] = name
	}

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
	botUser, _ := session.User(s.config.DiscBotId)
	s.updateCpLeaderboard(ctx, session, botUser)

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

	s.updateAllGoogleSheets(ctx)

	// Delete the slash commands the bot creates
	//session.ApplicationCommandDelete(session.State.User.ID, s.config.DiscGuildId, "")
	//s.removeSlashCommands(session)

	// Stop the cron scheduler
	s.scheduler.Stop()
}

func (s *Service) determineSpeedCategory(ctx context.Context) {
	for bossName, speedInfo := range s.speed {
		if _, ok := s.speedCategory[speedInfo.Category]; ok {
			s.speedCategory[speedInfo.Category] = append(s.speedCategory[speedInfo.Category], bossName)
		} else {
			s.speedCategory[speedInfo.Category] = []string{bossName}
		}
	}
}

func (s *Service) determineMainAndAlts(ctx context.Context) {
	// Key: Discord ID, Value: Name
	altsBeforeMainIsFound := make(map[int][]string)
	mainIdToName := make(map[int]string)

	for memberName, member := range s.members {
		if member.Main {
			// If it exists within alts, move all the alt names to the mainAndAlts with the main as the key
			if _, ok := altsBeforeMainIsFound[member.DiscordId]; ok {
				s.mainAndAlts[memberName] = altsBeforeMainIsFound[member.DiscordId]
				mainIdToName[member.DiscordId] = memberName
				delete(altsBeforeMainIsFound, member.DiscordId)
			} else {
				// If it doesn't exist, just set the mainKcs object with the member object
				mainIdToName[member.DiscordId] = memberName
				s.mainAndAlts[memberName] = []string{}
			}
		} else {
			// If it exists within the main, determine the main name by using the discord id and append it to the array
			if _, ok := mainIdToName[member.DiscordId]; ok {
				main := mainIdToName[member.DiscordId]
				s.mainAndAlts[main] = append(s.mainAndAlts[main], memberName)
			} else if _, ok := altsBeforeMainIsFound[member.DiscordId]; ok {
				// Check to see if an altsBeforeMainIsFound exists - if it does, add it to that
				altsBeforeMainIsFound[member.DiscordId] = append(altsBeforeMainIsFound[member.DiscordId], memberName)
			} else {
				// Otherwise, just add it to altsBeforeMainIsFound
				altsBeforeMainIsFound[member.DiscordId] = []string{memberName}
			}
		}
	}
}

func (s *Service) updateAllGoogleSheets(ctx context.Context) {
	logger := flume.FromContext(ctx)
	// Once the program is interrupted, update the Google Sheet clan points & screenshot sheets
	logger.Debug("Running cp sheets updates...")
	s.sheets.UpdateCpSheet(ctx, s.cp)
	logger.Debug("Running cp sc sheets updates...")
	s.sheets.UpdateCpScreenshotsSheet(ctx, s.cpscreenshots)
	logger.Debug("Running speed updates...")
	s.sheets.UpdateSpeedSheet(ctx, s.speed, s.speedCategory)
	logger.Debug("Running speed sc sheets updates...")
	s.sheets.UpdateSpeedScreenshotsSheet(ctx, s.speedscreenshots)
	logger.Debug("Running tid sheets updates...")
	s.sheets.UpdateTIDFromSheet(ctx, s.tid)
	logger.Debug("Running members sheets updates...")
	s.sheets.UpdateMembersSheet(ctx, s.members)
	logger.Info("Finished running sheets updates")
}

// initCron will instantiate the HallOfFameRequestInfos and kick off the cron job
func (s *Service) initCron(ctx context.Context, session *discordgo.Session) {
	s.log.Info("Initializing Hall Of Fame Cron Job...")

	// Kick off a scheduled job at a configured time
	job, err := s.scheduler.Every(1).Day().At(s.config.CronKickoffTime).Do(func() {
		s.log.Debug("Running Cron Job to update the Hall Of Fame, Collection Log, and Leagues...")
		botUser, _ := session.User(s.config.DiscBotId)
		s.updateKcHOF(ctx, session, botUser)
		s.updateSpeedHOF(ctx, session, botUser, "TzHaar", "Slayer", "Nightmare", "Nex", "Solo Bosses", "Chambers Of Xeric", "Theatre Of Blood", "Tombs Of Amascut", "Agility")
		s.updateColLog(ctx, session, botUser)
		s.updateTempleMilestones(ctx, session, botUser)
		s.updateAllGoogleSheets(ctx)
		s.log.Debug("Finished running Cron Job to update the Hall Of Fame, Collection Log, and Google Sheets!")
	})
	if err != nil {
		// handle the error related to setting up the job
		fmt.Printf("Job: %#v, Error: %#v", job, err)
	}
	job.SingletonMode()
	s.scheduler.StartAsync()
}

func (s *Service) listenForAllChannelMessages(session *discordgo.Session, message *discordgo.MessageCreate) {
	// If the author of the message is not the webhook, ignore the message
	if message.Author.ID != s.config.DiscLootLogWebhook {
		return
	}

	// Run certain tasks depending on the channel the message was posted in
	switch channel := message.ChannelID; channel {
	case s.config.DiscLootLogChan:
		ctx := flume.WithLogger(context.Background(), s.log.With("transactionID", s.tid).With("user", message.Author.Username))
		defer func() { s.tid++ }()
		s.listenForLootLog(ctx, session, message)
	default:
		// Return if the message was not posted in one of the channels we are handling
		return
	}
}
