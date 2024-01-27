package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"osrs-disc-bot/util"
	"time"

	"github.com/gemalto/flume"
)

type ImgurClient struct {
	client     *http.Client
	uploadURL  string
	refreshURL string
}

func NewImgurClient() *ImgurClient {
	client := new(ImgurClient)
	client.client = &http.Client{Timeout: 30 * time.Second}
	client.uploadURL = "https://api.imgur.com/3/image"
	client.refreshURL = "https://api.imgur.com/oauth2/token"
	return client
}

// Upload will upload the image provided into imgur and return back the imgur url
func (i ImgurClient) Upload(ctx context.Context, AccessToken string, image io.Reader) string {
	logger := flume.FromContext(ctx)

	var buf = new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	part, _ := writer.CreateFormFile("image", "dont care about name")
	io.Copy(part, image)

	writer.Close()
	req, _ := http.NewRequest(http.MethodPost, i.uploadURL, buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+AccessToken)

	res, _ := i.client.Do(req)
	defer res.Body.Close()

	// Read the whole body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error("Error reading body: ", err)
	}

	dec := json.NewDecoder(bytes.NewReader(body))
	var img util.ImageInfoDataWrapper
	if err = dec.Decode(&img); err != nil {
		logger.Error("Error decoding body: ", err)
	}

	return img.Ii.Link
}

/*
GetNewAccessToken will take the imgur refresh token along with client information to respond with
the access token required to make cp to the imgur API
*/
func (i ImgurClient) GetNewAccessToken(ctx context.Context, RefreshToken string, ClientID string, ClientSecret string) (string, error) {
	logger := flume.FromContext(ctx)
	rawBody, err := json.Marshal(
		util.GenerateAccessTokenRequest{
			RefreshToken: RefreshToken,
			ClientID:     ClientID,
			ClientSecret: ClientSecret,
			GrantType:    "refresh_token",
		})
	logger.Info("Initiating retrieval of new imgur access token.")

	req, err := http.NewRequest(http.MethodPost, i.refreshURL, bytes.NewBuffer(rawBody))
	if err != nil {
		logger.Error("Error while creating access token request: ", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := i.client.Do(req)
	if err != nil {
		logger.Error("Error while executing access token API call: ", err.Error())
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Error while reading access token response: ", err.Error())
	}
	defer resp.Body.Close()

	response := util.GenerateAccessTokenResponse{}
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	if err = decoder.Decode(&response); err != nil {
		return "", errors.New("failed to generate token - retry")
	}

	return response.AccessToken, nil
}
