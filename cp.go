package main

import (
	"fmt"
	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
)

// updateLeaderboard will update the cp-leaderboard channel in discord with a new ranking of everyone in the clan
func updateLeaderboard(session *discordgo.Session) {
	// Update the #cp-leaderboard
	keys := make([]string, 0, len(submissions))
	for key := range submissions {
		keys = append(keys, key)
	}

	// Sort the map based on the values
	sort.SliceStable(keys, func(i, j int) bool {
		return submissions[keys[i]] > submissions[keys[j]]
	})

	// Create the leaderboard message that will be sent
	leaderboard := ""
	for placement, k := range keys {
		leaderboard = leaderboard + strconv.Itoa(placement+1) + ") " + k + ": " + strconv.Itoa(submissions[k]) + "\n"
	}

	// Retrieve the one channel message and delete it in the leaderboard channel
	messages, err := session.ChannelMessages(config.DiscLeaderboardChan, 1, "", "", "")
	if err != nil {
		return
	}
	err = session.ChannelMessageDelete(config.DiscLeaderboardChan, messages[0].ID)
	if err != nil {
		return
	}

	_, err = session.ChannelMessageSendEmbed(config.DiscLeaderboardChan, embed.NewEmbed().
		SetTitle("Ponies Clan Points Leaderboard").
		SetDescription(fmt.Sprintf(leaderboard)).
		SetColor(0x1c1c1c).SetThumbnail("https://i.imgur.com/O4NzB95.png").MessageEmbed)
	if err != nil {
		return
	}
}

func downloadSubmissionScreenshot(submissionLink string) {
	// Build fileName from fullPath
	fileURL, err := url.Parse(submissionLink)
	if err != nil {
		log.Fatal(err)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	// Create blank file
	file, err := os.Create("submissions/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	resp, err := client.Get(submissionLink)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)

	defer file.Close()
}
