package handler

import "fmt"

type InvalidTrackCollectionIdError struct {
	trackCollectionId string
}

func (invalidTrackCollectionIdError InvalidTrackCollectionIdError) Error() string {
	return fmt.Sprintf("Invalid track collection ID: %s", invalidTrackCollectionIdError.trackCollectionId)
}

func NewInvalidTrackCollectionIdError(trackCollectionId string) InvalidTrackCollectionIdError {
	return InvalidTrackCollectionIdError{
		trackCollectionId,
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

type TrackCollectionNotFoundError struct {
	trackCollectionId string
}

func (trackCollectionNotFoundError TrackCollectionNotFoundError) Error() string {
	return fmt.Sprintf("Track collection not found: %s", trackCollectionNotFoundError.trackCollectionId)
}

func NewTrackCollectionNotFoundError(trackCollectionId string) TrackCollectionNotFoundError {
	return TrackCollectionNotFoundError{
		trackCollectionId,
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
