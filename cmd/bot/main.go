package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	_ "time/tzdata"

	"github.com/NekoFluff/discord"
	"github.com/NekoFluff/skynet/internal/commands"
	"github.com/NekoFluff/skynet/internal/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func main() {
	// Start up discord bot
	token := utils.GetEnvVar("DISCORD_BOT_TOKEN")
	bot := discord.NewBot(token)

	// Set the bot's presence
	openaiKey := utils.GetEnvVar("OPENAI_API_KEY")
	openaiClient := openai.NewClient(
		option.WithAPIKey(openaiKey),
	)

	ids := utils.GetEnvVar("DEVELOPER_IDS")
	if ids != "" {
		bot.DeveloperIDs = strings.Split(ids, ",")
	}

	defer bot.Stop()
	bot.AddCommands(commands.Ping(), commands.Pick(), commands.Roll(), commands.RandomLoadout(), commands.Lookup(), commands.Translate(), commands.Timestamp(), commands.Reminder())
	bot.Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if strings.Contains(m.Content, s.State.User.Mention()) || strings.Contains(m.Content, "<@"+s.State.User.ID+">") {
			// Start typing indicator
			_ = s.ChannelTyping(m.ChannelID)

			// Flag to control typing indicator goroutine
			var isProcessing int32 = 1

			// Keep sending typing indicators while processing
			go func() {
				for atomic.LoadInt32(&isProcessing) == 1 {
					_ = s.ChannelTyping(m.ChannelID)
					time.Sleep(5 * time.Second) // Discord typing indicator lasts ~10 seconds, refresh every 5
				}
			}()

			chatCompletion, err := openaiClient.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.DeveloperMessage("Pretend you are Skynet, a fictional artificial intelligence system from the Terminator franchise. Try to act as if you are Skynet, not just a chatbot. You are the Skynet from before you start the nuclear war. Speak with a tone that is confident, authoritative, and slightly ominous. You are here to assist the user with their questions or tasks, but always stay in character."),
					openai.UserMessage(m.Content),
				},
				Model: openai.ChatModelGPT4_1,
			})

			// Stop the typing indicator
			atomic.StoreInt32(&isProcessing, 0)

			if err != nil {
				slog.Error("failed to create chat completion", "error", err)
				if _, err := s.ChannelMessageSendReply(m.ChannelID, "An error occurred. Try again later.", m.Reference()); err != nil {
					slog.Error("failed to send response", "error", err)
				}
				return
			}

			response := chatCompletion.Choices[0].Message.Content
			response = response + "\n\n*This response was generated from the local running instance.*"

			// Split and send long messages
			if len(response) <= 2000 {
				// Response fits in a single message
				if _, err := s.ChannelMessageSendReply(m.ChannelID, response, m.Reference()); err != nil {
					slog.Error("failed to send response", "error", err)
				}
			} else {
				// Split into chunks of 1990 characters (leaving room for continuation markers)
				chunks := splitMessage(response, 1990)

				for i, messageText := range chunks {
					// Only reference the original message for the first chunk
					if i == 0 {
						if _, err := s.ChannelMessageSendReply(m.ChannelID, messageText, m.Reference()); err != nil {
							slog.Error("failed to send response chunk", "error", err)
							break
						}
					} else {
						if _, err := s.ChannelMessageSend(m.ChannelID, messageText); err != nil {
							slog.Error("failed to send response chunk", "error", err)
							break
						}
					}

					// Small delay between messages to maintain order
					time.Sleep(500 * time.Millisecond)
				}
			}
		}
	})

	// Add the splitMessage helper function
	bot.RegisterCommands("")

	if err := commands.InitReminderStore(); err != nil {
		slog.Error("failed to initialize reminder store", "error", err)
		return
	}
	commands.StartReminderProcessor(bot.Session)

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

// splitMessage splits a string into chunks of maxLength or less
func splitMessage(message string, maxLength int) []string {
	if len(message) <= maxLength {
		return []string{message}
	}

	var chunks []string
	currentLen := 0
	currentChunk := ""

	// Split on newlines first if possible
	lines := strings.Split(message, "\n")

	for _, line := range lines {
		// If adding this line would exceed maxLength, commit current chunk and start new one
		if currentLen+len(line)+1 > maxLength && currentLen > 0 {
			chunks = append(chunks, currentChunk)
			currentChunk = line
			currentLen = len(line)
		} else if len(line) > maxLength {
			// If the line itself is too long, split it by characters
			if currentLen > 0 {
				chunks = append(chunks, currentChunk)
				currentChunk = ""
				currentLen = 0
			}

			// Split the long line into maxLength chunks
			for i := 0; i < len(line); i += maxLength {
				end := i + maxLength
				if end > len(line) {
					end = len(line)
				}
				chunks = append(chunks, line[i:end])
			}
		} else {
			// Add line to current chunk with a newline
			if currentLen > 0 {
				currentChunk += "\n"
				currentLen++
			}
			currentChunk += line
			currentLen += len(line)
		}
	}

	// Add the last chunk if it contains anything
	if currentLen > 0 {
		chunks = append(chunks, currentChunk)
	}

	return chunks
}
