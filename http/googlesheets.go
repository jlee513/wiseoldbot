package http

import (
	"context"
	"github.com/gemalto/flume"
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
	cpSheet      string
	cpScSheet    string
	speedSheet   string
	speedScSheet string
	tidSheet     string
	members      string
}

func NewGoogleSheetsClient(config *util.Config) *GoogleSheetsClient {
	client := new(GoogleSheetsClient)
	client.client = &http.Client{Timeout: 30 * time.Second}
	client.cpSheet = config.SheetsCp
	client.cpScSheet = config.SheetsCpSC
	client.speedSheet = config.SheetsSpeed
	client.speedScSheet = config.SheetsSpeedSC
	client.tidSheet = config.SheetsTid
	client.members = config.SheetsMembers
	return client
}

/*
prepGoogleSheet will read the credentials in the client_secret.json in order to retrieve the JWT that
is used to fetch the spreadsheet. Once fetched, it is returned to whichever function needs it
*/
func (g *GoogleSheetsClient) prepGoogleSheet(ctx context.Context, sheetId string) *spreadsheet.Sheet {
	// Create the client with the correct JWT configuration
	data, err := os.ReadFile("config/client_secret.json")
	checkError(ctx, err)
	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	checkError(ctx, err)
	client := conf.Client(context.TODO())

	// Fetch the clan points google sheet
	service := spreadsheet.NewServiceWithClient(client)
	googlesheet, err := service.FetchSpreadsheet(sheetId)
	checkError(ctx, err)
	sheet, err := googlesheet.SheetByIndex(0)
	checkError(ctx, err)

	return sheet
}

/*
InitializeCpFromSheet will take all the clan points from the Cp Google Sheet and populate the
cp map for use within the bot
*/
func (g *GoogleSheetsClient) InitializeCpFromSheet(ctx context.Context, cp map[string]int) {
	sheet := g.prepGoogleSheet(ctx, g.cpSheet)

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

func (g *GoogleSheetsClient) InitializeSpeedsFromSheet(ctx context.Context, speeds map[string]util.SpeedInfo) {
	sheet := g.prepGoogleSheet(ctx, g.speedSheet)

	header := true

	// Set the in memory cp map with the Google sheets information
	for _, row := range sheet.Rows {
		if header {
			header = false
			continue
		}
		category := ""
		bossName := ""
		speedTime := ""
		players := ""
		url := ""
		for _, cell := range row {
			if len(category) == 0 {
				category = cell.Value
				continue
			}
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
			case 1:
				c, _ := strconv.Atoi(splitTime)
				t = t.Add(time.Duration(c) * time.Minute)
			case 2:
				if strings.Contains(splitTime, ".") {
					milliAndSeconds := strings.Split(splitTime, ".")
					c, _ := strconv.Atoi(milliAndSeconds[0])
					c2, _ := strconv.Atoi(milliAndSeconds[1])
					t = t.Add(time.Duration(c) * time.Second)
					t = t.Add(time.Duration(c2) * time.Millisecond * 10)
				} else {
					c, _ := strconv.Atoi(splitTime)
					t = t.Add(time.Duration(c) * time.Second)
				}
			}
		}

		speeds[bossName] = util.SpeedInfo{Time: t, PlayersInvolved: players, URL: url, Category: category}
	}
}

func (g *GoogleSheetsClient) InitializeTIDFromSheet(ctx context.Context) int {
	sheet := g.prepGoogleSheet(ctx, g.tidSheet)
	tid, _ := strconv.Atoi(sheet.Rows[0][0].Value)
	return tid
}

func (g *GoogleSheetsClient) InitializeMembersFromSheet(ctx context.Context, members map[string]util.MemberInfo) {
	sheet := g.prepGoogleSheet(ctx, g.members)

	header := true

	missingFeedback := make(map[string]string)
	foundFeedback := make(map[string]string)

	// Set the in memory cp map with the Google sheets information
	for _, row := range sheet.Rows {
		if header {
			header = false
			continue
		}
		player := ""
		id := 0
		discordName := ""
		feedback := ""
		main := false
		for col, cell := range row {
			switch col {
			case 0:
				player = cell.Value
			case 1:
				id, _ = strconv.Atoi(strings.Replace(cell.Value, "clan", "", -1))
			case 2:
				discordName = cell.Value
			case 3:
				feedback = strings.Replace(cell.Value, "clan", "", -1)
				if len(feedback) == 0 {
					missingFeedback[discordName] = player
				} else {
					foundFeedback[discordName] = feedback
				}
			case 4:
				main, _ = strconv.ParseBool(cell.Value)
			}
		}
		members[player] = util.MemberInfo{
			DiscordId:   id,
			DiscordName: discordName,
			Feedback:    feedback,
			Main:        main,
		}
	}

	// If there are entries in the clan members sheets that don't have feedback entries, check the ones that found the
	// entries and set it (this is for alts)
	for discordName, player := range missingFeedback {
		if _, ok := foundFeedback[discordName]; ok {
			members[player] = util.MemberInfo{
				DiscordId:   members[player].DiscordId,
				DiscordName: members[player].DiscordName,
				Feedback:    foundFeedback[discordName],
				Main:        members[player].Main,
			}
		}
	}
}
func (g *GoogleSheetsClient) UpdateMembersSheet(ctx context.Context, members map[string]util.MemberInfo) {
	sheet := g.prepGoogleSheet(ctx, g.members)

	// Delete all the values in the sheet before proceeding with the insertion of clan points
	// We are deleting as this is an easier way of ensuring deleted people are removed from the
	// sheets without adding additional logic

	// Update the Google sheets information with the in memory cp map
	row := 1

	// First, delete everything from the sheets
	err := sheet.DeleteRows(1, len(sheet.Rows))
	checkError(ctx, err)

	// Then insert in all the information from members
	for user, memberInfo := range members {
		sheet.Update(row, 0, user)
		sheet.Update(row, 1, "clan"+strconv.Itoa(memberInfo.DiscordId))
		sheet.Update(row, 2, memberInfo.DiscordName)
		sheet.Update(row, 3, "clan"+memberInfo.Feedback)
		sheet.Update(row, 4, strconv.FormatBool(memberInfo.Main))
		row++
	}

	// Make sure call Synchronize to reflect the changes
	err = sheet.Synchronize()
	checkError(ctx, err)
}

func (g *GoogleSheetsClient) UpdateTIDFromSheet(ctx context.Context, tid int) {
	sheet := g.prepGoogleSheet(ctx, g.tidSheet)
	sheet.Update(0, 0, strconv.Itoa(tid))
	// Make sure call Synchronize to reflect the changes
	err := sheet.Synchronize()
	checkError(ctx, err)
}

/*
UpdateCpSheet will take the cp map that was being locally updated and save it to the
Cp Google Sheets
*/
func (g *GoogleSheetsClient) UpdateCpSheet(ctx context.Context, cp map[string]int) {
	sheet := g.prepGoogleSheet(ctx, g.cpSheet)

	// Delete all the values in the sheet before proceeding with the insertion of clan points
	// We are deleting as this is an easier way of ensuring deleted people are removed from the
	// sheets without adding additional logic

	// Update the Google sheets information with the in memory cp map
	row := 1

	// First, delete everything from the sheets
	err := sheet.DeleteRows(1, len(sheet.Rows))
	checkError(ctx, err)

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
	err = sheet.Synchronize()
	checkError(ctx, err)
}

/*
UpdateCpScreenshotsSheet will take all the cpscreenshots map and store the imgur link along with the
people who got that item in the Google Sheet
*/
func (g *GoogleSheetsClient) UpdateCpScreenshotsSheet(ctx context.Context, cpscreenshots map[string]util.CpScInfo) {
	// If no screenshots need to be uploaded, skip
	if len(cpscreenshots) == 0 {
		return
	}
	sheet := g.prepGoogleSheet(ctx, g.cpScSheet)

	// Append new rows into the sheets
	startingRow := len(sheet.Rows)

	for submissionTime, scInfo := range cpscreenshots {
		sheet.Update(startingRow, 0, submissionTime)
		sheet.Update(startingRow, 1, scInfo.URL)
		sheet.Update(startingRow, 2, scInfo.PlayersInvolved)
		startingRow++
	}

	// Make sure call Synchronize to reflect the changes
	err := sheet.Synchronize()
	checkError(ctx, err)
}

func checkError(ctx context.Context, err error) {
	logger := flume.FromContext(ctx)
	if err != nil {
		logger.Error("Failed sheets call: " + err.Error())
	}
}

func (g *GoogleSheetsClient) UpdateSpeedSheet(ctx context.Context, speed map[string]util.SpeedInfo, speedCategory map[string][]string) {
	// If no screenshots need to be uploaded, skip
	if len(speed) == 0 {
		return
	}
	sheet := g.prepGoogleSheet(ctx, g.speedSheet)

	// Overwrite the rows
	startingRow := 1

	// First, delete everything from the sheets
	err := sheet.DeleteRows(1, len(sheet.Rows))
	checkError(ctx, err)

	for category := range util.HofSpeedCategories {
		bosses := speedCategory[category]
		sort.Strings(bosses)
		for _, boss := range bosses {
			sheet.Update(startingRow, 0, speed[boss].Category)
			sheet.Update(startingRow, 1, boss)
			sheet.Update(startingRow, 2, speed[boss].Time.Format("15:04:05.00"))
			sheet.Update(startingRow, 3, speed[boss].PlayersInvolved)
			sheet.Update(startingRow, 4, speed[boss].URL)
			startingRow++
		}
	}

	// Make sure call Synchronize to reflect the changes
	err = sheet.Synchronize()
	checkError(ctx, err)
}

func (g *GoogleSheetsClient) UpdateSpeedScreenshotsSheet(ctx context.Context, speedscreenshots map[string]util.SpeedScInfo) {
	// If no screenshots need to be uploaded, skip
	if len(speedscreenshots) == 0 {
		return
	}
	sheet := g.prepGoogleSheet(ctx, g.speedScSheet)

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
	checkError(ctx, err)
}
