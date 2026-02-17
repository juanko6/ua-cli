package uacloud

import (
	"encoding/json"
	"time"

	"ua-cli/internal/domain/schedule"
)

// UACloudEvent represents the raw JSON event from the UACloud API (FullCalendar format).
type UACloudEvent struct {
	Title     string `json:"title"`
	Start     string `json:"start"`     // ISO 8601: "2026-02-17T09:00:00"
	End       string `json:"end"`       // ISO 8601: "2026-02-17T11:00:00"
	ClassName string `json:"className"` // CSS class, may indicate type
	Color     string `json:"color"`
	ID        string `json:"id"`
	Rendering string `json:"rendering"`
}

// ParseSchedule parses the JSON response from UACloud's ObtenerEventosCalendarioJson.
func ParseSchedule(data []byte) ([]schedule.Event, error) {
	var rawEvents []UACloudEvent
	if err := json.Unmarshal(data, &rawEvents); err != nil {
		return nil, err
	}

	events := make([]schedule.Event, 0, len(rawEvents))
	for _, raw := range rawEvents {
		// Skip background events (e.g. rendering="background")
		if raw.Rendering == "background" {
			continue
		}

		start, _ := time.Parse("2006-01-02T15:04:05", raw.Start)
		end, _ := time.Parse("2006-01-02T15:04:05", raw.End)

		eventType := classifyEvent(raw.ClassName, raw.Color)

		events = append(events, schedule.Event{
			ID:       raw.ID,
			Title:    raw.Title,
			Start:    start,
			End:      end,
			Location: "", // UACloud JSON may not include location; we fill it if available
			Type:     eventType,
		})
	}

	return events, nil
}

// classifyEvent tries to infer the event type from CSS class or color.
func classifyEvent(className string, color string) schedule.EventType {
	switch {
	case color == "#6aa84f" || color == "green":
		return schedule.TypeTheory
	case color == "#3d85c6" || color == "blue":
		return schedule.TypePractice
	case color == "#e69138" || color == "orange":
		return schedule.TypeSeminar
	default:
		return schedule.TypeUnknown
	}
}
