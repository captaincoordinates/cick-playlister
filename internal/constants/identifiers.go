package constants

type RequestType int

const (
	PlaylistRequestType RequestType = iota
	TrackRequestType
)

var RequestTypeNames = map[RequestType]string{
	PlaylistRequestType: "playlist",
	TrackRequestType:    "track",
}

const PlaylistIdentifierParam = "playlistIdentifier"
const TrackIdentifierParam = "trackIdentifier"
