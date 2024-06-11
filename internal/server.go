package internal

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/captaincoordinates/cick-playlister/internal/config"
	"github.com/captaincoordinates/cick-playlister/internal/constants"
	"github.com/captaincoordinates/cick-playlister/internal/handler"
	"github.com/captaincoordinates/cick-playlister/internal/handler/spotify"

	"github.com/gorilla/mux"
)

//go:embed docs
var docsDirectory embed.FS

//go:embed client/dist
var clientDirectory embed.FS

//go:embed client/assets
var assetsDirectory embed.FS

func ConfigureRouter(
	newReleaseDays uint,
) *mux.Router {
	router := mux.NewRouter()
	router.Use(corsMiddleware)
	credentialsConfig := config.NewCredentialsConfig()
	for _, trackInfoHandler := range []handler.TrackInfoHandler{
		spotify.NewSpotifyHandler(
			credentialsConfig.Spotify.ClientID,
			credentialsConfig.Spotify.ClientSecret,
			newReleaseDays,
		),
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
		if albumHandler, ok := trackInfoHandler.(handler.TrackInfoAlbumHandler); ok {
			router.HandleFunc(
				fmt.Sprintf(
					"/%s/%s/{%s:.+}",
					trackInfoHandler.Identifier(),
					constants.RequestTypeNames[constants.AlbumRequestType],
					constants.AlbumIdentifierParam,
				),
				createHandlerFunctionClosure(albumHandler.Album),
			)
			handlerCapabilities = append(handlerCapabilities, constants.RequestTypeNames[constants.TrackRequestType])
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
	}
	router.PathPrefix("/docs/").Handler(http.FileServer(http.FS(fs.FS(docsDirectory))))
	router.PathPrefix("/client/dist/").Handler(http.FileServer(http.FS(fs.FS(clientDirectory))))
	router.PathPrefix("/client/assets/").Handler(http.FileServer(http.FS(fs.FS(assetsDirectory))))
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
	if _, ok := err.(handler.InvalidTrackCollectionIdError); ok {
		return http.StatusBadRequest, err.Error()
	}
	if _, ok := err.(handler.HandlerAuthenticationError); ok {
		return http.StatusUnauthorized, err.Error()
	}
	if _, ok := err.(handler.TrackCollectionNotFoundError); ok {
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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length")
		if request.Method == "OPTIONS" {
			writer.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(writer, request)
	})
}
