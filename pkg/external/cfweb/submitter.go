package cfweb

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Pre-compiled regexes for parsing submission results
var (
	reTimeMs   = regexp.MustCompile(`(\d+)\s*ms`)
	reMemoryKB = regexp.MustCompile(`(\d+)\s*KB`)
	reMemoryMB = regexp.MustCompile(`(\d+)\s*MB`)
)

// Submitter handles solution submission to CF
type Submitter struct {
	session *Session
}

// NewSubmitter creates a new submitter with an authenticated session
// Requires session to have cf_clearance + session cookies set
func NewSubmitter(session *Session) (*Submitter, error) {
	if !session.IsReadyForSubmission() {
		return nil, fmt.Errorf("session not ready for submission - need cf_clearance + session cookies + handle")
	}
	return &Submitter{session: session}, nil
}

// SubmissionResult contains the result of a submission
type SubmissionResult struct {
	SubmissionID int64
	ContestID    int
	ProblemIndex string
	Verdict      string
	Time         time.Duration
	Memory       int64 // bytes
	PassedTests  int
	SubmittedAt  time.Time
	Status       string
}

// Submit submits a solution to a problem
func (s *Submitter) Submit(contestID int, problemIndex string, langID int, sourceCode string) (*SubmissionResult, error) {
	// Construct submit URL
	submitURL := fmt.Sprintf("%s/contest/%d/submit", BaseURL, contestID)

	// Get the submit page first to extract CSRF token
	resp, err := s.get(submitURL)
	if err != nil {
		return nil, fmt.Errorf("get submit page: %w", err)
	}
	defer resp.Body.Close()

	// Use bounded reader to prevent OOM from large responses
	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxPageSize))
	if err != nil {
		return nil, fmt.Errorf("read submit page: %w", err)
	}

	// Extract CSRF token
	csrfToken := extractCSRFToken(string(body))
	if csrfToken == "" {
		return nil, fmt.Errorf("csrf token not found")
	}

	// Extract FTAA and BFAA
	ftaa := extractHiddenInput(string(body), "ftaa")
	bfaa := extractHiddenInput(string(body), "bfaa")

	// Prepare submission form
	form := url.Values{}
	form.Set("csrf_token", csrfToken)
	form.Set("action", "submitSolutionFormSubmitted")
	form.Set("submittedProblemIndex", problemIndex)
	form.Set("programTypeId", strconv.Itoa(langID))
	form.Set("source", sourceCode)
	form.Set("tabSize", "4")
	form.Set("sourceFile", "")
	if ftaa != "" {
		form.Set("ftaa", ftaa)
	}
	if bfaa != "" {
		form.Set("bfaa", bfaa)
	}

	// Submit the solution
	req, err := http.NewRequest(http.MethodPost, submitURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create submit request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Referer", submitURL)
	req.Header.Set("Origin", BaseURL)

	resp, err = s.session.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("submit solution: %w", err)
	}
	defer resp.Body.Close()

	// Check for redirect to my submissions page (indicates success)
	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusSeeOther {
		location := resp.Header.Get("Location")
		if strings.Contains(location, "/my") || strings.Contains(location, "/status") {
			// Extract submission ID from my submissions page
			return s.getLatestSubmission(contestID, problemIndex)
		}
	}

	// Read response to check for errors (bounded to prevent OOM)
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, MaxPageSize))
	respStr := string(respBody)

	// Check for common error messages
	if strings.Contains(respStr, "You have submitted exactly the same code before") {
		return nil, fmt.Errorf("duplicate submission: you have submitted exactly the same code before")
	}
	if strings.Contains(respStr, "Source code is too long") {
		return nil, fmt.Errorf("source code is too long")
	}
	if strings.Contains(respStr, "You are not allowed to submit") {
		return nil, fmt.Errorf("you are not allowed to submit to this contest")
	}
	if strings.Contains(respStr, "Contest is over") {
		return nil, fmt.Errorf("contest is over")
	}

	// Try to extract submission ID from response
	if resp.StatusCode == http.StatusOK {
		return s.getLatestSubmission(contestID, problemIndex)
	}

	return nil, fmt.Errorf("submission failed (status %d)", resp.StatusCode)
}

// SubmitToGym submits a solution to a gym problem
func (s *Submitter) SubmitToGym(gymID int, problemIndex string, langID int, sourceCode string) (*SubmissionResult, error) {
	submitURL := fmt.Sprintf("%s/gym/%d/submit", BaseURL, gymID)

	// Similar logic to Submit, but for gym
	resp, err := s.get(submitURL)
	if err != nil {
		return nil, fmt.Errorf("get gym submit page: %w", err)
	}
	defer resp.Body.Close()

	// Use bounded reader to prevent OOM from large responses
	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxPageSize))
	if err != nil {
		return nil, fmt.Errorf("read gym submit page: %w", err)
	}

	csrfToken := extractCSRFToken(string(body))
	if csrfToken == "" {
		return nil, fmt.Errorf("csrf token not found")
	}

	ftaa := extractHiddenInput(string(body), "ftaa")
	bfaa := extractHiddenInput(string(body), "bfaa")

	form := url.Values{}
	form.Set("csrf_token", csrfToken)
	form.Set("action", "submitSolutionFormSubmitted")
	form.Set("submittedProblemIndex", problemIndex)
	form.Set("programTypeId", strconv.Itoa(langID))
	form.Set("source", sourceCode)
	form.Set("tabSize", "4")
	if ftaa != "" {
		form.Set("ftaa", ftaa)
	}
	if bfaa != "" {
		form.Set("bfaa", bfaa)
	}

	req, err := http.NewRequest(http.MethodPost, submitURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create gym submit request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Referer", submitURL)
	req.Header.Set("Origin", BaseURL)

	resp, err = s.session.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("submit gym solution: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusSeeOther {
		return s.getLatestGymSubmission(gymID, problemIndex)
	}

	return nil, fmt.Errorf("gym submission failed (status %d)", resp.StatusCode)
}

// getLatestSubmission fetches the latest submission from my submissions
func (s *Submitter) getLatestSubmission(contestID int, problemIndex string) (*SubmissionResult, error) {
	myURL := fmt.Sprintf("%s/contest/%d/my", BaseURL, contestID)

	resp, err := s.get(myURL)
	if err != nil {
		return nil, fmt.Errorf("get my submissions: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse submissions page: %w", err)
	}

	// Find the first submission row
	row := doc.Find("table.status-frame-datatable tr[data-submission-id]").First()
	if row.Length() == 0 {
		return nil, fmt.Errorf("no submissions found")
	}

	return parseSubmissionRow(row, contestID)
}

// getLatestGymSubmission fetches the latest gym submission
func (s *Submitter) getLatestGymSubmission(gymID int, problemIndex string) (*SubmissionResult, error) {
	myURL := fmt.Sprintf("%s/gym/%d/my", BaseURL, gymID)

	resp, err := s.get(myURL)
	if err != nil {
		return nil, fmt.Errorf("get gym submissions: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse gym submissions page: %w", err)
	}

	row := doc.Find("table.status-frame-datatable tr[data-submission-id]").First()
	if row.Length() == 0 {
		return nil, fmt.Errorf("no gym submissions found")
	}

	return parseSubmissionRow(row, gymID)
}

// WaitForVerdict waits for the submission to be judged
func (s *Submitter) WaitForVerdict(submissionID int64, contestID int, timeout time.Duration) (*SubmissionResult, error) {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		result, err := s.GetSubmission(submissionID, contestID)
		if err != nil {
			return nil, err
		}

		// Check if judging is complete
		if result.Status != "In queue" && result.Status != "Running" {
			return result, nil
		}

		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("timeout waiting for verdict")
}

// GetSubmission gets a specific submission's status
func (s *Submitter) GetSubmission(submissionID int64, contestID int) (*SubmissionResult, error) {
	statusURL := fmt.Sprintf("%s/contest/%d/submission/%d", BaseURL, contestID, submissionID)

	resp, err := s.get(statusURL)
	if err != nil {
		return nil, fmt.Errorf("get submission status: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse submission page: %w", err)
	}

	// Find submission details
	verdictSpan := doc.Find(".verdict-accepted, .verdict-rejected, .verdict-waiting")
	verdict := strings.TrimSpace(verdictSpan.Text())

	result := &SubmissionResult{
		SubmissionID: submissionID,
		ContestID:    contestID,
		Verdict:      verdict,
		SubmittedAt:  time.Now(),
	}

	// Parse time and memory from the info table
	doc.Find(".datatable tr td").Each(func(i int, sel *goquery.Selection) {
		text := strings.TrimSpace(sel.Text())
		if strings.Contains(text, "ms") {
			result.Time = parseTime(text)
		}
		if strings.Contains(text, "KB") || strings.Contains(text, "MB") {
			result.Memory = parseMemory(text)
		}
	})

	// Determine status
	switch {
	case verdict == "" || strings.Contains(verdict, "queue"):
		result.Status = "In queue"
	case strings.Contains(verdict, "Running"):
		result.Status = "Running"
	case strings.Contains(verdict, "Accepted"):
		result.Status = "Accepted"
		result.Verdict = "OK"
	default:
		result.Status = "Judged"
	}

	return result, nil
}

// get makes a GET request
func (s *Submitter) get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	return s.session.client.Do(req)
}

// parseSubmissionRow parses a submission table row
func parseSubmissionRow(row *goquery.Selection, contestID int) (*SubmissionResult, error) {
	submissionIDStr, exists := row.Attr("data-submission-id")
	if !exists {
		return nil, fmt.Errorf("submission ID not found")
	}

	submissionID, err := strconv.ParseInt(submissionIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid submission ID: %w", err)
	}

	// Extract problem index
	problemCell := row.Find("td.id-cell").First()
	problemIndex := strings.TrimSpace(problemCell.Text())

	// Extract verdict
	verdictCell := row.Find("td.status-cell").First()
	verdict := strings.TrimSpace(verdictCell.Find(".verdict-accepted, .verdict-rejected").Text())
	if verdict == "" {
		verdict = strings.TrimSpace(verdictCell.Text())
	}

	// Extract time
	timeCell := row.Find("td.time-consumed-cell").First()
	timeText := strings.TrimSpace(timeCell.Text())

	// Extract memory
	memoryCell := row.Find("td.memory-consumed-cell").First()
	memoryText := strings.TrimSpace(memoryCell.Text())

	result := &SubmissionResult{
		SubmissionID: submissionID,
		ContestID:    contestID,
		ProblemIndex: problemIndex,
		Verdict:      normalizeVerdict(verdict),
		Time:         parseTime(timeText),
		Memory:       parseMemory(memoryText),
		SubmittedAt:  time.Now(),
	}

	// Determine status
	if verdict == "" || strings.Contains(verdict, "queue") {
		result.Status = "In queue"
	} else if strings.Contains(verdict, "Running") {
		result.Status = "Running"
	} else {
		result.Status = "Judged"
	}

	return result, nil
}

// parseTime parses time string like "46 ms"
func parseTime(text string) time.Duration {
	matches := reTimeMs.FindStringSubmatch(text)
	if len(matches) > 1 {
		ms, _ := strconv.Atoi(matches[1])
		return time.Duration(ms) * time.Millisecond
	}
	return 0
}

// parseMemory parses memory string like "1024 KB" or "1 MB"
func parseMemory(text string) int64 {
	if matches := reMemoryMB.FindStringSubmatch(text); len(matches) > 1 {
		mb, _ := strconv.ParseInt(matches[1], 10, 64)
		return mb * 1024 * 1024
	}

	if matches := reMemoryKB.FindStringSubmatch(text); len(matches) > 1 {
		kb, _ := strconv.ParseInt(matches[1], 10, 64)
		return kb * 1024
	}

	return 0
}

// normalizeVerdict converts CF verdict text to standard format
func normalizeVerdict(verdict string) string {
	verdict = strings.TrimSpace(verdict)

	switch {
	case strings.Contains(verdict, "Accepted"):
		return "OK"
	case strings.Contains(verdict, "Wrong answer"):
		return "WRONG_ANSWER"
	case strings.Contains(verdict, "Time limit"):
		return "TIME_LIMIT_EXCEEDED"
	case strings.Contains(verdict, "Memory limit"):
		return "MEMORY_LIMIT_EXCEEDED"
	case strings.Contains(verdict, "Runtime error"):
		return "RUNTIME_ERROR"
	case strings.Contains(verdict, "Compilation error"):
		return "COMPILATION_ERROR"
	case strings.Contains(verdict, "Presentation error"):
		return "PRESENTATION_ERROR"
	case strings.Contains(verdict, "Idleness"):
		return "IDLENESS_LIMIT_EXCEEDED"
	case strings.Contains(verdict, "Hacked"):
		return "CHALLENGED"
	default:
		return verdict
	}
}

// VerifySubmitPage checks if the submit page structure is valid
func (s *Submitter) VerifySubmitPage(contestID int) error {
	submitURL := fmt.Sprintf("%s/contest/%d/submit", BaseURL, contestID)

	resp, err := s.get(submitURL)
	if err != nil {
		return fmt.Errorf("get submit page: %w", err)
	}
	defer resp.Body.Close()

	// Use bounded reader to prevent OOM from large responses
	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxPageSize))
	if err != nil {
		return fmt.Errorf("read submit page: %w", err)
	}

	bodyStr := string(body)

	// Check for required elements
	checks := []struct {
		name    string
		pattern string
	}{
		{"csrf_token", `name="csrf_token"`},
		{"problem_index", `name="submittedProblemIndex"`},
		{"language_select", `name="programTypeId"`},
		{"source_code", `name="source"`},
		{"submit_button", `type="submit"`},
	}

	var missing []string
	for _, check := range checks {
		if !strings.Contains(bodyStr, check.pattern) {
			missing = append(missing, check.name)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("submit page missing elements: %s", strings.Join(missing, ", "))
	}

	return nil
}
