package presenter

import (
	"fmt"
	"strings"
	"time"

	"github.com/juanko6/ua-cli/internal/service/grades"
	"github.com/charmbracelet/lipgloss"
)

// TextGradesPresenter presents grades in a text table format
type TextGradesPresenter struct {
	style gradesPresenterStyle
}

type gradesPresenterStyle struct {
	header    lipgloss.Style
	subject   lipgloss.Style
	grade     lipgloss.Style
	location  lipgloss.Style
	status    lipgloss.Style
	newGrade  lipgloss.Style
}

// NewTextGradesPresenter creates a new text grades presenter
func NewTextGradesPresenter() *TextGradesPresenter {
	return &TextGradesPresenter{
		style: gradesPresenterStyle{
			header:    lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")),
			subject:   lipgloss.NewStyle().Foreground(lipgloss.Color("15")),
			grade:     lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true),
			location:  lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
			status:    lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
			newGrade:  lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true),
		},
	}
}

// Present formats and displays grades in text format
func (p *TextGradesPresenter) Present(result *grades.GradesResult) string {
	if result == nil {
		return "❌ No hay datos de calificaciones disponibles."
	}

	var builder strings.Builder

	// Header
	builder.WriteString(p.style.header.Render("📊 TUS CALIFICACIONES ACTUALES") + "\n")
	builder.WriteString(p.style.header.Render(strings.Repeat("═", 50)) + "\n\n")

	// Summary
	if result.TotalSubjects > 0 {
		summary := fmt.Sprintf("📈 %d asignaturas | Promedio: %.2f | Última actualización: %s",
			result.TotalSubjects,
			result.AverageGrade,
			result.LastCheck.Format("02/01/2006 15:04"))
		if result.HasChanges {
			summary += fmt.Sprintf" | 🆕 %d nuevos cambios", len(result.NewGrades))
		}
		builder.WriteString(p.style.status.Render(summary) + "\n\n")
	}

	// Grades table
	if len(result.Grades) > 0 {
		builder.WriteString(p.renderGradesTable(result.Grades, result.NewGrades) + "\n\n")
	}

	// New grades section
	if len(result.NewGrades) > 0 {
		builder.WriteString(p.style.newGrade.Render("🆕 NUEVAS CALIFICACIONES") + "\n")
		builder.WriteString(p.style.newGrade.Render(strings.Repeat("─", 30)) + "\n\n")
		for _, grade := range result.NewGrades {
			builder.WriteString(p.renderNewGrade(grade) + "\n")
		}
		builder.WriteString("\n")
	}

	// Message
	if result.Message != "" {
		builder.WriteString("ℹ️ " + result.Message + "\n")
	}

	// Empty state
	if result.TotalSubjects == 0 {
		builder.WriteString("📋 No hay calificaciones disponibles en este momento.\n")
		builder.WriteString("🔍 Esto puede deberse a que aún no se han publicado notas o no estás inscrito en asignaturas.\n")
	}

	return builder.String()
}

// renderGradesTable creates a formatted table of grades
func (p *TextGradesPresenter) renderGradesTable(grades []grades.Grade, newGrades []grades.Grade) string {
	var builder strings.Builder

	// Table header
	builder.WriteString(p.style.header.Render(fmt.Sprintf("%-40s %-8s %-8s %-12s", 
		"ASIGNATURA", "CALIFICACIÓN", "PROMEDIO", "ESTADO")) + "\n")
	builder.WriteString(p.style.header.Render(strings.Repeat("─", 70)) + "\n")

	// Create lookup for new grades
	newGradeMap := make(map[string]bool)
	for _, grade := range newGrades {
		newGradeMap[grade.SubjectID] = true
	}

	// Table rows
	for _, grade := range grades {
		isNew := newGradeMap[grade.SubjectID]
		subject := p.truncateString(grade.SubjectName, 38)
		currentGrade := p.formatGrade(grade.CurrentGrade)
		average := fmt.Sprintf("%.2f", grade.Average)
		status := p.formatStatus(grade.Status, isNew)

		builder.WriteString(fmt.Sprintf("%-40s %-8s %-8s %-12s\n",
			p.style.subject.Render(subject),
			p.style.grade.Render(currentGrade),
			average,
			status))
	}

	return builder.String()
}

// renderNewGrade formats a single new grade notification
func (p *TextGradesPresenter) renderNewGrade(grade grades.Grade) string {
	return fmt.Sprintf("📚 %s: %s (Promedio: %.2f) - %s",
		p.style.subject.Render(grade.SubjectName),
		p.style.grade.Render(grade.CurrentGrade),
		grade.Average,
		p.formatStatus(grade.Status, true))
}

// formatGrade formats a grade for display
func (p *TextGradesPresenter) formatGrade(grade string) string {
	if grade == "No calificación" || grade == "" {
		return "N/A"
	}
	return grade
}

// formatStatus formats the grade status with appropriate styling
func (p *TextGradesPresenter) formatStatus(status grades.GradeStatus, isNew bool) string {
	var baseStyle lipgloss.Style
	var emoji string

	switch status {
	case grades.StatusApproved:
		baseStyle = p.style.status
		emoji = "✅"
	case grades.StatusPending:
		baseStyle = p.style.status
		emoji = "⏳"
	case grades.StatusNeedsAttention:
		baseStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")) // Red
		emoji = "⚠️"
	default:
		baseStyle = p.style.status
		emoji = "❓"
	}

	text := baseStyle.Render(string(status))
	if isNew {
		text = fmt.Sprintf("%s %s", p.style.newGrade.Render("🆕"), text)
	}

	return emoji + " " + text
}

// truncateString truncates a string to fit the specified width
func (p *TextGradesPresenter) truncateString(s string, width int) string {
	if len(s) <= width {
		return s
	}
	return s[:width-3] + "..."
}

// formatTimestamp formats a timestamp for display
func (p *TextGradesPresenter) formatTimestamp(t time.Time) string {
	return t.Format("02/01/2006 15:04")
}