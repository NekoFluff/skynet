//go:generate mockgen -package=mydiscord -source=session.go -destination=session_mock.go
package mydiscord

import "github.com/bwmarrin/discordgo"

type Session interface {
	ChannelMessageSend(string, string) (*discordgo.Message, error)
	ChannelMessageSendEmbed(string, *discordgo.MessageEmbed) (*discordgo.Message, error)
	InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse) (err error)
}
