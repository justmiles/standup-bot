# standup-bot

## Getting started

- [install go](https://golang.org/doc/install)
- export any required environment variables
- run with `make dev`

## Environment Configs

If you set `SSM_PATH` config values are pull from SSM. Values in SSM are the same as below.

- `SSM_PATH` - Pull configs from SSM using this SSM path
- `SLACK_TOKEN` - API token used to auth against Slack as a bot user
- `OAUTH_SLACK_TOKEN` - API token used to auth against Slack using oauth. Required for non-bot API scope
- `SLACK_GROUP` - Name of Slack group containing participating users
- `SLACK_CHANNEL` - Channel ID to post standup notes in
- `CRON_SOLICIT_STANDUP` - cron expression in UTC to solicit for standup notes. example, "0 0 13 * * mon-fri", // 8AM central
- `CRON_SHARE_STANDUP` - cron expression in UTC to share standup notes to the main channel. example, "0 0 13 * * mon-fri", // 8AM central
- `STANDUP_MESSAGE` - optionally customize the message sent to users when soliciting for standup notes
- `SHAME_PARTICIPANTS` - if set to "true" post a message when user doesn't particiapte in the standup

## Roadmap

- [x] instead of sourcing partificpans by Slack Group, just message everyone in the `SLACK_CHANNEL`
- [x] add property `STANDUP_MESSAGE` to customize the message sent to users when solociting standup notes
