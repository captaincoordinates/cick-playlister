package spotify

import (
	"net/http"

	"github.com/gorilla/mux"
)

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

type SpotifyAlbumData struct {
	Name                 string `json:"name"`
	ReleaseDate          string `json:"release_date"`
	ReleaseDatePrecision string `json:"release_date_precision"`
	Tracks               struct {
		Next  string `json:"next"`
		Items []struct {
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
			Name string `json:"name"`
		} `json:"items"`
	} `json:"tracks"`
}

type SpotifyHandler struct {
	clientId             string
	clientSecret         string
	token                string
	tokenExpiryTimeMilli int64
	pathParamsProvider   func(*http.Request) map[string]string
	newReleaseDays       uint
}

func NewSpotifyHandler(
	clientId string,
	clientSecret string,
	newReleaseDays uint,
) *SpotifyHandler {
	return &SpotifyHandler{
		clientId:             clientId,
		clientSecret:         clientSecret,
		token:                "",
		tokenExpiryTimeMilli: 0,
		pathParamsProvider:   mux.Vars,
		newReleaseDays:       newReleaseDays,
	}
}
