package commands

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Pick() DiscordCommand {
	prefix := getPrefix()
	command := "pick"

	return NewDiscordCommand(
		prefix,
		command,
		fmt.Sprintf("Pick a random value (e.g. `%s%s optionA optionB optionC`)", prefix, command),
		func(s Session, m *discordgo.MessageCreate) {
			args := strings.Split(m.Content, " ")[1:]
			if len(args) == 0 {
				_, err := s.ChannelMessageSend(m.ChannelID, "No arguments provided")
				if err != nil {
					log.Println(err)
				}
				return
			}

			result := pick(args)
			_, err := s.ChannelMessageSend(m.ChannelID, result)
			if err != nil {
				log.Println(err)
			}
		},
	)
}

func pick(options []string) string {
	return options[rand.Intn(len(options))]
}
