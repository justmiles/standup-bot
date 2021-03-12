package standup

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	c "github.com/justmiles/standup-bot/lib/configs"

	"github.com/go-chat-bot/bot"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/slack-go/slack"
)

var api = slack.New(c.Get("SLACK_TOKEN", "", true))            // Bot User OAuth Access Token
var oauthAPI = slack.New(c.Get("OAUTH_SLACK_TOKEN", "", true)) // OAuth Access Token
var standupChannel = c.Get("SLACK_CHANNEL", "", true)
var cronSolicitStandup = c.Get("CRON_SOLICIT_STANDUP", "", true)   // example, "0 0 13 * * mon-fri", // 8AM central
var cronShareStandup = c.Get("CRON_SHARE_STANDUP", "", true)       // example, "0 0 15 * * mon-fri", // 10AM central
var shameParticipants = c.Get("SHAME_PARTICIPANTS", "true", false) // example, "0 0 15 * * mon-fri", // 10AM central
var solicitStandupMessage = c.Get("STANDUP_MESSAGE", `Hello! Could you share your standup notes with me? I'll post them in the #standup channel at 10:00AM. Just post them here anytime before then. Consider the following questions:
:point_left: What happened yesterday?
:bell: Is there anything others should be aware of?
:fire: What you're hoping to accomplish today?
:construction: Any blockers?
`, true)

var header = "Asynchronous Standups! Less time in meetings means more time getting things done. Keep the channel clean by using threads! If you would like to be removed from this list, just leave this channel"

func init() {
	userIDs, err := getUserIDsInSlackChannel(standupChannel)
	if err != nil {
		log.Fatal(err)
	}

	// Reach out to each user to gather feedback info
	bot.RegisterPeriodicCommand("solicit_standup", bot.PeriodicConfig{
		CronSpec: cronSolicitStandup,
		Channels: userIDs,
		CmdFunc:  solicitStandupNotes,
	})

	// Share the team member standup notes in the main slack channel
	bot.RegisterPeriodicCommand("start_standup", bot.PeriodicConfig{
		CronSpec: cronShareStandup,
		Channels: []string{standupChannel},
		CmdFunc:  shareStandupNotes,
	})

	// bot.RegisterPassiveCommand("standup", log)
	bot.RegisterCommand(
		"standup",
		"Write the standup results into this channel",
		"",
		shareStandupNotesRequest,
	)
	// bot.RegisterPassiveCommand("standup", log)
	bot.RegisterCommand(
		"solicit",
		"Solicit participants for their standup notes",
		"",
		solicitStandupNotesRequest,
	)

}

func solicitStandupNotes(channel string) (string, error) {
	return solicitStandupMessage, nil
}

func solicitStandupNotesRequest(command *bot.Cmd) (msg string, err error) {
	return solicitStandupNotes(command.Channel)
}

func shareStandupNotesRequest(command *bot.Cmd) (msg string, err error) {
	return shareStandupNotes(command.Channel)
}

func shareStandupNotes(standupChannel string) (string, error) {

	_, _, err := api.PostMessage(
		standupChannel,
		slack.MsgOptionText("", false),
		slack.MsgOptionAttachments(slack.Attachment{Pretext: header}))

	if err != nil {
		fmt.Printf("Error posting standup header: %s\n", err)
	}

	conversations, err := getConversations()
	if err != nil {
		fmt.Println("Error getting conversations: ", err)
	}

	users, err := getUsersInSlackChannel(standupChannel)
	if err != nil {
		fmt.Println("Error getting users: ", err)
	}

	// seed our random colors
	rand.Seed(time.Now().UTC().UnixNano())

	for userID, user := range users {
		if channel, ok := conversations[userID]; ok {
			fmt.Println(channel.ID, user.Name)
			notes, err := getStandupNote(channel.ID)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// Shame people in the channel that aren't participating
			if len(notes) == 0 {
				if shameParticipants == "true" {
					notes = append(notes, fmt.Sprintf(":poop: %s has no standup notes", user.Name))
				} else {
					continue
				}
			}

			attachments := []slack.Attachment{
				slack.Attachment{
					Title: fmt.Sprintf("%s standup notes:", user.RealName),
					Text:  fmt.Sprintf(strings.Join(reverse(notes), "\n")),
					Color: colorful.FastHappyColor().Hex(),
				}}

			_, _, err = api.PostMessage(
				standupChannel,
				slack.MsgOptionText("", false),
				slack.MsgOptionEnableLinkUnfurl(),
				slack.MsgOptionAttachments(attachments...))

			if err != nil {
				fmt.Printf("Error posting standup report: %s\n", err)
			}

		} else {
			fmt.Printf("channel %s is orphaned", channel.ID)
		}
	}
	return "", nil
}

func getConversations() (map[string]slack.Channel, error) {
	conversations := make(map[string]slack.Channel)
	r, _, err := api.GetConversations(&slack.GetConversationsParameters{
		ExcludeArchived: "true",
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

func getUserIDsInSlackChannel(channelID string) ([]string, error) {
	members, _, err := api.GetUsersInConversation(&slack.GetUsersInConversationParameters{ChannelID: channelID})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if len(members) == 0 {
		return nil, fmt.Errorf("there are no users in channel %s", channelID)
	}

	return members, err
}

func getUsersInSlackChannel(channelID string) (map[string]*slack.User, error) {
	userIDs, err := getUserIDsInSlackChannel(standupChannel)
	if err != nil {
		return nil, err
	}

	return getUsersByID(userIDs)
}

func getUserIDsByGroupName(groupName string) ([]string, error) {
	groupID, err := getGroupIDByGroupName(groupName)
	if err != nil {
		return nil, err
	}

	userIDs, err := oauthAPI.GetUserGroupMembers(groupID)
	if err != nil {
		return userIDs, fmt.Errorf("error getting users: %s", err)
	}
	return userIDs, err
}

func getUsersByID(userIDs []string) (map[string]*slack.User, error) {
	var err error
	users := make(map[string]*slack.User)

	for _, userID := range userIDs {
		users[userID], err = api.GetUserInfo(userID)
		if users[userID].IsBot {
			delete(users, userID)
		}
	}

	return users, err
}

func getUsersByGroupName(groupName string) (map[string]*slack.User, error) {
	userIDs, err := getUserIDsByGroupName(groupName)
	if err != nil {
		return nil, err
	}

	return getUsersByID(userIDs)
}

func getGroupIDByGroupName(groupName string) (groupID string, err error) {
	ug, err := oauthAPI.GetUserGroups()
	if err != nil {
		return "", err
	}
	for _, group := range ug {
		if groupName == group.Name {
			groupID = group.ID
			return groupID, nil
		}
	}
	return "", fmt.Errorf("group %s does not exist", groupName)
}

func getStandupNote(channelID string) ([]string, error) {
	conversationHistory, err := api.GetConversationHistory(&slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		// Oldest:    fmt.Sprintf("%d.000001", (now.Add(time.Hour * time.Duration(-1))).Unix()),
	})
	if err != nil {
		return nil, fmt.Errorf("error getting standup info from channel %s: %s", channelID, err)
	}

	var txt []string

	for _, m := range conversationHistory.Messages {

		if m.Text == solicitStandupMessage {
			break
		}
		txt = append(txt, m.Text)
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
