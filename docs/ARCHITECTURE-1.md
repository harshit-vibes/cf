# DSA Prep - Architecture 1 (Local CLI)

> Codeforces API wrapper with beautiful TUI, AI-powered assistance, and offline capability

---

## Table of Contents

1. [Overview](#overview)
2. [Development Phases](#development-phases)
3. [Phase 1: Codeforces Wrapper](#phase-1-codeforces-wrapper)
4. [Phase 2: AI-Powered Features](#phase-2-ai-powered-features)
5. [Phase 3: Offline Mode](#phase-3-offline-mode)
6. [CLI Structure](#cli-structure)
7. [TUI Application](#tui-application)
8. [Installation](#installation)
9. [Development](#development)
10. [Configuration](#configuration)
11. [Next Steps](#next-steps)
12. [References](#references)

---

## Overview

### What We're Building

A beautiful terminal-based Codeforces experience. Start simple—wrap the Codeforces API with an exceptional TUI. Then layer on AI assistance and offline capabilities.

### Philosophy

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        BUILD ITERATIVELY                                │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   Phase 1                Phase 2                Phase 3                 │
│   ────────               ────────               ────────                │
│                                                                         │
│   ┌─────────┐           ┌─────────┐           ┌─────────┐              │
│   │   TUI   │           │   TUI   │           │   TUI   │              │
│   │    +    │    ──►    │    +    │    ──►    │    +    │              │
│   │   API   │           │   API   │           │   API   │              │
│   └─────────┘           │    +    │           │    +    │              │
│                         │   AI    │           │   AI    │              │
│   "Make it work"        └─────────┘           │    +    │              │
│                                               │ Offline │              │
│                         "Make it smart"       └─────────┘              │
│                                                                         │
│                                               "Make it resilient"       │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### Key Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Approach | Online-first | Start simple, add complexity later |
| CLI Framework | Cobra | Industry standard, subcommands |
| TUI Framework | Bubble Tea | Full terminal control, beautiful UI |
| Problem Source | Codeforces API | 8000+ problems, free, no auth needed |
| AI Provider | Amazon Bedrock | Claude models, pay-per-use, no API key mgmt |
| Local Storage | SQLite | Only added in Phase 3, embedded |
| Distribution | Go binary | Single file, cross-platform |

### What Each Phase Delivers

| Phase | User Value | Technical Scope |
|-------|------------|-----------------|
| **Phase 1** | Browse & practice CF problems with beautiful TUI | Cobra + Bubble Tea + HTTP client |
| **Phase 2** | Get AI hints, explanations, code review | + Amazon Bedrock integration |
| **Phase 3** | Work without internet, sync when online | + SQLite cache, sync logic |

---

## Development Phases

### Phase 1: Codeforces Wrapper 🎯

**Goal**: Best-in-class terminal experience for Codeforces

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         PHASE 1 ARCHITECTURE                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│                              TERMINAL                                   │
│   ┌─────────────────────────────────────────────────────────────────┐  │
│   │                                                                  │  │
│   │                      DSA PREP CLI                                │  │
│   │                      (Go Binary)                                 │  │
│   │                                                                  │  │
│   │   ┌────────────────────────────────────────────────────────┐    │  │
│   │   │                    COBRA COMMANDS                       │    │  │
│   │   │                                                         │    │  │
│   │   │   dsaprep              → Launch TUI                     │    │  │
│   │   │   dsaprep problem      → Quick problem lookup           │    │  │
│   │   │   dsaprep random       → Random problem by filters      │    │  │
│   │   │   dsaprep stats        → Your CF stats                  │    │  │
│   │   └────────────────────────────────────────────────────────┘    │  │
│   │                              │                                   │  │
│   │                              ▼                                   │  │
│   │   ┌────────────────────────────────────────────────────────┐    │  │
│   │   │                   BUBBLE TEA TUI                        │    │  │
│   │   │                                                         │    │  │
│   │   │   ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌─────────┐   │    │  │
│   │   │   │Dashboard │ │ Problems │ │ Practice │ │  Stats  │   │    │  │
│   │   │   └──────────┘ └──────────┘ └──────────┘ └─────────┘   │    │  │
│   │   └────────────────────────────────────────────────────────┘    │  │
│   │                              │                                   │  │
│   └──────────────────────────────┼───────────────────────────────────┘  │
│                                  │                                      │
│   ┌──────────────────────────────┼───────────────────────────────────┐  │
│   │                              ▼                                   │  │
│   │   ┌────────────────────────────────────────────────────────┐    │  │
│   │   │              CODEFORCES API CLIENT                      │    │  │
│   │   │                                                         │    │  │
│   │   │   • Rate limiting (5 req/sec)                          │    │  │
│   │   │   • Response caching (in-memory)                       │    │  │
│   │   │   • Error handling & retries                           │    │  │
│   │   └────────────────────────────────────────────────────────┘    │  │
│   │                         HTTP CLIENT                              │  │
│   └──────────────────────────────────────────────────────────────────┘  │
│                                  │                                      │
└──────────────────────────────────┼──────────────────────────────────────┘
                                   │ HTTPS
                                   ▼
                    ┌───────────────────────────────┐
                    │       CODEFORCES API          │
                    │   codeforces.com/api/         │
                    │                               │
                    │   • problemset.problems       │
                    │   • user.info                 │
                    │   • user.status               │
                    │   • user.rating               │
                    └───────────────────────────────┘
```

**Features**:
- Browse all 8000+ Codeforces problems
- Filter by difficulty (800-3500), tags, solved count
- View problem in terminal (rendered markdown)
- Practice timer with tracking
- Open problem in browser
- View your CF profile stats
- Keyboard-driven navigation

**Tech Stack**:
```go
// go.mod (Phase 1)
require (
    github.com/spf13/cobra          v1.8.0   // CLI
    github.com/charmbracelet/bubbletea v0.26.0 // TUI
    github.com/charmbracelet/bubbles   v0.18.0 // Components
    github.com/charmbracelet/lipgloss  v0.11.0 // Styling
    github.com/charmbracelet/glamour   v0.7.0  // Markdown
)
```

---

### Phase 2: AI-Powered Features 🧠

**Goal**: Make practicing smarter with AI assistance

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         PHASE 2 ARCHITECTURE                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│                              TERMINAL                                   │
│   ┌─────────────────────────────────────────────────────────────────┐  │
│   │                                                                  │  │
│   │                      DSA PREP CLI                                │  │
│   │                                                                  │  │
│   │   ┌────────────────────────────────────────────────────────┐    │  │
│   │   │                   BUBBLE TEA TUI                        │    │  │
│   │   │                                                         │    │  │
│   │   │   Practice View                                         │    │  │
│   │   │   ┌─────────────────────────────────────────────────┐  │    │  │
│   │   │   │  Problem: Two Sum                               │  │    │  │
│   │   │   │  Rating: 800 │ Tags: math, implementation       │  │    │  │
│   │   │   │                                                 │  │    │  │
│   │   │   │  [h] Hint  [e] Explain  [r] Review Code        │  │    │  │
│   │   │   │  [s] Similar  [p] Patterns                     │  │    │  │
│   │   │   └─────────────────────────────────────────────────┘  │    │  │
│   │   │                                                         │    │  │
│   │   │   AI Assistant Panel                                    │    │  │
│   │   │   ┌─────────────────────────────────────────────────┐  │    │  │
│   │   │   │  💡 Hint 1/3:                                   │  │    │  │
│   │   │   │  Consider using a hash map to store seen        │  │    │  │
│   │   │   │  numbers. What would you store as the value?    │  │    │  │
│   │   │   │                                                 │  │    │  │
│   │   │   │  [Enter] Next hint  [Esc] Close                │  │    │  │
│   │   │   └─────────────────────────────────────────────────┘  │    │  │
│   │   └────────────────────────────────────────────────────────┘    │  │
│   │                              │                                   │  │
│   └──────────────────────────────┼───────────────────────────────────┘  │
│                                  │                                      │
│          ┌───────────────────────┴───────────────────────┐             │
│          │                                               │             │
│          ▼                                               ▼             │
│   ┌─────────────────┐                         ┌─────────────────┐      │
│   │  CODEFORCES     │                         │  AMAZON         │      │
│   │  API CLIENT     │                         │  BEDROCK        │      │
│   │                 │                         │                 │      │
│   │  Problems,      │                         │  Claude 3.5     │      │
│   │  User data      │                         │  Sonnet         │      │
│   └────────┬────────┘                         └────────┬────────┘      │
│            │                                           │               │
└────────────┼───────────────────────────────────────────┼───────────────┘
             │ HTTPS                                     │ HTTPS
             ▼                                           ▼
┌───────────────────────────┐             ┌───────────────────────────┐
│     CODEFORCES API        │             │     AWS BEDROCK           │
│  codeforces.com/api/      │             │  bedrock-runtime          │
└───────────────────────────┘             │                           │
                                          │  Models:                  │
                                          │  • Claude 3.5 Sonnet      │
                                          │  • Claude 3 Haiku         │
                                          └───────────────────────────┘
```

**AI Features**:

| Feature | Description | Trigger |
|---------|-------------|---------|
| **Progressive Hints** | 3-level hints (gentle → direct → solution approach) | `h` key |
| **Explain Solution** | Explain optimal approach after solving | `e` key |
| **Code Review** | Paste your code, get feedback | `r` key |
| **Pattern Recognition** | Identify algorithm patterns that apply | `p` key |
| **Similar Problems** | Suggest similar problems to practice | `s` key |
| **Complexity Analysis** | Analyze time/space complexity of approach | `c` key |

**Bedrock Integration**:
```go
// Additional dependencies for Phase 2
require (
    github.com/aws/aws-sdk-go-v2          v1.26.0
    github.com/aws/aws-sdk-go-v2/config   v1.27.0
    github.com/aws/aws-sdk-go-v2/service/bedrockruntime v1.7.0
)
```

**AWS Free Tier Note**:
- Bedrock is **pay-per-use** (no free tier)
- Claude 3 Haiku: ~$0.00025/1K input tokens, ~$0.00125/1K output
- Typical hint request: ~$0.001 (fraction of a cent)
- Make AI features **opt-in** with clear pricing info

---

### Phase 3: Offline Mode 📴

**Goal**: Work without internet, sync when connected

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         PHASE 3 ARCHITECTURE                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│                              TERMINAL                                   │
│   ┌─────────────────────────────────────────────────────────────────┐  │
│   │                                                                  │  │
│   │                      DSA PREP CLI                                │  │
│   │                                                                  │  │
│   │   ┌────────────────────────────────────────────────────────┐    │  │
│   │   │                   BUBBLE TEA TUI                        │    │  │
│   │   │                                                         │    │  │
│   │   │   Status Bar: [●] Online │ Last sync: 2 hours ago       │    │  │
│   │   │               [○] Offline │ 4521 problems cached        │    │  │
│   │   │                                                         │    │  │
│   │   └────────────────────────────────────────────────────────┘    │  │
│   │                              │                                   │  │
│   └──────────────────────────────┼───────────────────────────────────┘  │
│                                  │                                      │
│   ┌──────────────────────────────┼───────────────────────────────────┐  │
│   │                              ▼                                   │  │
│   │   ┌────────────────────────────────────────────────────────┐    │  │
│   │   │                   DATA LAYER                            │    │  │
│   │   │                                                         │    │  │
│   │   │   ┌─────────────────┐      ┌─────────────────┐         │    │  │
│   │   │   │  Online Mode    │      │  Offline Mode   │         │    │  │
│   │   │   │                 │      │                 │         │    │  │
│   │   │   │  CF API ────────┼──────▶ SQLite Cache   │         │    │  │
│   │   │   │       │         │      │       │         │         │    │  │
│   │   │   │       ▼         │      │       ▼         │         │    │  │
│   │   │   │  Fresh data     │      │  Cached data    │         │    │  │
│   │   │   └─────────────────┘      └─────────────────┘         │    │  │
│   │   │                                                         │    │  │
│   │   └────────────────────────────────────────────────────────┘    │  │
│   │                                                                  │  │
│   │   ┌────────────────────────────────────────────────────────┐    │  │
│   │   │                   SQLite Database                       │    │  │
│   │   │                   ~/.dsaprep/data.db                    │    │  │
│   │   │                                                         │    │  │
│   │   │   Tables:                                               │    │  │
│   │   │   • problems        (cached from CF)                    │    │  │
│   │   │   • tags            (problem tags)                      │    │  │
│   │   │   • user_progress   (local practice tracking)           │    │  │
│   │   │   • sessions        (practice sessions)                 │    │  │
│   │   │   • daily_stats     (streaks, daily counts)             │    │  │
│   │   │   • sync_meta       (last sync timestamps)              │    │  │
│   │   └────────────────────────────────────────────────────────┘    │  │
│   │                          LOCAL STORAGE                           │  │
│   └──────────────────────────────────────────────────────────────────┘  │
│                                  │                                      │
│                                  │ When online                          │
│                                  ▼                                      │
│   ┌──────────────────────────────────────────────────────────────────┐  │
│   │                        SYNC ENGINE                                │  │
│   │                                                                   │  │
│   │   On Launch:                                                      │  │
│   │   1. Check connectivity                                          │  │
│   │   2. If online + stale cache → background sync                   │  │
│   │   3. If offline → use cache                                      │  │
│   │                                                                   │  │
│   │   Manual: `dsaprep sync`                                         │  │
│   └──────────────────────────────────────────────────────────────────┘  │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**Offline Capabilities**:

| Feature | Online | Offline |
|---------|--------|---------|
| Browse problems | ✓ Fresh | ✓ Cached |
| Filter/search | ✓ | ✓ |
| Practice timer | ✓ | ✓ |
| Track progress | ✓ | ✓ (local) |
| AI hints | ✓ | ✗ |
| View CF profile | ✓ Fresh | ✓ Cached |
| Sync progress | ✓ | Queued |

**Additional Dependencies**:
```go
// Phase 3 additions
require (
    github.com/mattn/go-sqlite3  v1.14.22  // SQLite driver
    github.com/jmoiron/sqlx      v1.4.0    // SQL extensions
)
```

---

## Phase 1: Codeforces Wrapper

### Project Structure

```
dsaprep/
├── cmd/
│   └── dsaprep/
│       └── main.go                 # Entry point
│
├── internal/
│   ├── cmd/                        # Cobra commands
│   │   ├── root.go                 # Root → launches TUI
│   │   ├── problem.go              # problem subcommand
│   │   ├── random.go               # random problem
│   │   ├── stats.go                # user stats
│   │   └── version.go              # version info
│   │
│   ├── tui/                        # Bubble Tea TUI
│   │   ├── app.go                  # Main model
│   │   ├── keymap.go               # Key bindings
│   │   ├── views/
│   │   │   ├── dashboard.go        # Home view
│   │   │   ├── problems.go         # Problem browser
│   │   │   ├── practice.go         # Practice session
│   │   │   ├── problem_detail.go   # Single problem view
│   │   │   └── stats.go            # Statistics
│   │   ├── components/
│   │   │   ├── header.go           # App header
│   │   │   ├── footer.go           # Help footer
│   │   │   ├── table.go            # Problem table
│   │   │   ├── timer.go            # Practice timer
│   │   │   └── filter.go           # Filter bar
│   │   └── styles/
│   │       └── theme.go            # Lip Gloss styles
│   │
│   ├── codeforces/                 # CF API client
│   │   ├── client.go               # HTTP client
│   │   ├── types.go                # API types
│   │   ├── problems.go             # Problem endpoints
│   │   ├── user.go                 # User endpoints
│   │   └── cache.go                # In-memory cache
│   │
│   └── config/                     # Configuration
│       └── config.go               # Viper setup
│
├── .goreleaser.yml
├── Makefile
├── go.mod
└── go.sum
```

### Codeforces Client

```go
// internal/codeforces/client.go
package codeforces

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"

    "golang.org/x/time/rate"
)

const (
    BaseURL     = "https://codeforces.com/api"
    RateLimit   = 5 // requests per second
    CacheTTL    = 5 * time.Minute
)

type Client struct {
    httpClient *http.Client
    limiter    *rate.Limiter
    cache      *Cache
}

func NewClient() *Client {
    return &Client{
        httpClient: &http.Client{Timeout: 30 * time.Second},
        limiter:    rate.NewLimiter(rate.Every(200*time.Millisecond), 1),
        cache:      NewCache(CacheTTL),
    }
}

// GetProblems fetches all problems from Codeforces
func (c *Client) GetProblems(ctx context.Context) (*ProblemsResult, error) {
    // Check cache first
    if cached, ok := c.cache.Get("problems"); ok {
        return cached.(*ProblemsResult), nil
    }

    // Rate limit
    if err := c.limiter.Wait(ctx); err != nil {
        return nil, err
    }

    resp, err := c.httpClient.Get(BaseURL + "/problemset.problems")
    if err != nil {
        return nil, fmt.Errorf("fetch problems: %w", err)
    }
    defer resp.Body.Close()

    var apiResp APIResponse[ProblemsResult]
    if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
        return nil, fmt.Errorf("decode response: %w", err)
    }

    if apiResp.Status != "OK" {
        return nil, fmt.Errorf("API error: %s", apiResp.Comment)
    }

    // Cache the result
    c.cache.Set("problems", &apiResp.Result)

    return &apiResp.Result, nil
}

// GetUserInfo fetches user information
func (c *Client) GetUserInfo(ctx context.Context, handle string) (*User, error) {
    cacheKey := "user:" + handle
    if cached, ok := c.cache.Get(cacheKey); ok {
        return cached.(*User), nil
    }

    if err := c.limiter.Wait(ctx); err != nil {
        return nil, err
    }

    url := fmt.Sprintf("%s/user.info?handles=%s", BaseURL, handle)
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var apiResp APIResponse[[]User]
    if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
        return nil, err
    }

    if apiResp.Status != "OK" || len(apiResp.Result) == 0 {
        return nil, fmt.Errorf("user not found: %s", handle)
    }

    user := &apiResp.Result[0]
    c.cache.Set(cacheKey, user)

    return user, nil
}

// GetUserSubmissions fetches user's submission history
func (c *Client) GetUserSubmissions(ctx context.Context, handle string, count int) ([]Submission, error) {
    if err := c.limiter.Wait(ctx); err != nil {
        return nil, err
    }

    url := fmt.Sprintf("%s/user.status?handle=%s&count=%d", BaseURL, handle, count)
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var apiResp APIResponse[[]Submission]
    if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
        return nil, err
    }

    return apiResp.Result, nil
}
```

```go
// internal/codeforces/types.go
package codeforces

// APIResponse is the generic Codeforces API response wrapper
type APIResponse[T any] struct {
    Status  string `json:"status"`
    Comment string `json:"comment,omitempty"`
    Result  T      `json:"result"`
}

// ProblemsResult contains problems and their statistics
type ProblemsResult struct {
    Problems          []Problem          `json:"problems"`
    ProblemStatistics []ProblemStatistic `json:"problemStatistics"`
}

// Problem represents a Codeforces problem
type Problem struct {
    ContestID int      `json:"contestId"`
    Index     string   `json:"index"`
    Name      string   `json:"name"`
    Type      string   `json:"type"`
    Rating    int      `json:"rating,omitempty"`
    Tags      []string `json:"tags"`
}

// ProblemStatistic contains solve counts
type ProblemStatistic struct {
    ContestID   int `json:"contestId"`
    Index       string `json:"index"`
    SolvedCount int    `json:"solvedCount"`
}

// User represents a Codeforces user
type User struct {
    Handle       string `json:"handle"`
    Rating       int    `json:"rating"`
    MaxRating    int    `json:"maxRating"`
    Rank         string `json:"rank"`
    MaxRank      string `json:"maxRank"`
    Country      string `json:"country,omitempty"`
    Organization string `json:"organization,omitempty"`
    Contribution int    `json:"contribution"`
    FriendOf     int    `json:"friendOfCount"`
}

// Submission represents a user's submission
type Submission struct {
    ID                  int     `json:"id"`
    ContestID           int     `json:"contestId"`
    CreationTimeSeconds int64   `json:"creationTimeSeconds"`
    Problem             Problem `json:"problem"`
    Verdict             string  `json:"verdict"`
    PassedTestCount     int     `json:"passedTestCount"`
}

// Helper methods

func (p Problem) ID() string {
    return fmt.Sprintf("%d%s", p.ContestID, p.Index)
}

func (p Problem) URL() string {
    return fmt.Sprintf("https://codeforces.com/problemset/problem/%d/%s",
        p.ContestID, p.Index)
}

func (p Problem) DifficultyColor() string {
    switch {
    case p.Rating < 1200:
        return "#808080" // Gray (Newbie)
    case p.Rating < 1400:
        return "#008000" // Green (Pupil)
    case p.Rating < 1600:
        return "#03A89E" // Cyan (Specialist)
    case p.Rating < 1900:
        return "#0000FF" // Blue (Expert)
    case p.Rating < 2100:
        return "#AA00AA" // Violet (Candidate Master)
    case p.Rating < 2400:
        return "#FF8C00" // Orange (Master)
    default:
        return "#FF0000" // Red (Grandmaster+)
    }
}
```

```go
// internal/codeforces/cache.go
package codeforces

import (
    "sync"
    "time"
)

type cacheItem struct {
    value      interface{}
    expiration time.Time
}

type Cache struct {
    items map[string]cacheItem
    mu    sync.RWMutex
    ttl   time.Duration
}

func NewCache(ttl time.Duration) *Cache {
    c := &Cache{
        items: make(map[string]cacheItem),
        ttl:   ttl,
    }
    go c.cleanup()
    return c
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    item, exists := c.items[key]
    if !exists || time.Now().After(item.expiration) {
        return nil, false
    }
    return item.value, true
}

func (c *Cache) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.items[key] = cacheItem{
        value:      value,
        expiration: time.Now().Add(c.ttl),
    }
}

func (c *Cache) cleanup() {
    ticker := time.NewTicker(c.ttl)
    for range ticker.C {
        c.mu.Lock()
        for key, item := range c.items {
            if time.Now().After(item.expiration) {
                delete(c.items, key)
            }
        }
        c.mu.Unlock()
    }
}
```

### TUI Application

```go
// internal/tui/app.go
package tui

import (
    "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/yourusername/dsaprep/internal/codeforces"
)

type View int

const (
    DashboardView View = iota
    ProblemsView
    PracticeView
    ProblemDetailView
    StatsView
)

type Model struct {
    // Current view
    currentView View
    prevView    View

    // Sub-models
    dashboard     DashboardModel
    problems      ProblemsModel
    practice      PracticeModel
    problemDetail ProblemDetailModel
    stats         StatsModel

    // Shared state
    cfClient *codeforces.Client
    cfHandle string

    // Terminal dimensions
    width  int
    height int

    // Loading state
    loading bool
    err     error
}

func NewModel(cfHandle string) Model {
    client := codeforces.NewClient()

    return Model{
        currentView: DashboardView,
        cfClient:    client,
        cfHandle:    cfHandle,
        dashboard:   NewDashboardModel(),
        problems:    NewProblemsModel(client),
        practice:    NewPracticeModel(),
        stats:       NewStatsModel(client, cfHandle),
    }
}

func (m Model) Init() tea.Cmd {
    return tea.Batch(
        tea.EnterAltScreen,
        m.dashboard.Init(),
        m.loadProblems,
    )
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Global keybindings
        switch msg.String() {
        case "ctrl+c":
            return m, tea.Quit
        case "q", "esc":
            if m.currentView == DashboardView {
                return m, tea.Quit
            }
            m.currentView = m.prevView
            if m.currentView == 0 {
                m.currentView = DashboardView
            }
            return m, nil

        // Navigation (only when not in input mode)
        case "d":
            m.prevView = m.currentView
            m.currentView = DashboardView
        case "p":
            m.prevView = m.currentView
            m.currentView = ProblemsView
        case "s":
            m.prevView = m.currentView
            m.currentView = StatsView
        case "?":
            // Show help
        }

    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        // Propagate to sub-models
        m.dashboard = m.dashboard.SetSize(msg.Width, msg.Height)
        m.problems = m.problems.SetSize(msg.Width, msg.Height)

    case ProblemsLoadedMsg:
        m.problems = m.problems.SetProblems(msg.Problems, msg.Stats)
        m.loading = false

    case ErrorMsg:
        m.err = msg.Err
        m.loading = false

    case StartPracticeMsg:
        m.prevView = m.currentView
        m.currentView = PracticeView
        m.practice = NewPracticeModel().SetProblem(msg.Problem)
        return m, m.practice.Init()
    }

    // Update current view
    var cmd tea.Cmd
    switch m.currentView {
    case DashboardView:
        m.dashboard, cmd = m.dashboard.Update(msg)
    case ProblemsView:
        m.problems, cmd = m.problems.Update(msg)
    case PracticeView:
        m.practice, cmd = m.practice.Update(msg)
    case StatsView:
        m.stats, cmd = m.stats.Update(msg)
    }
    cmds = append(cmds, cmd)

    return m, tea.Batch(cmds...)
}

func (m Model) View() string {
    if m.loading {
        return m.renderLoading()
    }

    if m.err != nil {
        return m.renderError()
    }

    // Build layout
    header := m.renderHeader()
    footer := m.renderFooter()

    contentHeight := m.height - lipgloss.Height(header) - lipgloss.Height(footer)

    var content string
    switch m.currentView {
    case DashboardView:
        content = m.dashboard.View()
    case ProblemsView:
        content = m.problems.View()
    case PracticeView:
        content = m.practice.View()
    case StatsView:
        content = m.stats.View()
    }

    // Ensure content fills available space
    content = lipgloss.NewStyle().
        Height(contentHeight).
        Render(content)

    return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}

// Commands

func (m Model) loadProblems() tea.Msg {
    result, err := m.cfClient.GetProblems(context.Background())
    if err != nil {
        return ErrorMsg{Err: err}
    }
    return ProblemsLoadedMsg{
        Problems: result.Problems,
        Stats:    result.ProblemStatistics,
    }
}

func Run(cfHandle string) error {
    p := tea.NewProgram(
        NewModel(cfHandle),
        tea.WithAltScreen(),
        tea.WithMouseCellMotion(),
    )
    _, err := p.Run()
    return err
}
```

### Screen Designs

```
┌─────────────────────────────────────────────────────────────────────────┐
│  DSA Prep                                        tourist │ Expert │ 1847 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   ╭─ Welcome back! ────────────────────────────────────────────────────╮│
│   │                                                                    ││
│   │   🎯 Today's Goal: 5 problems                                      ││
│   │   ━━━━━━━━━━━━━━━━━━━━━━░░░░░░░░░░  2/5 completed                  ││
│   │                                                                    ││
│   │   📊 This Week                                                     ││
│   │   Mon  Tue  Wed  Thu  Fri  Sat  Sun                               ││
│   │    █    █    ░    █    ░    ░    ░                                ││
│   │    3    5    0    4    -    -    -                                ││
│   │                                                                    ││
│   ╰────────────────────────────────────────────────────────────────────╯│
│                                                                         │
│   ╭─ Quick Actions ────────────────────────────────────────────────────╮│
│   │                                                                    ││
│   │   [r] Random Problem    [p] Browse Problems    [s] Your Stats     ││
│   │                                                                    ││
│   ╰────────────────────────────────────────────────────────────────────╯│
│                                                                         │
│   ╭─ Recent Activity ──────────────────────────────────────────────────╮│
│   │                                                                    ││
│   │   ✓  1850A  Two Sum               800   math         2 hours ago  ││
│   │   ✓  1851B  Colored Segments     1000   greedy       5 hours ago  ││
│   │   ✗  1852C  Binary Search        1200   binary       yesterday    ││
│   │                                                                    ││
│   ╰────────────────────────────────────────────────────────────────────╯│
│                                                                         │
├─────────────────────────────────────────────────────────────────────────┤
│  [d]ashboard  [p]roblems  [s]tats  [?]help  [q]uit                     │
└─────────────────────────────────────────────────────────────────────────┘


┌─────────────────────────────────────────────────────────────────────────┐
│  DSA Prep › Problems                             tourist │ Expert │ 1847 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   Search: _______________     Difficulty: [800 ────────○── 1600]       │
│   Tags: dp ✕  greedy ✕  + add tag                                       │
│                                                                         │
│  ┌───────┬──────────────────────────────────┬────────┬──────────┬─────┐│
│  │   ID  │  Name                            │ Rating │ Tags     │  ✓  ││
│  ├───────┼──────────────────────────────────┼────────┼──────────┼─────┤│
│  │►1850A │  Two Sum                         │   800  │ math     │     ││
│  │ 1851B │  Colored Segments                │  1000  │ greedy   │  ✓  ││
│  │ 1852C │  Binary Search                   │  1200  │ binary   │  ✓  ││
│  │ 1853A │  Maximum Subarray                │  1100  │ dp       │     ││
│  │ 1854D │  Graph Traversal                 │  1400  │ graphs   │     ││
│  │ 1855B │  Tree Diameter                   │  1500  │ trees    │     ││
│  │ 1856A │  String Matching                 │   900  │ strings  │  ✓  ││
│  │ 1857C │  Number Theory                   │  1300  │ math     │     ││
│  └───────┴──────────────────────────────────┴────────┴──────────┴─────┘│
│                                                                         │
│   Showing 1-8 of 4521 problems                         Page 1 of 566   │
│                                                                         │
├─────────────────────────────────────────────────────────────────────────┤
│  [/]search [f]ilter [r]andom [Enter]start [←→]page [q]back             │
└─────────────────────────────────────────────────────────────────────────┘


┌─────────────────────────────────────────────────────────────────────────┐
│  DSA Prep › Practice                             tourist │ Expert │ 1847 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   ╭─ 1850A: Two Sum ───────────────────────────────────────────────────╮│
│   │                                                                    ││
│   │   Rating: 800  │  Tags: math, implementation  │  Solved by: 45.2k ││
│   │                                                                    ││
│   │   ─────────────────────────────────────────────────────────────── ││
│   │                                                                    ││
│   │   Given an array of integers nums and an integer target, return   ││
│   │   indices of the two numbers such that they add up to target.     ││
│   │                                                                    ││
│   │   You may assume that each input would have exactly one solution, ││
│   │   and you may not use the same element twice.                     ││
│   │                                                                    ││
│   │   Example 1:                                                       ││
│   │   Input: nums = [2,7,11,15], target = 9                           ││
│   │   Output: [0,1]                                                    ││
│   │   Explanation: nums[0] + nums[1] == 9, return [0, 1].             ││
│   │                                                                    ││
│   ╰────────────────────────────────────────────────────────────────────╯│
│                                                                         │
│   ╭─ Timer ─────╮  ╭─ Actions ─────────────────────────────────────────╮│
│   │             │  │                                                   ││
│   │   ⏱ 05:32   │  │  [o] Open in browser   [n] Next problem          ││
│   │             │  │  [✓] Mark solved       [s] Skip                  ││
│   ╰─────────────╯  ╰───────────────────────────────────────────────────╯│
│                                                                         │
├─────────────────────────────────────────────────────────────────────────┤
│  [o]pen [✓]solved [s]kip [n]ext [q]back                                │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Phase 2: AI-Powered Features

### Bedrock Integration

```go
// internal/ai/bedrock.go
package ai

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

const (
    ModelClaude35Sonnet = "anthropic.claude-3-5-sonnet-20241022-v2:0"
    ModelClaude3Haiku   = "anthropic.claude-3-haiku-20240307-v1:0"
)

type BedrockClient struct {
    client *bedrockruntime.Client
    model  string
}

func NewBedrockClient(ctx context.Context, model string) (*BedrockClient, error) {
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        return nil, fmt.Errorf("load AWS config: %w", err)
    }

    return &BedrockClient{
        client: bedrockruntime.NewFromConfig(cfg),
        model:  model,
    }, nil
}

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type ClaudeRequest struct {
    AnthropicVersion string    `json:"anthropic_version"`
    MaxTokens        int       `json:"max_tokens"`
    System           string    `json:"system,omitempty"`
    Messages         []Message `json:"messages"`
}

type ClaudeResponse struct {
    Content []struct {
        Text string `json:"text"`
    } `json:"content"`
}

func (c *BedrockClient) Complete(ctx context.Context, system, prompt string) (string, error) {
    req := ClaudeRequest{
        AnthropicVersion: "bedrock-2023-05-31",
        MaxTokens:        1024,
        System:           system,
        Messages: []Message{
            {Role: "user", Content: prompt},
        },
    }

    body, _ := json.Marshal(req)

    resp, err := c.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
        ModelId:     &c.model,
        ContentType: aws.String("application/json"),
        Body:        body,
    })
    if err != nil {
        return "", fmt.Errorf("invoke model: %w", err)
    }

    var result ClaudeResponse
    if err := json.Unmarshal(resp.Body, &result); err != nil {
        return "", err
    }

    if len(result.Content) == 0 {
        return "", fmt.Errorf("empty response")
    }

    return result.Content[0].Text, nil
}
```

### AI Features Implementation

```go
// internal/ai/assistant.go
package ai

import (
    "context"
    "fmt"
    "strings"

    "github.com/yourusername/dsaprep/internal/codeforces"
)

type Assistant struct {
    client *BedrockClient
}

func NewAssistant(client *BedrockClient) *Assistant {
    return &Assistant{client: client}
}

// GetHint returns a progressive hint (level 1-3)
func (a *Assistant) GetHint(ctx context.Context, problem codeforces.Problem, level int) (string, error) {
    system := `You are a helpful DSA tutor. Provide hints for competitive programming problems.
Be encouraging and guide the student to discover the solution themselves.
Format your response in markdown.`

    levelDesc := map[int]string{
        1: "Give a gentle nudge - just hint at the general approach without specifics",
        2: "Provide a more direct hint - mention the algorithm/data structure to use",
        3: "Give a detailed approach - explain the algorithm steps without full code",
    }

    prompt := fmt.Sprintf(`Problem: %s
Rating: %d
Tags: %s

%s

Provide a Level %d hint:
%s`,
        problem.Name,
        problem.Rating,
        strings.Join(problem.Tags, ", "),
        getProblemDescription(problem), // Fetch from CF or cache
        level,
        levelDesc[level],
    )

    return a.client.Complete(ctx, system, prompt)
}

// ExplainSolution explains the optimal approach after solving
func (a *Assistant) ExplainSolution(ctx context.Context, problem codeforces.Problem) (string, error) {
    system := `You are a DSA expert. Explain solutions to competitive programming problems.
Include:
1. The key insight
2. The algorithm used
3. Time and space complexity
4. Why this approach is optimal
Format in markdown.`

    prompt := fmt.Sprintf(`Explain the optimal solution for:

Problem: %s
Rating: %d
Tags: %s

%s`,
        problem.Name,
        problem.Rating,
        strings.Join(problem.Tags, ", "),
        getProblemDescription(problem),
    )

    return a.client.Complete(ctx, system, prompt)
}

// ReviewCode reviews user's code and provides feedback
func (a *Assistant) ReviewCode(ctx context.Context, problem codeforces.Problem, code string) (string, error) {
    system := `You are a code reviewer for competitive programming.
Review the code for:
1. Correctness - will it produce the right answer?
2. Edge cases - any cases it might miss?
3. Time complexity - is it efficient enough?
4. Style - any improvements for readability?

Be constructive and educational. Format in markdown.`

    prompt := fmt.Sprintf(`Review this solution:

Problem: %s (Rating: %d)

Code:
%s`,
        problem.Name,
        problem.Rating,
        code,
    )

    return a.client.Complete(ctx, system, prompt)
}

// IdentifyPatterns identifies algorithm patterns for a problem
func (a *Assistant) IdentifyPatterns(ctx context.Context, problem codeforces.Problem) (string, error) {
    system := `You are a competitive programming coach.
Identify the algorithm patterns and techniques relevant to this problem.
List 2-3 key patterns with brief explanations of why they apply.
Format as a numbered list in markdown.`

    prompt := fmt.Sprintf(`What patterns apply to:

Problem: %s
Rating: %d
Tags: %s

%s`,
        problem.Name,
        problem.Rating,
        strings.Join(problem.Tags, ", "),
        getProblemDescription(problem),
    )

    return a.client.Complete(ctx, system, prompt)
}

// SuggestSimilar suggests similar problems to practice
func (a *Assistant) SuggestSimilar(ctx context.Context, problem codeforces.Problem, allProblems []codeforces.Problem) ([]codeforces.Problem, error) {
    // This could be enhanced with embeddings later
    // For now, filter by similar tags and rating range

    var similar []codeforces.Problem
    for _, p := range allProblems {
        if p.ID() == problem.ID() {
            continue
        }
        if abs(p.Rating-problem.Rating) <= 200 && hasCommonTag(p.Tags, problem.Tags) {
            similar = append(similar, p)
            if len(similar) >= 5 {
                break
            }
        }
    }
    return similar, nil
}
```

### TUI AI Panel

```go
// internal/tui/views/ai_panel.go
package views

import (
    "github.com/charmbracelet/bubbles/viewport"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/glamour"
)

type AIPanelModel struct {
    viewport   viewport.Model
    content    string
    loading    bool
    visible    bool
    title      string
    width      int
    height     int
}

func NewAIPanelModel() AIPanelModel {
    return AIPanelModel{
        viewport: viewport.New(0, 0),
    }
}

func (m AIPanelModel) Update(msg tea.Msg) (AIPanelModel, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "esc" {
            m.visible = false
            return m, nil
        }
    case AIResponseMsg:
        m.loading = false
        m.content = msg.Content
        rendered, _ := glamour.Render(msg.Content, "dark")
        m.viewport.SetContent(rendered)
    }

    var cmd tea.Cmd
    m.viewport, cmd = m.viewport.Update(msg)
    return m, cmd
}

func (m AIPanelModel) View() string {
    if !m.visible {
        return ""
    }

    style := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("63")).
        Padding(1).
        Width(m.width - 4)

    titleStyle := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("63"))

    var content string
    if m.loading {
        content = "🤔 Thinking..."
    } else {
        content = m.viewport.View()
    }

    return style.Render(
        lipgloss.JoinVertical(lipgloss.Left,
            titleStyle.Render(m.title),
            "",
            content,
            "",
            lipgloss.NewStyle().Faint(true).Render("[↑↓] scroll  [Esc] close"),
        ),
    )
}

func (m *AIPanelModel) Show(title string) {
    m.visible = true
    m.loading = true
    m.title = title
    m.content = ""
}

func (m *AIPanelModel) SetSize(w, h int) {
    m.width = w
    m.height = h / 2
    m.viewport.Width = w - 6
    m.viewport.Height = m.height - 6
}
```

### Practice View with AI

```
┌─────────────────────────────────────────────────────────────────────────┐
│  DSA Prep › Practice                             tourist │ Expert │ 1847 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   ╭─ 1850A: Two Sum ───────────────────────────────────────────────────╮│
│   │                                                                    ││
│   │   Rating: 800  │  Tags: math, implementation  │  Solved by: 45.2k ││
│   │   ─────────────────────────────────────────────────────────────── ││
│   │   Given an array of integers nums and an integer target...        ││
│   │                                                                    ││
│   ╰────────────────────────────────────────────────────────────────────╯│
│                                                                         │
│   ╭─ 💡 Hint (Level 1/3) ──────────────────────────────────────────────╮│
│   │                                                                    ││
│   │   Think about what information you need to track as you iterate   ││
│   │   through the array. For each number, what are you looking for?   ││
│   │                                                                    ││
│   │   Consider: if you're at number `x` and the target is `T`, what   ││
│   │   number would you need to have seen before?                       ││
│   │                                                                    ││
│   │   [↑↓] scroll  [Enter] next hint  [Esc] close                     ││
│   ╰────────────────────────────────────────────────────────────────────╯│
│                                                                         │
│   ╭─ Timer ─────╮  ╭─ AI Actions ──────────────────────────────────────╮│
│   │   ⏱ 05:32   │  │  [h] Hint  [e] Explain  [p] Patterns  [r] Review ││
│   ╰─────────────╯  ╰───────────────────────────────────────────────────╯│
│                                                                         │
├─────────────────────────────────────────────────────────────────────────┤
│  [h]int [e]xplain [p]atterns [r]eview [o]pen [✓]solved [q]back         │
└─────────────────────────────────────────────────────────────────────────┘
```

### AWS Setup for Bedrock

```bash
# 1. Configure AWS credentials
aws configure
# Enter your AWS Access Key ID, Secret, and region (e.g., us-east-1)

# 2. Request model access in AWS Console
# Go to: Amazon Bedrock → Model access → Request access for Claude models

# 3. Test access
aws bedrock-runtime invoke-model \
    --model-id anthropic.claude-3-haiku-20240307-v1:0 \
    --body '{"anthropic_version":"bedrock-2023-05-31","max_tokens":100,"messages":[{"role":"user","content":"Hello"}]}' \
    --content-type application/json \
    output.json
```

### Cost Estimation

| Model | Input (per 1M tokens) | Output (per 1M tokens) | Typical Hint Cost |
|-------|----------------------|------------------------|-------------------|
| Claude 3 Haiku | $0.25 | $1.25 | ~$0.0005 |
| Claude 3.5 Sonnet | $3.00 | $15.00 | ~$0.005 |

**Recommendation**: Use Haiku for hints (fast, cheap), Sonnet for explanations (thorough)

---

## Phase 3: Offline Mode

### SQLite Schema

```sql
-- migrations/001_init.sql

-- Problems cache
CREATE TABLE problems (
    id TEXT PRIMARY KEY,              -- e.g., "1850A"
    contest_id INTEGER NOT NULL,
    problem_index TEXT NOT NULL,
    name TEXT NOT NULL,
    rating INTEGER,
    tags TEXT,                        -- JSON array
    solved_count INTEGER DEFAULT 0,
    url TEXT,
    cached_at INTEGER NOT NULL,       -- Unix timestamp
    UNIQUE(contest_id, problem_index)
);

CREATE INDEX idx_problems_rating ON problems(rating);
CREATE INDEX idx_problems_cached ON problems(cached_at);

-- User progress (local tracking)
CREATE TABLE progress (
    problem_id TEXT PRIMARY KEY REFERENCES problems(id),
    status TEXT DEFAULT 'unseen',     -- unseen, attempted, solved, skipped
    attempts INTEGER DEFAULT 0,
    time_spent INTEGER DEFAULT 0,     -- seconds
    first_seen_at INTEGER,
    solved_at INTEGER,
    updated_at INTEGER NOT NULL
);

CREATE INDEX idx_progress_status ON progress(status);

-- Practice sessions
CREATE TABLE sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    started_at INTEGER NOT NULL,
    ended_at INTEGER,
    problems_viewed INTEGER DEFAULT 0,
    problems_solved INTEGER DEFAULT 0
);

-- Daily statistics (for streaks)
CREATE TABLE daily_stats (
    date TEXT PRIMARY KEY,            -- YYYY-MM-DD
    problems_solved INTEGER DEFAULT 0,
    problems_attempted INTEGER DEFAULT 0,
    time_spent INTEGER DEFAULT 0,
    streak_day INTEGER DEFAULT 0
);

-- Sync metadata
CREATE TABLE sync_meta (
    key TEXT PRIMARY KEY,
    value TEXT,
    updated_at INTEGER
);

-- Store last sync time, CF handle, etc.
INSERT INTO sync_meta (key, value, updated_at) VALUES
    ('problems_synced_at', NULL, NULL),
    ('cf_handle', NULL, NULL);
```

### Data Layer

```go
// internal/storage/sqlite.go
package storage

import (
    "context"
    "database/sql"
    "embed"
    "encoding/json"
    "path/filepath"
    "time"

    "github.com/jmoiron/sqlx"
    _ "github.com/mattn/go-sqlite3"

    "github.com/yourusername/dsaprep/internal/codeforces"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Store struct {
    db *sqlx.DB
}

func New(dataDir string) (*Store, error) {
    dbPath := filepath.Join(dataDir, "data.db")

    db, err := sqlx.Connect("sqlite3", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
    if err != nil {
        return nil, err
    }

    store := &Store{db: db}
    if err := store.migrate(); err != nil {
        return nil, err
    }

    return store, nil
}

func (s *Store) migrate() error {
    schema, err := migrations.ReadFile("migrations/001_init.sql")
    if err != nil {
        return err
    }
    _, err = s.db.Exec(string(schema))
    return err
}

// Problems

func (s *Store) CacheProblems(ctx context.Context, problems []codeforces.Problem, stats []codeforces.ProblemStatistic) error {
    tx, err := s.db.BeginTxx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    now := time.Now().Unix()

    stmt, err := tx.PrepareContext(ctx, `
        INSERT OR REPLACE INTO problems
        (id, contest_id, problem_index, name, rating, tags, solved_count, url, cached_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `)
    if err != nil {
        return err
    }
    defer stmt.Close()

    // Create stats lookup
    statsMap := make(map[string]int)
    for _, stat := range stats {
        key := fmt.Sprintf("%d%s", stat.ContestID, stat.Index)
        statsMap[key] = stat.SolvedCount
    }

    for _, p := range problems {
        tagsJSON, _ := json.Marshal(p.Tags)
        solvedCount := statsMap[p.ID()]

        _, err := stmt.ExecContext(ctx,
            p.ID(), p.ContestID, p.Index, p.Name, p.Rating,
            string(tagsJSON), solvedCount, p.URL(), now,
        )
        if err != nil {
            return err
        }
    }

    // Update sync metadata
    _, err = tx.ExecContext(ctx,
        `UPDATE sync_meta SET value = ?, updated_at = ? WHERE key = 'problems_synced_at'`,
        now, now)
    if err != nil {
        return err
    }

    return tx.Commit()
}

func (s *Store) GetProblems(ctx context.Context, filter ProblemFilter) ([]Problem, error) {
    query := `
        SELECT p.*, pr.status, pr.attempts, pr.time_spent
        FROM problems p
        LEFT JOIN progress pr ON p.id = pr.problem_id
        WHERE 1=1
    `
    args := []interface{}{}

    if filter.MinRating > 0 {
        query += " AND p.rating >= ?"
        args = append(args, filter.MinRating)
    }
    if filter.MaxRating > 0 {
        query += " AND p.rating <= ?"
        args = append(args, filter.MaxRating)
    }
    if filter.Search != "" {
        query += " AND p.name LIKE ?"
        args = append(args, "%"+filter.Search+"%")
    }

    query += " ORDER BY p.rating ASC LIMIT ? OFFSET ?"
    args = append(args, filter.Limit, filter.Offset)

    var problems []Problem
    err := s.db.SelectContext(ctx, &problems, query, args...)
    return problems, err
}

func (s *Store) IsCacheStale(ctx context.Context, maxAge time.Duration) bool {
    var syncedAt sql.NullInt64
    err := s.db.GetContext(ctx, &syncedAt,
        `SELECT value FROM sync_meta WHERE key = 'problems_synced_at'`)

    if err != nil || !syncedAt.Valid {
        return true
    }

    return time.Since(time.Unix(syncedAt.Int64, 0)) > maxAge
}

// Progress

func (s *Store) UpdateProgress(ctx context.Context, problemID string, status string, timeSpent int) error {
    now := time.Now().Unix()

    _, err := s.db.ExecContext(ctx, `
        INSERT INTO progress (problem_id, status, attempts, time_spent, first_seen_at, updated_at)
        VALUES (?, ?, 1, ?, ?, ?)
        ON CONFLICT(problem_id) DO UPDATE SET
            status = excluded.status,
            attempts = attempts + 1,
            time_spent = time_spent + excluded.time_spent,
            solved_at = CASE WHEN excluded.status = 'solved' THEN ? ELSE solved_at END,
            updated_at = ?
    `, problemID, status, timeSpent, now, now, now, now)

    return err
}

func (s *Store) GetStats(ctx context.Context) (*Stats, error) {
    var stats Stats

    err := s.db.GetContext(ctx, &stats, `
        SELECT
            COUNT(*) FILTER (WHERE status = 'solved') as solved,
            COUNT(*) FILTER (WHERE status = 'attempted') as attempted,
            COUNT(*) FILTER (WHERE status = 'skipped') as skipped,
            COALESCE(SUM(time_spent), 0) as total_time
        FROM progress
    `)

    return &stats, err
}

func (s *Store) Close() error {
    return s.db.Close()
}
```

### Sync Engine

```go
// internal/sync/engine.go
package sync

import (
    "context"
    "net/http"
    "time"

    "github.com/yourusername/dsaprep/internal/codeforces"
    "github.com/yourusername/dsaprep/internal/storage"
)

type Engine struct {
    cfClient *codeforces.Client
    store    *storage.Store
    online   bool
}

func NewEngine(client *codeforces.Client, store *storage.Store) *Engine {
    return &Engine{
        cfClient: client,
        store:    store,
    }
}

// CheckConnectivity tests if we can reach Codeforces
func (e *Engine) CheckConnectivity(ctx context.Context) bool {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    req, _ := http.NewRequestWithContext(ctx, "HEAD", "https://codeforces.com", nil)
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        e.online = false
        return false
    }
    resp.Body.Close()

    e.online = resp.StatusCode == 200
    return e.online
}

// SyncProblems fetches problems from CF and caches locally
func (e *Engine) SyncProblems(ctx context.Context) error {
    if !e.CheckConnectivity(ctx) {
        return ErrOffline
    }

    result, err := e.cfClient.GetProblems(ctx)
    if err != nil {
        return err
    }

    return e.store.CacheProblems(ctx, result.Problems, result.ProblemStatistics)
}

// SyncIfStale syncs only if cache is older than maxAge
func (e *Engine) SyncIfStale(ctx context.Context, maxAge time.Duration) error {
    if !e.store.IsCacheStale(ctx, maxAge) {
        return nil
    }
    return e.SyncProblems(ctx)
}

// GetProblems returns problems, using cache when offline
func (e *Engine) GetProblems(ctx context.Context, filter storage.ProblemFilter) ([]storage.Problem, error) {
    // Try to sync in background if online and stale
    if e.CheckConnectivity(ctx) {
        go e.SyncIfStale(context.Background(), 24*time.Hour)
    }

    // Always return from cache
    return e.store.GetProblems(ctx, filter)
}

func (e *Engine) IsOnline() bool {
    return e.online
}

var ErrOffline = fmt.Errorf("offline mode: no internet connection")
```

### Offline Status in TUI

```go
// internal/tui/components/status.go
package components

import (
    "fmt"
    "time"

    "github.com/charmbracelet/lipgloss"
)

type StatusBar struct {
    online      bool
    lastSync    time.Time
    problemCount int
}

func (s StatusBar) View() string {
    var status string
    var statusColor lipgloss.Color

    if s.online {
        status = "● Online"
        statusColor = lipgloss.Color("42") // Green
    } else {
        status = "○ Offline"
        statusColor = lipgloss.Color("243") // Gray
    }

    statusStyle := lipgloss.NewStyle().
        Foreground(statusColor)

    syncInfo := ""
    if !s.lastSync.IsZero() {
        ago := time.Since(s.lastSync)
        if ago < time.Hour {
            syncInfo = fmt.Sprintf("Synced %dm ago", int(ago.Minutes()))
        } else if ago < 24*time.Hour {
            syncInfo = fmt.Sprintf("Synced %dh ago", int(ago.Hours()))
        } else {
            syncInfo = fmt.Sprintf("Synced %dd ago", int(ago.Hours()/24))
        }
    }

    cacheInfo := fmt.Sprintf("%d problems cached", s.problemCount)

    return lipgloss.JoinHorizontal(lipgloss.Center,
        statusStyle.Render(status),
        " │ ",
        syncInfo,
        " │ ",
        cacheInfo,
    )
}
```

---

## CLI Structure

### Commands

```
dsaprep
├── (default)              # Launch TUI
├── problem <id>           # Show problem details
├── random                 # Random problem
│   ├── --min-rating      # Minimum difficulty
│   ├── --max-rating      # Maximum difficulty
│   └── --tags            # Filter by tags
├── stats                  # Your statistics
├── sync                   # Force sync problems
├── config                 # Configuration
│   ├── set <key> <val>   # Set config value
│   ├── get <key>         # Get config value
│   └── init              # Initialize config
├── completion             # Shell completions
└── version                # Version info
```

### Root Command

```go
// internal/cmd/root.go
package cmd

import (
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"

    "github.com/yourusername/dsaprep/internal/tui"
)

var cfgFile string

var rootCmd = &cobra.Command{
    Use:   "dsaprep",
    Short: "Practice DSA problems in your terminal",
    Long: `DSA Prep - A beautiful terminal interface for practicing
Data Structures and Algorithms problems from Codeforces.

Launch without arguments to start the interactive TUI.`,
    Run: func(cmd *cobra.Command, args []string) {
        handle := viper.GetString("cf_handle")
        if err := tui.Run(handle); err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(1)
        }
    },
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    cobra.OnInitialize(initConfig)

    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
        "config file (default: ~/.dsaprep/config.yaml)")

    rootCmd.AddCommand(problemCmd)
    rootCmd.AddCommand(randomCmd)
    rootCmd.AddCommand(statsCmd)
    rootCmd.AddCommand(syncCmd)
    rootCmd.AddCommand(configCmd)
    rootCmd.AddCommand(versionCmd)
}

func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, _ := os.UserHomeDir()
        configDir := filepath.Join(home, ".dsaprep")

        // Create config dir if not exists
        os.MkdirAll(configDir, 0755)

        viper.AddConfigPath(configDir)
        viper.SetConfigName("config")
        viper.SetConfigType("yaml")
    }

    viper.AutomaticEnv()
    viper.ReadInConfig()
}
```

---

## Installation

### Quick Install

```bash
# Go install
go install github.com/yourusername/dsaprep@latest

# Homebrew (macOS/Linux)
brew install yourusername/tap/dsaprep

# Binary download
curl -sSL https://dsaprep.dev/install.sh | bash
```

### First Run

```bash
# Launch TUI
dsaprep

# On first run, you'll be prompted to:
# 1. Enter your Codeforces handle (optional)
# 2. Set difficulty preferences
# 3. Sync problems (requires internet)
```

---

## Configuration

### Config File

```yaml
# ~/.dsaprep/config.yaml

# Your Codeforces handle (for stats)
cf_handle: "tourist"

# Difficulty preferences
difficulty:
  min: 800
  max: 1600

# Daily goal
daily_goal: 5

# AI features (Phase 2)
ai:
  enabled: false
  model: "haiku"  # haiku (cheap) or sonnet (better)

# Offline mode (Phase 3)
offline:
  auto_sync: true
  sync_interval: "24h"
```

---

## Next Steps

### Phase 1: Codeforces Wrapper ⏱️ ~2 weeks
1. Initialize Go project with Cobra + Bubble Tea
2. Implement Codeforces API client with caching
3. Build TUI views: Dashboard, Problems, Practice, Stats
4. Add keyboard navigation and styling
5. Cross-platform build with GoReleaser

### Phase 2: AI Features ⏱️ ~1 week
6. Set up AWS Bedrock access
7. Implement AI assistant (hints, explanations)
8. Add AI panel to Practice view
9. Make AI opt-in with cost awareness

### Phase 3: Offline Mode ⏱️ ~1 week
10. Add SQLite storage layer
11. Implement sync engine
12. Add offline status indicators
13. Test offline → online transitions

---

## References

### Core Libraries
- [Cobra CLI](https://cobra.dev/) - CLI framework
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Glamour](https://github.com/charmbracelet/glamour) - Markdown rendering

### APIs & Services
- [Codeforces API](https://codeforces.com/apiHelp) - Problem source
- [Amazon Bedrock](https://aws.amazon.com/bedrock/) - AI models
- [Claude on Bedrock](https://docs.aws.amazon.com/bedrock/latest/userguide/model-parameters-claude.html) - Claude integration

### Build & Distribution
- [GoReleaser](https://goreleaser.com/) - Release automation
- [Homebrew](https://brew.sh/) - Package manager

---

## Summary

| Phase | Focus | Key Tech | Effort |
|-------|-------|----------|--------|
| **1** | Beautiful CF wrapper | Cobra + Bubble Tea + CF API | 2 weeks |
| **2** | AI assistance | Amazon Bedrock (Claude) | 1 week |
| **3** | Offline support | SQLite + Sync engine | 1 week |

Start simple. Ship Phase 1. Iterate.
