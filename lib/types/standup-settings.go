package types

import "fmt"

type StandupSettings struct {
	ChannelID       string   `json:"channelID"`
	SolicitCronSpec string   `json:"solicitCronSpec"`
	ShareCronSpec   string   `json:"shareCronSpec"`
	SolicitMsg      string   `json:"solicitMsg"`
	Shame           bool     `json:"shame"`
	Freeze          bool     `json:"freeze"`
	Participants    []string `json:"participants"` // List of slack user ids
}

func NewStandupSettings(channelID string, channelName string) *StandupSettings {
	return &StandupSettings{
		ChannelID:       channelID,
		SolicitCronSpec: "00 13 * * MON,TUE,WED,THU,FRI",
		ShareCronSpec:   "00 14 * * MON,TUE,WED,THU,FRI",
		SolicitMsg: fmt.Sprintf(`Hello! Could you share your standup notes with me? I'll post them in the #%s channel at the configured time. Just post them here anytime before then. Use a thread if you attend multiple standups and want to send specific notes to them. Consider the following questions:
:point_left: What happened yesterday?
:bell: Is there anything others should be aware of?
:fire: What you're hoping to accomplish today?
:construction: Any blockers?`, channelName),
	}
}
