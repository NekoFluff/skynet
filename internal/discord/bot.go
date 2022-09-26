package discord

import (
	"fmt"
	"log"

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
