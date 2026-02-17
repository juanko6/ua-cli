package uacloud

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
    "net/url"
    "strings"
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

    // Parse raw cookie string "key=value; key2=value2"
    baseURLStr := "https://cvnet.cpd.ua.es"
    u, _ := url.Parse(baseURLStr)
    
    var cookies []*http.Cookie
    rawCookies := strings.Split(cookieValue, ";")
    for _, raw := range rawCookies {
        parts := strings.SplitN(strings.TrimSpace(raw), "=", 2)
        if len(parts) == 2 {
            cookies = append(cookies, &http.Cookie{
                Name:  parts[0],
                Value: parts[1],
            })
        }
    }
    jar.SetCookies(u, cookies)

	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	return &UACloudAdapter{
		client:  client,
		baseURL: baseURLStr,
	}, nil
}

func (a *UACloudAdapter) GetWeekEvents(ctx context.Context, date time.Time) ([]schedule.Event, error) {
	// Calculate week boundaries (Monday 00:00 → Sunday 23:59) in Unix timestamps
	weekday := date.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	monday := date.AddDate(0, 0, -int(weekday-time.Monday))
	monday = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, date.Location())
	sunday := monday.AddDate(0, 0, 7) // Next Monday 00:00 as exclusive end

	startUnix := monday.Unix()
	endUnix := sunday.Unix()

	// Real UACloud JSON API
	url := fmt.Sprintf(
		"%s/uaHorarios/Home/ObtenerEventosCalendarioJson?calendario=docenciaalu&start=%d&end=%d",
		a.baseURL, startUnix, endUnix,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) ua-cli/1.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("UACloud API returned status %d: %s", resp.StatusCode, string(body[:min(len(body), 200)]))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return ParseSchedule(body)
}

func (a *UACloudAdapter) SaveEvents(ctx context.Context, events []schedule.Event) error {
	// This adapter is read-only (Fetch from Cloud)
	return fmt.Errorf("UACloud adapter is read-only")
}
