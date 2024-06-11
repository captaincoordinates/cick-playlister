package constants

type RequestType int

const (
	PlaylistRequestType RequestType = iota
	AlbumRequestType
	TrackRequestType
)

var RequestTypeNames = map[RequestType]string{
	PlaylistRequestType: "playlist",
	AlbumRequestType:    "album",
	TrackRequestType:    "track",
}

const PlaylistIdentifierParam = "playlistIdentifier"
const AlbumIdentifierParam = "albumIdentifier"
const TrackIdentifierParam = "trackIdentifier"
