//go:generate mockgen -package=commands -source=command.go -destination=command_mock_test.go
package commands

import (
	"strings"

	"warden/internal/utils"

	"github.com/bwmarrin/discordgo"
)

type Session interface {
	ChannelMessageSend(string, string) (*discordgo.Message, error)
	ChannelMessageSendEmbed(string, *discordgo.MessageEmbed) (*discordgo.Message, error)
}

type DiscordCommand struct {
	Command     string
	Description string
	Execute     func(s Session, m *discordgo.MessageCreate)
	prefix      string
}

func NewDiscordCommand(prefix string, command string, description string, execute func(s Session, m *discordgo.MessageCreate)) DiscordCommand {
	return DiscordCommand{
		Command:     prefix + command,
		Description: description,
		Execute:     execute,
		prefix:      prefix,
	}
}

type CommandsManager struct {
	Commands []DiscordCommand
}

func NewCommandsManager() *CommandsManager {
	c := &CommandsManager{}
	return c
}

func getPrefix() string {
	prefix := utils.GetEnvVar("COMMAND_PREFIX")
	if prefix == "" {
		prefix = "!"
	}
	return prefix
}

// This function will be called every time a new
// message is created on any channel that the authenticated bot has access to.
func (c *CommandsManager) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	cmd := strings.Split(strings.ToLower(m.Message.Content), " ")[0]
	for _, c := range c.Commands {
		if cmd == c.Command {
			c.Execute(s, m)
		}
	}
}
