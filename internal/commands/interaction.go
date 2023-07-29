package commands

import (
	"github.com/NekoFluff/skynet/internal/mydiscord"
	"github.com/bwmarrin/discordgo"
)

func respondToInteraction(s mydiscord.Session, i *discordgo.Interaction, msg string) (err error) {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}

func respondToInteractionWithEmbed(s mydiscord.Session, i *discordgo.Interaction, e *discordgo.MessageEmbed) (err error) {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{e},
		},
	})
}
