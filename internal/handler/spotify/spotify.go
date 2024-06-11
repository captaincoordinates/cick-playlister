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

func (spotifyHandler *SpotifyHandler) Album(request *http.Request) (albumInfo handler.TrackCollectionInfo, err error) {
	vars := spotifyHandler.pathParamsProvider(request)
	albumParamValue := vars[constants.AlbumIdentifierParam]
	token, err := spotifyHandler.getToken(spotifyHandler.clientId, spotifyHandler.clientSecret)
	if token == "" || err != nil {
		return handler.EmptyTrackCollectionInfo, handler.NewHandlerAuthenticationError(handler.ApplicationCredentials)
	}
	nextUrl := fmt.Sprintf(
		"https://api.spotify.com/v1/albums/%s",
		albumParamValue,
	)
	trackInfos := make([]handler.TrackInfo, 0)
	for nextUrl != "" {
		req, err := http.NewRequest(http.MethodGet, nextUrl, nil)
		if err != nil {
			return handler.EmptyTrackCollectionInfo, handler.NewInternalError(err.Error())
		}
		addAuthHeader(req, token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return handler.EmptyTrackCollectionInfo, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			switch resp.StatusCode {
			case http.StatusBadRequest:
				return handler.EmptyTrackCollectionInfo, handler.NewInvalidTrackCollectionIdError(albumParamValue)
			case http.StatusNotFound:
				return handler.EmptyTrackCollectionInfo, handler.NewTrackCollectionNotFoundError(albumParamValue)
			default:
				return handler.EmptyTrackCollectionInfo, fmt.Errorf("spotify API returned status: %d", resp.StatusCode)
			}
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return handler.EmptyTrackCollectionInfo, err
		}
		var data SpotifyAlbumData
		err = json.Unmarshal(body, &data)
		if err != nil {
			return handler.EmptyTrackCollectionInfo, err
		}
		for _, entry := range data.Tracks.Items {
			artistNames := make([]string, len(entry.Artists))
			for i, artist := range entry.Artists {
				artistNames[i] = artist.Name
			}
			artists := strings.Join(artistNames, ", ")
			trackInfos = append(
				trackInfos,
				handler.NewTrackInfo(
					artists,
					entry.Name,
					data.Name,
					false,
					spotifyHandler.trackIsNew(data.ReleaseDate, data.ReleaseDatePrecision),
				),
			)
		}
		nextUrl = data.Tracks.Next
	}
	return handler.NewTrackCollectionInfo(trackInfos, albumParamValue), nil
}

func (spotifyHandler *SpotifyHandler) Playlist(request *http.Request) (playlistInfo handler.TrackCollectionInfo, err error) {
	vars := spotifyHandler.pathParamsProvider(request)
	playlistParamValue := vars[constants.PlaylistIdentifierParam]
	token, err := spotifyHandler.getToken(spotifyHandler.clientId, spotifyHandler.clientSecret)
	if token == "" || err != nil {
		return handler.EmptyTrackCollectionInfo, handler.NewHandlerAuthenticationError(handler.ApplicationCredentials)
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
			return handler.EmptyTrackCollectionInfo, handler.NewInternalError(err.Error())
		}
		addAuthHeader(req, token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return handler.EmptyTrackCollectionInfo, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			switch resp.StatusCode {
			case http.StatusBadRequest:
				return handler.EmptyTrackCollectionInfo, handler.NewInvalidTrackCollectionIdError(playlistParamValue)
			case http.StatusNotFound:
				return handler.EmptyTrackCollectionInfo, handler.NewTrackCollectionNotFoundError(playlistParamValue)
			default:
				return handler.EmptyTrackCollectionInfo, fmt.Errorf("spotify API returned status: %d", resp.StatusCode)
			}
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return handler.EmptyTrackCollectionInfo, err
		}
		var data SpotifyPlaylistData
		err = json.Unmarshal(body, &data)
		if err != nil {
			return handler.EmptyTrackCollectionInfo, err
		}
		for _, entry := range data.Items {
			trackInfos = append(
				trackInfos,
				spotifyHandler.trackInfoFromSpotifyTrackData(entry.Track),
			)
		}
		nextUrl = data.Next
	}
	return handler.NewTrackCollectionInfo(trackInfos, playlistParamValue), nil
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
		spotifyHandler.trackIsNew(spotifyTrackData.Album.ReleaseDate, spotifyTrackData.Album.ReleaseDatePrecision),
	)
}

func (spotifyHandler *SpotifyHandler) trackIsNew(releaseDate string, releaseDatePrecision string) bool {
	dateLocation, _ := time.LoadLocation("UTC")
	var date time.Time
	var err error
	switch releaseDatePrecision {
	case "day":
		date, err = time.ParseInLocation(
			time.DateOnly,
			releaseDate,
			dateLocation,
		)
	case "month":
		date, err = time.ParseInLocation(
			time.DateOnly,
			fmt.Sprintf(
				"%s-01",
				releaseDate,
			),
			dateLocation,
		)
	case "year":
		date, err = time.ParseInLocation(
			time.DateOnly,
			fmt.Sprintf(
				"%s-01-01",
				releaseDate,
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
