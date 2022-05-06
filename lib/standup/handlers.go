package standup

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/justmiles/standup-bot/lib/types"
	"github.com/justmiles/standup-bot/lib/views"
	logger "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func handleSlashCommand(cmd slack.SlashCommand) {
	switch cmd.Text {
	case "":
		openSettingsModal(cmd)
	case "solicit":
		// mock standup settings if none exist with user issuing the command else use the channels settings
		settings := getStandupSettings(cmd.ChannelID, cmd.ChannelName)
		if len(settings.Participants) == 0 {
			settings.Participants = []string{cmd.UserID}
			settings.Shame = true
		}
		getSolicitStandupFunc(settings)()
	case "share":
		// mock standup settings if none exist with user issuing the command else use the channels settings
		settings := getStandupSettings(cmd.ChannelID, cmd.ChannelName)
		if len(settings.Participants) == 0 {
			settings.Participants = []string{cmd.UserID}
			settings.Shame = true
		}
		getShareStandupFunc(settings)()
	default:
		socketMode.Debugf("Ignored: %v", cmd)
	}
}

func handleSettingsModalActions(payload slack.InteractionCallback) {

	for _, action := range payload.ActionCallback.BlockActions {
		switch action.ActionID {
		case views.ModalStandupChannelActionId:
			handleSettingsModalSelectChannelAction(payload)
		default:
			socketMode.Debugf("Ignored: settings modal action with ActionID: %v", payload)
		}
	}
}

func openSettingsModal(cmd slack.SlashCommand) {

	settings := getStandupSettings(cmd.ChannelID, cmd.ChannelName)

	resp, err := webApi.OpenView(cmd.TriggerID, views.GetSettingsModal(settings))
	if err != nil {
		log.Printf("Failed to opemn a modal: %v", err)
	}
	socketMode.Debugf("views.open response: %v", resp)
}

func submitSettingsModal(payload slack.InteractionCallback) {
	socketMode.Debugf("submit settings modal payload.View.State.Values: %v", payload.View.State.Values)

	channelID := payload.View.State.Values[views.ModalStandupChannelBlockId][views.ModalStandupChannelActionId].SelectedChannel
	participants := payload.View.State.Values[views.ModalParticipantsBlockId][views.ModalParticipantsActionId].SelectedUsers
	solicitTime := strings.Split(payload.View.State.Values[views.ModalSolicitTimeBlockId][views.ModalSolicitTimeActionId].SelectedTime, ":")
	shareTime := strings.Split(payload.View.State.Values[views.ModalShareTimeBlockId][views.ModalShareTimeActionId].SelectedTime, ":")
	meetingDays := payload.View.State.Values[views.ModalMeetingDaysBlockId][views.ModalMeetingDaysActionId].SelectedOptions
	shame, _ := strconv.ParseBool(payload.View.State.Values[views.ModalShameBlockId][views.ModalShameActionId].SelectedOption.Value)
	solicitMsg := payload.View.State.Values[views.ModalSolicitMessgaeBlockId][views.ModalSolicitMessageActionId].Value

	settings := types.StandupSettings{
		ChannelID:       channelID,
		SolicitCronSpec: fmt.Sprintf("%s %s * * %s", solicitTime[1], solicitTime[0], meetingDaysToCron(meetingDays)),
		ShareCronSpec:   fmt.Sprintf("%s %s * * %s", shareTime[1], shareTime[0], meetingDaysToCron(meetingDays)),
		SolicitMsg:      solicitMsg,
		Shame:           shame,
		Participants:    participants,
	}

	socketMode.Debugf("standup submission settings: %v", settings)

	err := RegisterStandup(settings)
	if err != nil {
		logger.Errorf("Error Registering Standup: %v", err)
		return
	}

	// Save brain after successful submission
	err = brain.writeToS3()
	if err != nil {
		logger.Errorf("Error to writing to s3: %v", err)
	}
}

func handleSettingsModalSelectChannelAction(payload slack.InteractionCallback) {
	socketMode.Debugf("select channel: %v", payload)

	channelID := payload.View.State.Values[views.ModalStandupChannelBlockId][views.ModalStandupChannelActionId].SelectedChannel
	channelInfo, _ := webApi.GetConversationInfo(channelID, false)

	socketMode.Debugf("select channel action ID: %v", channelID)
	socketMode.Debugf("selected name: %v", channelInfo.Name)

	settings := getStandupSettings(channelID, channelInfo.Name)

	// Only way to dynamically update view form values
	// 	1. Update the view without no blocks
	//  2. Update the view with new blocks using the 'initial_value' config

	// Update view by removing blocks - no other way to dynamically change view values.
	resp, err := webApi.UpdateView(views.GetSettingsModal(types.StandupSettings{}), payload.View.ExternalID, payload.View.Hash, payload.View.ID)
	if err != nil {
		log.Printf("Failed to update a modal: %v", err)
	}
	socketMode.Debugf("views.update response: %v", resp)

	// Update view with new initial values
	resp, err = webApi.UpdateView(views.GetSettingsModal(settings), resp.ExternalID, resp.Hash, resp.ID)
	if err != nil {
		log.Printf("Failed to opemn a modal: %v", err)
	}
	socketMode.Debugf("views.update response: %v", resp)

}

func meetingDaysToCron(meetingDays []slack.OptionBlockObject) string {
	values := []string{}
	for _, day := range meetingDays {
		values = append(values, day.Value)
	}
	return strings.Join(values, ",")
}
