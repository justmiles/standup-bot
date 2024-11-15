package standup

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/chrispruitt/go-slackbot/lib/bot"
	"github.com/justmiles/standup-bot/lib/types"
	"github.com/lucasb-eyer/go-colorful"
	logger "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

var header = "Asynchronous Standups! Less time in meetings means more time getting things done. Keep the channel clean by using threads! To manage times, participants, etc, just type the `/standup` command."

func RegisterStandup(settings types.StandupSettings) error {

	logger.Infof("\n\nRegistering Standup \n %v\n\n", settings)

	err := bot.RegisterPeriodicScript(bot.PeriodicScript{
		Name:     fmt.Sprintf("standup-solicit-%s", settings.ChannelID),
		CronSpec: settings.SolicitCronSpec,
		Function: getSolicitStandupFunc(settings),
	})

	if err != nil {
		return err
	}

	bot.RegisterPeriodicScript(bot.PeriodicScript{
		Name:     fmt.Sprintf("standup-share-%s", settings.ChannelID),
		CronSpec: settings.ShareCronSpec,
		Function: getShareStandupFunc(settings),
	})

	brain.Standups[settings.ChannelID] = settings

	return nil
}

func getSolicitStandupFunc(settings types.StandupSettings) func() {
	return func() {
		// declare settings within function scope
		settings := settings

		if settings.Freeze {
			return
		}
		for _, userId := range settings.Participants {
			bot.PostMessage(userId, settings.SolicitMsg)
		}
	}
}

func getShareStandupFunc(settings types.StandupSettings) func() {
	return func() {
		// declare settings within function scope
		settings := settings

		if settings.Freeze {
			return
		}

		bot.PostMessage(settings.ChannelID, header)

		users, err := bot.SlackClient.GetUsersInfo(settings.Participants...)
		if err != nil {
			logger.Error("Error getting users: ", err)
		}

		conversations, err := getConversations()
		if err != nil {
			logger.Error("Error getting conversations: ", err)
		}

		// seed our random colors and randomize share order
		rand.Seed(time.Now().UTC().UnixNano())
		rand.Shuffle(len(*users), func(i, j int) { (*users)[i], (*users)[j] = (*users)[j], (*users)[i] })

		for _, user := range *users {

			if channel, ok := conversations[user.ID]; ok {
				notes, err := getStandupNote(channel.ID, settings.SolicitMsg, user.ID)
				if err != nil {
					logger.Error(err)
					continue
				}

				if len(notes) == 0 {
					if settings.Shame {
						notes = append(notes, fmt.Sprintf(":poop: %s has no standup notes", strings.Split(user.RealName, " ")[0]))
					} else {
						continue
					}
				}

				attachments := []slack.Attachment{
					{
						Title: fmt.Sprintf("%s standup notes:", user.RealName),
						Text:  fmt.Sprintf(strings.Join(reverse(notes), "\n")),
						Color: colorful.FastHappyColor().Hex(),
					},
				}

				_, _, err = bot.SlackClient.PostMessage(
					settings.ChannelID,
					slack.MsgOptionText("", false),
					slack.MsgOptionEnableLinkUnfurl(),
					slack.MsgOptionAttachments(attachments...))

				if err != nil {
					logger.Errorf("Error posting standup report: %s\n", err)
				}

			} else {
				logger.Errorf("channel %s is orphaned", channel.ID)
			}
		}
	}
}

func getConversations() (map[string]slack.Channel, error) {
	conversations := make(map[string]slack.Channel)
	r, _, err := bot.SlackClient.GetConversations(&slack.GetConversationsParameters{
		ExcludeArchived: true,
		Types:           []string{"im"},
	})
	if err != nil {
		return nil, err
	}
	for _, channel := range r {
		conversations[channel.User] = channel
	}
	return conversations, nil
}

func getStandupNote(channelID string, solicitStandupMessage string, userID string) ([]string, error) {
	conversationHistory, err := bot.SlackClient.GetConversationHistory(&slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		// Oldest:    fmt.Sprintf("%d.000001", (now.Add(time.Hour * time.Duration(-1))).Unix()),
	})
	if err != nil {
		return nil, fmt.Errorf("error getting standup info from channel %s: %s", channelID, err)
	}

	var txt []string

	for _, m := range conversationHistory.Messages {
		if m.Text == solicitStandupMessage {

			threadReplies, _, _, _ := bot.SlackClient.GetConversationReplies(&slack.GetConversationRepliesParameters{
				ChannelID: channelID,
				Timestamp: m.ThreadTimestamp,
			})

			userThreadReplies := []string{}
			for _, t := range threadReplies {
				if t.User == userID {
					txt = append(userThreadReplies, t.Text)
				}
			}

			// If user replied in thread, then use only thread message for standup notes
			if len(userThreadReplies) > 0 {
				txt = userThreadReplies
			}
			break
		}

		// filter out any other standup solicitations for users in multiple standups
		if m.User == userID {
			txt = append(txt, m.Text)
		}
	}

	return txt, nil
}

func reverse(ss []string) []string {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
	return ss
}
