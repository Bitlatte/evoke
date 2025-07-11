package logger

import (
	"os"

	"github.com/charmbracelet/log"
)

// Logger is the shared instance of the logger.
var Logger = log.NewWithOptions(os.Stderr, log.Options{
	Level: log.InfoLevel,
})
