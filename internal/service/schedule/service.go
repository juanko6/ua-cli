package service

import (
	"context"
	"fmt"
	"time"

	"github.com/juanko6/ua-cli/internal/domain/schedule"
)

// ScheduleService orchestrates the retrieval of schedule events.
type ScheduleService struct {
	repo    schedule.Repository
	adapter schedule.Repository // Cloud Adapter is also a Repository (Fetch)
}

// NewScheduleService creates a new instance of ScheduleService.
func NewScheduleService(repo schedule.Repository, adapter schedule.Repository) *ScheduleService {
	return &ScheduleService{
		repo:    repo,
		adapter: adapter,
	}
}

// GetScheduleForWeek returns the schedule for the week containing the given date.
// It implements the "Offline First" strategy with 24h TTL logic (TODO).
func (s *ScheduleService) GetScheduleForWeek(ctx context.Context, date time.Time, forceRefresh bool) ([]schedule.Event, error) {
	// 1. Try Cache (Repo)
	if !forceRefresh {
		events, err := s.repo.GetWeekEvents(ctx, date)
		if err == nil && len(events) > 0 {
			events = FilterEventsToWeek(events, date)
			events = DeduplicateEvents(events)
			if len(events) > 0 {
				return events, nil
			}
		}
	}

	// 2. Fetch from Cloud (Adapter)
	events, err := s.adapter.GetWeekEvents(ctx, date)
	if err != nil {
		return nil, err
	}

	// 3. Save to Cache
	_ = s.repo.SaveEvents(ctx, events)

	// 4. Filter and deduplicate
	events = FilterEventsToWeek(events, date)
	events = DeduplicateEvents(events)

	return events, nil
}

// FilterEventsToWeek returns only events whose Start falls within the
// Monday 00:00 – Sunday 23:59:59 range of the week containing `date`.
func FilterEventsToWeek(events []schedule.Event, date time.Time) []schedule.Event {
	monday, sunday := weekBounds(date)
	filtered := make([]schedule.Event, 0, len(events))
	for _, e := range events {
		if !e.Start.Before(monday) && e.Start.Before(sunday) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

// DeduplicateEvents removes duplicate events using composite key Title+Start+End.
func DeduplicateEvents(events []schedule.Event) []schedule.Event {
	seen := make(map[string]bool, len(events))
	result := make([]schedule.Event, 0, len(events))
	for _, e := range events {
		key := fmt.Sprintf("%s|%d|%d", e.Title, e.Start.Unix(), e.End.Unix())
		if !seen[key] {
			seen[key] = true
			result = append(result, e)
		}
	}
	return result
}

// weekBounds returns Monday 00:00:00 and next Monday 00:00:00 for the week containing date.
func weekBounds(date time.Time) (monday, nextMonday time.Time) {
	weekday := date.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	monday = date.AddDate(0, 0, -int(weekday-time.Monday))
	monday = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, date.Location())
	nextMonday = monday.AddDate(0, 0, 7)
	return monday, nextMonday
}
