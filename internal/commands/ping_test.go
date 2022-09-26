package commands

import (
	"fmt"
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
	gomock "github.com/golang/mock/gomock"
)

func TestServer_Ping(t *testing.T) {
	os.Setenv("COMMAND_PREFIX", "!")

	tests := []struct {
		name      string
		setupMock func(*MockSession)
	}{
		{
			name: "successfully pinged server",
			setupMock: func(session *MockSession) {
				session.EXPECT().InteractionRespond(gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
		},
		{
			name: "failed to ping server",
			setupMock: func(session *MockSession) {
				session.EXPECT().InteractionRespond(gomock.Any(), gomock.Any()).Times(1).Return(fmt.Errorf("random error"))
			},
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		session := NewMockSession(ctrl)
		tt.setupMock(session)

		ping := Ping()
		ping.Handler(session, &discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				ChannelID: "test-channel-id",
			},
		})
	}
}
