package uacloud

import (
	"encoding/json"


	"ua-cli/internal/domain/schedule"
)

// ParseSchedule parses the raw response from UACloud.
// Note: This is a placeholder. The actual parsing logic depends heavily on the response format.
// Assuming for now it returns a JSON-like structure inside HTML or pure JSON if requested correctly.
func ParseSchedule(data []byte) ([]schedule.Event, error) {
	// TODO: Implement actual parsing logic.
	// This requires inspecting the real response from 'cvnet'.
	// For now, we return existing dummy data or try to unmarshal if it's JSON.
	
	var events []schedule.Event
	// Attempt JSON unmarshal
	if err := json.Unmarshal(data, &events); err == nil {
		return events, nil
	}

	// If fail, maybe it's HTML and we need regex?
	// Returning empty list for now until we have real data samples.
	return []schedule.Event{}, nil
}
