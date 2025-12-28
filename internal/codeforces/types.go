package codeforces

import "fmt"

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

// ID returns a unique identifier for the problem
func (p Problem) ID() string {
	return fmt.Sprintf("%d%s", p.ContestID, p.Index)
}

// URL returns the Codeforces URL for this problem
func (p Problem) URL() string {
	return fmt.Sprintf("https://codeforces.com/problemset/problem/%d/%s",
		p.ContestID, p.Index)
}

// DifficultyColor returns a hex color based on problem rating
func (p Problem) DifficultyColor() string {
	switch {
	case p.Rating == 0:
		return "#808080" // Unrated
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

// RankName returns the rank name based on problem rating
func (p Problem) RankName() string {
	switch {
	case p.Rating == 0:
		return "Unrated"
	case p.Rating < 1200:
		return "Newbie"
	case p.Rating < 1400:
		return "Pupil"
	case p.Rating < 1600:
		return "Specialist"
	case p.Rating < 1900:
		return "Expert"
	case p.Rating < 2100:
		return "Candidate Master"
	case p.Rating < 2400:
		return "Master"
	case p.Rating < 2600:
		return "International Master"
	case p.Rating < 3000:
		return "Grandmaster"
	default:
		return "Legendary Grandmaster"
	}
}

// ProblemStatistic contains solve counts for a problem
type ProblemStatistic struct {
	ContestID   int    `json:"contestId"`
	Index       string `json:"index"`
	SolvedCount int    `json:"solvedCount"`
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

// RankColor returns a hex color based on user rank
func (u User) RankColor() string {
	switch u.Rank {
	case "newbie":
		return "#808080"
	case "pupil":
		return "#008000"
	case "specialist":
		return "#03A89E"
	case "expert":
		return "#0000FF"
	case "candidate master":
		return "#AA00AA"
	case "master":
		return "#FF8C00"
	case "international master":
		return "#FF8C00"
	case "grandmaster":
		return "#FF0000"
	case "international grandmaster":
		return "#FF0000"
	case "legendary grandmaster":
		return "#FF0000"
	default:
		return "#808080"
	}
}

// Submission represents a user's submission
type Submission struct {
	ID                  int64   `json:"id"`
	ContestID           int     `json:"contestId"`
	CreationTimeSeconds int64   `json:"creationTimeSeconds"`
	RelativeTimeSeconds int64   `json:"relativeTimeSeconds"`
	Problem             Problem `json:"problem"`
	Author              Party   `json:"author"`
	ProgrammingLanguage string  `json:"programmingLanguage"`
	Verdict             string  `json:"verdict"`
	Testset             string  `json:"testset"`
	PassedTestCount     int     `json:"passedTestCount"`
	TimeConsumedMillis  int     `json:"timeConsumedMillis"`
	MemoryConsumedBytes int64   `json:"memoryConsumedBytes"`
}

// Party represents participants in a contest
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

// Member represents a team member
type Member struct {
	Handle string `json:"handle"`
}

// RatingChange represents a rating change after a contest
type RatingChange struct {
	ContestID               int    `json:"contestId"`
	ContestName             string `json:"contestName"`
	Handle                  string `json:"handle"`
	Rank                    int    `json:"rank"`
	RatingUpdateTimeSeconds int64  `json:"ratingUpdateTimeSeconds"`
	OldRating               int    `json:"oldRating"`
	NewRating               int    `json:"newRating"`
}

// Contest represents a Codeforces contest
type Contest struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Type                string `json:"type"`
	Phase               string `json:"phase"`
	Frozen              bool   `json:"frozen"`
	DurationSeconds     int    `json:"durationSeconds"`
	StartTimeSeconds    int64  `json:"startTimeSeconds,omitempty"`
	RelativeTimeSeconds int64  `json:"relativeTimeSeconds,omitempty"`
	PreparedBy          string `json:"preparedBy,omitempty"`
	WebsiteURL          string `json:"websiteUrl,omitempty"`
	Description         string `json:"description,omitempty"`
	Difficulty          int    `json:"difficulty,omitempty"`
	Kind                string `json:"kind,omitempty"`
}
