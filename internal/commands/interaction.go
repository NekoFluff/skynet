package commands

import (
	"fmt"
	"strings"

	"github.com/NekoFluff/discord"
	"github.com/bwmarrin/discordgo"
)

func respondToInteraction(s discord.Session, i *discordgo.Interaction, msg string) (err error) {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}

func respondToInteractionWithEmbed(s discord.Session, i *discordgo.Interaction, e *discordgo.MessageEmbed) (err error) {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{e},
		},
	})
}

type containsInteractionResponseMatcher struct {
	msg string
}

func (e containsInteractionResponseMatcher) Matches(input interface{}) bool {
	response, ok := input.(*discordgo.InteractionResponse)
	if !ok {
		return false
	}

	return strings.Contains(response.Data.Content, e.msg)
}

func (e containsInteractionResponseMatcher) String() string {
	return fmt.Sprintf("to contain msg '%v'", e.msg)
}

func (e containsInteractionResponseMatcher) Got(input interface{}) string {
	response, ok := input.(*discordgo.InteractionResponse)
	if !ok {
		return "not an interaction response"
	}

	return response.Data.Content
}
