package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"warden/internal/commands"
	"warden/internal/discord"
	"warden/internal/utils"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Initialize global pseudo random generator
	rand.Seed(time.Now().Unix())

	// Start up discord bot
	token := utils.GetEnvVar("DISCORD_BOT_TOKEN")
	cmdMgr := discord.NewCommandsManager()
	bot := discord.NewBot(token)

	defer bot.Stop()

	// Generate Commands
	cmdMgr.AddCommands(commands.Ping(), commands.Pick(), commands.Roll())
	bot.Session.AddHandler(cmdMgr.HandleInteractionCreate)
	bot.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// registeredCommands := make([]*discordgo.ApplicationCommand, len(cmdMgr.Commands))
	for _, cmd := range cmdMgr.Commands {
		cmd, err := bot.Session.ApplicationCommandCreate(bot.Session.State.User.ID, "", &cmd.Command)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", cmd.Name, err)
		}
		// registeredCommands[i] = cmd
	}

	go handleSignalExit()

	// Bind to port
	port := utils.GetEnvVar("PORT")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Warden is online."))
	})
	fmt.Printf("Serving on port %s\n", port)

	_ = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

// Wait until CTRL-C or other term signal is received.
func handleSignalExit() {
	fmt.Println("Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	os.Exit(1)
}
