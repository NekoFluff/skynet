//go:generate mockgen -package=commands -source=command.go -destination=command_mock_test.go
package commands

import (
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Command discordgo.ApplicationCommand
	Handler func(s Session, m *discordgo.InteractionCreate)
}
