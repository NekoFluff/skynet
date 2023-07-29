package commands

import (
	"log"

	"github.com/NekoFluff/discord"
	"github.com/NekoFluff/skynet/internal/mydiscord"

	"github.com/bwmarrin/discordgo"
)

func Ping() discord.Command {
	command := "ping"

	return discord.Command{
		Command: discordgo.ApplicationCommand{
			Name:        command,
			Description: "Is the bot online?",
		},
		Handler: func(s mydiscord.Session, i *discordgo.InteractionCreate) {
			err := respondToInteraction(s, i.Interaction, "Pong!")
			if err != nil {
				log.Println("An error occurred while pinging the server")
				log.Println(err)
			}
		},
	}
}
