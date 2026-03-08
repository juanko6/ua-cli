package main

import (
	"strings"
	"testing"
	"time"

	"github.com/juanko6/ua-cli/internal/domain/grades"
)

func TestGradesCommandStructure(t *testing.T) {
	// Test that grades command exists and has correct flags
	t.Run("command_exists", func(t *testing.T) {
		if gradesCmd == nil {
			t.Fatal("grades command is nil")
		}
		
		expectedUse := "grades"
		if gradesCmd.Use != expectedUse {
			t.Errorf("expected use %s, got %s", expectedUse, gradesCmd.Use)
		}
		
		expectedShort := "Show your current grades across all subjects"
		if gradesCmd.Short != expectedShort {
			t.Errorf("expected short %s, got %s", expectedShort, gradesCmd.Short)
		}
	})
	
	t.Run("has_json_flag", func(t *testing.T) {
		found := false
		for _, flag := range gradesCmd.LocalFlags().FlagNames() {
			if flag == "json" || flag == "j" {
				found = true
				break
			}
		}
		if !found {
			t.Error("json flag not found")
		}
	})
	
	t.Run("has_status_flags", func(t *testing.T) {
		requiredFlags := []string{"approved", "pending", "attention"}
		for _, flagName := range requiredFlags {
			found := false
			for _, flag := range gradesCmd.LocalFlags().FlagNames() {
				if flag == flagName {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("%s flag not found", flagName)
			}
		}
	})
}

func TestGradeEntityValidation(t *testing.T) {
	t.Run("valid_grade", func(t *testing.T) {
		grade := grades.Grade{
			SubjectID:   "34041",
			SubjectName: "Test Subject",
			CurrentGrade: "8.5",
			Average:     8.5,
			Status:      grades.StatusApproved,
			LastUpdated: time.Now(),
		}
		
		if !grade.IsValid() {
			t.Error("expected grade to be valid")
		}
		
		if !grade.IsApproved() {
			t.Error("expected grade to be approved")
		}
		
		if grade.IsPending() {
			t.Error("expected grade not to be pending")
		}
	})
	
	t.Run("invalid_grade_empty_id", func(t *testing.T) {
		grade := grades.Grade{
			SubjectID:   "",
			SubjectName: "Test Subject",
			CurrentGrade: "8.5",
			Average:     8.5,
			Status:      grades.StatusApproved,
			LastUpdated: time.Now(),
		}
		
		if grade.IsValid() {
			t.Error("expected grade with empty ID to be invalid")
		}
	})
	
	t.Run("invalid_grade_empty_name", func(t *testing.T) {
		grade := grades.Grade{
			SubjectID:   "34041",
			SubjectName: "",
			CurrentGrade: "8.5",
			Average:     8.5,
			Status:      grades.StatusApproved,
			LastUpdated: time.Now(),
		}
		
		if grade.IsValid() {
			t.Error("expected grade with empty name to be invalid")
		}
	})
	
	t.Run("invalid_grade_empty_status", func(t *testing.T) {
		grade := grades.Grade{
			SubjectID:   "34041",
			SubjectName: "Test Subject",
			CurrentGrade: "8.5",
			Average:     8.5,
			Status:      "",
			LastUpdated: time.Now(),
		}
		
		if grade.IsValid() {
			t.Error("expected grade with empty status to be invalid")
		}
	})
}

func TestGradeStatusMethods(t *testing.T) {
	t.Run("status_methods", func(t *testing.T) {
		approved := grades.Grade{Status: grades.StatusApproved}
		pending := grades.Grade{Status: grades.StatusPending}
		attention := grades.Grade{Status: grades.StatusNeedsAttention}
		
		if !approved.IsApproved() {
			t.Error("approved grade should return IsApproved() = true")
		}
		if approved.IsPending() {
			t.Error("approved grade should return IsPending() = false")
		}
		if approved.NeedsAttention() {
			t.Error("approved grade should return NeedsAttention() = false")
		}
		
		if !pending.IsPending() {
			t.Error("pending grade should return IsPending() = true")
		}
		if pending.IsApproved() {
			t.Error("pending grade should return IsApproved() = false")
		}
		if pending.NeedsAttention() {
			t.Error("pending grade should return NeedsAttention() = false")
		}
		
		if !attention.NeedsAttention() {
			t.Error("attention grade should return NeedsAttention() = true")
		}
		if attention.IsApproved() {
			t.Error("attention grade should return IsApproved() = false")
		}
		if attention.IsPending() {
			t.Error("attention grade should return IsPending() = false")
		}
	})
	
	t.Run("status_emojis", func(t *testing.T) {
		testCases := []struct {
			status grades.GradeStatus
			emoji  string
		}{
			{grades.StatusApproved, "✅"},
			{grades.StatusPending, "⏳"},
			{grades.StatusNeedsAttention, "⚠️"},
		}
		
		for _, tc := range testCases {
			if tc.status.DisplayEmoji() != tc.emoji {
				t.Errorf("expected emoji %s for status %s, got %s", 
					tc.emoji, tc.status, tc.status.DisplayEmoji())
			}
		}
	})
}

func TestGradesCommandHelp(t *testing.T) {
	t.Run("help_contains_examples", func(t *testing.T) {
		if gradesCmd.Long == "" {
			t.Error("command long description is empty")
		}
		
		requiredExamples := []string{
			"ua grades",
			"ua grades --json",
			"ua grades --approved",
		}
		
		for _, example := range requiredExamples {
			if !strings.Contains(gradesCmd.Long, example) {
				t.Errorf("help text missing example: %s", example)
			}
		}
	})
}