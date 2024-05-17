package handler

import "fmt"

type InvalidPlaylistIdError struct {
	playlistId string
}

func (invalidPlaylistIdError InvalidPlaylistIdError) Error() string {
	return fmt.Sprintf("Invalid playlist ID: %s", invalidPlaylistIdError.playlistId)
}

func NewInvalidPlaylistIdError(playlistId string) InvalidPlaylistIdError {
	return InvalidPlaylistIdError{
		playlistId,
	}
}

type HandlerCredentialType int

const ApplicationCredentials HandlerCredentialType = iota

type HandlerAuthenticationError struct {
	credentialType HandlerCredentialType
}

func (handlerAuthenticationError HandlerAuthenticationError) Error() string {
	switch handlerAuthenticationError.credentialType {
	case ApplicationCredentials:
		return "Failed to authenticate with application credentials"
	}
	return "Failed to authenticate, reason unknown"
}

func NewHandlerAuthenticationError(credentialType HandlerCredentialType) HandlerAuthenticationError {
	return HandlerAuthenticationError{
		credentialType,
	}
}

type PlaylistNotFoundError struct {
	playlistId string
}

func (playlistNotFoundError PlaylistNotFoundError) Error() string {
	return fmt.Sprintf("Playlist not found: %s", playlistNotFoundError.playlistId)
}

func NewPlaylistNotFoundError(playlistId string) PlaylistNotFoundError {
	return PlaylistNotFoundError{
		playlistId,
	}
}

type TrackNotFoundError struct {
	trackId string
}

func (trackNotFoundError TrackNotFoundError) Error() string {
	return fmt.Sprintf("Track not found: %s", trackNotFoundError.trackId)
}

func NewTrackNotFoundError(trackId string) TrackNotFoundError {
	return TrackNotFoundError{
		trackId,
	}
}

type InternalError struct {
	reason string
}

func (internalError InternalError) Error() string {
	return fmt.Sprintf("Internal error: %s", internalError.reason)
}

func NewInternalError(reason string) InternalError {
	return InternalError{
		reason,
	}
}

type InvalidTrackIdError struct {
	trackId string
}

func (invalidTrackIdError InvalidTrackIdError) Error() string {
	return fmt.Sprintf("Invalid track ID: %s", invalidTrackIdError.trackId)
}

func NewInvalidTrackIdError(trackId string) InvalidTrackIdError {
	return InvalidTrackIdError{
		trackId,
	}
}
