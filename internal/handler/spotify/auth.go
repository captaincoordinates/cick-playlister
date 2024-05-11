package spotify

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type SpotifyTokenData struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (spotifyHandler *SpotifyHandler) getToken() (string, error) {
	nowMilli := time.Now().UTC().UnixMilli()
	if spotifyHandler.tokenExpiryTimeMilli-nowMilli <= 5000 {
		/*
			Need to think about how best to provide these values.
			May want to compile them into the binary for simplicity if the binary
			will not be publicly distributed.
		*/
		clientId := os.Getenv("SPOTIFY_CLIENT_ID")
		clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
		data := url.Values{}
		data.Set("grant_type", "client_credentials")
		req, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
		if err != nil {
			return "", err
		}
		req.Header.Add(
			"Authorization",
			fmt.Sprintf(
				"Basic %s",
				base64.StdEncoding.EncodeToString(
					[]byte(
						fmt.Sprintf(
							"%s:%s",
							clientId,
							clientSecret,
						),
					),
				),
			),
		)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		var tokenResponse SpotifyTokenData
		err = json.Unmarshal(body, &tokenResponse)
		if err != nil {
			return "", err
		}
		spotifyHandler.token = tokenResponse.AccessToken
		spotifyHandler.tokenExpiryTimeMilli = nowMilli + int64(tokenResponse.ExpiresIn)*1000
	}
	return spotifyHandler.token, nil
}
