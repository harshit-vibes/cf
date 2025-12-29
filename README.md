# DSA Prep

A command-line tool and Go SDK for practicing competitive programming problems from Codeforces.

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green.svg)
[![CI](https://github.com/harshit-vibes/dsaprep/actions/workflows/ci.yml/badge.svg)](https://github.com/harshit-vibes/dsaprep/actions/workflows/ci.yml)

## Features

- **Codeforces API Client** - Full-featured API client with caching and rate limiting
- **Web Scraper** - Parse problem statements, samples, and metadata from Codeforces
- **Workspace Management** - Organize problems locally with versioned schemas
- **Health Checks** - Verify system configuration and connectivity
- **Cross-Platform** - Works on Linux, macOS, and Windows

## Installation

### Using Go (requires Go 1.22+)

```bash
go install github.com/harshit-vibes/dsaprep/cmd/dsaprep@latest
```

### Download Binary

Download pre-built binaries from the [Releases](https://github.com/harshit-vibes/dsaprep/releases) page.

**macOS (Apple Silicon)**
```bash
curl -Lo dsaprep https://github.com/harshit-vibes/dsaprep/releases/latest/download/dsaprep-darwin-arm64
chmod +x dsaprep
sudo mv dsaprep /usr/local/bin/
```

**macOS (Intel)**
```bash
curl -Lo dsaprep https://github.com/harshit-vibes/dsaprep/releases/latest/download/dsaprep-darwin-amd64
chmod +x dsaprep
sudo mv dsaprep /usr/local/bin/
```

**Linux (amd64)**
```bash
curl -Lo dsaprep https://github.com/harshit-vibes/dsaprep/releases/latest/download/dsaprep-linux-amd64
chmod +x dsaprep
sudo mv dsaprep /usr/local/bin/
```

**Linux (arm64)**
```bash
curl -Lo dsaprep https://github.com/harshit-vibes/dsaprep/releases/latest/download/dsaprep-linux-arm64
chmod +x dsaprep
sudo mv dsaprep /usr/local/bin/
```

### Build from Source

```bash
git clone https://github.com/harshit-vibes/dsaprep.git
cd dsaprep
go build -o dsaprep ./cmd/dsaprep
sudo mv dsaprep /usr/local/bin/
```

## Quick Start

```bash
# Initialize a workspace
dsaprep init ~/dsa-practice

# Check system health
dsaprep health

# Parse a problem from Codeforces
dsaprep parse 1 A

# Show version
dsaprep version
```

## CLI Commands

### `dsaprep init [path]`

Initialize a new workspace for storing problems and tracking progress.

```bash
# Initialize in current directory
dsaprep init

# Initialize in a specific path
dsaprep init ~/competitive-programming
```

This creates:
```
workspace/
├── workspace.yaml      # Workspace manifest
├── problems/           # Problem metadata and statements
├── templates/          # Code templates
├── submissions/        # Your solutions
└── stats/              # Progress tracking
```

### `dsaprep health`

Check system configuration and connectivity.

```bash
dsaprep health
```

Output:
```
Health Check Report
══════════════════════════════════════════════════════════════════════

✓ Environment File       ~/.dsaprep.env exists
✓ Configuration          Config loaded successfully
✓ Workspace              Workspace initialized
✓ Schema Version         Schema version 1.0 compatible
✓ CF API                 API responding (latency: 156ms)
✓ CF Web                 Page structure verified
✓ CF Handle              Handle 'tourist' verified (rating: 3979)

══════════════════════════════════════════════════════════════════════
Overall: Healthy | Duration: 1.2s
```

### `dsaprep parse <contest_id> <problem_index>`

Parse a problem from Codeforces and save it to your workspace.

```bash
# Parse problem A from contest 1
dsaprep parse 1 A

# Parse problem B from contest 1800
dsaprep parse 1800 B
```

### `dsaprep version`

Display version information.

```bash
dsaprep version
```

## Configuration

### Config File

Configuration is stored in `~/.dsaprep/config.yaml`:

```yaml
cf_handle: your_codeforces_handle
difficulty:
  min: 800
  max: 1600
daily_goal: 5
workspace_path: ~/dsa-practice
```

### Credentials File

Credentials are stored in `~/.dsaprep.env`:

```bash
# Your Codeforces handle (required)
CF_HANDLE=your_handle

# API credentials (optional - for authenticated requests)
CF_API_KEY=your_api_key
CF_API_SECRET=your_api_secret

# Session cookies (optional - required for submissions)
CF_JSESSIONID=your_session_id
CF_39CE7=your_cookie_value
CF_CLEARANCE=cloudflare_clearance_token
CF_CLEARANCE_EXPIRES=1234567890
CF_CLEARANCE_UA=Mozilla/5.0...
```

To get API credentials:
1. Go to https://codeforces.com/settings/api
2. Create a new API key
3. Copy the key and secret to your `.dsaprep.env` file

To get session cookies (for submissions):
1. Log in to Codeforces in your browser
2. Open DevTools (F12) > Application > Cookies > codeforces.com
3. Copy the cookie values to your `.dsaprep.env` file

---

## SDK Usage

DSA Prep can also be used as a Go library in your own applications.

### Installation

```bash
go get github.com/harshit-vibes/dsaprep
```

### Package Overview

| Package | Description |
|---------|-------------|
| `pkg/external/cfapi` | Codeforces API client |
| `pkg/external/cfweb` | Web scraper for problem parsing |
| `pkg/internal/config` | Configuration management |
| `pkg/internal/workspace` | Workspace and problem storage |
| `pkg/internal/schema/v1` | Data schemas for problems, submissions |
| `pkg/internal/health` | Health check framework |
| `pkg/internal/errors` | Structured error handling |

---

### Codeforces API Client

The `cfapi` package provides a full-featured client for the Codeforces API.

#### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/harshit-vibes/dsaprep/pkg/external/cfapi"
)

func main() {
    // Create a client (no authentication required for most operations)
    client := cfapi.NewClient()
    ctx := context.Background()

    // Ping the API
    if err := client.Ping(ctx); err != nil {
        log.Fatal("API unreachable:", err)
    }

    // Get user info
    users, err := client.GetUserInfo(ctx, []string{"tourist"})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("User: %s, Rating: %d, Rank: %s\n",
        users[0].Handle, users[0].Rating, users[0].Rank)
}
```

#### Client Options

```go
// With API credentials (for authenticated requests)
client := cfapi.NewClient(
    cfapi.WithAPICredentials("your_api_key", "your_api_secret"),
)

// With custom HTTP client
httpClient := &http.Client{Timeout: 60 * time.Second}
client := cfapi.NewClient(
    cfapi.WithHTTPClient(httpClient),
)

// With custom cache TTL
client := cfapi.NewClient(
    cfapi.WithCacheTTL(10 * time.Minute),
)

// Combine multiple options
client := cfapi.NewClient(
    cfapi.WithAPICredentials(key, secret),
    cfapi.WithCacheTTL(15 * time.Minute),
)
```

#### Fetching Problems

```go
// Get all problems
problems, err := client.GetProblems(ctx, nil)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Total problems: %d\n", len(problems.Problems))

// Get problems with specific tags
dpProblems, err := client.GetProblems(ctx, []string{"dp"})

// Filter problems by criteria
filtered, err := client.FilterProblems(ctx,
    800,   // minRating
    1200,  // maxRating
    []string{"greedy", "math"},  // tags (nil for any)
    true,  // excludeSolved
    "your_handle",  // handle for exclusion
)
```

#### User Data

```go
// Get user info (supports multiple handles)
users, err := client.GetUserInfo(ctx, []string{"tourist", "Petr"})
for _, u := range users {
    fmt.Printf("%s: %d (%s)\n", u.Handle, u.Rating, u.Rank)
}

// Get user submissions
submissions, err := client.GetUserSubmissions(ctx, "tourist", 1, 10)
for _, s := range submissions {
    fmt.Printf("%d%s: %s\n", s.Problem.ContestID, s.Problem.Index, s.Verdict)
}

// Get rating history
ratings, err := client.GetUserRating(ctx, "tourist")
for _, r := range ratings {
    fmt.Printf("%s: %d -> %d\n", r.ContestName, r.OldRating, r.NewRating)
}

// Get solved problems
solved, err := client.GetSolvedProblems(ctx, "tourist")
fmt.Printf("Solved: %d problems\n", len(solved))
```

#### Contests

```go
// Get all contests
contests, err := client.GetContests(ctx, false) // false = exclude gym

// Get specific contest
contest, err := client.GetContest(ctx, 1)
fmt.Printf("Contest: %s, Phase: %s\n", contest.Name, contest.Phase)

// Get contest standings
standings, err := client.GetContestStandings(ctx,
    1,     // contestID
    1,     // from (1-indexed)
    10,    // count
    nil,   // handles filter
    false, // showUnofficial
)
```

#### Get Specific Problem

```go
// Get a single problem by contest ID and index
problem, err := client.GetProblem(ctx, 1, "A")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Problem: %s\n", problem.Name)
fmt.Printf("Rating: %d\n", problem.Rating)
fmt.Printf("Tags: %v\n", problem.Tags)
```

---

### Web Scraper

The `cfweb` package scrapes problem statements and metadata from the Codeforces website.

#### Parsing Problems

```go
package main

import (
    "fmt"
    "log"

    "github.com/harshit-vibes/dsaprep/pkg/external/cfweb"
)

func main() {
    // Create a session
    session, err := cfweb.NewSession()
    if err != nil {
        log.Fatal(err)
    }

    // Create a parser
    parser := cfweb.NewParser(session)

    // Parse a problem
    problem, err := parser.ParseProblem(1, "A")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Name: %s\n", problem.Name)
    fmt.Printf("Time Limit: %s\n", problem.TimeLimit)
    fmt.Printf("Memory Limit: %s\n", problem.MemoryLimit)
    fmt.Printf("Rating: %d\n", problem.Rating)
    fmt.Printf("Tags: %v\n", problem.Tags)

    // Print sample test cases
    for i, sample := range problem.Samples {
        fmt.Printf("\nSample %d:\n", i+1)
        fmt.Printf("Input:\n%s\n", sample.Input)
        fmt.Printf("Output:\n%s\n", sample.Output)
    }
}
```

#### Parse from Problemset URL

```go
// Parse from /problemset/problem/1/A instead of /contest/1/problem/A
problem, err := parser.ParseProblemset(1, "A")
```

#### Parse All Problems in a Contest

```go
problems, err := parser.ParseContestProblems(1)
for _, p := range problems {
    fmt.Printf("%s: %s\n", p.Index, p.Name)
}
```

#### Convert to Schema Format

```go
// Parse and convert to internal schema format
parsed, err := parser.ParseProblem(1, "A")
if err != nil {
    log.Fatal(err)
}

// Convert to v1.Problem schema (for storage)
schemaProblem := parsed.ToSchemaProblem()
fmt.Printf("ID: %s\n", schemaProblem.ID)          // "1A"
fmt.Printf("Platform: %s\n", schemaProblem.Platform)  // "codeforces"
```

#### Session with Authentication

```go
session, _ := cfweb.NewSession()

// Set Cloudflare clearance (for bypassing protection)
session.SetCFClearance(
    "clearance_token",
    "Mozilla/5.0...",  // User-Agent
    time.Now().Add(1*time.Hour),  // Expiration
)

// Set full authentication (for submissions)
session.SetFullAuth(
    "cf_clearance_token",
    "user_agent_string",
    time.Now().Add(1*time.Hour),
    "jsessionid_value",
    "39ce7_cookie_value",
    "your_handle",
)

// Check authentication status
if session.IsAuthenticated() {
    fmt.Println("Logged in as:", session.Handle())
}
```

#### Verify Page Structure

```go
// Verify that CSS selectors still work (useful for detecting site changes)
err := parser.VerifyPageStructure()
if err != nil {
    fmt.Println("Warning: Codeforces page structure may have changed")
}
```

---

### Configuration Management

The `config` package handles application settings.

```go
package main

import (
    "fmt"
    "log"

    "github.com/harshit-vibes/dsaprep/pkg/internal/config"
)

func main() {
    // Initialize config with workspace path
    err := config.Init("/path/to/workspace")
    if err != nil {
        log.Fatal(err)
    }

    // Get config
    cfg := config.Get()
    fmt.Printf("Handle: %s\n", cfg.CFHandle)
    fmt.Printf("Difficulty: %d-%d\n", cfg.Difficulty.Min, cfg.Difficulty.Max)
    fmt.Printf("Daily Goal: %d\n", cfg.DailyGoal)

    // Update settings
    config.SetCFHandle("your_handle")
    config.SetDifficulty(1000, 1600)
    config.SetDailyGoal(5)

    // Load credentials from ~/.dsaprep.env
    creds, err := config.LoadCredentials()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("API Key: %s\n", creds.APIKey)
}
```

---

### Workspace Management

The `workspace` package manages local problem storage.

```go
package main

import (
    "fmt"
    "log"

    "github.com/harshit-vibes/dsaprep/pkg/internal/workspace"
    v1 "github.com/harshit-vibes/dsaprep/pkg/internal/schema/v1"
)

func main() {
    // Create workspace manager
    ws := workspace.New("/path/to/workspace")

    // Initialize new workspace
    err := ws.Init("My DSA Practice", "your_handle")
    if err != nil {
        log.Fatal(err)
    }

    // Or load existing workspace
    err = ws.Load()
    if err != nil {
        log.Fatal(err)
    }

    // Check if workspace exists
    if ws.Exists() {
        fmt.Println("Workspace found!")
    }

    // Validate workspace integrity
    if err := ws.Validate(); err != nil {
        fmt.Printf("Validation failed: %v\n", err)
    }

    // Save a problem
    problem := &v1.Problem{
        ID:        "1A",
        Platform:  "codeforces",
        ContestID: 1,
        Index:     "A",
        Name:      "Theatre Square",
        URL:       "https://codeforces.com/problemset/problem/1/A",
        Limits: v1.ProblemLimits{
            TimeLimit:   "1 second",
            MemoryLimit: "256 megabytes",
        },
        Metadata: v1.ProblemMetadata{
            Rating: 1000,
            Tags:   []string{"math"},
        },
        Samples: []v1.Sample{
            {Index: 1, Input: "6 6 4", Output: "4"},
        },
    }
    err = ws.SaveProblem(problem)

    // Load a problem
    loaded, err := ws.LoadProblem("codeforces", 1, "A")
    fmt.Printf("Loaded: %s\n", loaded.Name)

    // List all problems
    problems, err := ws.ListProblems()
    for _, p := range problems {
        fmt.Printf("%s: %s\n", p.ID, p.Name)
    }

    // Check if problem exists
    exists := ws.ProblemExists("codeforces", 1, "A")

    // Get paths
    fmt.Println("Problems:", ws.ProblemsPath())
    fmt.Println("Templates:", ws.TemplatesPath())
    fmt.Println("Submissions:", ws.SubmissionsPath())
}
```

---

### Health Checks

The `health` package provides system diagnostics.

```go
package main

import (
    "context"
    "fmt"

    "github.com/harshit-vibes/dsaprep/pkg/internal/health"
    externalHealth "github.com/harshit-vibes/dsaprep/pkg/external/health"
    "github.com/harshit-vibes/dsaprep/pkg/external/cfapi"
)

func main() {
    ctx := context.Background()

    // Create checker
    checker := health.NewChecker()

    // Add internal checks
    checker.AddCheck(&health.EnvFileCheck{})
    checker.AddCheck(&health.ConfigCheck{})

    // Add external checks
    client := cfapi.NewClient()
    checker.AddCheck(externalHealth.NewCFAPICheck(client))
    checker.AddCheck(externalHealth.NewCFHandleCheck(client))

    // Run all checks
    report := checker.RunAll(ctx)

    // Process results
    fmt.Printf("Overall Status: %s\n", report.OverallStatus)
    fmt.Printf("Can Proceed: %v\n", report.CanProceed)
    fmt.Printf("Duration: %s\n", report.Duration)

    for _, result := range report.Results {
        status := "✓"
        if result.Status != health.StatusHealthy {
            status = "✗"
        }
        fmt.Printf("%s %s: %s\n", status, result.Name, result.Message)
    }
}
```

---

### Data Schemas

The `schema/v1` package defines data structures for problems, submissions, and progress.

```go
import v1 "github.com/harshit-vibes/dsaprep/pkg/internal/schema/v1"

// Problem structure
problem := &v1.Problem{
    ID:        "1A",
    Platform:  "codeforces",
    ContestID: 1,
    Index:     "A",
    Name:      "Theatre Square",
    URL:       "https://codeforces.com/problemset/problem/1/A",
    Limits: v1.ProblemLimits{
        TimeLimit:   "1 second",
        MemoryLimit: "256 megabytes",
    },
    Metadata: v1.ProblemMetadata{
        Rating: 1000,
        Tags:   []string{"math", "implementation"},
    },
    Samples: []v1.Sample{
        {Index: 1, Input: "6 6 4", Output: "4"},
    },
    Practice: v1.PracticeData{
        Status:       v1.StatusUnseen,  // unseen, attempted, solved
        AttemptCount: 0,
        TimeSpent:    0,
    },
    Notes: v1.UserNotes{
        Difficulty: "easy",
        Approach:   "greedy",
        CustomTags: []string{"math"},
    },
}

// Submission structure
submission := v1.NewSubmission(
    123456,           // submissionID
    1,                // contestID
    "A",              // problemIndex
    v1.VerdictAccepted,
    500,              // timeMs
    4096,             // memoryBytes
    "cpp",            // language
    "source_hash",
)

// Progress tracking
progress := v1.NewProgress()
progress.AddSolved(1000)  // Add solved problem with rating 1000
progress.AddAttempted()
fmt.Printf("Solved: %d, Attempted: %d\n", progress.TotalSolved, progress.TotalAttempted)
fmt.Printf("Current Streak: %d days\n", progress.CurrentStreak)
```

---

### Error Handling

The `errors` package provides structured error types.

```go
import "github.com/harshit-vibes/dsaprep/pkg/internal/errors"

// Create structured errors
err := errors.New(
    errors.ErrHandleNotSet,
    errors.CatUser,
    "CF handle not configured",
    errors.ActionUserPrompt,
)

// Check error type
if errors.Is(err, errors.ErrHandleNotSet) {
    fmt.Println("Please set your Codeforces handle")
}

// Get recovery action
if dsaErr, ok := err.(*errors.Error); ok {
    switch dsaErr.Action {
    case errors.ActionAutoFix:
        fmt.Println("This can be fixed automatically")
    case errors.ActionUserPrompt:
        fmt.Println("User input required")
    case errors.ActionRetry:
        fmt.Println("Please try again")
    }
}
```

---

## Architecture

```
dsaprep/
├── cmd/dsaprep/           # CLI entry point
│   └── main.go
├── pkg/
│   ├── cmd/               # CLI commands (Cobra)
│   │   └── root.go
│   ├── external/          # External integrations
│   │   ├── cfapi/         # Codeforces API client
│   │   ├── cfweb/         # Web scraper
│   │   └── health/        # External health checks
│   └── internal/          # Internal packages
│       ├── config/        # Configuration management
│       ├── errors/        # Error types
│       ├── health/        # Health check framework
│       ├── schema/        # Data schemas
│       │   └── v1/        # Schema version 1
│       └── workspace/     # Workspace management
├── .github/workflows/     # CI/CD
├── go.mod
└── README.md
```

## Development

```bash
# Clone the repository
git clone https://github.com/harshit-vibes/dsaprep.git
cd dsaprep

# Install dependencies
go mod download

# Run tests
go test ./...

# Run tests with race detection
go test -race ./...

# Run integration tests (requires network)
go test -tags=integration ./...

# Build
go build -o dsaprep ./cmd/dsaprep

# Run linter
golangci-lint run

# Build for all platforms
GOOS=linux GOARCH=amd64 go build -o dsaprep-linux-amd64 ./cmd/dsaprep
GOOS=darwin GOARCH=arm64 go build -o dsaprep-darwin-arm64 ./cmd/dsaprep
GOOS=windows GOARCH=amd64 go build -o dsaprep-windows-amd64.exe ./cmd/dsaprep
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
