package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func Help(commands *[]DiscordCommand) DiscordCommand {
	prefix := getPrefix()
	command := "help"

	return NewDiscordCommand(
		prefix,
		command,
		"Display all available commands",
		func(s Session, m *discordgo.MessageCreate) {
			embedFields := []*discordgo.MessageEmbedField{}

			// Setup MessageEmbedField
			embedFields = append(embedFields, &discordgo.MessageEmbedField{
				Name:   "Setup",
				Value:  "No setup necesary!",
				Inline: false,
			})

			// Build all the commands into MessageEmbedFields
			for _, c := range *commands {
				embedField := &discordgo.MessageEmbedField{
					Name:   c.Command,
					Value:  c.Description,
					Inline: false,
				}
				embedFields = append(embedFields, embedField)
			}

			// Build the embed
			embed := &discordgo.MessageEmbed{
				Type:   discordgo.EmbedTypeRich,
				Title:  "Help Page",
				Fields: embedFields,
			}

			_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				log.Println("An error occurred while sending the help embed")
				log.Println(err)
			}
		},
	)
}
