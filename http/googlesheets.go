package http

import (
	"context"
	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
	"net/http"
	"os"
	"strconv"
	"time"
)

type GoogleSheetsClient struct {
	client    *http.Client
	cpSheet   string
	cpScSheet string
}

func NewGoogleSheetsClient(cpSheet string, cpScSheet string) *GoogleSheetsClient {
	client := new(GoogleSheetsClient)
	client.client = &http.Client{Timeout: 30 * time.Second}
	client.cpSheet = cpSheet
	client.cpScSheet = cpScSheet
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
InitializeSubmissionsFromSheet will take all the clan points from the CP Google Sheet and populate the
submissions map for use within the bot
*/
func (g GoogleSheetsClient) InitializeSubmissionsFromSheet(ctx context.Context, submissions map[string]int) {
	sheet := g.prepGoogleSheet(g.cpSheet)

	// Set the in memory submissions map with the Google sheets information
	for _, row := range sheet.Rows {
		isPlayer := true
		player := ""
		cp := 0
		for _, cell := range row {
			if isPlayer {
				player = cell.Value
			} else {
				cp, _ = strconv.Atoi(cell.Value)
				break
			}
			isPlayer = false
		}
		submissions[player] = cp
	}
}

/*
UpdateCpSheet will take the submissions map that was being locally updated and save it to the
CP Google Sheets
*/
func (g GoogleSheetsClient) UpdateCpSheet(ctx context.Context, submissions map[string]int) {
	sheet := g.prepGoogleSheet(g.cpSheet)

	// Delete all the values in the sheet before proceeding with the insertion of clan points
	// We are deleting as this is an easier way of ensuring deleted people are removed from the
	// sheets without adding additional logic

	// Update the Google sheets information with the in memory submissions map
	row := 0
	for player, cp := range submissions {
		sheet.Update(row, 0, player)
		sheet.Update(row, 1, strconv.Itoa(cp))
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
