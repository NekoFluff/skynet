package commands

import (
	"fmt"
	"log"

	"github.com/NekoFluff/discord"
	"github.com/NekoFluff/skynet/internal/utils"
	"github.com/bwmarrin/discordgo"
)

func Timestamp() discord.Command {
	command := "timestamp"

	return discord.Command{
		Command: discordgo.ApplicationCommand{
			Name:        command,
			Description: fmt.Sprintf("Translates a generic date time string to a unix timestamp (e.g. `%s October 8, 2024 4PM MST` or `2006-01-02T15:04:05Z07:00`)", command),
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "date time",
					Description: "The date and time to convert to a unix timestamp",
					Required:    true,
				},
			},
		},
		Handler: func(s discord.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			dateTime := fmt.Sprint(optionMap["date time"].Value)

			timestamp, err := utils.ConvertToUnixTimestamp(dateTime)
			if err != nil {
				err := respondToInteraction(s, i.Interaction, "Could not convert the date time to a unix timestamp")
				if err != nil {
					log.Println(err)
				}
				return
			}

			err = respondToInteraction(s, i.Interaction, fmt.Sprintf("<t:%d:F> -> `<t:%d:F>`\n<t:%d:R> -> `<t:%d:R>`", timestamp, timestamp, timestamp, timestamp))
			if err != nil {
				log.Println(err)
			}
		},
	}
}
