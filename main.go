package main

import (
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

var config Config
var submissions map[string]int
var cpscreenshots map[string]string

func init() {
	// Initialize the Viper configuration ingestion and unmarshal
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Failed to read environment variable, exiting now...")
		return
	}
	err = viper.Unmarshal(&config)

	initCpSheet()
	cpscreenshots = make(map[string]string)
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

		// Once the program is interrupted, update the google clan points sheet
		updateCpSheet()
		updateCpScreenshotsSheet()
	}()

	startDiscordIRC()
}
