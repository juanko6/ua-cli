package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	adaptauth "ua-cli/internal/adapters/auth"
	domainauth "ua-cli/internal/domain/auth"
	authservice "ua-cli/internal/service/auth"
)

var (
	loginManual bool
	loginStatus bool
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with UACloud",
	Long:  "Opens UACloud login flow and securely stores your session cookie for future commands.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cookiePath, err := defaultCookiePath()
		if err != nil {
			return err
		}
		store := adaptauth.NewFileCredentialStore(cookiePath)
		svc := authservice.NewAuthService(store, adaptauth.ValidateCookie)

		if loginStatus {
			return printSessionStatus(svc)
		}

		cookie, err := promptCookie(loginManual)
		if err != nil {
			return err
		}
		if err := svc.StoreCookie(cookie); err != nil {
			return fmt.Errorf(errorStyle.Render("✗ Login failed:")+" %w", err)
		}

		fmt.Println(successStyle.Render("✓ Login successful. Session cookie saved."))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().BoolVar(&loginManual, "cookie", false, "Manual mode: paste cookie directly without opening browser")
	loginCmd.Flags().BoolVar(&loginStatus, "status", false, "Check saved session status and exit")
}

func defaultCookiePath() (string, error) {
	return adaptauth.DefaultCookiePath()
}

func promptCookie(manual bool) (string, error) {
	if !manual {
		fmt.Println(infoStyle.Render("Opening UACloud in your browser..."))
		if err := adaptauth.OpenInBrowser("https://cvnet.cpd.ua.es/uaCloud"); err != nil {
			fmt.Fprintf(os.Stderr, "Could not open browser automatically: %v\n", err)
		}
		fmt.Println(infoStyle.Render("After signing in, copy your Cookie header value and paste it below."))
	} else {
		fmt.Println(infoStyle.Render("Manual mode enabled. Paste your UACloud cookie below."))
	}

	fmt.Print("Cookie (timeout 60s): ")
	input := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			errCh <- err
			return
		}
		input <- strings.TrimSpace(line)
	}()

	select {
	case err := <-errCh:
		return "", err
	case val := <-input:
		if val == "" {
			return "", fmt.Errorf("empty cookie input")
		}
		return val, nil
	case <-time.After(60 * time.Second):
		return "", fmt.Errorf("input timeout after 60 seconds")
	}
}

func printSessionStatus(svc *authservice.AuthService) error {
	s, err := svc.CheckSession()
	if err != nil && s == nil {
		return err
	}
	if s == nil {
		return fmt.Errorf("unable to determine session status")
	}

	formatAge := func(age time.Duration) string {
		if age <= 0 {
			return "unknown"
		}
		return age.Round(time.Minute).String()
	}

	switch s.Status {
	case domainauth.SessionMissing:
		fmt.Println(errorStyle.Render("✗ Session missing") + " — run 'ua login' to authenticate.")
	case domainauth.SessionExpired:
		fmt.Println(errorStyle.Render("✗ Session expired") + " — run 'ua login' to refresh.")
		fmt.Println("Age:", formatAge(s.Age))
	case domainauth.SessionInvalid:
		fmt.Println(errorStyle.Render("✗ Session invalid") + " — check cookie file permissions/content.")
		if s.Path != "" {
			fmt.Println("Path:", s.Path)
		}
		fmt.Println("Age:", formatAge(s.Age))
	case domainauth.SessionValid:
		fmt.Println(successStyle.Render("✓ Session active"))
		fmt.Println("Age:", formatAge(s.Age))
	default:
		fmt.Println("Session status:", s.Status)
	}
	return err
}
