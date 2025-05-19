package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Reminder structure to store reminder data
type Reminder struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ChannelID   string    `json:"channel_id"`
	Message     string    `json:"message"`
	RemindTime  time.Time `json:"remind_time"`
	CreatedTime time.Time `json:"created_time"`
}

// ReminderStore manages reminder persistence
type ReminderStore struct {
	reminders []Reminder
	filePath  string
	mutex     sync.Mutex
}

// NewReminderStore creates a new reminder store with persistence
func NewReminderStore(dataDir string) (*ReminderStore, error) {
	// Create reminders directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create reminders data directory: %w", err)
	}

	filePath := filepath.Join(dataDir, "reminders.json")
	
	store := &ReminderStore{
		filePath:  filePath,
		reminders: []Reminder{},
	}

	// Load existing reminders if file exists
	if _, err := os.Stat(filePath); err == nil {
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read reminders file: %w", err)
		}

		if len(data) > 0 {
			if err := json.Unmarshal(data, &store.reminders); err != nil {
				return nil, fmt.Errorf("failed to parse reminders file: %w", err)
			}
		}
	}

	return store, nil
}

// AddReminder adds a new reminder and persists it
func (rs *ReminderStore) AddReminder(reminder Reminder) error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	// Add the reminder
	rs.reminders = append(rs.reminders, reminder)

	// Save to file
	return rs.saveToFile()
}

// RemoveReminder removes a reminder by ID
func (rs *ReminderStore) RemoveReminder(id string) error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	for i, r := range rs.reminders {
		if r.ID == id {
			// Remove the reminder
			rs.reminders = append(rs.reminders[:i], rs.reminders[i+1:]...)
			return rs.saveToFile()
		}
	}

	return fmt.Errorf("reminder with ID %s not found", id)
}

// GetAllReminders returns all pending reminders
func (rs *ReminderStore) GetAllReminders() []Reminder {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	return rs.reminders
}

// GetPendingReminders returns reminders that should be sent now
func (rs *ReminderStore) GetPendingReminders() []Reminder {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	var pending []Reminder
	now := time.Now()
	
	for i := 0; i < len(rs.reminders); i++ {
		if rs.reminders[i].RemindTime.Before(now) || rs.reminders[i].RemindTime.Equal(now) {
			pending = append(pending, rs.reminders[i])
			
			// Remove the pending reminder
			rs.reminders = append(rs.reminders[:i], rs.reminders[i+1:]...)
			i--
		}
	}

	// Save if we found any pending reminders
	if len(pending) > 0 {
		_ = rs.saveToFile()
	}

	return pending
}

// saveToFile saves reminders to the file
func (rs *ReminderStore) saveToFile() error {
	data, err := json.MarshalIndent(rs.reminders, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal reminders: %w", err)
	}

	if err := ioutil.WriteFile(rs.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to save reminders file: %w", err)
	}

	return nil
}
