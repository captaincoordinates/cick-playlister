package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/captaincoordinates/cick-playlister/internal/handler"
	"github.com/captaincoordinates/cick-playlister/internal/handler/spotify"

	"github.com/gorilla/mux"
)

func ConfigureRouter() *mux.Router {
	router := mux.NewRouter()
	for _, eachHandler := range []handler.TrackInfoHandler{
		spotify.NewSpotifyHandler(),
	} {
		router.HandleFunc(
			fmt.Sprintf(
				"/%s/%s",
				eachHandler.Identifier(),
				eachHandler.PathPattern(),
			),
			createHandlerClosure(eachHandler),
		)
	}
	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("./swagger/"))))
	router.HandleFunc("/openapi.yml", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "openapi.yml")
	})
	router.HandleFunc("/healthz", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
	})
	// router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
	// 	http.Redirect(writer, request, "/docs/", http.StatusMovedPermanently)
	// })
	return router
}

func createHandlerClosure(handler handler.TrackInfoHandler) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		result, err := handler.Content(request)
		if err != nil {
			// add more nuanced error handling once likely error conditions are known
			http.Error(
				writer,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
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
