package commands

import (
	"reflect"
	"testing"

	"github.com/NekoFluff/discord"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/mock/gomock"
)

func TestServer_Reminder(t *testing.T) {
	tests := []struct {
		name       string
		setupMock  func(*discord.MockSession)
		message    string
		dateTime   string
		timezone   string
		shouldPass bool
	}{
		{
			name: "successfully set reminder with absolute time",
			setupMock: func(session *discord.MockSession) {
				session.EXPECT().InteractionRespond(gomock.Any(), discord.ContainsInteractionResponse("I'll remind you about")).Times(1).Return(nil)
			},
			message:    "Test reminder",
			dateTime:   "October 8, 2025 4PM",
			timezone:   "MST",
			shouldPass: true,
		},
		{
			name: "successfully set reminder with relative time",
			setupMock: func(session *discord.MockSession) {
				session.EXPECT().InteractionRespond(gomock.Any(), discord.ContainsInteractionResponse("I'll remind you about")).Times(1).Return(nil)
			},
			message:    "Test reminder",
			dateTime:   "1h30m",
			timezone:   "",
			shouldPass: true,
		},
		{
			name: "fail with past time",
			setupMock: func(session *discord.MockSession) {
				session.EXPECT().InteractionRespond(gomock.Any(), discord.ContainsInteractionResponse("Cannot set a reminder in the past")).Times(1).Return(nil)
			},
			message:    "Test reminder",
			dateTime:   "January 1, 2020 12PM",
			timezone:   "MST",
			shouldPass: false,
		},
		{
			name: "invalid time format",
			setupMock: func(session *discord.MockSession) {
				session.EXPECT().InteractionRespond(gomock.Any(), discord.ContainsInteractionResponse("Could not understand the reminder time format")).Times(1).Return(nil)
			},
			message:    "Test reminder",
			dateTime:   "invalid time",
			timezone:   "MST",
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			session := discord.NewMockSession(ctrl)
			tt.setupMock(session)

			reminder := Reminder()
			options := []*discordgo.ApplicationCommandInteractionDataOption{
				{
					Name:  "message",
					Value: tt.message,
				},
				{
					Name:  "datetime",
					Value: tt.dateTime,
				},
			}

			if tt.timezone != "" {
				options = append(options, &discordgo.ApplicationCommandInteractionDataOption{
					Name:  "timezone",
					Value: tt.timezone,
				})
			}

			reflect.ValueOf(reminder.Handler).Call([]reflect.Value{
				reflect.ValueOf(session),
				reflect.ValueOf(&discordgo.InteractionCreate{
					Interaction: &discordgo.Interaction{
						ChannelID: "test-channel-id",
						Type:      discordgo.InteractionApplicationCommand,
						Member: &discordgo.Member{
							User: &discordgo.User{
								ID: "test-user-id",
							},
						},
						Data: discordgo.ApplicationCommandInteractionData{
							Options: options,
						},
					},
				}),
			})
		})
	}
}
