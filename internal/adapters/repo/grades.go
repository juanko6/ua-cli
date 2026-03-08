package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/juanko6/ua-cli/internal/domain/grades"
	"github.com/juanko6/ua-cli/internal/adapters/uacloud"
)

// UACloudGradesRepository implements GradesRepository with UACloud integration
type UACloudGradesRepository struct {
	mu         sync.RWMutex
	grades     []grades.Grade
	lastCheck  time.Time
	storagePath string
	adapter    *UACloudGradesAdapter
}

// NewUACloudGradesRepository creates a new UACloud grades repository
func NewUACloudGradesRepository(storagePath string, adapter *UACloudGradesAdapter) *UACloudGradesRepository {
	repo := &UACloudGradesRepository{
		storagePath: storagePath,
		adapter:    adapter,
	}
	
	// Load existing data on initialization
	_ = repo.load()
	
	return repo
}

// FetchGrades implements GradesRepository interface
func (r *UACloudGradesRepository) FetchGrades(ctx context.Context, session interface{}) ([]grades.Grade, error) {
	// Extract cookies from session
	cookies, ok := session.([]*http.Cookie)
	if !ok {
		return nil, fmt.Errorf("invalid session type: expected []*http.Cookie, got %T", session)
	}

	// Fetch grades from UACloud
	grades, err := r.adapter.FetchGrades(ctx, cookies)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch grades from UACloud: %w", err)
	}

	// Validate the fetched grades
	if err := r.adapter.ValidateGradeData(grades); err != nil {
		return nil, fmt.Errorf("invalid grade data: %w", err)
	}

	// Store grades in memory
	r.mu.Lock()
	r.grades = grades
	r.lastCheck = time.Now()
	r.mu.Unlock()

	// Persist to storage
	if err := r.save(); err != nil {
		// Don't fail the request if saving fails, but log it
		fmt.Fprintf(os.Stderr, "Warning: failed to save grades to storage: %v\n", err)
	}

	return grades, nil
}

// GetLastCheck implements GradesRepository interface
func (r *UACloudGradesRepository) GetLastCheck() (time.Time, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.lastCheck, nil
}

// SetLastCheck implements GradesRepository interface
func (r *UACloudGradesRepository) SetLastCheck(checkTime time.Time) error {
	r.mu.Lock()
	r.lastCheck = checkTime
	r.mu.Unlock()
	return r.save()
}

// DetectNewGrades implements GradesRepository interface
func (r *UACloudGradesRepository) DetectNewGrades(current []grades.Grade) ([]grades.Grade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.grades) == 0 {
		// No previous grades, all are new
		return current, nil
	}

	var newGrades []Grade
	
	// Create lookup map for previous grades
	previousMap := make(map[string]Grade)
	for _, grade := range r.grades {
		previousMap[grade.SubjectID] = grade
	}
	
	// Find new or updated grades
	for _, currentGrade := range current {
		if previous, exists := previousMap[currentGrade.SubjectID]; exists {
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

// GetGrades returns the current cached grades
func (r *UACloudGradesRepository) GetGrades() []grades.Grade {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.grades
}

// save persists grades to storage
func (r *UACloudGradesRepository) save() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Create storage directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(r.storagePath), 0755); err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Prepare data for saving
	data := struct {
		Grades    []grades.Grade `json:"grades"`
		LastCheck time.Time      `json:"last_check"`
	}{
		Grades:    r.grades,
		LastCheck: r.lastCheck,
	}

	// Marshal data to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal grades: %w", err)
	}

	// Write to file
	if err := os.WriteFile(r.storagePath, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to write grades file: %w", err)
	}

	return nil
}

// load loads grades from storage
func (r *UACloudGradesRepository) load() error {
	// Check if storage file exists
	if _, err := os.Stat(r.storagePath); os.IsNotExist(err) {
		// File doesn't exist, that's okay
		return nil
	}

	// Read file
	data, err := os.ReadFile(r.storagePath)
	if err != nil {
		return fmt.Errorf("failed to read grades file: %w", err)
	}

	// Parse JSON
	var storedData struct {
		Grades    []grades.Grade `json:"grades"`
		LastCheck time.Time      `json:"last_check"`
	}
	
	if err := json.Unmarshal(data, &storedData); err != nil {
		return fmt.Errorf("failed to unmarshal grades: %w", err)
	}

	// Load data into repository
	r.mu.Lock()
	r.grades = storedData.Grades
	r.lastCheck = storedData.LastCheck
	r.mu.Unlock()

	return nil
}

// ClearCache clears the grades cache and storage
func (r *UACloudGradesRepository) ClearCache() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.grades = nil
	r.lastCheck = time.Time{}

	// Remove storage file
	if err := os.Remove(r.storagePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove grades file: %w", err)
	}

	return nil
}