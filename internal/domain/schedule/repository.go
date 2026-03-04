package schedule

import (
	"context"
	"time"
)

// Repository defines the interface for storing and retrieving schedule events.
type Repository interface {
	// GetWeekEvents retrieves events for a specific week (defined by a date within that week).
	GetWeekEvents(ctx context.Context, date time.Time) ([]Event, error)

	// SaveEvents persists a list of events.
	SaveEvents(ctx context.Context, events []Event) error
}
