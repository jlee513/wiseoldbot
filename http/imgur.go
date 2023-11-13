package http

import (
	"bytes"
	"encoding/json"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"osrs-disc-bot/util"
	"time"
)

type ImgurClient struct {
	client *http.Client
}

func NewImgurClient() *ImgurClient {
	client := new(ImgurClient)
	client.client = &http.Client{Timeout: 30 * time.Second}
	return client
}

func (i ImgurClient) Upload(ctx context.Context, AccessToken string, image io.Reader) string {
	var buf = new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	part, _ := writer.CreateFormFile("image", "dont care about name")
	io.Copy(part, image)

	writer.Close()
	req, _ := http.NewRequest("POST", "https://api.imgur.com/3/image", buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	req.Header.Set("Authorization", "Bearer "+AccessToken)

	client := &http.Client{}
	res, _ := client.Do(req)
	defer res.Body.Close()

	// Read the whole body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	dec := json.NewDecoder(bytes.NewReader(body))
	var img util.ImageInfoDataWrapper
	if err = dec.Decode(&img); err != nil {
		panic(err)
	}

	return img.Ii.Link
}

func (i ImgurClient) GetNewAccessToken(ctx context.Context, RefreshToken string, ClientID string, ClientSecret string) string {
	rawBody, err := json.Marshal(
		util.GenerateAccessTokenRequest{
			RefreshToken: RefreshToken,
			ClientID:     ClientID,
			ClientSecret: ClientSecret,
			GrantType:    "refresh_token",
		})

	req, err := http.NewRequest(http.MethodPost, "https://api.imgur.com/oauth2/token", bytes.NewBuffer(rawBody))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := i.client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	response := util.GenerateAccessTokenResponse{}
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	if err = decoder.Decode(&response); err != nil {
		panic(err)
	}

	return response.AccessToken
}
