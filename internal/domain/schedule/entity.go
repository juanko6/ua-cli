package schedule

import "time"

// EventType represents the type of academic event (Theory, Practice, etc.)
type EventType string

const (
	TypeTheory   EventType = "Teoria"
	TypePractice EventType = "Practica"
	TypeSeminar  EventType = "Seminario"
	TypeUnknown  EventType = "Desconocido"
)

// Event represents a single class or academic activity in the schedule.
type Event struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Location string    `json:"location"`
	Type     EventType `json:"type"`
}
