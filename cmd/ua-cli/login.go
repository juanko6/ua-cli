package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	adaptauth "github.com/juanko6/ua-cli/internal/adapters/auth"
	domainauth "github.com/juanko6/ua-cli/internal/domain/auth"
	authservice "github.com/juanko6/ua-cli/internal/service/auth"
)

var (
	loginManual bool
	loginStatus bool
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
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

		var cookie string
		if loginManual {
			cookie, err = promptCookieManual()
		} else {
			cookie, err = loginViaProxy()
		}
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

// loginViaProxy starts a local reverse proxy, opens the browser, and captures
// authentication cookies automatically after the user completes CAS login.
func loginViaProxy() (string, error) {
	proxy := &adaptauth.ProxyAdapter{}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Start the proxy in the background so we can get the port and open the browser.
	cookieCh := make(chan string, 1)
	errCh := make(chan error, 1)
	go func() {
		cookie, err := proxy.Capture(ctx)
		if err != nil {
			errCh <- err
			return
		}
		cookieCh <- cookie
	}()

	// Give the proxy a moment to bind its port.
	time.Sleep(200 * time.Millisecond)

	proxyURL := proxy.ProxyURL()
	fmt.Println(infoStyle.Render("🔐 Starting login proxy..."))
	fmt.Println(infoStyle.Render("Opening UACloud in your browser..."))

	if err := adaptauth.OpenInBrowser(proxyURL); err != nil {
		fmt.Fprintf(os.Stderr, "%s %v\n",
			dimStyle.Render("Could not open browser automatically:"), err)
		fmt.Println(infoStyle.Render("Open this URL manually: ") + proxyURL)
	}

	fmt.Println(dimStyle.Render("Waiting for CAS login (timeout 2 min)..."))

	select {
	case cookie := <-cookieCh:
		return cookie, nil
	case err := <-errCh:
		return "", fmt.Errorf(errorStyle.Render("✗ Login failed:")+" %w", err)
	}
}

// promptCookieManual asks the user to paste a cookie string manually.
func promptCookieManual() (string, error) {
	fmt.Println(infoStyle.Render("Manual mode enabled. Paste your UACloud cookie below."))
	fmt.Println(dimStyle.Render("(Use this if the automatic browser login doesn't work)"))

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
