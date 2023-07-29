package commands

import (
	"fmt"
	"testing"

	"github.com/NekoFluff/internal/mocks"

	"github.com/bwmarrin/discordgo"
	gomock "github.com/golang/mock/gomock"
)

func TestServer_Ping(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*mocks.MockSession)
	}{
		{
			name: "successfully pinged server",
			setupMock: func(session *mocks.MockSession) {
				session.EXPECT().InteractionRespond(gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
		},
		{
			name: "failed to ping server",
			setupMock: func(session *mocks.MockSession) {
				session.EXPECT().InteractionRespond(gomock.Any(), gomock.Any()).Times(1).Return(fmt.Errorf("random error"))
			},
		},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		session := mocks.NewMockSession(ctrl)
		tt.setupMock(session)

		ping := Ping()
		ping.Handler(session, &discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				ChannelID: "test-channel-id",
			},
		})
	}
}
