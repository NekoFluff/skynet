package commands

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/NekoFluff/discord"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/mock/gomock"
)

func TestServer_Ping(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*discord.MockSession)
	}{
		{
			name: "successfully pinged server",
			setupMock: func(session *discord.MockSession) {
				session.EXPECT().InteractionRespond(gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
		},
		{
			name: "failed to ping server",
			setupMock: func(session *discord.MockSession) {
				session.EXPECT().InteractionRespond(gomock.Any(), gomock.Any()).Times(1).Return(fmt.Errorf("random error"))
			},
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		session := discord.NewMockSession(ctrl)
		tt.setupMock(session)

		ping := Ping()
		reflect.ValueOf(ping.Handler).Call([]reflect.Value{
			reflect.ValueOf(session),
			reflect.ValueOf(&discordgo.InteractionCreate{
				Interaction: &discordgo.Interaction{
					ChannelID: "test-channel-id",
				},
			}),
		})

	}
}
