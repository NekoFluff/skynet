package commands

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/NekoFluff/discord"
	"github.com/bwmarrin/discordgo"
)

func Translate() discord.Command {
	command := "translate"

	return discord.Command{
		Command: discordgo.ApplicationCommand{
			Name:        command,
			Description: fmt.Sprintf("Translates to a d20 roll (e.g. `%s (d8+WeaponDamageBonus)+[4+[2*#OfUpCasts]]d8`)", command),
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "dice",
					Description: "The dice to roll",
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

			diceRoll := fmt.Sprint(optionMap["dice"].Value)

			translation := translate(diceRoll)
			err := respondToInteraction(s, i.Interaction, translation)

			if err != nil {
				log.Println(err)
			}
		},
	}
}

func translate(input string) string {
	// Define the character mapping
	charMap := map[byte]string{
		'(': "{",
		')': "&#125",
		'[': "[[",
		']': "]]",
	}

	// Replace characters based on the mapping
	translated := strings.Builder{}
	for i := 0; i < len(input); i++ {
		if replacement, ok := charMap[input[i]]; ok {
			translated.WriteString(replacement)
		} else {
			translated.WriteByte(input[i])
		}
	}

	tempString := translated.String()
	re := regexp.MustCompile(`#?[A-Za-z][A-Za-z]+`)
	tempString = re.ReplaceAllString(tempString, "?{$0&#125")

	return "[[" + tempString + "]]"
}
