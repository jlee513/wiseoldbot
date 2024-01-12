package mock

import (
	"context"
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

func (g *GoogleSheetsClientMock) InitializeCpFromSheet(ctx context.Context, cp map[string]int) {
	return
}

func (g *GoogleSheetsClientMock) UpdateCpSheet(ctx context.Context, cp map[string]int) {
	return
}

func (g *GoogleSheetsClientMock) UpdateCpScreenshotsSheet(ctx context.Context, cpscreenshots map[string]string) {
	return
}

func (g *GoogleSheetsClientMock) UpdateSpeedScreenshotsSheet(ctx context.Context, speedscreenshots map[string]string) {
	return
}
