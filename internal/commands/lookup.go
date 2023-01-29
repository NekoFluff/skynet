package commands

import (
	"fmt"
	"log"
	"warden/internal/discord"

	"github.com/bwmarrin/discordgo"
)

var gtfoItems map[string]string = map[string]string{
	"Pistol":              "https://gtfo.fandom.com/wiki/Shelling_S49_Pistol",
	"HEL Revolver":        "https://gtfo.fandom.com/wiki/Bataldo_Model_8_HEL_Revolver",
	"Machine Pistol":      "https://gtfo.fandom.com/wiki/Raptus_Treffen_2_Machine_Pistol",
	"SMG":                 "https://gtfo.fandom.com/wiki/Van_Auken_LTC_5_SMG",
	"PDW":                 "https://gtfo.fandom.com/wiki/Accrat_STB_PDW",
	"Carbine":             "https://gtfo.fandom.com/wiki/Van_Auken_CAB_F4_Carbine",
	"Assault Rifle":       "https://gtfo.fandom.com/wiki/Malatack_LX_Assault_Rifle",
	"Bullpup Rifle":       "https://gtfo.fandom.com/wiki/Accrat_Golok_DA_Bullpup_Rifle",
	"Burst Rifle":         "https://gtfo.fandom.com/wiki/Malatack_CH_4_Burst_Rifle",
	"DMR":                 "https://gtfo.fandom.com/wiki/TR22_Hanaway_DMR",
	"HEL Shotgun":         "https://gtfo.fandom.com/wiki/Bataldo_J_300_Hel_Shotgun",
	"Machine Gun V":       "https://gtfo.fandom.com/wiki/Techman_Arbalist_V_Machine_Gun",
	"Machine Gun XII":     "https://gtfo.fandom.com/wiki/Techman_Veruta_XII_Machine_Gun",
	"Heavy Assault Rifle": "https://gtfo.fandom.com/wiki/Malatack_HXC_Heavy_Assault_Rifle",
	"Shotgun":             "https://gtfo.fandom.com/wiki/Buckland_S870_Shotgun",
	"Combat Shotgun":      "https://gtfo.fandom.com/wiki/Buckland_AF6_Combat_Shotgun",
	"Choke Mod Shotgun":   "https://gtfo.fandom.com/wiki/Buckland_XDIST2_Choke_Mod_Shotgun",
	"Revolver":            "https://gtfo.fandom.com/wiki/Mastaba_R66_Revolver",
	"High Cal Pistol":     "https://gtfo.fandom.com/wiki/Shelling_Arid_5_High_Caliber_Pistol",
	"Precision Rifle":     "https://gtfo.fandom.com/wiki/Drekker_Del_P1_Precision_Rifle",
	"Sniper":              "https://gtfo.fandom.com/wiki/K%C3%B6ning_PR_11_Sniper",
	"Sledgehammer":        "https://gtfo.fandom.com/wiki/Santonian_HDH_Sledgehammer",
	"Knife":               "https://gtfo.fandom.com/wiki/Mastaba_Fixed_Blade_Knife",
	"Bat":                 "https://gtfo.fandom.com/wiki/Kovac_Peacekeeper_Bat",
	"Spear":               "https://gtfo.fandom.com/wiki/Maco_Drillhead_Spear",
	"Rifle":               "https://gtfo.fandom.com/wiki/Drekker_Pres_MOD_556_Rifle",
	"Double Tap Rifle":    "https://gtfo.fandom.com/wiki/Hanaway_PSB_Double_Tap_Rifle",
	"Autopistol":          "https://gtfo.fandom.com/wiki/Raptus_Steigro_Autopistol",
	"Heavy SMG":           "https://gtfo.fandom.com/wiki/Accrat_ND6_Heavy_SMG",
	"Sawed-off Shotgun":   "https://gtfo.fandom.com/wiki/Buckland_SBS_III_Sawed-off_Shotgun",
	"Burst Cannon":        "https://gtfo.fandom.com/wiki/Techman_Klust_6_Burst_Cannon",
	"Scattergun":          "https://gtfo.fandom.com/wiki/Drekker_INEX_Drei_Scattergun",
	"HEL Rifle":           "https://gtfo.fandom.com/wiki/Omneco_LRG_HEL_Rifle",
	"HEL Gun":             "https://gtfo.fandom.com/wiki/Omneco_EXP1_HEL_Gun",
	"Bio Tracker":         "https://gtfo.fandom.com/wiki/D-Tek_Optron_IV_Bio_Tracker",
	"C-Foam Launcher":     "https://gtfo.fandom.com/wiki/Stalwart_G2_C-Foam_Launcher",
	"Mine Deployer":       "https://gtfo.fandom.com/wiki/Krieger_04_Mine_Deployer",
	"Burst Sentry":        "https://gtfo.fandom.com/wiki/Mechatronic_SGB3_Burst_Sentry",
	"Shotgun Sentry":      "https://gtfo.fandom.com/wiki/Mechatronic_B5_LFR_Shotgun_Sentry",
	"Sniper Sentry":       "https://gtfo.fandom.com/wiki/AutoTek_51_RSG_Sniper_Sentry",
}

func Lookup() discord.Command {

	choices := []*discordgo.ApplicationCommandOptionChoice{}
	for name := range gtfoItems {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  name,
			Value: name,
		})
		if len(choices) == 25 {
			break
		}
	}

	command := "lookup"

	return discord.Command{
		Command: discordgo.ApplicationCommand{
			Name:        command,
			Description: fmt.Sprintf("Lookup a GTFO weaopn (e.g. `/%s HEL Revolver`)", command),
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "item",
					Description: "The item to lookup",
					Required:    true,
					Choices:     choices,
				},
			},
		},
		Handler: func(s discord.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			item := fmt.Sprint(optionMap["item"].Value)

			result := ""
			for _, choice := range choices {
				if choice.Name == item {
					result = gtfoItems[item]
					break
				}
			}
			err := respondToInteraction(s, i.Interaction, fmt.Sprint(result))

			if err != nil {
				log.Println(err)
			}
		},
	}
}
