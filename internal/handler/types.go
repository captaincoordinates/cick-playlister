package handler

import (
	"net/http"
)

type TrackInfoHandler interface {
	Identifier() string
}

type TrackInfoPlaylistHandler interface {
	Playlist(*http.Request) (TrackCollectionInfo, error)
}

type TrackInfoAlbumHandler interface {
	Album(*http.Request) (TrackCollectionInfo, error)
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

type TrackCollectionInfo struct {
	Tracks       []TrackInfo `json:"tracks"`
	CollectionId string      `json:"collectionId"`
}

func NewTrackCollectionInfo(tracks []TrackInfo, collectionId string) TrackCollectionInfo {
	return TrackCollectionInfo{
		Tracks:       tracks,
		CollectionId: collectionId,
	}
}

var EmptyTrackCollectionInfo = TrackCollectionInfo{
	Tracks:       []TrackInfo{},
	CollectionId: "",
}
