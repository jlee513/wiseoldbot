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
	s.sheets.UpdateCpSheet(ctx, s.cp)
	s.sheets.UpdateCpScreenshotsSheet(ctx, s.cpscreenshots)
	s.sheets.UpdateSpeedSheet(ctx, s.speed)
	s.sheets.UpdateSpeedScreenshotsSheet(ctx, s.speedscreenshots)

	// Delete the slash commands the bot creates
	s.removeSlashCommands(session)

	// Stop the cron scheduler
	s.scheduler.Stop()
}

// initCron will instantiate the HallOfFameRequestInfos and kick off the cron job
func (s *Service) initCron(ctx context.Context, session *discordgo.Session) {
	s.log.Info("Initializing Hall Of Fame Cron Job...")

	// HOF KC
	slayerBosses := util.HofRequestInfo{Name: "Slayer Bosses", DiscChan: s.config.DiscSlayerBossesChan, AfterId: "1194801106291785778", Bosses: util.HofSlayerBosses}
	gwd := util.HofRequestInfo{Name: "GWD Bosses", DiscChan: s.config.DiscGwdChan, AfterId: "1194801166429724884", Bosses: util.HofGWDBosses}
	wildy := util.HofRequestInfo{Name: "Wildy Bosses", DiscChan: s.config.DiscWildyChan, AfterId: "1194801335376285726", Bosses: util.HofWildyBosses}
	other := util.HofRequestInfo{Name: "Other Bosses", DiscChan: s.config.DiscOtherChan, AfterId: "1194801512870846535", Bosses: util.HofOtherBosses}
	misc := util.HofRequestInfo{Name: "Miscellaneous Bosses", DiscChan: s.config.DiscMiscChan, AfterId: "1194804397507620935", Bosses: util.HofMiscBosses}
	dt2 := util.HofRequestInfo{Name: "Desert Treasure 2 Bosses", DiscChan: s.config.DiscDT2Chan, AfterId: "1194802032855498832", Bosses: util.HofDT2Bosses}
	raids := util.HofRequestInfo{Name: "Raids Bosses", DiscChan: s.config.DiscRaidsChan, AfterId: "1194802206487089182", Bosses: util.HofRaidsBosses}
	pvp := util.HofRequestInfo{Name: "PVP", DiscChan: s.config.DiscPVPChan, AfterId: "1194802450209718272", Bosses: util.HofPVPBosses}
	clues := util.HofRequestInfo{Name: "Clues", DiscChan: s.config.DiscCluesChan, AfterId: "1194802590270103582", Bosses: util.HofCluesBosses}

	// HOF Speed
	tzhaar := util.SpeedsRequestInfo{Name: "TzHaar", DiscChan: s.config.DiscSpeedTzhaarChan, AfterId: "1194999599652425778", Bosses: util.HofSpeedTzhaar}
	slayer := util.SpeedsRequestInfo{Name: "Slayer", DiscChan: s.config.DiscSpeedSlayerChan, AfterId: "1194999714710573078", Bosses: util.HofSpeedSlayer}
	nightmare := util.SpeedsRequestInfo{Name: "Nightmare", DiscChan: s.config.DiscSpeedNightmareChan, AfterId: "1195000377288958023", Bosses: util.HofSpeedNightmare}
	nex := util.SpeedsRequestInfo{Name: "Nex", DiscChan: s.config.DiscSpeedNexChan, AfterId: "1195000695594684416", Bosses: util.HofSpeedNex}
	solo := util.SpeedsRequestInfo{Name: "Solo", DiscChan: s.config.DiscSpeedSoloChan, AfterId: "1195000959911350294", Bosses: util.HofSpeedSolo}
	cox := util.SpeedsRequestInfo{Name: "COX", DiscChan: s.config.DiscSpeedCOXChan, AfterId: "1195001187276161155", Bosses: util.HofSpeedCox}
	tob := util.SpeedsRequestInfo{Name: "TOB", DiscChan: s.config.DiscSpeedTOBChan, AfterId: "1195001367685779509", Bosses: util.HofSpeedTob}
	toa := util.SpeedsRequestInfo{Name: "TOA", DiscChan: s.config.DiscSpeedTOAChan, AfterId: "1195001626604355656", Bosses: util.HofSpeedToa}
	agility := util.SpeedsRequestInfo{Name: "Agility", DiscChan: s.config.DiscSpeedAgilityChan, AfterId: "1195002755132174368", Bosses: util.HofSpeedAgility}

	//s.updateKcHOF(ctx, session, slayerBosses, gwd, wildy, other, misc, dt2, raids, pvp, clues)
	//s.updateSpeedHOF(ctx, session, tzhaar, slayer, nightmare, nex, solo, cox, tob, toa, agility)
	//s.updateColLog(ctx, session)

	// Kick off a scheduled job at a configured time
	job, err := s.scheduler.Every(1).Day().At(s.config.CronKickoffTime).Do(func() {
		s.log.Debug("Running Cron Job to update the Hall Of Fame, Collection Log, and Leagues...")
		s.updateKcHOF(ctx, session, slayerBosses, gwd, wildy, other, misc, dt2, raids, pvp, clues)
		s.updateSpeedHOF(ctx, session, tzhaar, slayer, nightmare, nex, solo, cox, tob, toa, agility)
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

	// Run certain tasks depending on the channel the message was posted in
	switch channel := message.ChannelID; channel {
	case s.config.DiscLootLogChan:
		s.listenForLootLog(session, message)
	default:
		// Return if the message was not posted in one of the channels we are handling
		return
	}
}
