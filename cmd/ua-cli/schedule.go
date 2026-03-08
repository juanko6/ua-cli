package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	adaptauth "github.com/juanko6/ua-cli/internal/adapters/auth"
	"github.com/juanko6/ua-cli/internal/adapters/presenter"
	"github.com/juanko6/ua-cli/internal/adapters/repo"
	"github.com/juanko6/ua-cli/internal/adapters/uacloud"
	service "github.com/juanko6/ua-cli/internal/service/schedule"
)

var (
	scheduleNextWeek bool
	schedulePrevWeek bool
	scheduleJSON     bool
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "View your weekly class schedule",
	Long: `Displays the schedule for the current week.
Use --next/--prev to navigate weeks.
Use --json for raw output.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		home, _ := os.UserHomeDir()

		cachePath := filepath.Join(home, ".ua-cli", "cache", "schedule.json")
		repoAdapter, err := repo.NewJSONFileRepo(cachePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing cache: %v\n", err)
			os.Exit(1)
		}

		cookiePath, err := adaptauth.DefaultCookiePath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error resolving cookie path: %v\n", err)
			os.Exit(1)
		}
		cookieStore := adaptauth.NewFileCredentialStore(cookiePath)
		cookieVal, err := cookieStore.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading cookie: %v\n", err)
			os.Exit(1)
		}

		cloudAdapter, err := uacloud.NewUACloudAdapter(cookieVal)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing cloud adapter: %v\n", err)
			os.Exit(1)
		}

		svc := service.NewScheduleService(repoAdapter, cloudAdapter)

		targetDate := time.Now()
		if scheduleNextWeek {
			targetDate = targetDate.AddDate(0, 0, 7)
		} else if schedulePrevWeek {
			targetDate = targetDate.AddDate(0, 0, -7)
		}

		events, err := svc.GetScheduleForWeek(ctx, targetDate, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching schedule: %v\n", err)
			os.Exit(1)
		}

		isTTY := term.IsTerminal(int(os.Stdout.Fd()))
		if scheduleJSON {
			presenter.RenderJSON(os.Stdout, events)
		} else if isTTY {
			presenter.RenderTUI(events)
		} else {
			presenter.RenderTextTable(os.Stdout, events)
		}
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)

	scheduleCmd.Flags().BoolVarP(&scheduleNextWeek, "next", "n", false, "Show next week's schedule")
	scheduleCmd.Flags().BoolVarP(&schedulePrevWeek, "prev", "p", false, "Show previous week's schedule")
	scheduleCmd.Flags().BoolVar(&scheduleJSON, "json", false, "Output in JSON format")
}
