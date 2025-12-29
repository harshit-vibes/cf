package health

import (
	"testing"
)

func TestStatus_String(t *testing.T) {
	tests := []struct {
		status Status
		want   string
	}{
		{StatusHealthy, "healthy"},
		{StatusDegraded, "degraded"},
		{StatusCritical, "critical"},
		{Status(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("Status.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAction_String(t *testing.T) {
	tests := []struct {
		action Action
		want   string
	}{
		{ActionNone, "none"},
		{ActionAutoFix, "autofix"},
		{ActionUserPrompt, "prompt"},
		{ActionManualFix, "manual"},
		{ActionRetry, "retry"},
		{ActionFatal, "fatal"},
		{Action(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.action.String(); got != tt.want {
				t.Errorf("Action.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusConstants(t *testing.T) {
	// Verify constants are defined correctly
	if StatusHealthy != 0 {
		t.Errorf("StatusHealthy = %v, want 0", StatusHealthy)
	}
	if StatusDegraded != 1 {
		t.Errorf("StatusDegraded = %v, want 1", StatusDegraded)
	}
	if StatusCritical != 2 {
		t.Errorf("StatusCritical = %v, want 2", StatusCritical)
	}
}

func TestActionConstants(t *testing.T) {
	// Verify constants are defined correctly
	if ActionNone != 0 {
		t.Errorf("ActionNone = %v, want 0", ActionNone)
	}
	if ActionAutoFix != 1 {
		t.Errorf("ActionAutoFix = %v, want 1", ActionAutoFix)
	}
	if ActionUserPrompt != 2 {
		t.Errorf("ActionUserPrompt = %v, want 2", ActionUserPrompt)
	}
	if ActionManualFix != 3 {
		t.Errorf("ActionManualFix = %v, want 3", ActionManualFix)
	}
	if ActionRetry != 4 {
		t.Errorf("ActionRetry = %v, want 4", ActionRetry)
	}
	if ActionFatal != 5 {
		t.Errorf("ActionFatal = %v, want 5", ActionFatal)
	}
}

func TestResult_Fields(t *testing.T) {
	result := Result{
		Name:        "Test Check",
		Category:    "internal",
		Status:      StatusHealthy,
		Message:     "All good",
		Details:     "No issues found",
		Recoverable: true,
		Action:      ActionNone,
	}

	if result.Name != "Test Check" {
		t.Errorf("Result.Name = %v, want %v", result.Name, "Test Check")
	}
	if result.Category != "internal" {
		t.Errorf("Result.Category = %v, want %v", result.Category, "internal")
	}
	if result.Status != StatusHealthy {
		t.Errorf("Result.Status = %v, want %v", result.Status, StatusHealthy)
	}
	if !result.Recoverable {
		t.Error("Result.Recoverable should be true")
	}
}

func TestReport_Fields(t *testing.T) {
	report := Report{
		OverallStatus:         StatusHealthy,
		CanProceed:            true,
		CurrentSchemaVersion:  "1.0.0",
		DetectedSchemaVersion: "1.0.0",
		NeedsMigration:        false,
		Warnings:              []string{"warning1"},
		Errors:                []string{},
	}

	if report.OverallStatus != StatusHealthy {
		t.Errorf("Report.OverallStatus = %v, want %v", report.OverallStatus, StatusHealthy)
	}
	if !report.CanProceed {
		t.Error("Report.CanProceed should be true")
	}
	if len(report.Warnings) != 1 {
		t.Errorf("len(Report.Warnings) = %v, want %v", len(report.Warnings), 1)
	}
	if len(report.Errors) != 0 {
		t.Errorf("len(Report.Errors) = %v, want %v", len(report.Errors), 0)
	}
}
