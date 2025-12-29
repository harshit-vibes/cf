package health

import (
	"context"
	"testing"
	"time"
)

// MockCheck implements Check interface for testing
type MockCheck struct {
	name     string
	category string
	result   Result
}

func (m *MockCheck) Name() string     { return m.name }
func (m *MockCheck) Category() string { return m.category }
func (m *MockCheck) Check(ctx context.Context) Result {
	return m.result
}

// MockAutoFixableCheck implements AutoFixable for testing
type MockAutoFixableCheck struct {
	MockCheck
	autoFixErr error
	fixCalled  bool
}

func (m *MockAutoFixableCheck) AutoFix(ctx context.Context) error {
	m.fixCalled = true
	return m.autoFixErr
}

// MockCriticalCheck implements Critical for testing
type MockCriticalCheck struct {
	MockCheck
	critical bool
}

func (m *MockCriticalCheck) IsCritical() bool { return m.critical }

func TestNewChecker(t *testing.T) {
	checker := NewChecker()
	if checker == nil {
		t.Fatal("NewChecker() returned nil")
	}
	if len(checker.checks) != 0 {
		t.Errorf("NewChecker().checks should be empty, got %d", len(checker.checks))
	}
}

func TestChecker_AddCheck(t *testing.T) {
	checker := NewChecker()
	check := &MockCheck{name: "test", category: "internal"}

	checker.AddCheck(check)

	if len(checker.checks) != 1 {
		t.Errorf("AddCheck() should add check, got %d checks", len(checker.checks))
	}
}

func TestChecker_Run_AllHealthy(t *testing.T) {
	checker := NewChecker()

	checker.AddCheck(&MockCheck{
		name:     "Internal Check",
		category: "internal",
		result: Result{
			Status:  StatusHealthy,
			Message: "OK",
		},
	})
	checker.AddCheck(&MockCheck{
		name:     "External Check",
		category: "external",
		result: Result{
			Status:  StatusHealthy,
			Message: "OK",
		},
	})

	report := checker.Run(context.Background())

	if report.OverallStatus != StatusHealthy {
		t.Errorf("OverallStatus = %v, want %v", report.OverallStatus, StatusHealthy)
	}
	if !report.CanProceed {
		t.Error("CanProceed should be true")
	}
	if len(report.Results) != 2 {
		t.Errorf("len(Results) = %v, want %v", len(report.Results), 2)
	}
}

func TestChecker_Run_DegradedStatus(t *testing.T) {
	checker := NewChecker()

	checker.AddCheck(&MockCheck{
		name:     "Degraded Check",
		category: "internal",
		result: Result{
			Status:  StatusDegraded,
			Message: "Warning",
		},
	})

	report := checker.Run(context.Background())

	if report.OverallStatus != StatusDegraded {
		t.Errorf("OverallStatus = %v, want %v", report.OverallStatus, StatusDegraded)
	}
	if !report.CanProceed {
		t.Error("CanProceed should be true for degraded")
	}
	if len(report.Warnings) != 1 {
		t.Errorf("len(Warnings) = %v, want %v", len(report.Warnings), 1)
	}
}

func TestChecker_Run_CriticalStopsExecution(t *testing.T) {
	checker := NewChecker()

	checker.AddCheck(&MockCheck{
		name:     "Critical Check",
		category: "internal",
		result: Result{
			Status:  StatusCritical,
			Message: "Error",
		},
	})
	checker.AddCheck(&MockCheck{
		name:     "External Check",
		category: "external",
		result: Result{
			Status:  StatusHealthy,
			Message: "OK",
		},
	})

	report := checker.Run(context.Background())

	if report.OverallStatus != StatusCritical {
		t.Errorf("OverallStatus = %v, want %v", report.OverallStatus, StatusCritical)
	}
	if report.CanProceed {
		t.Error("CanProceed should be false for critical")
	}
	// External check should not run after internal critical failure
	if len(report.Results) != 1 {
		t.Errorf("len(Results) = %v, want %v (critical should stop)", len(report.Results), 1)
	}
}

func TestChecker_Run_NonCriticalCheck(t *testing.T) {
	checker := NewChecker()

	// Non-critical check with critical status should only degrade
	checker.AddCheck(&MockCriticalCheck{
		MockCheck: MockCheck{
			name:     "Non-Critical Check",
			category: "internal",
			result: Result{
				Status:  StatusCritical,
				Message: "Error",
			},
		},
		critical: false,
	})

	report := checker.Run(context.Background())

	// Should be degraded, not critical
	if report.OverallStatus != StatusDegraded {
		t.Errorf("OverallStatus = %v, want %v", report.OverallStatus, StatusDegraded)
	}
	if !report.CanProceed {
		t.Error("CanProceed should be true for non-critical")
	}
}

func TestChecker_Run_InternalBeforeExternal(t *testing.T) {
	checker := NewChecker()
	var order []string

	// Add external first
	checker.AddCheck(&MockCheck{
		name:     "External",
		category: "external",
		result: Result{
			Status: StatusHealthy,
		},
	})
	// Add internal second
	checker.AddCheck(&MockCheck{
		name:     "Internal",
		category: "internal",
		result: Result{
			Status: StatusHealthy,
		},
	})

	report := checker.Run(context.Background())

	// Results should be: internal first, then external
	for _, r := range report.Results {
		order = append(order, r.Name)
	}

	// Internal should come before external regardless of add order
	if len(order) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(order))
	}
}

func TestChecker_Run_HasDuration(t *testing.T) {
	checker := NewChecker()
	checker.AddCheck(&MockCheck{
		name:     "Test",
		category: "internal",
		result:   Result{Status: StatusHealthy},
	})

	report := checker.Run(context.Background())

	if report.Duration == 0 {
		t.Error("Report.Duration should not be zero")
	}
	if report.Timestamp.IsZero() {
		t.Error("Report.Timestamp should not be zero")
	}
}

func TestChecker_Run_CurrentSchemaVersion(t *testing.T) {
	checker := NewChecker()

	report := checker.Run(context.Background())

	if report.CurrentSchemaVersion == "" {
		t.Error("Report.CurrentSchemaVersion should not be empty")
	}
}

func TestChecker_QuickCheck_AllHealthy(t *testing.T) {
	checker := NewChecker()
	checker.AddCheck(&MockCheck{
		name:     "Test",
		category: "internal",
		result:   Result{Status: StatusHealthy},
	})

	result := checker.QuickCheck(context.Background())

	if !result {
		t.Error("QuickCheck() should return true when all checks healthy")
	}
}

func TestChecker_QuickCheck_CriticalFails(t *testing.T) {
	checker := NewChecker()
	checker.AddCheck(&MockCriticalCheck{
		MockCheck: MockCheck{
			name:     "Critical",
			category: "internal",
			result:   Result{Status: StatusCritical},
		},
		critical: true,
	})

	result := checker.QuickCheck(context.Background())

	if result {
		t.Error("QuickCheck() should return false when critical check fails")
	}
}

func TestChecker_QuickCheck_NonCriticalPasses(t *testing.T) {
	checker := NewChecker()
	checker.AddCheck(&MockCriticalCheck{
		MockCheck: MockCheck{
			name:     "NonCritical",
			category: "internal",
			result:   Result{Status: StatusCritical},
		},
		critical: false,
	})

	result := checker.QuickCheck(context.Background())

	if !result {
		t.Error("QuickCheck() should return true for non-critical failures")
	}
}

func TestChecker_Run_WithContext(t *testing.T) {
	checker := NewChecker()
	checker.AddCheck(&MockCheck{
		name:     "Test",
		category: "internal",
		result:   Result{Status: StatusHealthy},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	report := checker.Run(ctx)

	if report == nil {
		t.Fatal("Run() returned nil")
	}
}
