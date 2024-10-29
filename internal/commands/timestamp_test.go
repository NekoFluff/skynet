package commands

import (
	"reflect"
	"testing"

	"github.com/NekoFluff/discord"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/mock/gomock"
)

func TestServer_Timestamp(t *testing.T) {
	tests := []struct {
		name              string
		setupMock         func(*discord.MockSession)
		dateTime          string
		expectedTimestamp string
	}{
		{
			name: "successfully converted to timestamp",
			setupMock: func(session *discord.MockSession) {
				session.EXPECT().InteractionRespond(gomock.Any(), discord.ContainsInteractionResponse("t:1728428400")).Times(1).Return(nil)
			},
			dateTime: "October 8, 2024 4PM",
		},
		{
			name: "successfully converted to timestamp",
			setupMock: func(session *discord.MockSession) {
				session.EXPECT().InteractionRespond(gomock.Any(), discord.ContainsInteractionResponse("t:1728428400")).Times(1).Return(nil)
			},
			dateTime: "Oct 8, 2024 4PM",
		},
		{
			name: "successfully converted to timestamp 3",
			setupMock: func(session *discord.MockSession) {
				session.EXPECT().InteractionRespond(gomock.Any(), discord.ContainsInteractionResponse("t:1728403200")).Times(1).Return(nil)
			},
			dateTime: "2024-10-08T16:00:00+00:00",
		},
		{
			name: "successfully converted to timestamp 4",
			setupMock: func(session *discord.MockSession) {
				session.EXPECT().InteractionRespond(gomock.Any(), discord.ContainsInteractionResponse("t:1730336400")).Times(1).Return(nil)
			},
			dateTime: "2024-10-31T01:00:00+00:00",
		},
		{
			name: "successfully converted to timestamp 5",
			setupMock: func(session *discord.MockSession) {
				session.EXPECT().InteractionRespond(gomock.Any(), discord.ContainsInteractionResponse("t:1728403200")).Times(1).Return(nil)
			},
			dateTime: "2024-10-08T16:00:00Z",
		},
		{
			name: "successfully converted to timestamp 6",
			setupMock: func(session *discord.MockSession) {
				session.EXPECT().InteractionRespond(gomock.Any(), discord.ContainsInteractionResponse("t:1728370800")).Times(1).Return(nil)
			},
			dateTime: "10/8/2024", // month/day/year
		},
		{
			name: "failed to convert to timestamp",
			setupMock: func(session *discord.MockSession) {
				session.EXPECT().InteractionRespond(gomock.Any(), discord.ContainsInteractionResponse("Could not convert the date time to a unix timestamp")).Times(1).Return(nil)
			},
			dateTime: "invalid date time",
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		session := discord.NewMockSession(ctrl)
		tt.setupMock(session)

		timestamp := Timestamp()
		reflect.ValueOf(timestamp.Handler).Call([]reflect.Value{
			reflect.ValueOf(session),
			reflect.ValueOf(&discordgo.InteractionCreate{
				Interaction: &discordgo.Interaction{
					ChannelID: "test-channel-id",
					Type:      discordgo.InteractionApplicationCommand,
					Data: discordgo.ApplicationCommandInteractionData{
						Options: []*discordgo.ApplicationCommandInteractionDataOption{
							{
								Name:  "datetime",
								Value: tt.dateTime,
							},
							{
								Name:  "timezone",
								Value: "",
							},
						},
					},
				},
			}),
		})
	}
}
