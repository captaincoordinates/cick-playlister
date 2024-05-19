package internal

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/captaincoordinates/cick-playlister/internal/constants"
	"github.com/captaincoordinates/cick-playlister/internal/handler"
	"github.com/captaincoordinates/cick-playlister/internal/handler/spotify"

	"github.com/gorilla/mux"
)

//go:embed docs
var docsDirectory embed.FS

func ConfigureRouter(
	newReleaseDays uint,
) *mux.Router {
	router := mux.NewRouter()
	capabilitiesMap := make(map[string][]string)
	for _, trackInfoHandler := range []handler.TrackInfoHandler{
		spotify.NewSpotifyHandler(newReleaseDays),
	} {
		handlerCapabilities := make([]string, 0)
		if playlistHandler, ok := trackInfoHandler.(handler.TrackInfoPlaylistHandler); ok {
			router.HandleFunc(
				fmt.Sprintf(
					"/%s/%s/{%s:.+}",
					trackInfoHandler.Identifier(),
					constants.RequestTypeNames[constants.PlaylistRequestType],
					constants.PlaylistIdentifierParam,
				),
				createHandlerFunctionClosure(playlistHandler.Playlist),
			)
			handlerCapabilities = append(handlerCapabilities, constants.RequestTypeNames[constants.PlaylistRequestType])
		}
		if trackHandler, ok := trackInfoHandler.(handler.TrackInfoTrackHandler); ok {
			router.HandleFunc(
				fmt.Sprintf(
					"/%s/%s/{%s:.+}",
					trackInfoHandler.Identifier(),
					constants.RequestTypeNames[constants.TrackRequestType],
					constants.TrackIdentifierParam,
				),
				createHandlerFunctionClosure(trackHandler.Track),
			)
			handlerCapabilities = append(handlerCapabilities, constants.RequestTypeNames[constants.TrackRequestType])
		}
		capabilitiesMap[trackInfoHandler.Identifier()] = handlerCapabilities
	}
	router.HandleFunc("/capabilities", func(writer http.ResponseWriter, request *http.Request) {
		jsonResponseType(&writer)
		json.NewEncoder(writer).Encode(capabilitiesMap)
	})
	router.PathPrefix("/docs/").Handler(http.FileServer(http.FS(fs.FS(docsDirectory))))
	router.HandleFunc("/healthz", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
	})
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, "/docs/swagger/", http.StatusMovedPermanently)
	})
	return router
}

func createHandlerFunctionClosure[T any](handlerFunction func(*http.Request) (T, error)) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		result, err := handlerFunction(request)
		if err != nil {
			statusCode, message := statusCodeFromError(err)
			http.Error(
				writer,
				message,
				statusCode,
			)
			return
		}
		jsonResponseType(&writer)
		err = json.NewEncoder(writer).Encode(result)
		if err != nil {
			http.Error(
				writer,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}
	}
}

func statusCodeFromError(err error) (int, string) {
	if _, ok := err.(handler.InvalidPlaylistIdError); ok {
		return http.StatusBadRequest, err.Error()
	}
	if _, ok := err.(handler.HandlerAuthenticationError); ok {
		return http.StatusUnauthorized, err.Error()
	}
	if _, ok := err.(handler.PlaylistNotFoundError); ok {
		return http.StatusNotFound, err.Error()
	}
	if _, ok := err.(handler.TrackNotFoundError); ok {
		return http.StatusNotFound, err.Error()
	}
	if _, ok := err.(handler.InternalError); ok {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusInternalServerError, ""
}

func jsonResponseType(writer *http.ResponseWriter) {
	(*writer).Header().Set("Content-Type", "application/json")
}
