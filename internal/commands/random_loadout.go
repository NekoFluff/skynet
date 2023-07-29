package commands

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/NekoFluff/discord"
	"github.com/bwmarrin/discordgo"
)

func RandomLoadout() discord.Command {
	command := "randomloadout"

	return discord.Command{
		Command: discordgo.ApplicationCommand{
			Name:        command,
			Description: "Randomize your GTFO loadout",
		},
		Handler: func(s discord.Session, i *discordgo.InteractionCreate) {
			main := randomStringFromList([]string{"Pistol", "HEL Revolver", "Machine Pistol", "SMG", "PDW", "Carbine", "Assault Rifle", "Bullpup Rifle", "Burst Rifle", "DMR", "HEL Shotgun"})
			special := randomStringFromList([]string{"Machine Gun V", "Machine Gun XII", "Heavy Assault Rifle", "Shotgun", "Combat Shotgun", "Choke Mod Shotgun", "Revolver", "High Cal Pistol", "Precision Rifle", "Sniper"})
			tool := randomStringFromList([]string{"Bio Tracker", "C-Foam Launcher", "Mine Deployer", "Burst Sentry", "Shotgun Sentry", "Sniper Sentry"})
			melee := randomStringFromList([]string{"Sledgehammer", "Knife", "Bat", "Spear"})

			embed := discordgo.MessageEmbed{
				Title: "Random Loadout",
				Color: 15548997, // Red (https://gist.github.com/thomasbnt/b6f455e2c7d743b796917fa3c205f812)
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Main Weapon",
						Value: fmt.Sprintf("[%s](%s)", main, gtfoItems[main]),
					},
					{
						Name:  "Special Weapon",
						Value: fmt.Sprintf("[%s](%s)", special, gtfoItems[special]),
					},
					{
						Name:  "Tool",
						Value: fmt.Sprintf("[%s](%s)", tool, gtfoItems[tool]),
					},
					{
						Name:  "Melee Weapon",
						Value: fmt.Sprintf("[%s](%s)", melee, gtfoItems[melee]),
					},
				},
			}

			err := respondToInteractionWithEmbed(s, i.Interaction, &embed)

			if err != nil {
				log.Println(err)
			}
		},
	}
}

func randomStringFromList(list []string) string {
	return list[rand.Intn(len(list))]
}
