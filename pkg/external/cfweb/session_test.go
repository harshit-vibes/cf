package cfweb

import (
	"strings"
	"testing"
	"time"
)

func TestNewSession(t *testing.T) {
	session, err := NewSession()
	if err != nil {
		t.Fatalf("NewSession() failed: %v", err)
	}

	if session == nil {
		t.Fatal("NewSession() returned nil")
	}

	if session.client == nil {
		t.Error("session.client should not be nil")
	}

	if session.jar == nil {
		t.Error("session.jar should not be nil")
	}

	if session.Handle() != "" {
		t.Error("new session handle should be empty")
	}

	if session.GetCSRFToken() != "" {
		t.Error("new session CSRF token should be empty")
	}

	if session.IsCFClearanceValid() {
		t.Error("new session should not have valid cf_clearance")
	}

	if session.HasSessionCookies() {
		t.Error("new session should not have session cookies")
	}
}

func TestNewSessionWithUserAgent(t *testing.T) {
	customUA := "Custom User Agent/1.0"
	session, err := NewSessionWithUserAgent(customUA)
	if err != nil {
		t.Fatalf("NewSessionWithUserAgent() failed: %v", err)
	}

	if session == nil {
		t.Fatal("NewSessionWithUserAgent() returned nil")
	}

	if session.GetUserAgent() != customUA {
		t.Errorf("GetUserAgent() = %s, want %s", session.GetUserAgent(), customUA)
	}
}

func TestSession_SetCFClearance(t *testing.T) {
	session, err := NewSession()
	if err != nil {
		t.Fatalf("NewSession() failed: %v", err)
	}

	clearance := "test_cf_clearance_value"
	userAgent := "Test User Agent"
	expiresAt := time.Now().Add(30 * time.Minute)

	session.SetCFClearance(clearance, userAgent, expiresAt)

	if session.GetCFClearance() != clearance {
		t.Errorf("GetCFClearance() = %s, want %s", session.GetCFClearance(), clearance)
	}

	if session.GetUserAgent() != userAgent {
		t.Errorf("GetUserAgent() = %s, want %s", session.GetUserAgent(), userAgent)
	}

	if !session.IsCFClearanceValid() {
		t.Error("IsCFClearanceValid() should return true")
	}
}

func TestSession_SetSessionCookies(t *testing.T) {
	session, err := NewSession()
	if err != nil {
		t.Fatalf("NewSession() failed: %v", err)
	}

	jsessionID := "test-jsessionid-12345"
	ce7Cookie := "test-39ce7-cookie"
	handle := "testuser"

	session.SetSessionCookies(jsessionID, ce7Cookie, handle)

	if session.Handle() != handle {
		t.Errorf("Handle() = %s, want %s", session.Handle(), handle)
	}

	if !session.HasSessionCookies() {
		t.Error("HasSessionCookies() should return true after setting cookies")
	}

	if session.ExpiresAt().IsZero() {
		t.Error("ExpiresAt() should be set after setting session cookies")
	}
}

func TestSession_SetFullAuth(t *testing.T) {
	session, err := NewSession()
	if err != nil {
		t.Fatalf("NewSession() failed: %v", err)
	}

	cfClearance := "test_clearance"
	userAgent := "Test UA"
	cfClearExp := time.Now().Add(30 * time.Minute)
	jsessionID := "jsession123"
	ce7Cookie := "ce7cookie456"
	handle := "testhandle"

	session.SetFullAuth(cfClearance, userAgent, cfClearExp, jsessionID, ce7Cookie, handle)

	if !session.IsCFClearanceValid() {
		t.Error("cf_clearance should be valid")
	}

	if !session.HasSessionCookies() {
		t.Error("session cookies should be set")
	}

	if !session.IsAuthenticated() {
		t.Error("session should be authenticated")
	}

	if !session.IsReadyForSubmission() {
		t.Error("session should be ready for submission")
	}
}

func TestSession_IsCFClearanceValid_NotSet(t *testing.T) {
	session, err := NewSession()
	if err != nil {
		t.Fatalf("NewSession() failed: %v", err)
	}

	if session.IsCFClearanceValid() {
		t.Error("IsCFClearanceValid() should return false when not set")
	}
}

func TestSession_IsCFClearanceValid_Expired(t *testing.T) {
	session, err := NewSession()
	if err != nil {
		t.Fatalf("NewSession() failed: %v", err)
	}

	// Set expired clearance
	session.SetCFClearance("expired_clearance", "UA", time.Now().Add(-1*time.Hour))

	if session.IsCFClearanceValid() {
		t.Error("IsCFClearanceValid() should return false for expired clearance")
	}
}

func TestSession_CFClearanceExpiresIn(t *testing.T) {
	session, err := NewSession()
	if err != nil {
		t.Fatalf("NewSession() failed: %v", err)
	}

	// No clearance set - should return 0
	if session.CFClearanceExpiresIn() != 0 {
		t.Error("CFClearanceExpiresIn() should return 0 when not set")
	}

	// Set future clearance
	session.SetCFClearance("test", "UA", time.Now().Add(30*time.Minute))
	expiresIn := session.CFClearanceExpiresIn()
	if expiresIn <= 0 {
		t.Error("CFClearanceExpiresIn() should return positive for future expiry")
	}
}

func TestSession_HasSessionCookies_Empty(t *testing.T) {
	session, err := NewSession()
	if err != nil {
		t.Fatalf("NewSession() failed: %v", err)
	}

	if session.HasSessionCookies() {
		t.Error("HasSessionCookies() should return false for new session")
	}
}

func TestSession_HasSessionCookies_OnlyJSessionID(t *testing.T) {
	session, err := NewSession()
	if err != nil {
		t.Fatalf("NewSession() failed: %v", err)
	}

	session.SetSessionCookies("jsessionid123", "", "handle")

	if !session.HasSessionCookies() {
		t.Error("HasSessionCookies() should return true with only JSESSIONID")
	}
}

func TestSession_HasSessionCookies_OnlyCE7(t *testing.T) {
	session, err := NewSession()
	if err != nil {
		t.Fatalf("NewSession() failed: %v", err)
	}

	session.SetSessionCookies("", "ce7cookie123", "handle")

	if !session.HasSessionCookies() {
		t.Error("HasSessionCookies() should return true with only 39ce7")
	}
}

func TestSession_IsAuthenticated(t *testing.T) {
	tests := []struct {
		name         string
		cfClearance  string
		cfClearExp   time.Time
		jsessionID   string
		ce7Cookie    string
		wantAuth     bool
	}{
		{
			name:     "no cookies at all",
			wantAuth: false,
		},
		{
			name:        "only cf_clearance",
			cfClearance: "test",
			cfClearExp:  time.Now().Add(30 * time.Minute),
			wantAuth:    false,
		},
		{
			name:       "only session cookies",
			jsessionID: "test",
			wantAuth:   false,
		},
		{
			name:        "cf_clearance + session cookies",
			cfClearance: "test",
			cfClearExp:  time.Now().Add(30 * time.Minute),
			jsessionID:  "test",
			wantAuth:    true,
		},
		{
			name:        "expired cf_clearance + session cookies",
			cfClearance: "test",
			cfClearExp:  time.Now().Add(-1 * time.Hour),
			jsessionID:  "test",
			wantAuth:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, _ := NewSession()

			if tt.cfClearance != "" {
				session.SetCFClearance(tt.cfClearance, "UA", tt.cfClearExp)
			}
			if tt.jsessionID != "" || tt.ce7Cookie != "" {
				session.SetSessionCookies(tt.jsessionID, tt.ce7Cookie, "handle")
			}

			if session.IsAuthenticated() != tt.wantAuth {
				t.Errorf("IsAuthenticated() = %v, want %v", session.IsAuthenticated(), tt.wantAuth)
			}
		})
	}
}

func TestSession_IsReadyForSubmission(t *testing.T) {
	tests := []struct {
		name        string
		cfClearance string
		cfClearExp  time.Time
		jsessionID  string
		handle      string
		wantReady   bool
	}{
		{
			name:      "nothing set",
			wantReady: false,
		},
		{
			name:        "authenticated but no handle",
			cfClearance: "test",
			cfClearExp:  time.Now().Add(30 * time.Minute),
			jsessionID:  "test",
			handle:      "",
			wantReady:   false,
		},
		{
			name:        "fully configured",
			cfClearance: "test",
			cfClearExp:  time.Now().Add(30 * time.Minute),
			jsessionID:  "test",
			handle:      "testuser",
			wantReady:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, _ := NewSession()

			if tt.cfClearance != "" {
				session.SetCFClearance(tt.cfClearance, "UA", tt.cfClearExp)
			}
			if tt.jsessionID != "" {
				session.SetSessionCookies(tt.jsessionID, "", tt.handle)
			}

			if session.IsReadyForSubmission() != tt.wantReady {
				t.Errorf("IsReadyForSubmission() = %v, want %v", session.IsReadyForSubmission(), tt.wantReady)
			}
		})
	}
}

func TestSession_Handle(t *testing.T) {
	session, err := NewSession()
	if err != nil {
		t.Fatalf("NewSession() failed: %v", err)
	}

	if session.Handle() != "" {
		t.Error("new session handle should be empty")
	}

	session.SetSessionCookies("jsession", "ce7", "testuser")
	if session.Handle() != "testuser" {
		t.Errorf("Handle() = %s, want testuser", session.Handle())
	}
}

func TestSession_ExpiresAt(t *testing.T) {
	session, err := NewSession()
	if err != nil {
		t.Fatalf("NewSession() failed: %v", err)
	}

	// New session should have zero time
	if !session.ExpiresAt().IsZero() {
		t.Error("new session ExpiresAt should be zero")
	}

	// Set session cookies to trigger expiry setting
	session.SetSessionCookies("jsession", "ce7", "handle")

	if session.ExpiresAt().IsZero() {
		t.Error("ExpiresAt should be set after setting session cookies")
	}
}

func TestSession_CFClearanceExpiresAt(t *testing.T) {
	session, err := NewSession()
	if err != nil {
		t.Fatalf("NewSession() failed: %v", err)
	}

	expiry := time.Now().Add(30 * time.Minute)
	session.SetCFClearance("test", "UA", expiry)

	if !session.CFClearanceExpiresAt().Equal(expiry) {
		t.Errorf("CFClearanceExpiresAt() = %v, want %v", session.CFClearanceExpiresAt(), expiry)
	}
}

func TestSession_Client(t *testing.T) {
	session, err := NewSession()
	if err != nil {
		t.Fatalf("NewSession() failed: %v", err)
	}

	client := session.Client()
	if client == nil {
		t.Error("Client() should not return nil")
	}

	if client != session.client {
		t.Error("Client() should return the internal client")
	}
}

func TestSession_GetCSRFToken(t *testing.T) {
	session, err := NewSession()
	if err != nil {
		t.Fatalf("NewSession() failed: %v", err)
	}

	if session.GetCSRFToken() != "" {
		t.Error("new session CSRF token should be empty")
	}

	session.csrfToken = "test-token-12345"
	if session.GetCSRFToken() != "test-token-12345" {
		t.Errorf("GetCSRFToken() = %s, want test-token-12345", session.GetCSRFToken())
	}
}

func TestSession_GetUserAgent_Default(t *testing.T) {
	session, err := NewSession()
	if err != nil {
		t.Fatalf("NewSession() failed: %v", err)
	}

	if session.GetUserAgent() != UserAgent {
		t.Errorf("GetUserAgent() = %s, want default %s", session.GetUserAgent(), UserAgent)
	}
}

func TestExtractCSRFToken(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
	}{
		{
			name: "input field name first",
			html: `<input type="hidden" name="csrf_token" value="abc123def456"/>`,
			want: "abc123def456",
		},
		{
			name: "input field value first",
			html: `<input type="hidden" value="xyz789abc123" name="csrf_token"/>`,
			want: "xyz789abc123",
		},
		{
			name: "meta tag",
			html: `<meta name="X-Csrf-Token" content="meta123token"/>`,
			want: "meta123token",
		},
		{
			name: "javascript variable",
			html: `<script>Codeforces.getCsrfToken = function() { return "js_token_123"; }</script>`,
			want: "js_token_123",
		},
		{
			name: "no token",
			html: `<html><body>No token here</body></html>`,
			want: "",
		},
		{
			name: "empty input",
			html: "",
			want: "",
		},
		{
			name: "real CF-like HTML",
			html: `<!DOCTYPE html>
<html>
<head><title>Codeforces</title></head>
<body>
<form>
<input type="hidden" name="csrf_token" value="8a9b0c1d2e3f4g5h6i7j8k9l0m"/>
<input type="text" name="username"/>
</form>
</body>
</html>`,
			want: "8a9b0c1d2e3f4g5h6i7j8k9l0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractCSRFToken(tt.html)
			if got != tt.want {
				t.Errorf("extractCSRFToken() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractHiddenInput(t *testing.T) {
	tests := []struct {
		name  string
		html  string
		field string
		want  string
	}{
		{
			name:  "ftaa name first",
			html:  `<input type="hidden" name="ftaa" value="ftaa_value_123"/>`,
			field: "ftaa",
			want:  "ftaa_value_123",
		},
		{
			name:  "bfaa value first",
			html:  `<input type="hidden" value="bfaa_value_456" name="bfaa"/>`,
			field: "bfaa",
			want:  "bfaa_value_456",
		},
		{
			name:  "field not found",
			html:  `<input type="hidden" name="other" value="something"/>`,
			field: "ftaa",
			want:  "",
		},
		{
			name:  "empty value",
			html:  `<input type="hidden" name="ftaa" value=""/>`,
			field: "ftaa",
			want:  "",
		},
		{
			name:  "empty html",
			html:  "",
			field: "ftaa",
			want:  "",
		},
		{
			name:  "multiple inputs",
			html:  `<input name="a" value="1"/><input name="ftaa" value="found"/><input name="b" value="2"/>`,
			field: "ftaa",
			want:  "found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractHiddenInput(tt.html, tt.field)
			if got != tt.want {
				t.Errorf("extractHiddenInput() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetHTMLDocument(t *testing.T) {
	html := `<!DOCTYPE html><html><head><title>Test</title></head><body><p>Hello</p></body></html>`
	reader := strings.NewReader(html)

	doc, err := GetHTMLDocument(reader)
	if err != nil {
		t.Fatalf("GetHTMLDocument() failed: %v", err)
	}

	if doc == nil {
		t.Error("GetHTMLDocument() returned nil")
	}
}

func TestGetHTMLDocument_Empty(t *testing.T) {
	reader := strings.NewReader("")

	doc, err := GetHTMLDocument(reader)
	// Empty HTML should still parse (to an empty document)
	if err != nil {
		t.Fatalf("GetHTMLDocument() failed: %v", err)
	}

	if doc == nil {
		t.Error("GetHTMLDocument() returned nil")
	}
}

func TestConstants(t *testing.T) {
	if BaseURL != "https://codeforces.com" {
		t.Errorf("BaseURL = %s, want https://codeforces.com", BaseURL)
	}
	if SessionExpiry != 24*time.Hour {
		t.Errorf("SessionExpiry = %v, want 24h", SessionExpiry)
	}
	if UserAgent == "" {
		t.Error("UserAgent should not be empty")
	}
}

func TestSession_Validate_NoClearance(t *testing.T) {
	session, err := NewSession()
	if err != nil {
		t.Fatalf("NewSession() failed: %v", err)
	}

	err = session.Validate()
	if err == nil {
		t.Error("Validate() should return error when cf_clearance not set")
	}
}

func TestSession_Validate_NoHandle(t *testing.T) {
	session, err := NewSession()
	if err != nil {
		t.Fatalf("NewSession() failed: %v", err)
	}

	session.SetCFClearance("cf_clearance_value", "UA", time.Now().Add(30*time.Minute))

	err = session.Validate()
	if err == nil {
		t.Error("Validate() should return error when handle not set")
	}
}

func TestSession_GetUserAgent_Custom(t *testing.T) {
	customUA := "Custom/1.0"
	session, err := NewSession()
	if err != nil {
		t.Fatalf("NewSession() failed: %v", err)
	}

	session.userAgent = customUA
	if session.GetUserAgent() != customUA {
		t.Errorf("GetUserAgent() = %s, want %s", session.GetUserAgent(), customUA)
	}
}

func TestSession_GetUserAgent_Empty(t *testing.T) {
	session, err := NewSession()
	if err != nil {
		t.Fatalf("NewSession() failed: %v", err)
	}

	session.userAgent = ""
	if session.GetUserAgent() != UserAgent {
		t.Errorf("GetUserAgent() = %s, want default %s", session.GetUserAgent(), UserAgent)
	}
}
