package views

import "github.com/slack-go/slack"

func buildChannelSelect(label string, blockID string, actionID string, initialChannelId string) *slack.ActionBlock {
	element := slack.NewOptionsGroupSelectBlockElement(
		slack.OptTypeChannels,
		slack.NewTextBlockObject(
			"plain_text",
			label,
			false,
			false,
		),
		actionID,
	)

	element.InitialChannel = initialChannelId

	return slack.NewActionBlock(
		blockID,
		element,
	)
}

func buildUserMultiSelect(label string, blockID, actionID string, initialParticipants []string) *slack.InputBlock {
	element := slack.NewOptionsGroupMultiSelectBlockElement(
		slack.MultiOptTypeUser,
		slack.NewTextBlockObject(
			"plain_text",
			label,
			false,
			false,
		),
		actionID,
	)

	element.InitialUsers = initialParticipants

	return slack.NewInputBlock(
		blockID,
		slack.NewTextBlockObject(
			"plain_text",
			label,
			false,
			false,
		),
		element,
	)
}

func buildCheckboxGroup(label string, blockID, actionID string, options map[string]string, initialOptions []string) *slack.InputBlock {
	element := slack.NewCheckboxGroupsBlockElement(actionID)

	for key, val := range options {
		option := slack.NewOptionBlockObject(
			val,
			slack.NewTextBlockObject(
				"plain_text",
				key,
				false,
				false,
			),
			nil,
		)
		element.Options = append(element.Options, option)
		if contains(initialOptions, val) {
			element.InitialOptions = append(element.InitialOptions, option)
		}
	}

	return slack.NewInputBlock(
		blockID,
		slack.NewTextBlockObject(
			"plain_text",
			label,
			false,
			false,
		),
		element,
	)
}

func buildRadioGroup(label string, blockID, actionID string, options map[string]string, initialOption string) *slack.InputBlock {
	element := slack.NewRadioButtonsBlockElement(actionID)

	for key, val := range options {
		option := slack.NewOptionBlockObject(
			val,
			slack.NewTextBlockObject(
				"plain_text",
				key,
				false,
				false,
			),
			nil,
		)
		element.Options = append(element.Options, option)

		if initialOption == val {
			element.InitialOption = option
		}
	}

	return slack.NewInputBlock(
		blockID,
		slack.NewTextBlockObject(
			"plain_text",
			label,
			false,
			false,
		),
		element,
	)
}

func buildTimePicker(label string, blockID, actionID string, initialValue string) *slack.InputBlock {
	element := slack.NewTimePickerBlockElement(actionID)

	element.InitialTime = initialValue

	return slack.NewInputBlock(
		blockID,
		slack.NewTextBlockObject(
			"plain_text",
			label,
			false,
			false,
		),
		element,
	)
}

func buildTextInput(label string, blockID, actionID string, multiline bool, placeholder string, initialValue string) *slack.InputBlock {
	element := slack.NewPlainTextInputBlockElement(
		slack.NewTextBlockObject(
			"plain_text",
			placeholder,
			false,
			false,
		),
		actionID,
	)

	element.InitialValue = initialValue
	element.Multiline = multiline

	return slack.NewInputBlock(
		blockID,
		slack.NewTextBlockObject(
			"plain_text",
			label,
			false,
			false,
		),
		element,
	)
}

func contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
