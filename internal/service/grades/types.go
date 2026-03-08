package grades

import (
	"fmt"
	"strings"
	"time"
)

// NewGrade represents a newly detected grade change
type NewGrade struct {
	SubjectID     string    `json:"subject_id"`
	SubjectName   string    `json:"subject_name"`
	OldGrade      string    `json:"old_grade,omitempty"`
	NewGrade      string    `json:"new_grade"`
	ChangeTime    time.Time `json:"change_time"`
	ChangeType    ChangeType `json:"change_type"`
}

// ChangeType represents the type of grade change
type ChangeType string

const (
	ChangeTypeNew      ChangeType = "new"
	ChangeTypeUpdated  ChangeType = "updated"
	ChangeTypeImproved ChangeType = "improved"
	ChangeTypeDeclined ChangeType = "declined"
)

// GradeFilter represents filtering criteria for grades
type GradeFilter struct {
	Status      []GradeStatus `json:"status,omitempty"`
	MinAverage  float64       `json:"min_average,omitempty"`
	MaxAverage  float64       `json:"max_average,omitempty"`
	SubjectName string        `json:"subject_name,omitempty"`
}

// GradeStatistics provides statistical information about grades
type GradeStatistics struct {
	TotalSubjects     int       `json:"total_subjects"`
	ApprovedCount     int       `json:"approved_count"`
	PendingCount      int       `json:"pending_count"`
	NeedsAttention    int       `json:"needs_attention"`
	AverageGrade      float64   `json:"average_grade"`
	HighestGrade      float64   `json:"highest_grade"`
	LowestGrade       float64   `json:"lowest_grade"`
	LastUpdated       time.Time `json:"last_updated"`
}

// GradeExport represents grade data for export purposes
type GradeExport struct {
	ExportDate    string        `json:"export_date"`
	StudentInfo   StudentInfo   `json:"student_info"`
	Statistics   GradeStatistics `json:"statistics"`
	Grades       []Grade       `json:"grades"`
}

// StudentInfo represents student information for export
type StudentInfo struct {
	Name        string `json:"name"`
	StudentID   string `json:"student_id"`
	Faculty     string `json:"faculty"`
	Program     string `json:"program"`
	Semester    string `json:"semester"`
}

// Validate validates the grade filter
func (f *GradeFilter) Validate() error {
	if f.MinAverage > f.MaxAverage && f.MaxAverage > 0 {
		return fmt.Errorf("min_average cannot be greater than max_average")
	}
	return nil
}

// Apply applies the filter to a list of grades
func (f *GradeFilter) Apply(grades []Grade) []Grade {
	var filtered []Grade
	
	for _, grade := range grades {
		// Filter by status
		if len(f.Status) > 0 {
			statusMatch := false
			for _, status := range f.Status {
				if grade.Status == status {
					statusMatch = true
					break
				}
			}
			if !statusMatch {
				continue
			}
		}
		
		// Filter by average
		if f.MinAverage > 0 && grade.Average < f.MinAverage {
			continue
		}
		if f.MaxAverage > 0 && grade.Average > f.MaxAverage {
			continue
		}
		
		// Filter by subject name (case insensitive)
		if f.SubjectName != "" {
			if !strings.Contains(strings.ToLower(grade.SubjectName), strings.ToLower(f.SubjectName)) {
				continue
			}
		}
		
		filtered = append(filtered, grade)
	}
	
	return filtered
}