package uacloud

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"

	"ua-cli/internal/domain/schedule"
)

type UACloudAdapter struct {
	client  *http.Client
	baseURL string
}

func NewUACloudAdapter(cookieValue string) (*UACloudAdapter, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	// Set initial cookie (manual login)
	// In a real app, this should be more robust
	// For now, we assume basic cookie auth via configuration
	return &UACloudAdapter{
		client:  client,
		baseURL: "https://cvnet.cpd.ua.es",
	}, nil
}

func (a *UACloudAdapter) GetWeekEvents(ctx context.Context, date time.Time) ([]schedule.Event, error) {
	// Calculate week start/end based on date
	// UACloud specific logic to fetch 'horario'
	// For MVP, we'll implement a basic fetch that parses the HTML/JSON response
	// The URL construction logic is complex and needs to be correct.
	// Example URL: https://cvnet.cpd.ua.es/uaip/horario/index
	
	url := fmt.Sprintf("%s/uaip/horario/index", a.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add necessary headers
	req.Header.Add("User-Agent", "ua-cli/1.0")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("UACloud API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response using our parser
	return ParseSchedule(body)
}

func (a *UACloudAdapter) SaveEvents(ctx context.Context, events []schedule.Event) error {
	// This adapter is read-only (Fetch from Cloud)
	return fmt.Errorf("UACloud adapter is read-only")
}
