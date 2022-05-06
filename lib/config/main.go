package config

import (
	"fmt"
	"os"
	"strconv"
)

var (
	SlackBotToken string
	SlackAppToken string
	Debug         bool
	S3BrainBucket string
	S3BrainKey    string
)

func init() {
	SlackBotToken = getenv("SLACK_BOT_TOKEN", "", true).(string)
	SlackAppToken = getenv("SLACK_APP_TOKEN", "", true).(string)
	S3BrainBucket = getenv("S3_BRAIN_BUCKET", "", false).(string)
	S3BrainKey = getenv("S3_BRAIN_KEY", "", false).(string)
	Debug = getenv("DEBUG", false, false).(bool)
}

func getenv(key string, fallback interface{}, required bool) interface{} {
	value := os.Getenv(key)
	if len(value) == 0 {
		if required {
			panic(fmt.Sprintf("Missing required environment variable: '%s'", key))
		}
		return fallback
	}

	switch fallback.(type) {
	case string:
		v := os.Getenv(key)
		if len(value) == 0 {
			return fallback
		}
		return v
	case int:
		s := os.Getenv(key)
		v, err := strconv.Atoi(s)
		if err != nil {
			return fallback
		}
		return v

	case bool:
		s := os.Getenv(key)
		v, err := strconv.ParseBool(s)
		if err != nil {
			return fallback
		}
		return v
	default:
		return value
	}
}
