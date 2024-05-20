package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/captaincoordinates/cick-playlister/internal/constants"
	"github.com/captaincoordinates/cick-playlister/internal/handler"
)

func (spotifyHandler *SpotifyHandler) Identifier() string {
	return "spotify"
}

func (spotifyHandler *SpotifyHandler) Track(request *http.Request) (trackInfo handler.TrackInfo, err error) {
	vars := spotifyHandler.pathParamsProvider(request)
	trackParamValue := vars[constants.TrackIdentifierParam]
	token, err := spotifyHandler.getToken(spotifyHandler.clientId, spotifyHandler.clientSecret)
	if token == "" || err != nil {
		return handler.EmptyTrackInfo, handler.NewHandlerAuthenticationError(handler.ApplicationCredentials)
	}
	url := fmt.Sprintf(
		"https://api.spotify.com/v1/tracks/%s",
		trackParamValue,
	)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return handler.EmptyTrackInfo, handler.NewInternalError(err.Error())
	}
	addAuthHeader(req, token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return handler.EmptyTrackInfo, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return handler.EmptyTrackInfo, handler.NewInvalidTrackIdError(trackParamValue)
		case http.StatusNotFound:
			return handler.EmptyTrackInfo, handler.NewTrackNotFoundError(trackParamValue)
		default:
			return handler.EmptyTrackInfo, fmt.Errorf("spotify API returned status: %d", resp.StatusCode)
		}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return handler.EmptyTrackInfo, err
	}
	var data SpotifyTrackData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return handler.EmptyTrackInfo, err
	}
	return spotifyHandler.trackInfoFromSpotifyTrackData(data), nil
}

func (spotifyHandler *SpotifyHandler) Playlist(request *http.Request) (playlistInfo handler.PlaylistInfo, err error) {
	vars := spotifyHandler.pathParamsProvider(request)
	playlistParamValue := vars[constants.PlaylistIdentifierParam]
	token, err := spotifyHandler.getToken(spotifyHandler.clientId, spotifyHandler.clientSecret)
	if token == "" || err != nil {
		return handler.EmptyPlaylistInfo, handler.NewHandlerAuthenticationError(handler.ApplicationCredentials)
	}
	nextUrl := fmt.Sprintf(
		"https://api.spotify.com/v1/playlists/%s/tracks?fields=%s",
		playlistParamValue,
		"next,items(track(name,artists(name),album(name,album_type,release_date,release_date_precision))",
	)
	trackInfos := make([]handler.TrackInfo, 0)
	for nextUrl != "" {
		req, err := http.NewRequest(http.MethodGet, nextUrl, nil)
		if err != nil {
			return handler.EmptyPlaylistInfo, handler.NewInternalError(err.Error())
		}
		addAuthHeader(req, token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return handler.EmptyPlaylistInfo, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			switch resp.StatusCode {
			case http.StatusBadRequest:
				return handler.EmptyPlaylistInfo, handler.NewInvalidPlaylistIdError(playlistParamValue)
			case http.StatusNotFound:
				return handler.EmptyPlaylistInfo, handler.NewPlaylistNotFoundError(playlistParamValue)
			default:
				return handler.EmptyPlaylistInfo, fmt.Errorf("spotify API returned status: %d", resp.StatusCode)
			}
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
			trackInfos = append(
				trackInfos,
				spotifyHandler.trackInfoFromSpotifyTrackData(entry.Track),
			)
		}
		nextUrl = data.Next
	}
	return handler.NewPlaylistInfo(trackInfos, playlistParamValue), nil
}

func (spotifyHandler *SpotifyHandler) trackInfoFromSpotifyTrackData(spotifyTrackData SpotifyTrackData) handler.TrackInfo {
	artistNames := make([]string, len(spotifyTrackData.Artists))
	for i, artist := range spotifyTrackData.Artists {
		artistNames[i] = artist.Name
	}
	artists := strings.Join(artistNames, ", ")
	return handler.NewTrackInfo(
		artists,
		spotifyTrackData.Name,
		spotifyTrackData.Album.Name,
		spotifyTrackData.Album.AlbumType == "single",
		spotifyHandler.trackIsNew(spotifyTrackData),
	)
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

func addAuthHeader(request *http.Request, token string) {
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
}
