package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/NekoFluff/discord"
	"github.com/NekoFluff/skynet/internal/commands"
	"github.com/NekoFluff/skynet/internal/utils"
)

func main() {
	// Start up discord bot
	token := utils.GetEnvVar("DISCORD_BOT_TOKEN")
	bot := discord.NewBot(token)

	ids := utils.GetEnvVar("DEVELOPER_IDS")
	if ids != "" {
		bot.DeveloperIDs = strings.Split(ids, ",")
	}

	defer bot.Stop()

	// Generate Commands
	bot.AddCommands(commands.Ping(), commands.Pick(), commands.Roll(), commands.RandomLoadout(), commands.Lookup())
	bot.RegisterCommands()

	go handleSignalExit()

	// Bind to port
	port := utils.GetEnvVar("PORT")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Skynet is online."))
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
