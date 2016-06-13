package stormpath

import (
	"log"
	"os"

	"github.com/hashicorp/logutils"
)

//Logger library wide logger
var Logger *log.Logger
var logLevel string
var configured = false

func InitLog() {
	if !configured {
		logLevel = os.Getenv("STORMPATH_LOG_LEVEL")

		if logLevel == "" {
			logLevel = "ERROR"
		}

		filter := &logutils.LevelFilter{
			Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR", "NONE"},
			MinLevel: logutils.LogLevel(logLevel),
			Writer:   os.Stderr,
		}

		Logger = log.New(filter, "", log.Ldate|log.Ltime|log.Lshortfile)
		configured = true
	}
}
