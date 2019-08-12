module github.com/justmiles/standup-bot

replace github.com/justmiles/standup-bot/lib/plugins => ./lib/plugins

replace github.com/justmiles/standup-bot/lib/slack => ./lib/slack

replace github.com/justmiles/standup-bot/lib/configs => ./lib/configs

replace github.com/justmiles/standup-bot/lib/standup => ./lib/standup

require (
	github.com/aws/aws-sdk-go v1.19.15
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dustin/go-humanize v1.0.0
	github.com/go-chat-bot/bot v0.0.0-20181007231045-cf9880602203
	github.com/go-chat-bot/plugins v0.0.0-20181008204345-f5bd5fb31f12
	github.com/lucasb-eyer/go-colorful v1.0.2
	github.com/lusis/go-slackbot v0.0.0-20180109053408-401027ccfef5 // indirect
	github.com/lusis/slack-test v0.0.0-20190408224659-6cf59653add2 // indirect
	github.com/nlopes/slack v0.4.0
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.2.2 // indirect
	golang.org/x/net v0.0.0-20190311183353-d8887717615a // indirect
)
