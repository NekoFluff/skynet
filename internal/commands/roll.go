package commands

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
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

			r, _ := regexp.Compile(`(?P<count>\d+)d(?P<size>\d+)`)
			matches := r.FindStringSubmatch(args[0])

			matchMap := make(map[string]string)
			for i, name := range r.SubexpNames() {
				if i != 0 && name != "" {
					matchMap[name] = matches[i]
				}
			}

			count, _ := strconv.Atoi(matchMap["count"])
			size, _ := strconv.Atoi(matchMap["size"])
			result := roll(count, size)
			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprint(result))
			if err != nil {
				log.Println(err)
			}
		},
	)
}

func roll(count int, size int) int {
	if size <= 0 {
		return 0
	}

	total := 0
	for i := 0; i < count; i++ {
		total += rand.Intn(size) + 1
	}
	return total
}
