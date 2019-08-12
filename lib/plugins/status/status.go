package status

import (
	"fmt"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/go-chat-bot/bot"
)

var startTime time.Time

func status(command *bot.Cmd) (msg string, err error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	msg = fmt.Sprintf("booted up %s on host `%s`\n", humanize.Time(startTime), hostname)
	return
}

func init() {
	startTime = time.Now()
	bot.RegisterCommand(
		"status",
		"returns uptime information",
		"",
		status)
}
