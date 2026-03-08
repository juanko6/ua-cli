package grades

import (
	"time"
)

// GradeTracker handles change detection and tracking of grade updates
type GradeTracker struct {
	storage Storage
}

// Storage defines the interface for persisting grade tracking data
type Storage interface {
	GetLastCheck() (time.Time, error)
	SaveLastCheck(checkTime time.Time) error
	SavePreviousGrades(grades []Grade) error
	GetPreviousGrades() ([]Grade, error)
}

// NewGradeTracker creates a new grade tracker
func NewGradeTracker(storage Storage) *GradeTracker {
	return &GradeTracker{
		storage: storage,
	}
}

// DetectNewGrades compares current grades with previously stored grades
func (t *GradeTracker) DetectNewGrades(current []Grade) ([]Grade, error) {
	previous, err := t.storage.GetPreviousGrades()
	if err != nil {
		// If we can't get previous grades, assume all are new
		return current, nil
	}

	if len(previous) == 0 {
		// No previous grades, all are new
		return current, nil
	}

	var newGrades []Grade
	
	// Create lookup map for previous grades
	previousMap := make(map[string]Grade)
	for _, grade := range previous {
		previousMap[grade.SubjectID] = grade
	}
	
	// Compare current grades with previous ones
	for _, currentGrade := range current {
		if previous, exists := previousMap[currentGrade.SubjectID]; exists {
			// Check if this grade was updated since last check
			if currentGrade.LastUpdated.After(previous.LastUpdated) {
				newGrades = append(newGrades, currentGrade)
			}
		} else {
			// This is a new subject that wasn't there before
			newGrades = append(newGrades, currentGrade)
		}
	}
	
	return newGrades, nil
}

// UpdateTracking updates the tracking storage with current grades and timestamp
func (t *GradeTracker) UpdateTracking(current []Grade) error {
	now := time.Now()
	
	// Update last check timestamp
	if err := t.storage.SaveLastCheck(now); err != nil {
		return err
	}
	
	// Save current grades as previous for next comparison
	if err := t.storage.SavePreviousGrades(current); err != nil {
		return err
	}
	
	return nil
}

// GetLastCheckTime returns the last time grades were checked
func (t *GradeTracker) GetLastCheckTime() (time.Time, error) {
	return t.storage.GetLastCheck()
}

// HasPreviousGrades checks if there are previously stored grades
func (t *GradeTracker) HasPreviousGrades() (bool, error) {
	previous, err := t.storage.GetPreviousGrades()
	if err != nil {
		return false, err
	}
	return len(previous) > 0, nil
}