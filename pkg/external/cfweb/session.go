package cfweb

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

// Pre-compiled regexes for CSRF token extraction
var (
	reCSRFInput1   = regexp.MustCompile(`<input[^>]+name="csrf_token"[^>]+value="([^"]+)"`)
	reCSRFInput2   = regexp.MustCompile(`<input[^>]+value="([^"]+)"[^>]+name="csrf_token"`)
	reCSRFMeta     = regexp.MustCompile(`<meta[^>]+name="X-Csrf-Token"[^>]+content="([^"]+)"`)
	reCSRFJS       = regexp.MustCompile(`Codeforces\.getCsrfToken[^"]*"([^"]+)"`)
	reHiddenInput1 = regexp.MustCompile(`<input[^>]+name="([^"]+)"[^>]+value="([^"]*)"`)
	reHiddenInput2 = regexp.MustCompile(`<input[^>]+value="([^"]*)"[^>]+name="([^"]+)"`)
)

const (
	BaseURL         = "https://codeforces.com"
	UserAgent       = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	SessionExpiry   = 24 * time.Hour
	MaxPageSize     = 5 * 1024 * 1024 // 5MB max page size to prevent OOM
)

// Session manages CF web authentication using pre-extracted cookies
// No login functionality - users extract cookies from browser
type Session struct {
	client      *http.Client
	jar         *cookiejar.Jar
	csrfToken   string
	handle      string
	expiresAt   time.Time
	cfClearance string    // cf_clearance cookie for Cloudflare bypass
	cfClearExp  time.Time // cf_clearance expiration
	userAgent   string    // User-Agent tied to cf_clearance (MUST match)
}

// NewSession creates a new CF session
func NewSession() (*Session, error) {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, fmt.Errorf("create cookie jar: %w", err)
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	return &Session{
		client:    client,
		jar:       jar,
		userAgent: UserAgent,
	}, nil
}

// NewSessionWithUserAgent creates a new CF session with custom User-Agent
// This is important for Cloudflare bypass - the User-Agent MUST match the one
// used to obtain the cf_clearance cookie
func NewSessionWithUserAgent(ua string) (*Session, error) {
	session, err := NewSession()
	if err != nil {
		return nil, err
	}
	session.userAgent = ua
	return session, nil
}

// SetCFClearance sets the cf_clearance cookie for Cloudflare bypass
// The userAgent parameter MUST match the User-Agent used to obtain the cookie
func (s *Session) SetCFClearance(clearance, userAgent string, expiresAt time.Time) {
	s.cfClearance = clearance
	s.cfClearExp = expiresAt
	s.userAgent = userAgent

	cfURL, _ := url.Parse(BaseURL)
	cookies := []*http.Cookie{
		{
			Name:     "cf_clearance",
			Value:    clearance,
			Path:     "/",
			Domain:   ".codeforces.com",
			Expires:  expiresAt,
			HttpOnly: true,
			Secure:   true,
		},
	}
	s.jar.SetCookies(cfURL, cookies)
}

// SetSessionCookies sets session cookies directly (extracted from browser)
func (s *Session) SetSessionCookies(jsessionID, ce7Cookie, handle string) {
	cfURL, _ := url.Parse(BaseURL)

	var cookies []*http.Cookie

	if jsessionID != "" {
		cookies = append(cookies, &http.Cookie{
			Name:     "JSESSIONID",
			Value:    jsessionID,
			Path:     "/",
			Domain:   "codeforces.com",
			HttpOnly: true,
		})
	}

	if ce7Cookie != "" {
		cookies = append(cookies, &http.Cookie{
			Name:     "39ce7",
			Value:    ce7Cookie,
			Path:     "/",
			Domain:   ".codeforces.com",
			Expires:  time.Now().Add(365 * 24 * time.Hour),
			HttpOnly: true,
		})
	}

	if len(cookies) > 0 {
		s.jar.SetCookies(cfURL, cookies)
	}

	s.handle = handle
	s.expiresAt = time.Now().Add(SessionExpiry)
}

// SetFullAuth sets all authentication cookies at once (cf_clearance + session)
// This is the primary method for cookie-based auth
func (s *Session) SetFullAuth(cfClearance, userAgent string, cfClearExp time.Time,
	jsessionID, ce7Cookie, handle string) {
	s.SetCFClearance(cfClearance, userAgent, cfClearExp)
	s.SetSessionCookies(jsessionID, ce7Cookie, handle)
}

// IsCFClearanceValid returns true if cf_clearance is set and not expired
func (s *Session) IsCFClearanceValid() bool {
	if s.cfClearance == "" {
		return false
	}
	return time.Now().Before(s.cfClearExp)
}

// HasSessionCookies returns true if session cookies are set
func (s *Session) HasSessionCookies() bool {
	cfURL, _ := url.Parse(BaseURL)
	cookies := s.jar.Cookies(cfURL)

	for _, c := range cookies {
		if c.Name == "JSESSIONID" || c.Name == "39ce7" {
			return true
		}
	}
	return false
}

// IsAuthenticated returns true if session has valid cf_clearance and session cookies
func (s *Session) IsAuthenticated() bool {
	return s.IsCFClearanceValid() && s.HasSessionCookies()
}

// IsReadyForSubmission returns true if session can submit solutions
func (s *Session) IsReadyForSubmission() bool {
	return s.IsAuthenticated() && s.handle != ""
}

// get makes a GET request with proper headers
func (s *Session) get(urlStr string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}

	ua := s.userAgent
	if ua == "" {
		ua = UserAgent
	}
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	return s.client.Do(req)
}

// GetCSRFToken returns the current CSRF token
func (s *Session) GetCSRFToken() string {
	return s.csrfToken
}

// RefreshCSRFToken fetches a fresh CSRF token from any CF page
func (s *Session) RefreshCSRFToken() error {
	resp, err := s.get(BaseURL)
	if err != nil {
		return fmt.Errorf("get page: %w", err)
	}
	defer resp.Body.Close()

	// Use bounded reader to prevent OOM from large responses
	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxPageSize))
	if err != nil {
		return fmt.Errorf("read page: %w", err)
	}

	csrfToken := extractCSRFToken(string(body))
	if csrfToken == "" {
		return fmt.Errorf("csrf token not found")
	}

	s.csrfToken = csrfToken
	return nil
}

// Validate checks if the session is still valid by verifying handle appears on CF
func (s *Session) Validate() error {
	if !s.IsCFClearanceValid() {
		return fmt.Errorf("cf_clearance expired or not set")
	}

	if s.handle == "" {
		return fmt.Errorf("handle not set")
	}

	resp, err := s.get(BaseURL)
	if err != nil {
		return fmt.Errorf("validation request failed: %w", err)
	}
	defer resp.Body.Close()

	// Use bounded reader to prevent OOM from large responses
	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxPageSize))
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	// Check if we're logged in by looking for handle in page
	if !strings.Contains(string(body), s.handle) {
		return fmt.Errorf("session invalid - handle not found on page")
	}

	// Refresh CSRF token while we're at it
	if csrfToken := extractCSRFToken(string(body)); csrfToken != "" {
		s.csrfToken = csrfToken
	}

	return nil
}

// Handle returns the configured handle
func (s *Session) Handle() string {
	return s.handle
}

// ExpiresAt returns when the session expires
func (s *Session) ExpiresAt() time.Time {
	return s.expiresAt
}

// CFClearanceExpiresAt returns when cf_clearance expires
func (s *Session) CFClearanceExpiresAt() time.Time {
	return s.cfClearExp
}

// CFClearanceExpiresIn returns time until cf_clearance expires
func (s *Session) CFClearanceExpiresIn() time.Duration {
	if s.cfClearance == "" {
		return 0
	}
	return time.Until(s.cfClearExp)
}

// Client returns the underlying HTTP client
func (s *Session) Client() *http.Client {
	return s.client
}

// GetUserAgent returns the User-Agent tied to this session
func (s *Session) GetUserAgent() string {
	if s.userAgent == "" {
		return UserAgent
	}
	return s.userAgent
}

// GetCFClearance returns the current cf_clearance value
func (s *Session) GetCFClearance() string {
	return s.cfClearance
}

// Helper functions

// extractCSRFToken extracts CSRF token from HTML
func extractCSRFToken(htmlStr string) string {
	// Try input field first
	if matches := reCSRFInput1.FindStringSubmatch(htmlStr); len(matches) > 1 {
		return matches[1]
	}

	// Try with value before name
	if matches := reCSRFInput2.FindStringSubmatch(htmlStr); len(matches) > 1 {
		return matches[1]
	}

	// Try meta tag
	if matches := reCSRFMeta.FindStringSubmatch(htmlStr); len(matches) > 1 {
		return matches[1]
	}

	// Try JavaScript variable
	if matches := reCSRFJS.FindStringSubmatch(htmlStr); len(matches) > 1 {
		return matches[1]
	}

	return ""
}

// extractHiddenInput extracts value of a hidden input field
func extractHiddenInput(htmlStr, name string) string {
	// Use pre-compiled regex with name capture, then filter by name
	for _, match := range reHiddenInput1.FindAllStringSubmatch(htmlStr, -1) {
		if len(match) > 2 && match[1] == name {
			return match[2]
		}
	}

	// Try with value before name
	for _, match := range reHiddenInput2.FindAllStringSubmatch(htmlStr, -1) {
		if len(match) > 2 && match[2] == name {
			return match[1]
		}
	}

	return ""
}

// GetHTMLDocument parses HTML into a document
func GetHTMLDocument(r io.Reader) (*html.Node, error) {
	return html.Parse(r)
}
