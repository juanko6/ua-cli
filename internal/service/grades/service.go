package grades

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// GradesService handles the business logic for grades operations
type GradesService struct {
	repo      GradesRepository
	tracker   *GradeTracker
	presenter GradesPresenter
}

// GradesOptions defines the options for fetching and presenting grades
type GradesOptions struct {
	JSONOutput    bool     `json:"json_output"`
	ShowApproved  bool     `json:"show_approved"`
	ShowPending   bool     `json:"show_pending"`
	ShowAttention bool     `json:"show_attention"`
}

// GradesResult contains the result of a grades operation
type GradesResult struct {
	Grades         []Grade    `json:"grades"`
	NewGrades      []Grade    `json:"new_grades"`
	LastCheck      time.Time  `json:"last_check"`
	TotalSubjects  int        `json:"total_subjects"`
	AverageGrade   float64    `json:"average_grade"`
	HasChanges     bool       `json:"has_changes"`
	Message        string     `json:"message"`
}

// GradesPresenter defines the interface for presenting grades data
type GradesPresenter interface {
	Present(result *GradesResult) string
}

// NewGradesService creates a new grades service
func NewGradesService(
	repo GradesRepository,
	tracker *GradeTracker,
	presenter GradesPresenter,
) *GradesService {
	return &GradesService{
		repo:      repo,
		tracker:   tracker,
		presenter: presenter,
	}
}

// GetGrades fetches and presents grades based on the provided options
func (s *GradesService) GetGrades(ctx context.Context, session interface{}, opts GradesOptions) (*GradesResult, error) {
	// Fetch all grades
	grades, err := s.repo.FetchGrades(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch grades: %w", err)
	}

	// Filter grades based on options
	filtered := s.filterGrades(grades, opts)

	// Calculate statistics
	result := &GradesResult{
		Grades:        filtered,
		TotalSubjects: len(filtered),
		AverageGrade:  s.calculateAverage(filtered),
	}

	// Get last check time
	if lastCheck, err := s.repo.GetLastCheck(); err == nil {
		result.LastCheck = lastCheck
	}

	// Detect new grades if we have previous data
	if newGrades, err := s.repo.DetectNewGrades(filtered); err == nil {
		result.NewGrades = newGrades
		result.HasChanges = len(newGrades) > 0
	}

	// Generate appropriate message
	result.Message = s.generateMessage(result, opts)

	return result, nil
}

// DetectNewGrades detects and returns newly published grades
func (s *GradesService) DetectNewGrades(ctx context.Context, session interface{}) ([]Grade, error) {
	// Fetch current grades
	current, err := s.repo.FetchGrades(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch current grades: %w", err)
	}

	// Detect new grades using the tracker
	newGrades, err := s.tracker.DetectNewGrades(current)
	if err != nil {
		return nil, fmt.Errorf("failed to detect new grades: %w", err)
	}

	return newGrades, nil
}

// filterGrades filters grades based on the provided options
func (s *GradesService) filterGrades(grades []Grade, opts GradesOptions) []Grade {
	if !opts.ShowApproved && !opts.ShowPending && !opts.ShowAttention {
		// Show all if no filters specified
		return grades
	}

	var filtered []Grade
	for _, grade := range grades {
		switch grade.Status {
		case StatusApproved:
			if opts.ShowApproved {
				filtered = append(filtered, grade)
			}
		case StatusPending:
			if opts.ShowPending {
				filtered = append(filtered, grade)
			}
		case StatusNeedsAttention:
			if opts.ShowAttention {
				filtered = append(filtered, grade)
			}
		}
	}
	
	return filtered
}

// calculateAverage calculates the average grade across all subjects
func (s *GradesService) calculateAverage(grades []Grade) float64 {
	if len(grades) == 0 {
		return 0.0
	}

	var total float64
	for _, grade := range grades {
		total += grade.Average
	}
	
	return total / float64(len(grades))
}

// generateMessage generates an appropriate message based on the results
func (s *GradesService) generateMessage(result *GradesResult, opts GradesOptions) string {
	if opts.JSONOutput {
		return ""
	}

	if result.TotalSubjects == 0 {
		return "📊 No hay asignaturas con calificaciones disponibles."
	}

	if result.HasChanges {
		return fmt.Sprintf("🆕 ¡Tienes %d nuevas calificaciones!", len(result.NewGrades))
	}

	return "✅ No hay nuevas calificaciones desde la última revisión."
}