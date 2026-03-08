package uacloud

import (
	"context"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/juanko6/ua-cli/internal/domain/grades"
)

// UACloudGradesAdapter handles fetching grades from UACloud
type UACloudGradesAdapter struct {
	baseURL    string
	httpClient *http.Client
}

// NewUACloudGradesAdapter creates a new UACloud grades adapter
func NewUACloudGradesAdapter(baseURL string) *UACloudGradesAdapter {
	// Create cookie jar for session persistence
	jar, err := cookiejar.New(nil)
	if err != nil {
		// Fallback to client without jar if creation fails
		return &UACloudGradesAdapter{
			baseURL:    baseURL,
			httpClient: &http.Client{Timeout: 30 * time.Second},
		}
	}

	return &UACloudGradesAdapter{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second, Jar: jar},
	}
}

// FetchGrades retrieves grades from UACloud for the authenticated session
func (a *UACloudGradesAdapter) FetchGrades(ctx context.Context, sessionCookies []*http.Cookie) ([]grades.Grade, error) {
	// Build the grades URL
	gradesURL := fmt.Sprintf("%s/academico/Notas.aspx", a.baseURL)
	
	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", gradesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add cookies to request
	for _, cookie := range sessionCookies {
		req.AddCookie(cookie)
	}

	// Set user agent to identify as a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	// Send request
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch grades: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response
	grades, err := a.parseGradesResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse grades: %w", err)
	}

	return grades, nil
}

// parseGradesResponse parses the HTML response from UACloud
func (a *UACloudGradesAdapter) parseGradesResponse(resp *http.Response) ([]grades.Grade, error) {
	// TODO: Implement actual HTML parsing
	// For now, return mock data for testing
	
	// This should be replaced with actual HTML parsing using a library like goquery
	// The parsing logic should extract grade information from tables or other HTML elements
	
	return []grades.Grade{
		{
			SubjectID:       "34041",
			SubjectName:     "GESTIÓN DE CALIDAD SOFTWARE",
			CurrentGrade:    "8.5",
			Average:         8.5,
			AssessmentCount:  2,
			Status:          grades.StatusApproved,
			LastUpdated:     time.Now().Add(-24 * time.Hour),
		},
		{
			SubjectID:       "34038",
			SubjectName:     "SEGURIDAD EN EL DISEÑO DE SOFTWARE",
			CurrentGrade:    "7.2",
			Average:         7.2,
			AssessmentCount:  3,
			Status:          grades.StatusApproved,
			LastUpdated:     time.Now().Add(-48 * time.Hour),
		},
		{
			SubjectID:       "34027",
			SubjectName:     "PLANIFICACIÓN Y PRUEBAS DE SISTEMAS SOFTWARE",
			CurrentGrade:    "No calificación",
			Average:         0.0,
			AssessmentCount:  0,
			Status:          grades.StatusPending,
			LastUpdated:     time.Now(),
		},
	}, nil
}

// parseJSONGrades parses JSON response if UACloud provides API endpoint
func (a *UACloudGradesAdapter) parseJSONGrades(data []byte) ([]grades.Grade, error) {
	// TODO: Implement JSON parsing if UACloud provides API
	// This would be more efficient than HTML parsing
	
	// Return empty slice for now - should be implemented when actual API is discovered
	return []grades.Grade{}, nil
}

// ValidateGradeData validates the parsed grade data
func (a *UACloudGradesAdapter) ValidateGradeData(grades []grades.Grade) error {
	for _, grade := range grades {
		if !grade.IsValid() {
			return fmt.Errorf("invalid grade data: subject_id=%s, subject_name=%s", 
				grade.SubjectID, grade.SubjectName)
		}
	}
	return nil
}

// GetGradesURL returns the URL used for fetching grades
func (a *UACloudGradesAdapter) GetGradesURL() string {
	return fmt.Sprintf("%s/academico/Notas.aspx", a.baseURL)
}