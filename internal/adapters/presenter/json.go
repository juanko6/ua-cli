package presenter

import (
	"encoding/json"
	"io"

	"github.com/juanko6/ua-cli/internal/domain/schedule"
)

// RenderJSON prints events as raw JSON.
func RenderJSON(w io.Writer, events []schedule.Event) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(events); err != nil {
		// En presenter no solemos loguear errores fatales si es stdout, pero
		// al menos documentamos o ignoramos explícitamente el return.
		_ = err
	}
}
