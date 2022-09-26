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
		ch, err := bot.Session.UserChannelCreate(developerId)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = bot.Session.ChannelMessageSend(ch.ID, message)
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
