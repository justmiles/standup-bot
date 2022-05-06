module github.com/justmiles/standup-bot

go 1.16

replace github.com/justmiles/standup-bot/lib/standup => ./standup

require (
	github.com/aws/aws-sdk-go v1.41.4
	github.com/chrispruitt/go-slackbot v0.3.2
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/slack-go/slack v0.9.4
)
