package log

import (
	"fmt"

	"github.com/captaincoordinates/cick-playlister/internal/constants"
	"github.com/sirupsen/logrus"
)

func AllLogLevels() []string {
	all := make([]string, len(logrus.AllLevels))
	for i, level := range logrus.AllLevels {
		all[i] = level.String()
	}
	return all
}

func NewLogger(logLevelString string) logrus.FieldLogger {
	logLevel := logLevelFromString(logLevelString)
	fmt.Printf("logging at level '%s'\n", logLevel.String())
	logger := logrus.New()
	logger.SetLevel(logLevel)
	return logger
}

func logLevelFromString(logLevelString string) logrus.Level {
	logLevel, err := logrus.ParseLevel(logLevelString)
	if err != nil {
		logLevel = constants.DefaultLogLevel
		fmt.Println(err.Error())
		fmt.Printf("defaulting to '%s' log level\n", logLevel.String())
	}
	return logLevel
}
