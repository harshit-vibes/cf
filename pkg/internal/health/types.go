// Package health provides startup health checks
package health

import (
	"context"
	"time"
)

// Status represents the health status
type Status int

const (
	StatusHealthy  Status = iota // Everything OK
	StatusDegraded               // Some features unavailable
	StatusCritical               // Cannot proceed
)

// String returns the status name
func (s Status) String() string {
	switch s {
	case StatusHealthy:
		return "healthy"
	case StatusDegraded:
		return "degraded"
	case StatusCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// Result represents a single check result
type Result struct {
	Name        string
	Category    string // "internal" or "external"
	Status      Status
	Message     string
	Details     string
	Recoverable bool
	Action      Action
	Duration    time.Duration
}

// Action represents recovery action
type Action int

const (
	ActionNone       Action = iota // Nothing to do
	ActionAutoFix                  // Can fix automatically
	ActionUserPrompt               // Need user input
	ActionManualFix                // User must fix
	ActionRetry                    // Retry may help
	ActionFatal                    // Cannot continue
)

// String returns the action name
func (a Action) String() string {
	switch a {
	case ActionNone:
		return "none"
	case ActionAutoFix:
		return "autofix"
	case ActionUserPrompt:
		return "prompt"
	case ActionManualFix:
		return "manual"
	case ActionRetry:
		return "retry"
	case ActionFatal:
		return "fatal"
	default:
		return "unknown"
	}
}

// Report is the complete health report
type Report struct {
	Timestamp     time.Time
	Duration      time.Duration
	OverallStatus Status
	Results       []Result
	CanProceed    bool
	Warnings      []string
	Errors        []string

	// Schema info
	CurrentSchemaVersion  string
	DetectedSchemaVersion string
	NeedsMigration        bool
}

// Check is the interface for health checks
type Check interface {
	Name() string
	Category() string
	Check(ctx context.Context) Result
}

// AutoFixable is implemented by checks that can auto-fix
type AutoFixable interface {
	Check
	AutoFix(ctx context.Context) error
}

// Critical is implemented by checks that are critical
type Critical interface {
	Check
	IsCritical() bool
}
