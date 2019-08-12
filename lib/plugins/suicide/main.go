package status

import (
	"os"
	"time"

	"github.com/go-chat-bot/bot"
)

func suicide(command *bot.Cmd) (msg string, err error) {

	go func() {
		time.Sleep(3 * time.Second)
		os.Exit(33)
	}()

	return "shutting down", nil
}

func init() {
	bot.RegisterCommand(
		"suicide",
		"kills the bot",
		"",
		suicide)
}
