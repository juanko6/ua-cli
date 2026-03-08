package presenter

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/juanko6/ua-cli/internal/domain/schedule"
)

// RenderTextTable prints events in a tabular format compatible with grep/awk.
func RenderTextTable(w io.Writer, events []schedule.Event) {
	if len(events) == 0 {
		fmt.Fprintln(w, "No classes scheduled for this week.")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "DAY\tTIME\tSUBJECT\tLOCATION\tTYPE")

	for _, e := range events {
		day := e.Start.Weekday().String()[0:3]
		timeRange := fmt.Sprintf("%s-%s", e.Start.Format("15:04"), e.End.Format("15:04"))
		// Basic shortening of title might be needed
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", day, timeRange, e.Title, e.Location, e.Type)
	}
	tw.Flush()
}
