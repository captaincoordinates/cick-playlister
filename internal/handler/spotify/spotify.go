package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/captaincoordinates/cick-playlister/internal/handler"
	"github.com/gorilla/mux"
)

const playlistIdentifierParam = "playlistIdentifier"

type SpotifyTrackData struct {
	Artists []struct {
		Name string `json:"name"`
	} `json:"artists"`
	Name  string `json:"name"`
	Album struct {
		Name                 string `json:"name"`
		ReleaseDate          string `json:"release_date"`
		ReleaseDatePrecision string `json:"release_date_precision"`
		AlbumType            string `json:"album_type"`
	} `json:"album"`
}

type SpotifyPlaylistData struct {
	Next  string `json:"next"`
	Items []struct {
		Track SpotifyTrackData `json:"track"`
	} `json:"items"`
}

type SpotifyHandler struct {
	token                string
	tokenExpiryTimeMilli int64
	pathParamsProvider   func(*http.Request) map[string]string
	newReleaseDays       uint
}

func NewSpotifyHandler(
	newReleaseDays uint,
) *SpotifyHandler {
	return &SpotifyHandler{
		token:                "",
		tokenExpiryTimeMilli: 0,
		pathParamsProvider:   mux.Vars,
		newReleaseDays:       newReleaseDays,
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
	nextUrl := fmt.Sprintf(
		"https://api.spotify.com/v1/playlists/%s/tracks?fields=%s",
		playlistId,
		"next,items(track(name,artists(name),album(name,album_type,release_date,release_date_precision))",
	)
	trackInfos := make([]handler.TrackInfo, 0)
	for nextUrl != "" {
		req, err := http.NewRequest(http.MethodGet, nextUrl, nil)
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
		var data SpotifyPlaylistData
		err = json.Unmarshal(body, &data)
		if err != nil {
			return handler.EmptyPlaylistInfo, err
		}
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
					entry.Track.Album.AlbumType == "single",
					spotifyHandler.trackIsNew(entry.Track),
				),
			)
		}
		nextUrl = data.Next
	}
	return handler.NewPlaylistInfo(trackInfos, playlistParamValue), nil
}

func (spotifyHandler *SpotifyHandler) trackIsNew(trackInfo SpotifyTrackData) bool {
	dateLocation, _ := time.LoadLocation("UTC")
	var date time.Time
	var err error
	switch trackInfo.Album.ReleaseDatePrecision {
	case "day":
		date, err = time.ParseInLocation(
			time.DateOnly,
			trackInfo.Album.ReleaseDate,
			dateLocation,
		)
	case "month":
		date, err = time.ParseInLocation(
			time.DateOnly,
			fmt.Sprintf(
				"%s-01",
				trackInfo.Album.ReleaseDate,
			),
			dateLocation,
		)
	case "year":
		date, err = time.ParseInLocation(
			time.DateOnly,
			fmt.Sprintf(
				"%s-01-01",
				trackInfo.Album.ReleaseDate,
			),
			dateLocation,
		)
	}
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return time.Now().UTC().AddDate(0, 0, -int(spotifyHandler.newReleaseDays)).Before(date)
}
