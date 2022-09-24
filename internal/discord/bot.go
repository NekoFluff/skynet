package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session *discordgo.Session
}

func NewBot(token string, handler func(*discordgo.Session, *discordgo.MessageCreate)) *Bot {
	session, err := createBot(token, handler)
	if err != nil {
		log.Fatal(err)
	}

	return &Bot{
		session: session,
	}
}

func (bot *Bot) Stop() {
	// Cleanly close down the Discord session.
	bot.session.Close()
}

func createBot(Token string, handler func(*discordgo.Session, *discordgo.MessageCreate)) (s *discordgo.Session, err error) {
	// Create a new Discord session using the provided bot token.
	s, err = discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the handler func as a callback for MessageCreate events.
	s.AddHandler(handler)

	s.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = s.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	return
}
