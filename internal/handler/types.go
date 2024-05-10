package handler

import (
	"net/http"
)

type TrackInfoHandler interface {
	Identifier() string
	PathPattern() string
	Content(*http.Request) (PlaylistInfo, error)
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

type PlaylistInfo struct {
	Tracks     []TrackInfo
	PlaylistId string
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
