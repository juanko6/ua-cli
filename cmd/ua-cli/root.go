package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	adaptauth "ua-cli/internal/adapters/auth"
	domainauth "ua-cli/internal/domain/auth"
	authservice "ua-cli/internal/service/auth"
)

var rootCmd = &cobra.Command{
	Use:   "ua",
	Short: "ua-cli - University of Alicante CLI",
	Long:  `Fast, privacy-focused CLI for University of Alicante services (Schedule, Grades, Campus). Includes smart login with session checks.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if isAuthExempt(cmd) {
			return nil
		}
		cookiePath, err := defaultCookiePathRoot()
		if err != nil {
			return err
		}
		store := adaptauth.NewFileCredentialStore(cookiePath)
		svc := authservice.NewAuthService(store, nil)
		session, err := svc.CheckSession()
		if err != nil && session == nil {
			return err
		}
		if session == nil {
			return fmt.Errorf("no active session. run 'ua login' to authenticate")
		}
		svc.WarnIfAged(session)

		switch session.Status {
		case domainauth.SessionMissing:
			return fmt.Errorf("no active session. run 'ua login' to authenticate")
		case domainauth.SessionExpired:
			return fmt.Errorf("session expired (%s old). run 'ua login' to refresh", session.Age.Round(time.Minute))
		case domainauth.SessionInvalid:
			return fmt.Errorf("invalid session cookie. run 'ua login' to re-authenticate")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func isAuthExempt(cmd *cobra.Command) bool {
	exempt := map[string]bool{
		"login":      true,
		"help":       true,
		"version":    true,
		"completion": true,
	}
	for c := cmd; c != nil && c.Name() != "ua"; c = c.Parent() {
		if exempt[c.Name()] {
			return true
		}
	}
	return false
}

func defaultCookiePathRoot() (string, error) {
	return adaptauth.DefaultCookiePath()
}
