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
	listenPort := flag.Uint("server-port", constants.DefaultPort, "Port the server listens on")
	logLevelStr := flag.String("log-level", "info", strings.Join(log.AllLogLevels(), " | "))
	newReleaseDays := flag.Uint("new-release-days", constants.DefaultNewReleaseDays, "Number of days to consider a release new")
	flag.Parse()
	logger := log.NewLogger(*logLevelStr)
	logger.Debug(fmt.Sprintf("Server port %d", *listenPort))
	listenAddress := fmt.Sprintf(":%d", *listenPort)
	err := http.ListenAndServe(listenAddress, internal.ConfigureRouter(*newReleaseDays))
	if err != nil {
		panic(err)
	}
}
