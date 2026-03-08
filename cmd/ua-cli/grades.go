package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/juanko6/ua-cli/internal/adapters/presenter/grades"
	"github.com/juanko6/ua-cli/internal/adapters/repo"
	"github.com/juanko6/ua-cli/internal/adapters/uacloud"
	"github.com/juanko6/ua-cli/internal/domain/grades"
	"github.com/juanko6/ua-cli/internal/service/grades"
	"github.com/spf13/cobra"
)

var (
	showJSON     bool
	showApproved bool
	showPending  bool
	showAttention bool
)

var gradesCmd = &cobra.Command{
	Use:   "grades",
	Short: "Show your current grades across all subjects",
	Long: `Show your current academic grades across all enrolled subjects.
	
This command fetches your latest grades from UACloud and displays them in a
formatted table. It can also show newly published grades since your last check.

Examples:
  ua grades                    # Show all grades
  ua grades --json             # Show grades in JSON format
  ua grades --approved        # Show only approved grades
  ua grades --pending          # Show only pending grades`,
	RunE: runGradesCmd,
}

func init() {
	rootCmd.AddCommand(gradesCmd)

	// Output format flags
	gradesCmd.Flags().BoolVarP(&showJSON, "json", "j", false, "Output in JSON format")
	gradesCmd.Flags().BoolVar(&showApproved, "approved", false, "Show only approved subjects")
	gradesCmd.Flags().BoolVar(&showPending, "pending", false, "Show only pending subjects")
	gradesCmd.Flags().BoolVar(&showAttention, "attention", false, "Show only subjects needing attention")
}

func runGradesCmd(cmd *cobra.Command, args []string) error {
	// Get authentication service
	authService, err := getAuthService()
	if err != nil {
		return fmt.Errorf("failed to get auth service: %w", err)
	}

	// Check if user is authenticated
	session, err := authService.GetSession()
	if err != nil {
		return fmt.Errorf("authentication required: %w", err)
	}

	if session.IsEmpty() {
		return fmt.Errorf("no active session. Please run 'ua login' first")
	}

	// Create grades service components
	storagePath := filepath.Join(getConfigDir(), "grades.json")
	adapter := uacloud.NewUACloudGradesAdapter(getUACloudBaseURL())
	repository := repo.NewUACloudGradesRepository(storagePath, adapter)

	// Determine presenter
	var presenter grades.GradesPresenter
	if showJSON {
		presenter = grades.NewJSONGradesPresenter()
	} else {
		presenter = grades.NewTextGradesPresenter()
	}

	// Create grades service
	gradeService := grades.NewGradesService(repository, nil, presenter)

	// Prepare options
	opts := grades.GradesOptions{
		JSONOutput:    showJSON,
		ShowApproved:  showApproved,
		ShowPending:   showPending,
		ShowAttention: showAttention,
	}

	// Fetch and display grades
	result, err := gradeService.GetGrades(cmd.Context(), session.Cookies(), opts)
	if err != nil {
		return fmt.Errorf("failed to get grades: %w", err)
	}

	// Output the result
	output := presenter.Present(result)
	fmt.Println(output)

	return nil
}

// Helper functions (these should be imported from existing code)
func getConfigDir() string {
	// This should come from the existing root command configuration
	// For now, return a default
	return filepath.Join(os.Getenv("HOME"), ".ua-cli")
}

func getUACloudBaseURL() string {
	// This should come from configuration
	return "https://ucloud.ua.es"
}

// getAuthService should return the existing authentication service
// This is a placeholder - should be imported from existing auth code
func getAuthService() (interface{}, error) {
	// Return a mock for now - should be replaced with actual auth service
	return &mockAuthService{}, nil
}

// mockAuthService is a placeholder for the real auth service
type mockAuthService struct{}

func (m *mockAuthService) GetSession() (*mockSession, error) {
	return &mockSession{
		cookies: []*http.Cookie{
			{Name: ".ASPXFORMSAUTHUACLOUD", Value: "mock-cookie"},
		},
	}, nil
}

type mockSession struct {
	cookies []*http.Cookie
}

func (m *mockSession) IsEmpty() bool {
	return len(m.cookies) == 0
}

func (m *mockSession) Cookies() []*http.Cookie {
	return m.cookies
}