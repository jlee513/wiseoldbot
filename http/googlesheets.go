package http

import (
	"context"
	"net/http"
	"os"
	"osrs-disc-bot/util"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

type GoogleSheetsClient struct {
	client       *http.Client
	feedback     string
	cpSheet      string
	cpScSheet    string
	speedSheet   string
	speedScSheet string
}

func NewGoogleSheetsClient(cpSheet string, cpScSheet string, speedSheet string, speedScSheet string, feedback string) *GoogleSheetsClient {
	client := new(GoogleSheetsClient)
	client.client = &http.Client{Timeout: 30 * time.Second}
	client.cpSheet = cpSheet
	client.cpScSheet = cpScSheet
	client.speedSheet = speedSheet
	client.speedScSheet = speedScSheet
	client.feedback = feedback
	return client
}

/*
prepGoogleSheet will read the credentials in the client_secret.json in order to retrieve the JWT that
is used to fetch the spreadsheet. Once fetched, it is returned to whichever function needs it
*/
func (g GoogleSheetsClient) prepGoogleSheet(sheetId string) *spreadsheet.Sheet {
	// Create the client with the correct JWT configuration
	data, err := os.ReadFile("config/client_secret.json")
	checkError(err)
	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	checkError(err)
	client := conf.Client(context.TODO())

	// Fetch the clan points google sheet
	service := spreadsheet.NewServiceWithClient(client)
	googlesheet, err := service.FetchSpreadsheet(sheetId)
	checkError(err)
	sheet, err := googlesheet.SheetByIndex(0)
	checkError(err)

	return sheet
}

/*
InitializeCpFromSheet will take all the clan points from the CP Google Sheet and populate the
cp map for use within the bot
*/
func (g GoogleSheetsClient) InitializeCpFromSheet(ctx context.Context, cp map[string]int) {
	sheet := g.prepGoogleSheet(g.cpSheet)

	header := true

	// Set the in memory cp map with the Google sheets information
	for _, row := range sheet.Rows {
		if header {
			header = false
			continue
		}
		isPlayer := true
		player := ""
		points := 0
		for _, cell := range row {
			if isPlayer {
				player = cell.Value
			} else {
				points, _ = strconv.Atoi(cell.Value)
				break
			}
			isPlayer = false
		}
		cp[player] = points
	}
}

func (g GoogleSheetsClient) InitializeSpeedsFromSheet(ctx context.Context, speeds map[string]util.SpeedInfo) {
	sheet := g.prepGoogleSheet(g.speedSheet)

	header := true

	// Set the in memory cp map with the Google sheets information
	for _, row := range sheet.Rows {
		if header {
			header = false
			continue
		}
		bossName := ""
		speedTime := ""
		players := ""
		url := ""
		for _, cell := range row {
			if len(bossName) == 0 {
				bossName = cell.Value
				continue
			} else if len(speedTime) == 0 {
				speedTime = cell.Value
				continue
			} else if len(players) == 0 {
				players = cell.Value
				continue
			} else if len(url) == 0 {
				url = cell.Value
				continue
			}
		}

		// Convert the time string into time
		var t time.Time
		speedTimeSplit := strings.Split(speedTime, ":")

		for index, splitTime := range speedTimeSplit {
			switch index {
			case 0:
				c, _ := strconv.Atoi(splitTime)
				t = t.Add(time.Duration(c) * time.Hour)
				break
			case 1:
				c, _ := strconv.Atoi(splitTime)
				t = t.Add(time.Duration(c) * time.Minute)
				break
			case 2:
				if strings.Contains(splitTime, ".") {
					milliAndSeconds := strings.Split(splitTime, ".")
					c, _ := strconv.Atoi(milliAndSeconds[0])
					c2, _ := strconv.Atoi(milliAndSeconds[1])
					t = t.Add(time.Duration(c) * time.Second)
					t = t.Add(time.Duration(c2) * time.Millisecond)
				} else {
					c, _ := strconv.Atoi(splitTime)
					t = t.Add(time.Duration(c) * time.Second)
				}
				break
			}
		}

		speeds[bossName] = util.SpeedInfo{Time: t, PlayersInvolved: players, URL: url}
	}
}

func (g GoogleSheetsClient) InitializeFeedbackFromSheet(ctx context.Context, feedback map[string]string) {
	sheet := g.prepGoogleSheet(g.feedback)

	header := true

	// Set the in memory cp map with the Google sheets information
	for _, row := range sheet.Rows {
		if header {
			header = false
			continue
		}
		isPlayer := true
		player := ""
		channel := ""
		for _, cell := range row {
			if isPlayer {
				player = cell.Value
			} else {
				channel = strings.Replace(cell.Value, "ponies", "", -1)
				break
			}
			isPlayer = false
		}
		feedback[player] = channel
	}
}

/*
UpdateFeedbackChannel will take the feedback map that was being locally updated and save it to the Google Sheets
*/
func (g GoogleSheetsClient) UpdateFeedbackChannel(ctx context.Context, feedback map[string]string) {
	sheet := g.prepGoogleSheet(g.feedback)

	// Delete all the values in the sheet before proceeding with the insertion of clan points
	// We are deleting as this is an easier way of ensuring deleted people are removed from the
	// sheets without adding additional logic

	// Update the Google sheets information with the in memory cp map
	row := 1

	for user, channelId := range feedback {
		sheet.Update(row, 0, user)
		sheet.Update(row, 1, "ponies"+channelId)
		row++
	}

	// Make sure call Synchronize to reflect the changes
	err := sheet.Synchronize()
	checkError(err)
}

/*
UpdateCpSheet will take the cp map that was being locally updated and save it to the
CP Google Sheets
*/
func (g GoogleSheetsClient) UpdateCpSheet(ctx context.Context, cp map[string]int) {
	sheet := g.prepGoogleSheet(g.cpSheet)

	// Delete all the values in the sheet before proceeding with the insertion of clan points
	// We are deleting as this is an easier way of ensuring deleted people are removed from the
	// sheets without adding additional logic

	// Update the Google sheets information with the in memory cp map
	row := 1

	// Sort based on number of clan points
	keys := make([]string, 0, len(cp))

	for key := range cp {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return cp[keys[i]] > cp[keys[j]]
	})

	for _, playerName := range keys {
		sheet.Update(row, 0, playerName)
		sheet.Update(row, 1, strconv.Itoa(cp[playerName]))
		row++
	}

	// Make sure call Synchronize to reflect the changes
	err := sheet.Synchronize()
	checkError(err)
}

/*
UpdateCpScreenshotsSheet will take all the cpscreenshots map and store the imgur link along with the
people who got that item in the Google Sheet
*/
func (g GoogleSheetsClient) UpdateCpScreenshotsSheet(ctx context.Context, cpscreenshots map[string]string) {
	// If no screenshots need to be uploaded, skip
	if len(cpscreenshots) == 0 {
		return
	}
	sheet := g.prepGoogleSheet(g.cpScSheet)

	// Append new rows into the sheets
	startingRow := len(sheet.Rows)

	for imgurUrl, players := range cpscreenshots {
		sheet.Update(startingRow, 0, imgurUrl)
		sheet.Update(startingRow, 1, players)
		startingRow++
	}

	// Make sure call Synchronize to reflect the changes
	err := sheet.Synchronize()
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func (g GoogleSheetsClient) UpdateSpeedSheet(ctx context.Context, speed map[string]util.SpeedInfo) {
	// If no screenshots need to be uploaded, skip
	if len(speed) == 0 {
		return
	}
	sheet := g.prepGoogleSheet(g.speedSheet)

	// Overwrite the rows
	startingRow := 1

	// To ensure order, we will make an array that will be used to get the key/value from the map
	order := []string{"TzHaar Fight Cave",
		"Inferno",
		"Alchemical Hydra",
		"Fragment of Seren",
		"Galvek",
		"The Gauntlet",
		"The Corrupted Gauntlet",
		"Grotesque Guardians",
		"Hespori",
		"Phantom Muspah",
		"Nex",
		"The Nightmare Solo",
		"The Nightmare 2",
		"The Nightmare 3",
		"The Nightmare 4",
		"The Nightmare 5",
		"The Nightmare 6+",
		"Phosani's Nightmare",
		"Vorkath",
		"Zulrah",
		"Tempoross",
		"Chambers of Xeric Solo",
		"Chambers of Xeric 2",
		"Chambers of Xeric 3",
		"Chambers of Xeric 4",
		"Chambers of Xeric 5",
		"Chambers of Xeric 6",
		"Chambers of Xeric 7",
		"Chambers of Xeric 8",
		"Chambers of Xeric 9",
		"Chambers of Xeric 10",
		"Chambers of Xeric 11-15",
		"Chambers of Xeric 16-23",
		"Chambers of Xeric 24+",
		"Chambers of Xeric - Challenge mode Solo",
		"Chambers of Xeric - Challenge mode 2",
		"Chambers of Xeric - Challenge mode 3",
		"Chambers of Xeric - Challenge mode 4",
		"Chambers of Xeric - Challenge mode 5",
		"Chambers of Xeric - Challenge mode 6",
		"Chambers of Xeric - Challenge mode 7",
		"Chambers of Xeric - Challenge mode 8",
		"Chambers of Xeric - Challenge mode 9",
		"Chambers of Xeric - Challenge mode 10",
		"Theatre of Blood - Entry Room 1",
		"Theatre of Blood - Entry Room 2",
		"Theatre of Blood - Entry Room 3",
		"Theatre of Blood - Entry Room 4",
		"Theatre of Blood - Entry Overall 1",
		"Theatre of Blood - Entry Overall 2",
		"Theatre of Blood - Entry Overall 3",
		"Theatre of Blood - Entry Overall 4",
		"Theatre of Blood Room 2",
		"Theatre of Blood Room 3",
		"Theatre of Blood Room 4",
		"Theatre of Blood Room 5",
		"Theatre of Blood Overall 2",
		"Theatre of Blood Overall 3",
		"Theatre of Blood Overall 4",
		"Theatre of Blood Overall 5",
		"Theatre of Blood - Hard Room 3",
		"Theatre of Blood - Hard Room 4",
		"Theatre of Blood - Hard Room 5",
		"Theatre of Blood - Hard Overall 3",
		"Theatre of Blood - Hard Overall 4",
		"Theatre of Blood - Hard Overall 5",
		"Tombs of Amascut - Entry Room Solo",
		"Tombs of Amascut - Entry Room 2",
		"Tombs of Amascut - Entry Room 4",
		"Tombs of Amascut - Entry Room 6",
		"Tombs of Amascut - Entry Overall Solo",
		"Tombs of Amascut - Entry Overall 2",
		"Tombs of Amascut - Entry Overall 4",
		"Tombs of Amascut - Entry Overall 6",
		"Tombs of Amascut Room Solo",
		"Tombs of Amascut Room 2",
		"Tombs of Amascut Room 3",
		"Tombs of Amascut Room 4",
		"Tombs of Amascut Room 5",
		"Tombs of Amascut Room 6",
		"Tombs of Amascut Room 7",
		"Tombs of Amascut Overall Solo",
		"Tombs of Amascut Overall 2",
		"Tombs of Amascut Overall 3",
		"Tombs of Amascut Overall 4",
		"Tombs of Amascut Overall 5",
		"Tombs of Amascut Overall 6",
		"Tombs of Amascut Overall 7",
		"Tombs of Amascut Expert Room Solo",
		"Tombs of Amascut Expert Room 2",
		"Tombs of Amascut Expert Room 3",
		"Tombs of Amascut Expert Room 4",
		"Tombs of Amascut Expert Room 5",
		"Tombs of Amascut Expert Room 6",
		"Tombs of Amascut Expert Room 7",
		"Tombs of Amascut Expert Room 8",
		"Tombs of Amascut Expert Overall Solo",
		"Tombs of Amascut Expert Overall 2",
		"Tombs of Amascut Expert Overall 3",
		"Tombs of Amascut Expert Overall 4",
		"Tombs of Amascut Expert Overall 5",
		"Tombs of Amascut Expert Overall 6",
		"Tombs of Amascut Expert Overall 7",
		"Tombs of Amascut Expert Overall 8",
		"Hallowed Sepulchre",
		"Prifddinas Agility Course",
	}

	for _, name := range order {
		sheet.Update(startingRow, 0, name)
		sheet.Update(startingRow, 1, speed[name].Time.Format("15:04:05.00"))
		sheet.Update(startingRow, 2, speed[name].PlayersInvolved)
		sheet.Update(startingRow, 3, speed[name].URL)
		startingRow++
	}

	// Make sure call Synchronize to reflect the changes
	err := sheet.Synchronize()
	checkError(err)
}

func (g GoogleSheetsClient) UpdateSpeedScreenshotsSheet(ctx context.Context, speedscreenshots map[string]util.SpeedScInfo) {
	// If no screenshots need to be uploaded, skip
	if len(speedscreenshots) == 0 {
		return
	}
	sheet := g.prepGoogleSheet(g.speedScSheet)

	// Sort the map on keys so that we can have submission screenshots saved in order
	keys := make([]string, 0, len(speedscreenshots))
	for k := range speedscreenshots {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	startingRow := len(sheet.Rows)

	// Append new rows into the sheets
	for _, k := range keys {
		sheet.Update(startingRow, 0, k)
		sheet.Update(startingRow, 1, speedscreenshots[k].URL)
		sheet.Update(startingRow, 2, speedscreenshots[k].BossName)
		sheet.Update(startingRow, 3, speedscreenshots[k].Time)
		sheet.Update(startingRow, 4, speedscreenshots[k].PlayersInvolved)
		startingRow++
	}

	// Make sure call Synchronize to reflect the changes
	err := sheet.Synchronize()
	checkError(err)
}
