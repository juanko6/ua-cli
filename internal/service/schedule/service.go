package service

import (
	"context"
	"time"

	"ua-cli/internal/domain/schedule"
)

// ScheduleService orchestrates the retrieval of schedule events.
type ScheduleService struct {
	repo    schedule.Repository
	adapter schedule.Repository // Cloud Adapter is also a Repository (Fetch)
}

// NewScheduleService creates a new instance of ScheduleService.
// We might need distinct interfaces for "CloudFetcher" vs "LocalCache" later,
// but for now, we can treat them as trusted sources.
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
			// TODO: Check TTL (We need metadata about when it was saved)
			// For MVP, if we have data, we return it.
			return events, nil
		}
	}

	// 2. Fetch from Cloud (Adapter)
	events, err := s.adapter.GetWeekEvents(ctx, date)
	if err != nil {
		return nil, err
	}

	// 3. Save to Cache
	_ = s.repo.SaveEvents(ctx, events)

	return events, nil
}
