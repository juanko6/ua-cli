package grades

import (
	"context"
	"time"
)

// GradesRepository defines the interface for grade data operations
type GradesRepository interface {
	// FetchGrades retrieves all grades for the current user
	FetchGrades(ctx context.Context, session interface{}) ([]Grade, error)
	
	// GetLastCheck returns the timestamp of the last grade check
	GetLastCheck() (time.Time, error)
	
	// SetLastCheck updates the timestamp of the last grade check
	SetLastCheck(checkTime time.Time) error
	
	// DetectNewGrades compares current grades with previous ones and returns new grades
	DetectNewGrades(current []Grade) ([]Grade, error)
}

// InMemoryGradesRepository is a test implementation of GradesRepository
type InMemoryGradesRepository struct {
	grades      []Grade
	lastCheck  time.Time
}

// NewInMemoryGradesRepository creates a new in-memory repository
func NewInMemoryGradesRepository() *InMemoryGradesRepository {
	return &InMemoryGradesRepository{
		lastCheck: time.Time{},
	}
}

// FetchGrades implements GradesRepository interface
func (r *InMemoryGradesRepository) FetchGrades(ctx context.Context, session interface{}) ([]Grade, error) {
	// Return stored grades for testing
	return r.grades, nil
}

// GetLastCheck implements GradesRepository interface
func (r *InMemoryGradesRepository) GetLastCheck() (time.Time, error) {
	return r.lastCheck, nil
}

// SetLastCheck implements GradesRepository interface
func (r *InMemoryGradesRepository) SetLastCheck(checkTime time.Time) error {
	r.lastCheck = checkTime
	return nil
}

// DetectNewGrades implements GradesRepository interface
func (r *InMemoryGradesRepository) DetectNewGrades(current []Grade) ([]Grade, error) {
	if len(r.grades) == 0 {
		// No previous grades, all are new
		return current, nil
	}

	var newGrades []Grade
	
	// Create a map for quick lookup of previous grades
	previousGrades := make(map[string]Grade)
	for _, grade := range r.grades {
		previousGrades[grade.SubjectID] = grade
	}
	
	// Find new or updated grades
	for _, currentGrade := range current {
		if previous, exists := previousGrades[currentGrade.SubjectID]; exists {
			// Check if grade was updated
			if currentGrade.LastUpdated.After(previous.LastUpdated) {
				newGrades = append(newGrades, currentGrade)
			}
		} else {
			// New subject
			newGrades = append(newGrades, currentGrade)
		}
	}
	
	return newGrades, nil
}

// AddGrade adds a grade to the repository (testing helper)
func (r *InMemoryGradesRepository) AddGrade(grade Grade) {
	r.grades = append(r.grades, grade)
}