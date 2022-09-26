package commands

import "github.com/bwmarrin/discordgo"

type CommandsManager struct {
	Commands map[string]Command
}

func NewCommandsManager() *CommandsManager {
	c := &CommandsManager{
		Commands: make(map[string]Command),
	}
	return c
}

func (c *CommandsManager) AddCommand(cmd Command) {
	c.Commands[cmd.Command.Name] = cmd
}

// This function will be called every time a new
// message is created on any channel that the authenticated bot has access to.
func (c *CommandsManager) HandleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Ignore all messages created by the bot itself
	// if i.Message.Author.ID == s.State.User.ID {
	// 	return
	// }

	if cmd, ok := c.Commands[i.ApplicationCommandData().Name]; ok {
		cmd.Handler(s, i)
	}
}

func respondToInteraction(s Session, i *discordgo.Interaction, msg string) (err error) {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}
