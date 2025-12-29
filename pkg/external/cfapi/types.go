// Package cfapi provides a client for the Codeforces API
package cfapi

import (
	"fmt"
	"time"
)

// API Response wrapper
type Response[T any] struct {
	Status  string `json:"status"`
	Result  T      `json:"result,omitempty"`
	Comment string `json:"comment,omitempty"`
}

// Problem represents a Codeforces problem
type Problem struct {
	ContestID      int      `json:"contestId"`
	ProblemsetName string   `json:"problemsetName,omitempty"`
	Index          string   `json:"index"`
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	Points         float64  `json:"points,omitempty"`
	Rating         int      `json:"rating,omitempty"`
	Tags           []string `json:"tags"`
}

// ProblemStatistics represents problem solve statistics
type ProblemStatistics struct {
	ContestID   int    `json:"contestId"`
	Index       string `json:"index"`
	SolvedCount int    `json:"solvedCount"`
}

// ProblemsResponse contains problems and their statistics
type ProblemsResponse struct {
	Problems          []Problem           `json:"problems"`
	ProblemStatistics []ProblemStatistics `json:"problemStatistics"`
}

// User represents a Codeforces user
type User struct {
	Handle                  string `json:"handle"`
	Email                   string `json:"email,omitempty"`
	FirstName               string `json:"firstName,omitempty"`
	LastName                string `json:"lastName,omitempty"`
	Country                 string `json:"country,omitempty"`
	City                    string `json:"city,omitempty"`
	Organization            string `json:"organization,omitempty"`
	Contribution            int    `json:"contribution"`
	Rank                    string `json:"rank"`
	Rating                  int    `json:"rating"`
	MaxRank                 string `json:"maxRank"`
	MaxRating               int    `json:"maxRating"`
	LastOnlineTimeSeconds   int64  `json:"lastOnlineTimeSeconds"`
	RegistrationTimeSeconds int64  `json:"registrationTimeSeconds"`
	FriendOfCount           int    `json:"friendOfCount"`
	Avatar                  string `json:"avatar"`
	TitlePhoto              string `json:"titlePhoto"`
}

// Submission represents a problem submission
type Submission struct {
	ID                  int64   `json:"id"`
	ContestID           int     `json:"contestId,omitempty"`
	CreationTimeSeconds int64   `json:"creationTimeSeconds"`
	RelativeTimeSeconds int64   `json:"relativeTimeSeconds"`
	Problem             Problem `json:"problem"`
	Author              Party   `json:"author"`
	ProgrammingLanguage string  `json:"programmingLanguage"`
	Verdict             string  `json:"verdict,omitempty"`
	Testset             string  `json:"testset"`
	PassedTestCount     int     `json:"passedTestCount"`
	TimeConsumedMillis  int64   `json:"timeConsumedMillis"`
	MemoryConsumedBytes int64   `json:"memoryConsumedBytes"`
	Points              float64 `json:"points,omitempty"`
}

// Party represents submission author info
type Party struct {
	ContestID        int      `json:"contestId,omitempty"`
	Members          []Member `json:"members"`
	ParticipantType  string   `json:"participantType"`
	TeamID           int      `json:"teamId,omitempty"`
	TeamName         string   `json:"teamName,omitempty"`
	Ghost            bool     `json:"ghost"`
	Room             int      `json:"room,omitempty"`
	StartTimeSeconds int64    `json:"startTimeSeconds,omitempty"`
}

// Member represents a party member
type Member struct {
	Handle string `json:"handle"`
	Name   string `json:"name,omitempty"`
}

// Contest represents a Codeforces contest
type Contest struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Type                string `json:"type"`
	Phase               string `json:"phase"`
	Frozen              bool   `json:"frozen"`
	DurationSeconds     int64  `json:"durationSeconds"`
	StartTimeSeconds    int64  `json:"startTimeSeconds,omitempty"`
	RelativeTimeSeconds int64  `json:"relativeTimeSeconds,omitempty"`
	PreparedBy          string `json:"preparedBy,omitempty"`
	WebsiteURL          string `json:"websiteUrl,omitempty"`
	Description         string `json:"description,omitempty"`
	Difficulty          int    `json:"difficulty,omitempty"`
	Kind                string `json:"kind,omitempty"`
	ICPCRegion          string `json:"icpcRegion,omitempty"`
	Country             string `json:"country,omitempty"`
	City                string `json:"city,omitempty"`
	Season              string `json:"season,omitempty"`
}

// RatingChange represents a user's rating change after a contest
type RatingChange struct {
	ContestID               int    `json:"contestId"`
	ContestName             string `json:"contestName"`
	Handle                  string `json:"handle"`
	Rank                    int    `json:"rank"`
	RatingUpdateTimeSeconds int64  `json:"ratingUpdateTimeSeconds"`
	OldRating               int    `json:"oldRating"`
	NewRating               int    `json:"newRating"`
}

// ContestStandings represents contest standings
type ContestStandings struct {
	Contest  Contest              `json:"contest"`
	Problems []Problem            `json:"problems"`
	Rows     []RanklistRow        `json:"rows"`
}

// RanklistRow represents a row in contest standings
type RanklistRow struct {
	Party                     Party               `json:"party"`
	Rank                      int                 `json:"rank"`
	Points                    float64             `json:"points"`
	Penalty                   int                 `json:"penalty"`
	SuccessfulHackCount       int                 `json:"successfulHackCount"`
	UnsuccessfulHackCount     int                 `json:"unsuccessfulHackCount"`
	ProblemResults            []ProblemResult     `json:"problemResults"`
}

// ProblemResult represents result for a specific problem in standings
type ProblemResult struct {
	Points                    float64 `json:"points"`
	Penalty                   int     `json:"penalty,omitempty"`
	RejectedAttemptCount      int     `json:"rejectedAttemptCount"`
	Type                      string  `json:"type"`
	BestSubmissionTimeSeconds int64   `json:"bestSubmissionTimeSeconds,omitempty"`
}

// Verdict constants
const (
	VerdictOK                  = "OK"
	VerdictFailed              = "FAILED"
	VerdictPartial             = "PARTIAL"
	VerdictCompilationError    = "COMPILATION_ERROR"
	VerdictRuntimeError        = "RUNTIME_ERROR"
	VerdictWrongAnswer         = "WRONG_ANSWER"
	VerdictPresentationError   = "PRESENTATION_ERROR"
	VerdictTimeLimitExceeded   = "TIME_LIMIT_EXCEEDED"
	VerdictMemoryLimitExceeded = "MEMORY_LIMIT_EXCEEDED"
	VerdictIdlenessLimitExc    = "IDLENESS_LIMIT_EXCEEDED"
	VerdictSecurityViolated    = "SECURITY_VIOLATED"
	VerdictCrashed             = "CRASHED"
	VerdictInputPrepCrashed    = "INPUT_PREPARATION_CRASHED"
	VerdictChallenged          = "CHALLENGED"
	VerdictSkipped             = "SKIPPED"
	VerdictTesting             = "TESTING"
	VerdictRejected            = "REJECTED"
)

// ContestPhase constants
const (
	PhaseBefore        = "BEFORE"
	PhaseCoding        = "CODING"
	PhasePendingTest   = "PENDING_SYSTEM_TEST"
	PhaseSystemTest    = "SYSTEM_TEST"
	PhaseFinished      = "FINISHED"
)

// Rank thresholds
var RankThresholds = map[string]int{
	"newbie":                 0,
	"pupil":                  1200,
	"specialist":             1400,
	"expert":                 1600,
	"candidate master":       1900,
	"master":                 2100,
	"international master":   2300,
	"grandmaster":            2400,
	"international grandmaster": 2600,
	"legendary grandmaster":  3000,
}

// Helper methods

// ProblemID returns a unique identifier for the problem
func (p *Problem) ProblemID() string {
	if p.ContestID > 0 {
		return fmt.Sprintf("%d%s", p.ContestID, p.Index)
	}
	return p.Index
}

// URL returns the problem URL
func (p *Problem) URL() string {
	if p.ContestID > 0 {
		return fmt.Sprintf("https://codeforces.com/problemset/problem/%d/%s", p.ContestID, p.Index)
	}
	return ""
}

// ContestURL returns the contest URL for this problem
func (p *Problem) ContestURL() string {
	if p.ContestID > 0 {
		return fmt.Sprintf("https://codeforces.com/contest/%d/problem/%s", p.ContestID, p.Index)
	}
	return ""
}

// IsAccepted returns true if the submission was accepted
func (s *Submission) IsAccepted() bool {
	return s.Verdict == VerdictOK
}

// SubmissionTime returns the submission time as time.Time
func (s *Submission) SubmissionTime() time.Time {
	return time.Unix(s.CreationTimeSeconds, 0)
}

// StartTime returns the contest start time
func (c *Contest) StartTime() time.Time {
	return time.Unix(c.StartTimeSeconds, 0)
}

// Duration returns the contest duration
func (c *Contest) Duration() time.Duration {
	return time.Duration(c.DurationSeconds) * time.Second
}

// IsRunning returns true if the contest is currently running
func (c *Contest) IsRunning() bool {
	return c.Phase == PhaseCoding
}

// IsFinished returns true if the contest has finished
func (c *Contest) IsFinished() bool {
	return c.Phase == PhaseFinished
}

// LastOnline returns when the user was last online
func (u *User) LastOnline() time.Time {
	return time.Unix(u.LastOnlineTimeSeconds, 0)
}

// RegistrationTime returns when the user registered
func (u *User) RegistrationTime() time.Time {
	return time.Unix(u.RegistrationTimeSeconds, 0)
}

// RatingDelta returns the rating change amount
func (rc *RatingChange) RatingDelta() int {
	return rc.NewRating - rc.OldRating
}
