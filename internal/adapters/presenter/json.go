package presenter

import (
	"encoding/json"
	"io"

	"ua-cli/internal/domain/schedule"
)

// RenderJSON prints events as raw JSON.
func RenderJSON(w io.Writer, events []schedule.Event) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(events)
}
