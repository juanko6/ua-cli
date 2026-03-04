package uacloud

import (
	"encoding/json"
	"strings"
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
	Aula      string `json:"uaIdAula"` // Classroom identifier, e.g. "A3/0007"
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
			Location: raw.Aula, // Classroom identifier from UACloud, e.g. "A3/0007"
			Type:     eventType,
		})
	}

	return events, nil
}

// theoryColors are hex codes used by UACloud for theory/lecture events.
var theoryColors = map[string]bool{
	"#6aa84f": true, "#274e13": true, "#38761d": true,
	"#93c47d": true, "#b6d7a8": true, "green": true,
}

// practiceColors are hex codes used by UACloud for practice/lab events.
var practiceColors = map[string]bool{
	"#3d85c6": true, "#0b5394": true, "#1155cc": true,
	"#6fa8dc": true, "#9fc5e8": true, "blue": true,
}

// seminarColors are hex codes used by UACloud for seminar/tutorial events.
var seminarColors = map[string]bool{
	"#e69138": true, "#b45f06": true, "#f6b26b": true,
	"#f9cb9c": true, "orange": true,
}

// classifyEvent infers the event type from CSS class name or color hex code.
func classifyEvent(className string, color string) schedule.EventType {
	lower := strings.ToLower(color)

	// 1. Try color-based classification
	if theoryColors[lower] {
		return schedule.TypeTheory
	}
	if practiceColors[lower] {
		return schedule.TypePractice
	}
	if seminarColors[lower] {
		return schedule.TypeSeminar
	}

	// 2. Fallback: className-based hints
	cn := strings.ToLower(className)
	switch {
	case strings.Contains(cn, "teor"):
		return schedule.TypeTheory
	case strings.Contains(cn, "pract") || strings.Contains(cn, "lab"):
		return schedule.TypePractice
	case strings.Contains(cn, "semin") || strings.Contains(cn, "tutor"):
		return schedule.TypeSeminar
	}

	return schedule.TypeUnknown
}
