package main

import (
	"bufio"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

type Config struct {
	DiscBotToken         string `mapstructure:"DISCORD_BOT_TOKEN"`
	DiscSubChan          string `mapstructure:"DISCORD_SUBMISSION_CHANNEL"`
	DiscSignUpChan       string `mapstructure:"DISCORD_SIGNUP_CHANNEL"`
	DiscLeaderboardChan  string `mapstructure:"DISCORD_LEADERBOARD_CHANNEL"`
	DiscDirectoryChan    string `mapstructure:"DISCORD_DIRECTORY_CHANNEL"`
	DiscSlayerBossesChan string `mapstructure:"DISCORD_SLAYER_BOSSES_CHANNEL"`
	DiscGwdChan          string `mapstructure:"DISCORD_GWD_CHANNEL"`
	DiscWildyChan        string `mapstructure:"DISCORD_WILDERNESS_CHANNEL"`
	DiscOtherChan        string `mapstructure:"DISCORD_OTHER_CHANNEL"`
	DiscMiscChan         string `mapstructure:"DISCORD_MISC_CHANNEL"`
	DiscDT2Chan          string `mapstructure:"DISCORD_DT2_CHANNEL"`
	DiscRaidsChan        string `mapstructure:"DISCORD_RAIDS_CHANNEL"`
	DiscPVPChan          string `mapstructure:"DISCORD_PVP_CHANNEL"`
	DiscCluesChan        string `mapstructure:"DISCORD_CLUES_CHANNEL"`
	TempleGroupId        string `mapstructure:"TEMPLE_GROUP_ID"`
	TempleGroupKey       string `mapstructure:"TEMPLE_GROUP_KEY"`
}

var config Config
var submissions map[string]int

func init() {
	// Initialize the Viper configuration ingestion and unmarshal
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Failed to read environment variable, exiting now...")
		return
	}
	err = viper.Unmarshal(&config)

	// Read the clan points file and break it up into an array then into a map
	f, err := os.ReadFile("submissions/clan_points.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Break up the names in the clan points file and store them on a map
	names := strings.Split(string(f), "\n")
	submissions = make(map[string]int)
	for _, nameAndCp := range names {
		if len(nameAndCp) == 0 {
			continue
		}
		parts := strings.Split(nameAndCp, ":")
		cp, _ := strconv.Atoi(parts[1])
		submissions[parts[0]] = cp
	}
}

func main() {
	// Throw a go func that will capture signal interrupts and will populate the submissions file
	go func() {
		sigchan := make(chan os.Signal)
		signal.Notify(sigchan,
			// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
			syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
			syscall.SIGINT,  // Ctrl+C
			syscall.SIGQUIT, // Ctrl-\
			syscall.SIGKILL, // "always fatal", "SIGKILL and SIGSTOP may not be caught by a program"
			syscall.SIGHUP,  // "terminal is disconnected"
		)
		<-sigchan

		// Once the program is interrupted, remake the clan_points.txt file with the locally stored submissions map
		f, _ := os.Create("submissions/clan_points.txt")
		datawriter := bufio.NewWriter(f)
		for name, cp := range submissions {
			cpToString := strconv.Itoa(cp)
			_, _ = datawriter.WriteString(name + ":" + cpToString + "\n")
		}
		datawriter.Flush()
		f.Close()

		os.Exit(0)
	}()

	startDiscordIRC()
}
