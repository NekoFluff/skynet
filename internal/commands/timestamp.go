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
			Description: "Translate a date time string to a unix timestamp",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "datetime",
					Description: "The date and time to convert to a unix timestamp",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "timezone",
					Description: "The timezone to use when converting the date time to a unix timestamp (default MST)",
					Required:    false,
				},
			},
		},
		Handler: func(s discord.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			dateTime := fmt.Sprint(optionMap["datetime"].Value)
			timezone := "MST"

			if optionMap["timezone"] != nil {
				timezone = fmt.Sprint(optionMap["timezone"].Value)
			}

			timestamp, err := utils.ConvertToUnixTimestamp(dateTime, timezone)
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
