package stormpath

import (
	"log"
	"os"

	"github.com/hashicorp/logutils"
)

func initLog() {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: "WARN",
		Writer:   os.Stderr,
	}

	log.SetOutput(filter)
}
