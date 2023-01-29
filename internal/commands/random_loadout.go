package commands

import (
	"fmt"
	"log"
	"math/rand"
	"warden/internal/discord"

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

			main := fmt.Sprintf("**Main Weapon**: %s", randomStringFromList([]string{"Pistol", "HEL Revolver", "Machine Pistol", "SMG", "PDW", "Carbine", "Assault Rifle", "Bullpup Rifle", "Burst Rifle", "DMR", "HEL Shotgun"}))
			special := fmt.Sprintf("**Special Weapon**: %s", randomStringFromList([]string{"Machine Gun V", "Machine Gun XII", "Heavy Assault Rifle", "Shotgun", "Combat Shotgun", "Choke Mod Shotgun", "Revolver", "High Cal Pistol", "Precision Rifle", "Sniper"}))
			tool := fmt.Sprintf("**Tool**: %s", randomStringFromList([]string{"Bio Tracker", "C-Foam Launcher", "Mine Deployer", "Burst Sentry", "Shotgun Sentry", "Sniper Sentry"}))
			melee := fmt.Sprintf("**Melee Weapon**: %s", randomStringFromList([]string{"Sledgehammer", "Knife", "Bat", "Spear"}))

			response := fmt.Sprintf("%s\n%s\n%s\n%s", main, special, tool, melee)

			err := respondToInteraction(s, i.Interaction, response)

			if err != nil {
				log.Println(err)
			}
		},
	}
}

func randomStringFromList(list []string) string {
	return list[rand.Intn(len(list))]
}
