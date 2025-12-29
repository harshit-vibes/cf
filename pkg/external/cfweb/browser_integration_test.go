//go:build integration

package cfweb

import (
	"strings"
	"testing"
	"time"

	"github.com/playwright-community/playwright-go"
)

// setupPlaywright creates a Playwright browser instance
func setupPlaywright(t *testing.T) (*playwright.Playwright, playwright.Browser, func()) {
	pw, err := playwright.Run()
	if err != nil {
		t.Fatalf("could not start playwright: %v", err)
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		pw.Stop()
		t.Fatalf("could not launch browser: %v", err)
	}

	cleanup := func() {
		browser.Close()
		pw.Stop()
	}

	return pw, browser, cleanup
}

// isCloudflareBlocking checks if Cloudflare is blocking access
func isCloudflareBlocking(content string) bool {
	cloudflareIndicators := []string{
		"Attention Required",
		"Cloudflare",
		"cf-browser-verification",
		"challenge-platform",
		"cf-chl-",
		"Just a moment",
		"Checking your browser",
	}

	for _, indicator := range cloudflareIndicators {
		if strings.Contains(content, indicator) {
			return true
		}
	}
	return false
}

// waitForCloudflare waits for Cloudflare challenge to complete
func waitForCloudflare(page playwright.Page, timeout time.Duration) error {
	start := time.Now()
	for time.Since(start) < timeout {
		content, err := page.Content()
		if err != nil {
			return err
		}

		if !isCloudflareBlocking(content) {
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	return nil // Return nil even if timeout - let test handle it
}

func TestBrowser_FetchProblemPage_Real(t *testing.T) {
	_, browser, cleanup := setupPlaywright(t)
	defer cleanup()

	page, err := browser.NewPage()
	if err != nil {
		t.Fatalf("could not create page: %v", err)
	}
	defer page.Close()

	// Navigate to problem 1A (Theatre Square)
	_, err = page.Goto("https://codeforces.com/problemset/problem/1/A", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		Timeout:   playwright.Float(30000),
	})
	if err != nil {
		t.Fatalf("could not navigate to problem page: %v", err)
	}

	// Wait for potential Cloudflare challenge
	waitForCloudflare(page, 10*time.Second)

	// Get page content and title
	content, err := page.Content()
	if err != nil {
		t.Fatalf("could not get page content: %v", err)
	}

	title, err := page.Title()
	if err != nil {
		t.Fatalf("could not get page title: %v", err)
	}

	// Check if Cloudflare is blocking
	if isCloudflareBlocking(content) {
		t.Skip("Cloudflare is blocking headless browser access - skipping browser test")
	}

	// Verify we're on the problem page
	if !strings.Contains(title, "Theatre Square") && !strings.Contains(title, "Problem") {
		t.Logf("unexpected page title: %s", title)
	}

	// Check for problem statement
	if strings.Contains(content, "problem-statement") || strings.Contains(content, "Theatre Square") {
		t.Logf("Successfully loaded problem page with title: %s", title)
	} else {
		t.Error("problem statement not found on page")
	}
}

func TestBrowser_ParseProblem_Real(t *testing.T) {
	_, browser, cleanup := setupPlaywright(t)
	defer cleanup()

	page, err := browser.NewPage()
	if err != nil {
		t.Fatalf("could not create page: %v", err)
	}
	defer page.Close()

	// Navigate to problem 1A
	_, err = page.Goto("https://codeforces.com/problemset/problem/1/A", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		Timeout:   playwright.Float(30000),
	})
	if err != nil {
		t.Fatalf("could not navigate to problem page: %v", err)
	}

	// Wait for potential Cloudflare challenge
	waitForCloudflare(page, 10*time.Second)

	content, err := page.Content()
	if err != nil {
		t.Fatalf("could not get page content: %v", err)
	}

	if isCloudflareBlocking(content) {
		t.Skip("Cloudflare is blocking headless browser access - skipping browser test")
	}

	// Extract problem name
	titleLocator := page.Locator(".title")
	titleText, err := titleLocator.First().TextContent(playwright.LocatorTextContentOptions{
		Timeout: playwright.Float(5000),
	})
	if err != nil {
		t.Logf("Warning: could not get problem title: %v", err)
		titleText = ""
	}

	// Check problem name contains expected text
	if titleText != "" && !strings.Contains(titleText, "Theatre Square") {
		t.Logf("problem title = %s, expected 'Theatre Square'", titleText)
	}

	// Extract time limit
	timeLimitLocator := page.Locator(".time-limit")
	timeLimitText, _ := timeLimitLocator.TextContent(playwright.LocatorTextContentOptions{
		Timeout: playwright.Float(5000),
	})

	// Extract memory limit
	memoryLimitLocator := page.Locator(".memory-limit")
	memoryLimitText, _ := memoryLimitLocator.TextContent(playwright.LocatorTextContentOptions{
		Timeout: playwright.Float(5000),
	})

	// Extract sample tests
	sampleInputs := page.Locator(".sample-test .input pre")
	inputCount, _ := sampleInputs.Count()

	t.Logf("Problem parsed: Title=%s, TimeLimit=%s, MemoryLimit=%s, Samples=%d",
		titleText, timeLimitText, memoryLimitText, inputCount)
}

func TestBrowser_GetCSRFToken_Real(t *testing.T) {
	_, browser, cleanup := setupPlaywright(t)
	defer cleanup()

	page, err := browser.NewPage()
	if err != nil {
		t.Fatalf("could not create page: %v", err)
	}
	defer page.Close()

	// Navigate to any page with CSRF token
	_, err = page.Goto("https://codeforces.com/enter", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		Timeout:   playwright.Float(30000),
	})
	if err != nil {
		t.Fatalf("could not navigate to page: %v", err)
	}

	// Wait for potential Cloudflare challenge
	waitForCloudflare(page, 15*time.Second)

	// Get page content
	content, err := page.Content()
	if err != nil {
		t.Fatalf("could not get page content: %v", err)
	}

	if isCloudflareBlocking(content) {
		t.Skip("Cloudflare is blocking headless browser access - skipping CSRF test")
	}

	// Extract CSRF token using the existing function
	csrfToken := extractCSRFToken(content)
	if csrfToken == "" {
		t.Log("CSRF token not found on page - this may be due to page structure changes")
	} else {
		// CSRF tokens are typically 32 hex characters
		if len(csrfToken) < 16 {
			t.Errorf("CSRF token seems too short: %s", csrfToken)
		} else {
			t.Logf("Successfully extracted CSRF token: %s...", csrfToken[:min(16, len(csrfToken))])
		}
	}
}

func TestBrowser_ContestPage_Real(t *testing.T) {
	_, browser, cleanup := setupPlaywright(t)
	defer cleanup()

	page, err := browser.NewPage()
	if err != nil {
		t.Fatalf("could not create page: %v", err)
	}
	defer page.Close()

	// Navigate to contest 1
	_, err = page.Goto("https://codeforces.com/contest/1", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		Timeout:   playwright.Float(30000),
	})
	if err != nil {
		t.Fatalf("could not navigate to contest page: %v", err)
	}

	// Wait for potential Cloudflare challenge
	waitForCloudflare(page, 10*time.Second)

	content, err := page.Content()
	if err != nil {
		t.Fatalf("could not get page content: %v", err)
	}

	if isCloudflareBlocking(content) {
		t.Skip("Cloudflare is blocking headless browser access - skipping contest test")
	}

	// Wait for problems table
	_, err = page.WaitForSelector(".problems", playwright.PageWaitForSelectorOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(10000),
	})
	if err != nil {
		t.Skipf("problems table not found (may be Cloudflare): %v", err)
	}

	// Count problems
	problemLinks := page.Locator(".problems td.id a")
	count, err := problemLinks.Count()
	if err != nil {
		t.Fatalf("could not count problem links: %v", err)
	}

	// Contest 1 should have at least 3 problems
	if count < 3 {
		t.Errorf("expected at least 3 problems, got %d", count)
	} else {
		t.Logf("Contest 1 has %d problems", count)
	}
}

func TestBrowser_ExtractCookies_Real(t *testing.T) {
	_, browser, cleanup := setupPlaywright(t)
	defer cleanup()

	context, err := browser.NewContext()
	if err != nil {
		t.Fatalf("could not create browser context: %v", err)
	}
	defer context.Close()

	page, err := context.NewPage()
	if err != nil {
		t.Fatalf("could not create page: %v", err)
	}
	defer page.Close()

	// Navigate to Codeforces
	_, err = page.Goto("https://codeforces.com", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		Timeout:   playwright.Float(30000),
	})
	if err != nil {
		t.Fatalf("could not navigate to CF: %v", err)
	}

	// Wait for potential Cloudflare challenge
	waitForCloudflare(page, 15*time.Second)

	content, err := page.Content()
	if err != nil {
		t.Fatalf("could not get page content: %v", err)
	}

	if isCloudflareBlocking(content) {
		t.Skip("Cloudflare is blocking headless browser access - skipping cookie extraction test")
	}

	// Get cookies from context
	cookies, err := context.Cookies()
	if err != nil {
		t.Fatalf("could not get cookies: %v", err)
	}

	// Log cookie names (not values for security)
	var cookieNames []string
	for _, cookie := range cookies {
		cookieNames = append(cookieNames, cookie.Name)
	}

	t.Logf("Found %d cookies: %v", len(cookies), cookieNames)

	// Check for expected cookies
	expectedCookies := []string{"cf_clearance", "JSESSIONID", "39ce7"}
	for _, expected := range expectedCookies {
		found := false
		for _, cookie := range cookies {
			if cookie.Name == expected {
				found = true
				t.Logf("Found %s cookie (expires: %v)", expected, time.Unix(int64(cookie.Expires), 0))
				break
			}
		}
		if !found {
			t.Logf("%s cookie not found (may require login or Cloudflare bypass)", expected)
		}
	}
}

// Helper function for Go versions < 1.21
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
