package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	adaptauth "ua-cli/internal/adapters/auth"
	"ua-cli/internal/adapters/repo"
	"ua-cli/internal/adapters/uacloud"
	"ua-cli/internal/domain/schedule"
	service "ua-cli/internal/service/schedule"
)

var nowJSON bool

var nowCmd = &cobra.Command{
	Use:   "now",
	Short: "Show what class you have right now or next",
	Long:  `Displays your current ongoing class or the immediate next one for today.`,
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
			// If not logged in, we can't reliably fetch a new schedule, but we could try cache.
			// Let's enforce login for now as the ScheduleService currently expects the cloudAdapter anyway.
			fmt.Fprintf(os.Stderr, "No valid session found. Please run 'ua login' first.\nError: %v\n", err)
			os.Exit(1)
		}

		cloudAdapter, err := uacloud.NewUACloudAdapter(cookieVal)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing cloud adapter: %v\n", err)
			os.Exit(1)
		}

		svc := service.NewScheduleService(repoAdapter, cloudAdapter)

		// Fetch schedule for current week
		now := time.Now()
		events, err := svc.GetScheduleForWeek(ctx, now, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching schedule: %v\n", err)
			os.Exit(1)
		}

		// Filter for today's events only
		var todaysEvents []schedule.Event
		for _, e := range events {
			y1, m1, d1 := e.Start.Date()
			y2, m2, d2 := now.Date()
			if y1 == y2 && m1 == m2 && d1 == d2 {
				todaysEvents = append(todaysEvents, e)
			}
		}

		// Sort events by chronological start time
		sort.Slice(todaysEvents, func(i, j int) bool {
			return todaysEvents[i].Start.Before(todaysEvents[j].Start)
		})

		var activeEvent *schedule.Event

		// Logic for deciding which class to show
		for i, e := range todaysEvents {
			if now.Before(e.Start) {
				// We haven't reached this class yet. It's the chronological next.
				activeEvent = &todaysEvents[i]
				break
			} else if now.After(e.Start) && now.Before(e.End) {
				// We are currently IN this class.
				timeLeft := e.End.Sub(now)
				if timeLeft <= 30*time.Minute {
					// 30 min threshold rule: look ahead to the next class
					if i+1 < len(todaysEvents) {
						activeEvent = &todaysEvents[i+1]
					} else {
						// No more classes today. Break and leave activeEvent as nil.
						break
					}
				} else {
					// Plentiful time remaining in current class.
					activeEvent = &todaysEvents[i]
				}
				break
			}
		}

		isTTY := term.IsTerminal(int(os.Stdout.Fd()))

		if activeEvent == nil {
			if nowJSON {
				fmt.Println("[]")
			} else {
				msgStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
				fmt.Println(msgStyle.Render("🎉 No tienes más clases por hoy. ¡A descansar!"))
			}
			return
		}

		if nowJSON {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode([]schedule.Event{*activeEvent}); err != nil {
				fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
			}
			return
		}

		if isTTY {
			renderNowCard(*activeEvent, now)
		} else {
			fmt.Printf("%s\t%s\t%s\t%s\n", activeEvent.Start.Format("15:04"), activeEvent.End.Format("15:04"), activeEvent.Location, activeEvent.Title)
		}
	},
}

func init() {
	rootCmd.AddCommand(nowCmd)
	nowCmd.Flags().BoolVar(&nowJSON, "json", false, "Output in JSON format")
}

func renderNowCard(e schedule.Event, now time.Time) {
	status := "PRÓXIMA"
	timeColor := "12" // light blue

	if now.After(e.Start) && now.Before(e.End) {
		status = "AHORA MISMO"
		timeColor = "10" // green
	}

	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(timeColor)).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	locStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("220"))

	timeRange := fmt.Sprintf("%s - %s", e.Start.Format("15:04"), e.End.Format("15:04"))

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 4).
		Render(
			fmt.Sprintf("%s\n\n%s\n%s %s\n%s %s\n%s %s",
				statusStyle.Render("● "+status),
				titleStyle.Render(e.Title),
				dimStyle.Render("🕒 Horario:"), timeRange,
				dimStyle.Render("📍 Aula:   "), locStyle.Render(e.Location),
				dimStyle.Render("📚 Tipo:   "), string(e.Type),
			),
		)

	fmt.Println(box)
}
