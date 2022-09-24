package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func (bot *Bot) SendEmbedMessage(channelName string, message *discordgo.MessageEmbed) {
	for _, guild := range bot.session.State.Guilds {
		// Get channels for this guild (a.k.a discord server)
		channels, _ := bot.session.GuildChannels(guild.ID)

		for _, c := range channels {
			// Ensure the channel is a guild text channel and not a voice or DM channel
			if c.Type != discordgo.ChannelTypeGuildText {
				continue
			}

			// Check if the channel name matches target name
			if c.Name != channelName {
				continue
			}

			// Send a message to the discord channel
			_, err := bot.session.ChannelMessageSendEmbed(
				c.ID,
				message,
			)
			if err != nil {
				log.Println("An error occurred while sending a message to a discord server")
				log.Println(err)
			}
		}
	}
}
