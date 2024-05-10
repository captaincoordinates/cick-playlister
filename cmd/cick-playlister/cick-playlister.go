package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/captaincoordinates/cick-playlister/internal"
	"github.com/captaincoordinates/cick-playlister/internal/constants"
	"github.com/captaincoordinates/cick-playlister/internal/log"
)

func main() {
	listenPort := flag.Int("server-port", int(constants.DefaultPort), "Port the server listens on")
	logLevelStr := flag.String("log-level", "info", strings.Join(log.AllLogLevels(), " | "))
	flag.Parse()
	logger := log.NewLogger(*logLevelStr)
	logger.Debug(fmt.Sprintf("Server port %d", *listenPort))
	listenAddress := fmt.Sprintf(":%d", *listenPort)
	err := http.ListenAndServe(listenAddress, internal.ConfigureRouter())
	if err != nil {
		panic(err)
	}
}
