package mock

import (
	"net/http"
	"time"
)

type GoogleSheetsClientMock struct {
	client *http.Client
}

func NewGoogleSheetsClient() *GoogleSheetsClientMock {
	client := new(GoogleSheetsClientMock)
	client.client = &http.Client{Timeout: 30 * time.Second}
	return client
}

func (g *GoogleSheetsClientMock) InitializeSubmissionsFromSheet(submissions map[string]int) {
	return
}

func (g *GoogleSheetsClientMock) UpdateCpSheet(submissions map[string]int) {
	return
}

func (g *GoogleSheetsClientMock) UpdateCpScreenshotsSheet(cpscreenshots map[string]string) {
	return
}
