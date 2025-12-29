package errors

import (
	"errors"
	"strings"
	"testing"
)

func TestCategory_String(t *testing.T) {
	tests := []struct {
		cat  Category
		want string
	}{
		{CatInternal, "internal"},
		{CatExternal, "external"},
		{CatUser, "user"},
		{CatNetwork, "network"},
		{Category(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.cat.String(); got != tt.want {
				t.Errorf("Category.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryConstants(t *testing.T) {
	if CatInternal != 0 {
		t.Errorf("CatInternal = %v, want 0", CatInternal)
	}
	if CatExternal != 1 {
		t.Errorf("CatExternal = %v, want 1", CatExternal)
	}
	if CatUser != 2 {
		t.Errorf("CatUser = %v, want 2", CatUser)
	}
	if CatNetwork != 3 {
		t.Errorf("CatNetwork = %v, want 3", CatNetwork)
	}
}

func TestRecoveryActionConstants(t *testing.T) {
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

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     *AppError
		want    string
		wantIn  string // substring to check
	}{
		{
			name: "with details",
			err: &AppError{
				Code:    "TEST_ERROR",
				Message: "Test error occurred",
				Details: "Detailed information",
			},
			want: "[TEST_ERROR] Test error occurred: Detailed information",
		},
		{
			name: "without details",
			err: &AppError{
				Code:    "TEST_ERROR",
				Message: "Test error occurred",
			},
			want: "[TEST_ERROR] Test error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("AppError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &AppError{
		Code:  "TEST",
		Cause: cause,
	}

	unwrapped := err.Unwrap()
	if unwrapped != cause {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, cause)
	}
}

func TestAppError_Unwrap_Nil(t *testing.T) {
	err := &AppError{
		Code: "TEST",
	}

	if err.Unwrap() != nil {
		t.Error("Unwrap() should return nil when no cause")
	}
}

func TestNew_KnownError(t *testing.T) {
	err := New(ErrEnvMissing)

	if err.Code != ErrEnvMissing {
		t.Errorf("Code = %v, want %v", err.Code, ErrEnvMissing)
	}
	if err.Category != CatUser {
		t.Errorf("Category = %v, want %v", err.Category, CatUser)
	}
	if err.Message == "" {
		t.Error("Message should not be empty for known error")
	}
	if err.Suggestion == "" {
		t.Error("Suggestion should not be empty for known error")
	}
	if err.Action != ActionAutoFix {
		t.Errorf("Action = %v, want %v", err.Action, ActionAutoFix)
	}
}

func TestNew_UnknownError(t *testing.T) {
	err := New("UNKNOWN_CODE")

	if err.Code != "UNKNOWN_CODE" {
		t.Errorf("Code = %v, want UNKNOWN_CODE", err.Code)
	}
	if err.Category != CatInternal {
		t.Errorf("Category = %v, want %v", err.Category, CatInternal)
	}
	if err.Message != "Unknown error" {
		t.Errorf("Message = %v, want 'Unknown error'", err.Message)
	}
}

func TestWrap(t *testing.T) {
	cause := errors.New("connection refused")
	err := Wrap(ErrCFAPIDown, cause)

	if err.Code != ErrCFAPIDown {
		t.Errorf("Code = %v, want %v", err.Code, ErrCFAPIDown)
	}
	if err.Cause != cause {
		t.Error("Cause should be set")
	}
	if err.Details != "connection refused" {
		t.Errorf("Details = %v, want 'connection refused'", err.Details)
	}
}

func TestAppError_WithDetails(t *testing.T) {
	err := New(ErrEnvMissing)
	result := err.WithDetails("file not found at /path/to/file")

	// Should return same error (for chaining)
	if result != err {
		t.Error("WithDetails() should return same error for chaining")
	}
	if err.Details != "file not found at /path/to/file" {
		t.Errorf("Details = %v", err.Details)
	}
}

func TestAppError_WithSuggestion(t *testing.T) {
	err := New(ErrEnvMissing)
	result := err.WithSuggestion("Try running as admin")

	// Should return same error (for chaining)
	if result != err {
		t.Error("WithSuggestion() should return same error for chaining")
	}
	if err.Suggestion != "Try running as admin" {
		t.Errorf("Suggestion = %v", err.Suggestion)
	}
}

func TestAppError_Chaining(t *testing.T) {
	err := New(ErrCFAPIDown).
		WithDetails("timeout after 30s").
		WithSuggestion("Check your internet connection")

	if err.Details != "timeout after 30s" {
		t.Errorf("Details = %v", err.Details)
	}
	if err.Suggestion != "Check your internet connection" {
		t.Errorf("Suggestion = %v", err.Suggestion)
	}
}

func TestRegistry_AllCodes(t *testing.T) {
	codes := []string{
		ErrEnvMissing,
		ErrEnvCorrupt,
		ErrHandleNotSet,
		ErrCredentialsMissing,
		ErrSessionExpired,
		ErrCFAPIDown,
		ErrCFAPIRateLimit,
		ErrCFWebChanged,
		ErrCFLoginFailed,
		ErrSchemaIncompatible,
		ErrWorkspaceNotFound,
		ErrNetworkOffline,
	}

	for _, code := range codes {
		t.Run(code, func(t *testing.T) {
			template, ok := Registry[code]
			if !ok {
				t.Errorf("Registry missing code %s", code)
				return
			}
			if template.Code != code {
				t.Errorf("Registry[%s].Code = %v", code, template.Code)
			}
			if template.Message == "" {
				t.Errorf("Registry[%s].Message is empty", code)
			}
		})
	}
}

func TestRegistry_Categories(t *testing.T) {
	// Verify error categories are correct
	userErrors := []string{
		ErrEnvMissing,
		ErrEnvCorrupt,
		ErrHandleNotSet,
		ErrCredentialsMissing,
		ErrSessionExpired,
		ErrWorkspaceNotFound,
	}
	for _, code := range userErrors {
		if Registry[code].Category != CatUser {
			t.Errorf("%s should be CatUser, got %v", code, Registry[code].Category)
		}
	}

	externalErrors := []string{
		ErrCFAPIDown,
		ErrCFAPIRateLimit,
		ErrCFWebChanged,
		ErrCFLoginFailed,
	}
	for _, code := range externalErrors {
		if Registry[code].Category != CatExternal {
			t.Errorf("%s should be CatExternal, got %v", code, Registry[code].Category)
		}
	}

	networkErrors := []string{
		ErrNetworkOffline,
	}
	for _, code := range networkErrors {
		if Registry[code].Category != CatNetwork {
			t.Errorf("%s should be CatNetwork, got %v", code, Registry[code].Category)
		}
	}
}

func TestErrorCodes(t *testing.T) {
	// Verify error codes are properly prefixed
	tests := []struct {
		code   string
		prefix string
	}{
		{ErrSchemaInvalid, "SCHEMA_"},
		{ErrSchemaIncompatible, "SCHEMA_"},
		{ErrWorkspaceCorrupt, "WORKSPACE_"},
		{ErrMigrationFailed, "MIGRATION_"},
		{ErrEnvMissing, "ENV_"},
		{ErrEnvCorrupt, "ENV_"},
		{ErrHandleNotSet, "HANDLE_"},
		{ErrCredentialsMissing, "CREDENTIALS_"},
		{ErrSessionExpired, "SESSION_"},
		{ErrCFAPIDown, "CF_"},
		{ErrCFAPIRateLimit, "CF_"},
		{ErrCFWebChanged, "CF_"},
		{ErrCFLoginFailed, "CF_"},
		{ErrCFSubmitFailed, "CF_"},
		{ErrCFParseFailed, "CF_"},
		{ErrNetworkOffline, "NETWORK_"},
		{ErrNetworkTimeout, "NETWORK_"},
		{ErrNetworkDNS, "NETWORK_"},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			if !strings.HasPrefix(tt.code, tt.prefix) {
				t.Errorf("%s should have prefix %s", tt.code, tt.prefix)
			}
		})
	}
}

func TestAppError_ImplementsError(t *testing.T) {
	var _ error = &AppError{}
}

func TestNew_DoesNotModifyRegistry(t *testing.T) {
	// Get original values
	original := Registry[ErrEnvMissing]
	originalMessage := original.Message

	// Create new error and modify it
	err := New(ErrEnvMissing)
	err.Message = "Modified message"

	// Registry should be unchanged
	if Registry[ErrEnvMissing].Message != originalMessage {
		t.Error("New() should not allow modification of Registry")
	}
}
