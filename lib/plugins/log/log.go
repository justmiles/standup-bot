package deploy

import (
	"fmt"

	"github.com/go-chat-bot/bot"
)

func log(command *bot.PassiveCmd) (msg string, err error) {
	fmt.Printf("[LOG] #%s (%s)\t%s (%s)\t%s\n", command.ChannelData.HumanName, command.Channel, command.User.Nick, command.User.ID, command.Raw)
	return
}

func init() {
	bot.RegisterPassiveCommand("log", log)
}
