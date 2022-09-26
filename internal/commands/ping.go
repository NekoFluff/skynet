package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func Ping() Command {
	command := "ping"

	return Command{
		Command: discordgo.ApplicationCommand{
			Name:        command,
			Description: "Is the bot online?",
		},
		Handler: func(s Session, i *discordgo.InteractionCreate) {
			err := respondToInteraction(s, i.Interaction, "Pong!")
			if err != nil {
				log.Println("An error occurred while pinging the server")
				log.Println(err)
			}
		},
	}
}
