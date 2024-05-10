package spotify

import (
	"net/http"

	"github.com/captaincoordinates/cick-playlister/internal/handler"
)

type SpotifyHandler struct{}

func NewSpotifyHandler() SpotifyHandler {
	return SpotifyHandler{}
}

func (spotifyHandler SpotifyHandler) Identifier() string {
	return "spotify"
}

func (spotifyHandler SpotifyHandler) PathPattern() string {
	return "{playlistIdentifier:.+}"
}

func (spotifyHandler SpotifyHandler) Content(request *http.Request) (trackInfos []handler.TrackInfo, err error) {
	trackInfos = append(
		trackInfos,
		handler.NewTrackInfo(
			"artist",
			"track",
			"album",
			true,
			true,
		),
	)
	return trackInfos, nil
}
