package commands

import (
	"fmt"
	"log"
	"log/slog"
	"path/filepath"
	"sync"
	"time"

	"github.com/NekoFluff/discord"
	"github.com/NekoFluff/skynet/internal/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

var (
	reminderStore   *utils.ReminderStore
	processorActive bool
	storeMutex      sync.Mutex
)

// initReminderStore initializes the reminder store if it's not already initialized
func InitReminderStore() error {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	if reminderStore != nil {
		return nil
	}

	// Create data directory in the current directory
	dataDir := filepath.Join(".", "data", "reminders")
	store, err := utils.NewReminderStore(dataDir)
	if err != nil {
		return fmt.Errorf("failed to initialize reminder store: %w", err)
	}

	reminderStore = store

	return nil
}

// startReminderProcessor starts a background goroutine to process reminders
func StartReminderProcessor(s discord.Session) {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	// Only start the processor once
	if processorActive {
		return
	}

	// Start the processor in a goroutine
	go func() {
		processorActive = true
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			<-ticker.C

			// Process any pending reminders
			pendingReminders := reminderStore.GetPendingReminders()

			for _, reminder := range pendingReminders {
				// Send the reminder message
				reminderMsg := fmt.Sprintf("<@%s> ⏰ **Reminder:** %s", reminder.UserID, reminder.Message)

				_, err := s.ChannelMessageSend(reminder.ChannelID, reminderMsg)
				if err != nil {
					slog.Error("failed to send reminder", "error", err, "channelId", reminder.ChannelID, "userId", reminder.UserID)
				}
			}
		}
	}()
}

func Reminder() discord.Command {
	command := "reminder"

	return discord.Command{
		Command: discordgo.ApplicationCommand{
			Name:        command,
			Description: "Set a reminder",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "The reminder message",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "datetime",
					Description: "When to send the reminder (e.g. 'Oct 8, 2025 4PM' or '1h30m')",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "timezone",
					Description: "The timezone to use (default MST)",
					Required:    false,
				},
			},
		},
		Handler: func(s discord.Session, i *discordgo.InteractionCreate) {
			// Initialize the reminder store if needed
			if err := InitReminderStore(); err != nil {
				slog.Error("failed to initialize reminder store", "error", err)
				err := respondToInteraction(s, i.Interaction, "Sorry, I couldn't set up the reminder system. Please try again later.")
				if err != nil {
					log.Println(err)
				}
				return
			}

			// Start the processor if this is the first reminder
			StartReminderProcessor(s)

			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			message := fmt.Sprint(optionMap["message"].Value)
			datetime := fmt.Sprint(optionMap["datetime"].Value)
			timezone := "MST"

			if optionMap["timezone"] != nil && optionMap["timezone"].Value != "" {
				timezone = fmt.Sprint(optionMap["timezone"].Value)
			}

			// First try to parse as a time duration (e.g., "1h30m")
			duration, err := time.ParseDuration(datetime)
			var timestamp int64
			var reminderTime time.Time

			if err == nil {
				// If it's a valid duration, use it relative to now
				reminderTime = time.Now().Add(duration)
				timestamp = reminderTime.Unix()
			} else {
				// Otherwise try to parse as a datetime string
				timestamp, err = utils.ConvertToUnixTimestamp(datetime, timezone)
				if err != nil {
					slog.Error("failed to parse reminder time", "error", err, "datetime", datetime, "timezone", timezone)
					err := respondToInteraction(s, i.Interaction, "Could not understand the reminder time format. Try something like '1h30m' for a relative time or 'Oct 8, 2025 4PM' for a specific date and time.")
					if err != nil {
						log.Println(err)
					}
					return
				}
				reminderTime = time.Unix(timestamp, 0)
			}

			// Check if the time is in the past
			if time.Now().After(reminderTime) {
				err := respondToInteraction(s, i.Interaction, "Cannot set a reminder in the past.")
				if err != nil {
					log.Println(err)
				}
				return
			}

			// Create and store the reminder
			userId := i.Member.User.ID
			channelId := i.ChannelID

			reminder := utils.Reminder{
				ID:          uuid.New().String(),
				UserID:      userId,
				ChannelID:   channelId,
				Message:     message,
				RemindTime:  reminderTime,
				CreatedTime: time.Now(),
			}

			err = reminderStore.AddReminder(reminder)
			if err != nil {
				slog.Error("failed to add reminder", "error", err)
				err := respondToInteraction(s, i.Interaction, "Sorry, I couldn't save your reminder. Please try again later.")
				if err != nil {
					log.Println(err)
				}
				return
			}

			// Format the response
			formattedTime := fmt.Sprintf("<t:%d:F>", timestamp)
			relativeTime := fmt.Sprintf("<t:%d:R>", timestamp)
			response := fmt.Sprintf("✓ I'll remind you about **%s** %s (%s)", message, formattedTime, relativeTime)

			err = respondToInteraction(s, i.Interaction, response)
			if err != nil {
				log.Println(err)
			}
		},
	}
}
