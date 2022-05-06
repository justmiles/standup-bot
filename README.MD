**Description**

This is a simple standup bot for slack.

**Setup**

1. Create standup bot with `app-manifest.yaml` file provided in the root of this directory
1. Get bot and bot user tokens and set environment variables defined in `lib/config/main.go`
1. Set up s3 bucket and provide AWS Creds if you would like standup settings to persists between restarts

**Docker**

```bash
# quickstart
docker run \
    -e SLACK_BOT_TOKEN=xoxb-blahblah \
    -e SLACK_APP_TOKEN=xapp-blahblah \
    -t justmiles/standup-bot

# with persistant brain
docker run \
    -e SLACK_BOT_TOKEN=xoxb-blahblah \
    -e SLACK_APP_TOKEN=xapp-blahblah \
    -e S3_BRAIN_BUCKET=some_s3_bucket \
    -e S3_BRAIN_KEY=brain.json \
    -e AWS_REGION \
	-e AWS_ACCESS_KEY_ID \
	-e AWS_SECRET_ACCESS_KEY \
	-e AWS_SESSION_TOKEN \
    -t justmiles/standup-bot
```

**Configure your bot**

In slack just type `/standup`

**Test your bot**

In a channel type the below commands - this will create a mock standup with you as the sole participant.
`/standup solicit`
`/standup share`

**Users in multiple standups Scenarios**

If User recieves a standup solicitation from multiple channel configurations at the same time, then:
- If User replies as normal, then both standups will receive the same notes.
- If user replies in thread, then only the thread response will be used for the standup notes.

**Roadmap**

- Provide terraform module for quick setup in fargate
- Update readme with a "how to"
