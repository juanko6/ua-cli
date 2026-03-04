package repo

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"ua-cli/internal/domain/schedule"
)

type JSONFileRepo struct {
	filePath string
}

func NewJSONFileRepo(path string) (*JSONFileRepo, error) {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}
	return &JSONFileRepo{filePath: path}, nil
}

type CachedSchedule struct {
	UpdatedAt time.Time        `json:"updated_at"`
	Events    []schedule.Event `json:"events"`
}

func (r *JSONFileRepo) GetWeekEvents(ctx context.Context, date time.Time) ([]schedule.Event, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No cache yet
		}
		return nil, err
	}
	defer file.Close()

	var cached CachedSchedule
	if err := json.NewDecoder(file).Decode(&cached); err != nil {
		return nil, err
	}

	// Check TTL here if we wanted to enforce strictly in repo,
	// but Service handles logic usually.
	// We just return what we have.
	return cached.Events, nil
}

func (r *JSONFileRepo) SaveEvents(ctx context.Context, events []schedule.Event) error {
	cached := CachedSchedule{
		UpdatedAt: time.Now(),
		Events:    events,
	}

	file, err := os.Create(r.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(cached)
}
