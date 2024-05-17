package handler

import (
	"net/http"
)

type TrackInfoHandler interface {
	Identifier() string
}

type TrackInfoPlaylistHandler interface {
	Playlist(*http.Request) (PlaylistInfo, error)
}

type TrackInfoTrackHandler interface {
	Track(*http.Request) (TrackInfo, error)
}

type TrackInfo struct {
	Artist   string `json:"artist"`
	Track    string `json:"track"`
	IsSingle bool   `json:"isSingle"`
	Album    string `json:"album"`
	IsNew    bool   `json:"isNew"`
}

func NewTrackInfo(artist, track, album string, isSingle, isNew bool) TrackInfo {
	return TrackInfo{
		Artist:   artist,
		Track:    track,
		IsSingle: isSingle,
		Album:    album,
		IsNew:    isNew,
	}
}

var EmptyTrackInfo = TrackInfo{}

type PlaylistInfo struct {
	Tracks     []TrackInfo `json:"tracks"`
	PlaylistId string      `json:"playlistId"`
}

func NewPlaylistInfo(tracks []TrackInfo, playlistId string) PlaylistInfo {
	return PlaylistInfo{
		Tracks:     tracks,
		PlaylistId: playlistId,
	}
}

var EmptyPlaylistInfo = PlaylistInfo{
	Tracks:     []TrackInfo{},
	PlaylistId: "",
}
