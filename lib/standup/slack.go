package standup

import (
	"fmt"
	"log"
	"os"

	"github.com/justmiles/standup-bot/lib/config"
	"github.com/justmiles/standup-bot/lib/views"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

var (
	webApi     *slack.Client
	socketMode *socketmode.Client
)

func Start() {
	listner()
}

func listner() {
	webApi = slack.New(
		config.SlackBotToken,
		slack.OptionAppLevelToken(config.SlackAppToken),
		slack.OptionDebug(config.Debug),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)
	socketMode = socketmode.New(
		webApi,
		socketmode.OptionDebug(config.Debug),
		socketmode.OptionLog(log.New(os.Stdout, "sm: ", log.Lshortfile|log.LstdFlags)),
	)
	_, authTestErr := webApi.AuthTest()
	if authTestErr != nil {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN is invalid: %v\n", authTestErr)
		os.Exit(1)
	}
	go func() {
		for envelope := range socketMode.Events {
			switch envelope.Type {
			case socketmode.EventTypeSlashCommand:
				socketMode.Ack(*envelope.Request)
				cmd, _ := envelope.Data.(slack.SlashCommand)
				socketMode.Debugf("Slash command received: %+v", cmd)

				handleSlashCommand(cmd)
			case socketmode.EventTypeInteractive:
				socketMode.Ack(*envelope.Request)
				payload, _ := envelope.Data.(slack.InteractionCallback)

				switch payload.Type {
				case slack.InteractionTypeBlockActions:
					switch payload.View.CallbackID {
					case views.SettingsModalCallBackId:
						handleSettingsModalActions(payload)
					default:
						socketMode.Debugf("Ignore Submission with CallbackID: %v", payload.View.CallbackID)
					}
				case slack.InteractionTypeViewSubmission:
					switch payload.View.CallbackID {
					case views.SettingsModalCallBackId:
						submitSettingsModal(payload)
					default:
						socketMode.Debugf("Ignore Submission with CallbackID: %v", payload.View.CallbackID)
					}
				default:
					socketMode.Debugf("Ignore Payload Type: %v", payload.Type)
				}
			default:
				socketMode.Debugf("Skipped: %v", envelope.Type)
			}
		}
	}()

	socketMode.Run()
}
