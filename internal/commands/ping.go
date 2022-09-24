package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func Ping() DiscordCommand {
	prefix := getPrefix()
	command := "ping"

	return NewDiscordCommand(
		prefix,
		command,
		"Is the warden watching?",
		func(s Session, m *discordgo.MessageCreate) {
			_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
			if err != nil {
				log.Println("An error occurred while pinging the server")
				log.Println(err)
			}
		},
	)
}
