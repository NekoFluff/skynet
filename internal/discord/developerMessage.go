package discord

import (
	"fmt"
	"strings"
	"warden/internal/utils"
)

func (bot *Bot) SendDeveloperMessage(message string) {
	developer_mode := utils.GetEnvVar("DEVELOPER_MODE")
	if developer_mode != "ON" && developer_mode != "1" {
		return
	}

	developerIds := getDeveloperIds()

	for _, developerId := range developerIds {
		ch, err := bot.session.UserChannelCreate(developerId)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = bot.session.ChannelMessageSend(ch.ID, message)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func getDeveloperIds() []string {
	ids := utils.GetEnvVar("DEVELOPER_IDS")
	if ids == "" {
		return []string{}
	}
	return strings.Split(ids, ",")
}
