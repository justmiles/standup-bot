package main

import (

	// community plugins
	_ "github.com/go-chat-bot/plugins/cmd"

	// internal plugins
	c "github.com/justmiles/standup-bot/lib/configs"
	_ "github.com/justmiles/standup-bot/lib/plugins/log"
	_ "github.com/justmiles/standup-bot/lib/plugins/standup"
	_ "github.com/justmiles/standup-bot/lib/plugins/status"
	_ "github.com/justmiles/standup-bot/lib/plugins/suicide"
	"github.com/justmiles/standup-bot/lib/slack"
)

func main() {
	slack.Run(c.Get("SLACK_TOKEN", "", true))
}
