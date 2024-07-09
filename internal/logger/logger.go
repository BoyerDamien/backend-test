package logger

import (
	"os"
	"time"

	charmLog "github.com/charmbracelet/log"
)

var Logger = charmLog.NewWithOptions(os.Stderr, charmLog.Options{
	Formatter:       charmLog.TextFormatter,
	ReportCaller:    true,
	ReportTimestamp: true,
	TimeFormat:      time.Kitchen,
	Prefix:          "🧑‍💻 backend-test",
	Level:           charmLog.DebugLevel,
})
