module github.com/justmiles/standup-bot

replace github.com/justmiles/standup-bot/lib/plugins => ./lib/plugins

replace github.com/justmiles/standup-bot/lib/slack => ./lib/slack

replace github.com/justmiles/standup-bot/lib/configs => ./lib/configs

replace github.com/justmiles/standup-bot/lib/standup => ./lib/standup

go 1.15

require (
	github.com/aws/aws-sdk-go v1.37.21
	github.com/dustin/go-humanize v1.0.0
	github.com/go-chat-bot/bot v0.0.0-20201004141219-763f9eeac7d5
	github.com/go-chat-bot/plugins v0.0.0-20201024114236-00ff43fcf77f
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/slack-go/slack v0.8.1
)
