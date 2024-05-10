package spotify

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/captaincoordinates/cick-playlister/internal/handler"
	"github.com/gorilla/mux"
)

const playlistIdentifierParam = "playlistIdentifier"

type SpotifyTokenData struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type PlaylistItemsResponse struct {
	Items []struct {
		Track struct {
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
			Name  string `json:"name"`
			Album struct {
				Name string `json:"name"`
			} `json:"album"`
		} `json:"track"`
	} `json:"items"`
}

type SpotifyHandler struct {
	token                string
	tokenExpiryTimeMilli int64
	pathParamsProvider   func(*http.Request) map[string]string
}

func NewSpotifyHandler() *SpotifyHandler {
	return &SpotifyHandler{
		token:                "",
		tokenExpiryTimeMilli: 0,
		pathParamsProvider:   mux.Vars,
	}
}

func (spotifyHandler *SpotifyHandler) Identifier() string {
	return "spotify"
}

func (spotifyHandler *SpotifyHandler) PathPattern() string {
	return fmt.Sprintf(`{%s:open\.spotify\.com/playlist/[a-zA-Z0-9]+}`, playlistIdentifierParam)
}

func (spotifyHandler *SpotifyHandler) Content(request *http.Request) (playlistInfo handler.PlaylistInfo, err error) {
	vars := spotifyHandler.pathParamsProvider(request)
	playlistParamValue := vars[playlistIdentifierParam]
	token, err := spotifyHandler.getToken()
	re := regexp.MustCompile(`playlist/([a-zA-Z0-9]+)$`)
	match := re.FindStringSubmatch(playlistParamValue)
	playlistId := ""
	if len(match) > 1 {
		playlistId = match[1]
	} else {
		return handler.EmptyPlaylistInfo, fmt.Errorf("invalid playlist identifier: %s", playlistParamValue)
	}
	req, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/playlists/"+playlistId+"/tracks", nil)
	if err != nil {
		return handler.EmptyPlaylistInfo, err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return handler.EmptyPlaylistInfo, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return handler.EmptyPlaylistInfo, fmt.Errorf("spotify API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return handler.EmptyPlaylistInfo, err
	}
	var data PlaylistItemsResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return handler.EmptyPlaylistInfo, err
	}

	trackInfos := make([]handler.TrackInfo, 0, len(data.Items))
	for _, entry := range data.Items {
		artistNames := make([]string, len(entry.Track.Artists))
		for i, artist := range entry.Track.Artists {
			artistNames[i] = artist.Name
		}
		artists := strings.Join(artistNames, ", ")
		trackInfos = append(
			trackInfos,
			handler.NewTrackInfo(
				artists,
				entry.Track.Name,
				entry.Track.Album.Name,
				false,
				false,
			),
		)
	}
	return handler.NewPlaylistInfo(trackInfos, playlistParamValue), nil
}

func (spotifyHandler *SpotifyHandler) getToken() (string, error) {
	nowMilli := time.Now().UTC().UnixMilli()
	if spotifyHandler.tokenExpiryTimeMilli-nowMilli <= 5000 {
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
