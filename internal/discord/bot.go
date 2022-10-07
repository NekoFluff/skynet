package discord

import (
	"fmt"
	"log"
	"strings"
	"warden/internal/utils"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session *discordgo.Session
}

func NewBot(token string) *Bot {
	session, err := createBot(token)
	if err != nil {
		log.Fatal(err)
	}

	return &Bot{
		Session: session,
	}
}

func (bot *Bot) Stop() {
	// Cleanly close down the Discord session.
	bot.Session.Close()
}

func createBot(Token string) (s *discordgo.Session, err error) {
	// Create a new Discord session using the provided bot token.
	s, err = discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	s.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = s.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	return
}

func (bot *Bot) SendChannelMessage(channelName string, message string) {
	for _, guild := range bot.Session.State.Guilds {
		// Get channels for this guild (a.k.a discord server)
		channels, _ := bot.Session.GuildChannels(guild.ID)

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
			_, err := bot.Session.ChannelMessageSend(
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

func (bot *Bot) SendDeveloperMessage(message string) {
	developer_mode := utils.GetEnvVar("DEVELOPER_MODE")
	if developer_mode != "ON" && developer_mode != "1" {
		return
	}

	developerIds := getDeveloperIds()

	for _, developerId := range developerIds {
		ch, err := bot.Session.UserChannelCreate(developerId)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = bot.Session.ChannelMessageSend(ch.ID, message)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func getDeveloperIds() []string {
	ids := utils.GetEnvVar("DEVELOPER_IDS")
	if ids == "" {
		return []string{}
	}
	return strings.Split(ids, ",")
}

func (bot *Bot) SendEmbedMessage(channelName string, message *discordgo.MessageEmbed) {
	for _, guild := range bot.Session.State.Guilds {
		// Get channels for this guild (a.k.a discord server)
		channels, _ := bot.Session.GuildChannels(guild.ID)

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
			_, err := bot.Session.ChannelMessageSendEmbed(
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
