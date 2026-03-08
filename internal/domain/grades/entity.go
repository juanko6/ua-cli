package grades

import (
	"time"
)

// Grade represents a student's grade for a specific subject
type Grade struct {
	SubjectID        string    `json:"subject_id"`
	SubjectName      string    `json:"subject_name"`
	CurrentGrade     string    `json:"current_grade"`
	Average          float64   `json:"average"`
	AssessmentCount  int       `json:"assessment_count"`
	Status           GradeStatus `json:"status"`
	LastUpdated      time.Time `json:"last_updated"`
}

// GradeStatus represents the academic status of a grade
type GradeStatus string

const (
	StatusApproved     GradeStatus = "approved"
	StatusPending      GradeStatus = "pending"
	StatusNeedsAttention GradeStatus = "needs_attention"
)

// IsValid checks if the grade data is valid
func (g *Grade) IsValid() bool {
	return g.SubjectID != "" && 
		   g.SubjectName != "" && 
		   g.Status != "" &&
		   !g.LastUpdated.IsZero()
}

// NeedsAttention returns true if the grade needs attention
func (g *Grade) NeedsAttention() bool {
	return g.Status == StatusNeedsAttention
}

// IsPending returns true if the grade is pending
func (g *Grade) IsPending() bool {
	return g.Status == StatusPending
}

// IsApproved returns true if the grade is approved
func (g *Grade) IsApproved() bool {
	return g.Status == StatusApproved
}

// String returns a string representation of the grade status
func (s GradeStatus) String() string {
	return string(s)
}

// DisplayEmoji returns an emoji representing the grade status
func (s GradeStatus) DisplayEmoji() string {
	switch s {
	case StatusApproved:
		return "✅"
	case StatusPending:
		return "⏳"
	case StatusNeedsAttention:
		return "⚠️"
	default:
		return "❓"
	}
}