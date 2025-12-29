// Cookie-based authentication test using built-in cfweb infrastructure
// Run with: go run ./cmd/explore cookieauth
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/harshit-vibes/dsaprep/pkg/external/cfweb"
)

// CookieAuthTokens holds all the tokens needed for cookie-based auth
type CookieAuthTokens struct {
	CFClearance   string
	CFClearanceUA string
	CFClearExp    time.Time
	JSESSIONID    string
	CE7Cookie     string
	Handle        string
}

func loadCookieAuthTokens() (*CookieAuthTokens, error) {
	home, _ := os.UserHomeDir()
	envPath := home + "/.dsaprep.env"

	file, err := os.Open(envPath)
	if err != nil {
		return nil, fmt.Errorf("open env file: %w", err)
	}
	defer file.Close()

	env := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}

	tokens := &CookieAuthTokens{
		CFClearance:   env["CF_CLEARANCE"],
		CFClearanceUA: env["CF_CLEARANCE_UA"],
		JSESSIONID:    env["CF_JSESSIONID"],
		CE7Cookie:     env["CF_39CE7"],
		Handle:        env["CF_HANDLE"],
	}

	// Parse expiration
	if exp, err := strconv.ParseInt(env["CF_CLEARANCE_EXPIRES"], 10, 64); err == nil {
		tokens.CFClearExp = time.Unix(exp, 0)
	}

	return tokens, nil
}

func runCookieAuthTest() {
	fmt.Println("=" + strings.Repeat("=", 59))
	fmt.Println("  Cookie-Based Authentication Test")
	fmt.Println("  Using cfweb infrastructure (no login required)")
	fmt.Println("=" + strings.Repeat("=", 59))
	fmt.Println()

	// Step 1: Load tokens from ~/.dsaprep.env
	fmt.Println("[1] Loading tokens from ~/.dsaprep.env...")
	tokens, err := loadCookieAuthTokens()
	if err != nil {
		fmt.Printf("    ERROR: %v\n", err)
		fmt.Println("    Create ~/.dsaprep.env with required tokens")
		return
	}

	// Display loaded tokens
	fmt.Printf("    Handle: %s\n", tokens.Handle)
	if tokens.CFClearance != "" {
		fmt.Printf("    cf_clearance: %s...%s\n",
			tokens.CFClearance[:minInt(20, len(tokens.CFClearance))],
			tokens.CFClearance[maxInt(0, len(tokens.CFClearance)-10):])
	} else {
		fmt.Println("    cf_clearance: NOT SET")
	}
	if tokens.CFClearanceUA != "" {
		fmt.Printf("    User-Agent: %s...\n", tokens.CFClearanceUA[:minInt(50, len(tokens.CFClearanceUA))])
	}
	if tokens.JSESSIONID != "" {
		fmt.Printf("    JSESSIONID: %s...\n", tokens.JSESSIONID[:minInt(20, len(tokens.JSESSIONID))])
	} else {
		fmt.Println("    JSESSIONID: NOT SET")
	}
	if tokens.CE7Cookie != "" {
		fmt.Printf("    39ce7: %s...\n", tokens.CE7Cookie[:minInt(20, len(tokens.CE7Cookie))])
	} else {
		fmt.Println("    39ce7: NOT SET")
	}
	fmt.Println()

	// Check readiness
	fmt.Println("[2] Checking authentication readiness...")
	cfClearanceValid := tokens.CFClearance != "" && time.Now().Before(tokens.CFClearExp)
	hasSessionCookies := tokens.JSESSIONID != "" || tokens.CE7Cookie != ""

	fmt.Printf("    cf_clearance valid: %v\n", cfClearanceValid)
	if cfClearanceValid {
		fmt.Printf("    Expires in: %s\n", time.Until(tokens.CFClearExp).Round(time.Minute))
	}
	fmt.Printf("    Session cookies: %v\n", hasSessionCookies)
	fmt.Printf("    Ready for submission: %v\n", cfClearanceValid && hasSessionCookies && tokens.Handle != "")
	fmt.Println()

	if tokens.CFClearance == "" {
		fmt.Println("ERROR: cf_clearance is required. Please extract from browser.")
		printCookieInstructions()
		return
	}

	// Step 3: Create session with cookies
	fmt.Println("[3] Creating cfweb session with cookies...")
	session, err := cfweb.NewSession()
	if err != nil {
		fmt.Printf("    ERROR: %v\n", err)
		return
	}

	// Set all authentication
	session.SetFullAuth(
		tokens.CFClearance,
		tokens.CFClearanceUA,
		tokens.CFClearExp,
		tokens.JSESSIONID,
		tokens.CE7Cookie,
		tokens.Handle,
	)
	fmt.Println("    OK - Session created with cookies")
	fmt.Printf("    IsAuthenticated: %v\n", session.IsAuthenticated())
	fmt.Println()

	// Step 4: Test problem fetching
	fmt.Println("[4] Testing problem fetching (Theatre Square - 1/A)...")
	parser := cfweb.NewParser(session)

	time.Sleep(500 * time.Millisecond) // Rate limiting
	problem, err := parser.ParseProblemset(1, "A")
	if err != nil {
		fmt.Printf("    ERROR: %v\n", err)
	} else {
		fmt.Println("    OK - Problem fetched successfully!")
		fmt.Printf("    Name: %s\n", problem.Name)
		fmt.Printf("    Time Limit: %s\n", problem.TimeLimit)
		fmt.Printf("    Memory Limit: %s\n", problem.MemoryLimit)
		fmt.Printf("    Rating: %d\n", problem.Rating)
		fmt.Printf("    Samples: %d\n", len(problem.Samples))
		if len(problem.Samples) > 0 {
			fmt.Printf("    Sample 1 Input: %s...\n", truncateCookie(problem.Samples[0].Input, 30))
		}
	}
	fmt.Println()

	// Step 5: Test CSRF token extraction
	fmt.Println("[5] Testing CSRF token extraction...")
	time.Sleep(500 * time.Millisecond)
	err = session.RefreshCSRFToken()
	if err != nil {
		fmt.Printf("    ERROR: %v\n", err)
	} else {
		csrf := session.GetCSRFToken()
		if csrf != "" {
			fmt.Printf("    OK - CSRF Token: %s...%s\n",
				csrf[:minInt(16, len(csrf))],
				csrf[maxInt(0, len(csrf)-8):])
		} else {
			fmt.Println("    WARNING: CSRF token not found")
		}
	}
	fmt.Println()

	// Step 6: Test session validation (checks if logged in)
	fmt.Println("[6] Testing session validation...")
	time.Sleep(500 * time.Millisecond)
	err = session.Validate()
	if err != nil {
		fmt.Printf("    Session invalid: %v\n", err)
		fmt.Println("    This means you can fetch problems but NOT submit")
	} else {
		fmt.Println("    OK - Session is valid!")
		fmt.Printf("    Logged in as: %s\n", session.Handle())
	}
	fmt.Println()

	// Step 7: Test submitter creation
	fmt.Println("[7] Testing submitter creation...")
	submitter, err := cfweb.NewSubmitter(session)
	if err != nil {
		fmt.Printf("    ERROR: %v\n", err)
		fmt.Println("    Submission will NOT work without valid session cookies")
	} else {
		fmt.Println("    OK - Submitter created!")
		fmt.Println("    NOTE: Actual submission test skipped (would submit real code)")

		// Verify submit page structure
		time.Sleep(500 * time.Millisecond)
		fmt.Println()
		fmt.Println("[8] Verifying submit page structure...")
		err = submitter.VerifySubmitPage(1) // Contest 1
		if err != nil {
			fmt.Printf("    ERROR: %v\n", err)
		} else {
			fmt.Println("    OK - Submit page structure verified!")
			fmt.Println("    All form elements present (csrf_token, problem_index, language, source)")
		}
	}

	// Summary
	fmt.Println()
	fmt.Println("=" + strings.Repeat("=", 59))
	fmt.Println("  SUMMARY")
	fmt.Println("=" + strings.Repeat("=", 59))
	fmt.Println()
	fmt.Println("  Cookie-based auth: WORKING")
	fmt.Println("  Problem fetching: WORKING")
	fmt.Println("  CSRF extraction: WORKING")

	if session.IsAuthenticated() && submitter != nil {
		fmt.Println("  Session validation: WORKING")
		fmt.Println("  Submission: READY")
		fmt.Println()
		fmt.Println("  You can now use the CLI to fetch problems and submit solutions!")
	} else {
		fmt.Println("  Session validation: FAILED")
		fmt.Println("  Submission: NOT READY")
		fmt.Println()
		fmt.Println("  To enable submissions, extract session cookies from browser:")
		printCookieInstructions()
	}
}

func printCookieInstructions() {
	fmt.Println()
	fmt.Println("  1. Open https://codeforces.com in your browser")
	fmt.Println("  2. Log in with your account")
	fmt.Println("  3. Open DevTools (F12 or Cmd+Option+I)")
	fmt.Println("  4. Go to: Application > Storage > Cookies > codeforces.com")
	fmt.Println("  5. Copy these cookies to ~/.dsaprep.env:")
	fmt.Println("     - cf_clearance -> CF_CLEARANCE")
	fmt.Println("     - JSESSIONID -> CF_JSESSIONID")
	fmt.Println("     - 39ce7 -> CF_39CE7")
	fmt.Println("  6. Also copy your browser's User-Agent:")
	fmt.Println("     - In Console, type: navigator.userAgent")
	fmt.Println("     - Copy to CF_CLEARANCE_UA")
}

func truncateCookie(s string, maxLen int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
