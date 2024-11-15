package views

import (
	"embed"
	"fmt"
	"log"
	"strconv"
	"strings"

	"encoding/json"

	"github.com/justmiles/standup-bot/lib/types"
	"github.com/slack-go/slack"
)

const (
	// Define Action_id as constant so we can refer to them in the controller
	SettingsModalCallBackId     = "modal_standup_settings"
	ModalMeetingDaysBlockId     = "meeting_days"
	ModalMeetingDaysActionId    = "meeting_days"
	ModalShameBlockId           = "shame"
	ModalShameActionId          = "shame"
	ModalFreezeBlockId          = "freeze"
	ModalFreezeActionId         = "freeze"
	ModalSolicitTimeBlockId     = "input_solicit_time"
	ModalSolicitTimeActionId    = "input_solicit_time"
	ModalShareTimeBlockId       = "input_share_time"
	ModalShareTimeActionId      = "input_share_time"
	ModalSolicitMessgaeBlockId  = "input_solicit_message"
	ModalSolicitMessageActionId = "input_solicit_message"
	ModalStandupChannelBlockId  = "standup_channel_select"
	ModalStandupChannelActionId = "standup_channel_select"
	ModalParticipantsBlockId    = "participants_select"
	ModalParticipantsActionId   = "participants_select"
)

//go:embed assets/*
var assets embed.FS

func GetSettingsModal(settings types.StandupSettings) slack.ModalViewRequest {

	// Any static blocks can be declared in the json file built using the slack block kit builder.
	str, err := assets.ReadFile("assets/settings-modal.json")
	if err != nil {
		log.Printf("Unable to read view `SettingsModal`: %v", err)
	}
	view := slack.ModalViewRequest{}
	json.Unmarshal([]byte(str), &view)

	checkBoxMeetingDaysOptions := map[string]string{
		"Monday":    "MON",
		"Tuesday":   "TUE",
		"Wednesday": "WED",
		"Thursday":  "THU",
		"Friday":    "FRI",
	}

	initialDays := []string{}
	solicitTimeValue := ""
	if solicitCronSpec := strings.Split(settings.SolicitCronSpec, " "); len(solicitCronSpec) >= 4 {
		initialDays = strings.Split(solicitCronSpec[4], ",")
		solicitTimeValue = fmt.Sprintf("%s:%s", solicitCronSpec[1], solicitCronSpec[0])
	}
	shareTimeValue := ""
	if shareCronSpec := strings.Split(settings.ShareCronSpec, " "); len(shareCronSpec) >= 4 {
		shareTimeValue = fmt.Sprintf("%s:%s", shareCronSpec[1], shareCronSpec[0])
	}

	dynamicBlocks := []slack.Block{
		buildChannelSelect("Standup Channel", ModalStandupChannelBlockId, ModalStandupChannelActionId, settings.ChannelID),
	}

	if settings.ChannelID != "" {
		dynamicBlocks = append(dynamicBlocks, []slack.Block{
			buildUserMultiSelect("Participants", ModalParticipantsBlockId, ModalParticipantsActionId, settings.Participants),
			buildCheckboxGroup("Meeting Days", ModalMeetingDaysBlockId, ModalMeetingDaysActionId, checkBoxMeetingDaysOptions, initialDays),
			buildRadioGroup("Shame", ModalShameBlockId, ModalShameActionId, map[string]string{"Yes": "true", "No": "false"}, strconv.FormatBool(settings.Shame)),
			buildRadioGroup("Freeze", ModalFreezeBlockId, ModalFreezeActionId, map[string]string{"Yes": "true", "No": "false"}, strconv.FormatBool(settings.Freeze)),
			buildTimePicker("Solicit (UTC)", ModalSolicitTimeBlockId, ModalSolicitTimeActionId, solicitTimeValue),
			buildTimePicker("Share (UTC)", ModalShareTimeBlockId, ModalShareTimeActionId, shareTimeValue),
			buildTextInput("Solicit Message", ModalSolicitMessgaeBlockId, ModalSolicitMessageActionId, true, "Message to prompt user to report their standup notes.", settings.SolicitMsg),
		}...)
	}

	view.Blocks.BlockSet = append(view.Blocks.BlockSet, dynamicBlocks...)

	return view
}
