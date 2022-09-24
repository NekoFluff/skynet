package commands

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Roll() DiscordCommand {
	prefix := getPrefix()
	command := "roll"

	return NewDiscordCommand(
		prefix,
		command,
		fmt.Sprintf("Roll some dice (e.g. `%s%s 3d20`)", prefix, command),
		func(s Session, m *discordgo.MessageCreate) {
			args := strings.Split(m.Content, " ")[1:]
			if len(args) == 0 {
				_, err := s.ChannelMessageSend(m.ChannelID, "No arguments provided")
				if err != nil {
					log.Println(err)
				}
				return
			}

			count := 1
			size := 10
			result := roll(count, size)
			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprint(result))
			if err != nil {
				log.Println(err)
			}
		},
	)
}

func roll(count int, size int) int {
	total := 0
	for i := 0; i < count; i++ {
		total += rand.Intn(size) + 1
	}
	return total
}
