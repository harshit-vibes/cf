package codeforces

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

const (
	// BaseURL is the Codeforces API base URL
	BaseURL = "https://codeforces.com/api"

	// DefaultCacheTTL is the default cache TTL (5 minutes)
	DefaultCacheTTL = 5 * time.Minute

	// DefaultTimeout is the default HTTP timeout
	DefaultTimeout = 30 * time.Second
)

// Client is a Codeforces API client
type Client struct {
	httpClient *http.Client
	limiter    *rate.Limiter
	cache      *Cache
	baseURL    string
}

// NewClient creates a new Codeforces API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: DefaultTimeout},
		// Rate limit: 5 requests per second (1 request every 200ms)
		limiter: rate.NewLimiter(rate.Every(200*time.Millisecond), 1),
		cache:   NewCache(DefaultCacheTTL),
		baseURL: BaseURL,
	}
}

// GetProblems fetches all problems from Codeforces
func (c *Client) GetProblems(ctx context.Context) (*ProblemsResult, error) {
	cacheKey := "problems:all"

	// Check cache first
	if cached, ok := c.cache.Get(cacheKey); ok {
		return cached.(*ProblemsResult), nil
	}

	// Wait for rate limiter
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter: %w", err)
	}

	// Make request
	url := c.baseURL + "/problemset.problems"
	resp, err := c.doGet(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode response
	var apiResp APIResponse[ProblemsResult]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if apiResp.Status != "OK" {
		return nil, fmt.Errorf("API error: %s", apiResp.Comment)
	}

	// Cache the result
	c.cache.Set(cacheKey, &apiResp.Result)

	return &apiResp.Result, nil
}

// GetProblemsWithTags fetches problems filtered by tags
func (c *Client) GetProblemsWithTags(ctx context.Context, tags []string) (*ProblemsResult, error) {
	cacheKey := fmt.Sprintf("problems:tags:%s", strings.Join(tags, ","))

	// Check cache first
	if cached, ok := c.cache.Get(cacheKey); ok {
		return cached.(*ProblemsResult), nil
	}

	// Wait for rate limiter
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter: %w", err)
	}

	// Make request
	url := fmt.Sprintf("%s/problemset.problems?tags=%s", c.baseURL, strings.Join(tags, ";"))
	resp, err := c.doGet(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode response
	var apiResp APIResponse[ProblemsResult]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if apiResp.Status != "OK" {
		return nil, fmt.Errorf("API error: %s", apiResp.Comment)
	}

	// Cache the result
	c.cache.Set(cacheKey, &apiResp.Result)

	return &apiResp.Result, nil
}

// GetUserInfo fetches user information
func (c *Client) GetUserInfo(ctx context.Context, handles ...string) ([]User, error) {
	if len(handles) == 0 {
		return nil, fmt.Errorf("at least one handle is required")
	}

	handleStr := strings.Join(handles, ";")
	cacheKey := fmt.Sprintf("user:info:%s", handleStr)

	// Check cache first
	if cached, ok := c.cache.Get(cacheKey); ok {
		return cached.([]User), nil
	}

	// Wait for rate limiter
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter: %w", err)
	}

	// Make request
	url := fmt.Sprintf("%s/user.info?handles=%s", c.baseURL, handleStr)
	resp, err := c.doGet(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode response
	var apiResp APIResponse[[]User]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if apiResp.Status != "OK" {
		return nil, fmt.Errorf("API error: %s", apiResp.Comment)
	}

	// Cache the result
	c.cache.Set(cacheKey, apiResp.Result)

	return apiResp.Result, nil
}

// GetUserStatus fetches user's submission history
func (c *Client) GetUserStatus(ctx context.Context, handle string, count int) ([]Submission, error) {
	cacheKey := fmt.Sprintf("user:status:%s:%d", handle, count)

	// Check cache first (shorter TTL for submissions)
	if cached, ok := c.cache.Get(cacheKey); ok {
		return cached.([]Submission), nil
	}

	// Wait for rate limiter
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter: %w", err)
	}

	// Make request
	url := fmt.Sprintf("%s/user.status?handle=%s&count=%d", c.baseURL, handle, count)
	resp, err := c.doGet(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode response
	var apiResp APIResponse[[]Submission]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if apiResp.Status != "OK" {
		return nil, fmt.Errorf("API error: %s", apiResp.Comment)
	}

	// Cache the result
	c.cache.Set(cacheKey, apiResp.Result)

	return apiResp.Result, nil
}

// GetUserRating fetches user's rating history
func (c *Client) GetUserRating(ctx context.Context, handle string) ([]RatingChange, error) {
	cacheKey := fmt.Sprintf("user:rating:%s", handle)

	// Check cache first
	if cached, ok := c.cache.Get(cacheKey); ok {
		return cached.([]RatingChange), nil
	}

	// Wait for rate limiter
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter: %w", err)
	}

	// Make request
	url := fmt.Sprintf("%s/user.rating?handle=%s", c.baseURL, handle)
	resp, err := c.doGet(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode response
	var apiResp APIResponse[[]RatingChange]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if apiResp.Status != "OK" {
		return nil, fmt.Errorf("API error: %s", apiResp.Comment)
	}

	// Cache the result
	c.cache.Set(cacheKey, apiResp.Result)

	return apiResp.Result, nil
}

// GetContestList fetches the list of contests
func (c *Client) GetContestList(ctx context.Context, gym bool) ([]Contest, error) {
	cacheKey := fmt.Sprintf("contest:list:%v", gym)

	// Check cache first
	if cached, ok := c.cache.Get(cacheKey); ok {
		return cached.([]Contest), nil
	}

	// Wait for rate limiter
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter: %w", err)
	}

	// Make request
	url := fmt.Sprintf("%s/contest.list?gym=%v", c.baseURL, gym)
	resp, err := c.doGet(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode response
	var apiResp APIResponse[[]Contest]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if apiResp.Status != "OK" {
		return nil, fmt.Errorf("API error: %s", apiResp.Comment)
	}

	// Cache the result
	c.cache.Set(cacheKey, apiResp.Result)

	return apiResp.Result, nil
}

// ClearCache clears the client's cache
func (c *Client) ClearCache() {
	c.cache.Clear()
}

// doGet performs a GET request with context
func (c *Client) doGet(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("User-Agent", "DSAPrep-CLI/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil
}
