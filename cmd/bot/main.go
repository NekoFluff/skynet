package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"warden/internal/commands"
	"warden/internal/discord"
	"warden/internal/utils"
)

func main() {
	// Initialize global pseudo random generator
	rand.Seed(time.Now().Unix())

	// Start up discord bot
	token := utils.GetEnvVar("DISCORD_BOT_TOKEN")
	cmdMgr := commands.NewCommandsManager()
	bot := discord.NewBot(token, cmdMgr.MessageCreate)
	defer bot.Stop()

	// Generate Commands
	commands := []commands.DiscordCommand{
		commands.Help(&cmdMgr.Commands),
		commands.Ping(),
		commands.Pick(),
		commands.Roll(),
	}
	cmdMgr.Commands = append(cmdMgr.Commands, commands...)

	go handleSignalExit()

	// Bind to port
	port := utils.GetEnvVar("PORT")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Warden is online."))
	})
	fmt.Printf("Warden is online. Serving on port %s\n", port)

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
