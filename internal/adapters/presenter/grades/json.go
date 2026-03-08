package presenter

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/juanko6/ua-cli/internal/service/grades"
)

// JSONGradesPresenter presents grades in JSON format
type JSONGradesPresenter struct{}

// NewJSONGradesPresenter creates a new JSON grades presenter
func NewJSONGradesPresenter() *JSONGradesPresenter {
	return &JSONGradesPresenter{}
}

// Present formats and displays grades in JSON format
func (p *JSONGradesPresenter) Present(result *grades.GradesResult) string {
	if result == nil {
		jsonResult := map[string]interface{}{
			"error": "No hay datos de calificaciones disponibles",
		}
		data, _ := json.MarshalIndent(jsonResult, "", "  ")
		return string(data)
	}

	// Build structured result for JSON output
	jsonResult := map[string]interface{}{
		"summary": map[string]interface{}{
			"total_subjects":   result.TotalSubjects,
			"average_grade":    result.AverageGrade,
			"last_check":       result.LastCheck.Format(time.RFC3339),
			"has_changes":     result.HasChanges,
			"new_grades_count": len(result.NewGrades),
		},
		"grades":    result.Grades,
		"new_grades": result.NewGrades,
		"message":   result.Message,
	}

	// Marshal with indentation for readability
	data, err := json.MarshalIndent(jsonResult, "", "  ")
	if err != nil {
		// Fallback to simple JSON if marshaling fails
		jsonResult["error"] = fmt.Sprintf("Error al formatear JSON: %v", err)
		data, _ = json.Marshal(jsonResult)
	}

	return string(data)
}