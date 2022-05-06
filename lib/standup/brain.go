package standup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/justmiles/standup-bot/lib/config"
	"github.com/justmiles/standup-bot/lib/types"

	logger "github.com/sirupsen/logrus"
)

var (
	brain = Brain{Standups: make(map[string]types.StandupSettings)}
	sess  = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	s3Client         = s3.New(sess)
	saveBrainEnabled = false
)

func init() {
	saveBrainEnabled = (len(config.S3BrainBucket) > 0 && len(config.S3BrainKey) > 0)

	if saveBrainEnabled {
		err := brain.readFromS3(config.S3BrainBucket, config.S3BrainKey)
		if err != nil {
			logger.Infof("Unable to read brain from s3. Starting bot with fresh brain: %v", err)
		} else {
			logger.Infof("Brain will persist at s3://%s/%s", config.S3BrainBucket, config.S3BrainKey)
		}

		for _, standup := range brain.Standups {
			err := RegisterStandup(standup)
			if err != nil {
				logger.Errorf("Error registering standup for channel '%s' : %v", standup.ChannelID, err)
			}
		}
	} else {
		logger.Warnf("Brain not configured. Settings will not persist between restarts.")
	}
}

type Brain struct {
	Standups map[string]types.StandupSettings `json:"standups"`
}

func getStandupSettings(channelID string, channelName string) types.StandupSettings {

	channelInfo, _ := webApi.GetConversationInfo(channelID, false)

	settings := types.StandupSettings{}

	if channelInfo != nil {
		settings = *types.NewStandupSettings(channelID, channelName)
		if value, ok := brain.Standups[channelID]; ok {
			settings = value
		}
	}

	return settings
}

func (b *Brain) writeToS3() error {

	// Convert struct to json formated byte array
	p, err := json.Marshal(b)
	if err != nil {
		return err
	}

	// Push to s3
	putObjectInput := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(bytes.NewReader(p)),
		Bucket: aws.String(config.S3BrainBucket),
		Key:    aws.String(config.S3BrainKey),
	}

	_, err = s3Client.PutObject(putObjectInput)

	if err != nil {
		return err
	}

	logger.Infof("Brain saved to s3://%s/%s", config.S3BrainBucket, config.S3BrainKey)

	return nil
}

func (b *Brain) readFromS3(bucket string, key string) error {
	result, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	defer result.Body.Close()
	body1, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}
	bodyString1 := fmt.Sprintf("%s", body1)

	decoder := json.NewDecoder(strings.NewReader(bodyString1))
	err = decoder.Decode(&b)
	if err != nil {
		return err
	}

	return nil
}
