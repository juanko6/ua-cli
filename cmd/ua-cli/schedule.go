package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
    "golang.org/x/term"

	"ua-cli/internal/adapters/presenter"
	"ua-cli/internal/adapters/repo"
	"ua-cli/internal/adapters/uacloud"
	"ua-cli/internal/service/schedule"
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

		// 1. Setup Dependencies
        // Use user home dir for cache
        home, _ := os.UserHomeDir()
        cachePath := fmt.Sprintf("%s/.ua-cli/cache/schedule.json", home)

		repoAdapter, err := repo.NewJSONFileRepo(cachePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing cache: %v\n", err)
			os.Exit(1)
		}

		// TODO: Load cookie from Keyring or Config
		// For MVP, if no cookie, we might fail or prompt.
		// Detailed auth logic is in a separate step, here we assume empty or configured.
		cookieVal := "TODO_LOAD_COOKIE" 
		cloudAdapter, err := uacloud.NewUACloudAdapter(cookieVal)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing cloud adapter: %v\n", err)
			os.Exit(1)
		}

		svc := service.NewScheduleService(repoAdapter, cloudAdapter)

		// 2. Calculate Date
		targetDate := time.Now()
		if scheduleNextWeek {
			targetDate = targetDate.AddDate(0, 0, 7)
		} else if schedulePrevWeek {
			targetDate = targetDate.AddDate(0, 0, -7)
		}

		// 3. Execute Service
		events, err := svc.GetScheduleForWeek(ctx, targetDate, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching schedule: %v\n", err)
			os.Exit(1)
		}

		// 4. Adaptive UI
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
