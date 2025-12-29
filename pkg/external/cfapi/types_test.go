package cfapi

import (
	"testing"
	"time"
)

func TestProblem_ProblemID(t *testing.T) {
	tests := []struct {
		name      string
		problem   Problem
		wantID    string
	}{
		{
			name:    "with contest ID",
			problem: Problem{ContestID: 1325, Index: "A"},
			wantID:  "1325A",
		},
		{
			name:    "problem with number suffix",
			problem: Problem{ContestID: 1500, Index: "B2"},
			wantID:  "1500B2",
		},
		{
			name:    "without contest ID",
			problem: Problem{Index: "A"},
			wantID:  "A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.problem.ProblemID(); got != tt.wantID {
				t.Errorf("ProblemID() = %v, want %v", got, tt.wantID)
			}
		})
	}
}

func TestProblem_URL(t *testing.T) {
	tests := []struct {
		name    string
		problem Problem
		wantURL string
	}{
		{
			name:    "with contest ID",
			problem: Problem{ContestID: 1325, Index: "A"},
			wantURL: "https://codeforces.com/problemset/problem/1325/A",
		},
		{
			name:    "without contest ID",
			problem: Problem{Index: "A"},
			wantURL: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.problem.URL(); got != tt.wantURL {
				t.Errorf("URL() = %v, want %v", got, tt.wantURL)
			}
		})
	}
}

func TestProblem_ContestURL(t *testing.T) {
	tests := []struct {
		name    string
		problem Problem
		wantURL string
	}{
		{
			name:    "with contest ID",
			problem: Problem{ContestID: 1325, Index: "A"},
			wantURL: "https://codeforces.com/contest/1325/problem/A",
		},
		{
			name:    "without contest ID",
			problem: Problem{Index: "A"},
			wantURL: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.problem.ContestURL(); got != tt.wantURL {
				t.Errorf("ContestURL() = %v, want %v", got, tt.wantURL)
			}
		})
	}
}

func TestSubmission_IsAccepted(t *testing.T) {
	tests := []struct {
		name       string
		submission Submission
		want       bool
	}{
		{
			name:       "accepted",
			submission: Submission{Verdict: VerdictOK},
			want:       true,
		},
		{
			name:       "wrong answer",
			submission: Submission{Verdict: VerdictWrongAnswer},
			want:       false,
		},
		{
			name:       "time limit exceeded",
			submission: Submission{Verdict: VerdictTimeLimitExceeded},
			want:       false,
		},
		{
			name:       "testing",
			submission: Submission{Verdict: VerdictTesting},
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.submission.IsAccepted(); got != tt.want {
				t.Errorf("IsAccepted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubmission_SubmissionTime(t *testing.T) {
	timestamp := int64(1609459200) // 2021-01-01 00:00:00 UTC
	submission := Submission{CreationTimeSeconds: timestamp}

	result := submission.SubmissionTime()
	expected := time.Unix(timestamp, 0)

	if !result.Equal(expected) {
		t.Errorf("SubmissionTime() = %v, want %v", result, expected)
	}
}

func TestContest_StartTime(t *testing.T) {
	timestamp := int64(1609459200)
	contest := Contest{StartTimeSeconds: timestamp}

	result := contest.StartTime()
	expected := time.Unix(timestamp, 0)

	if !result.Equal(expected) {
		t.Errorf("StartTime() = %v, want %v", result, expected)
	}
}

func TestContest_Duration(t *testing.T) {
	contest := Contest{DurationSeconds: 7200} // 2 hours

	result := contest.Duration()
	expected := 2 * time.Hour

	if result != expected {
		t.Errorf("Duration() = %v, want %v", result, expected)
	}
}

func TestContest_IsRunning(t *testing.T) {
	tests := []struct {
		name    string
		phase   string
		want    bool
	}{
		{"before", PhaseBefore, false},
		{"coding", PhaseCoding, true},
		{"pending", PhasePendingTest, false},
		{"system test", PhaseSystemTest, false},
		{"finished", PhaseFinished, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contest := Contest{Phase: tt.phase}
			if got := contest.IsRunning(); got != tt.want {
				t.Errorf("IsRunning() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContest_IsFinished(t *testing.T) {
	tests := []struct {
		name  string
		phase string
		want  bool
	}{
		{"before", PhaseBefore, false},
		{"coding", PhaseCoding, false},
		{"finished", PhaseFinished, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contest := Contest{Phase: tt.phase}
			if got := contest.IsFinished(); got != tt.want {
				t.Errorf("IsFinished() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_LastOnline(t *testing.T) {
	timestamp := int64(1609459200)
	user := User{LastOnlineTimeSeconds: timestamp}

	result := user.LastOnline()
	expected := time.Unix(timestamp, 0)

	if !result.Equal(expected) {
		t.Errorf("LastOnline() = %v, want %v", result, expected)
	}
}

func TestUser_RegistrationTime(t *testing.T) {
	timestamp := int64(1609459200)
	user := User{RegistrationTimeSeconds: timestamp}

	result := user.RegistrationTime()
	expected := time.Unix(timestamp, 0)

	if !result.Equal(expected) {
		t.Errorf("RegistrationTime() = %v, want %v", result, expected)
	}
}

func TestRatingChange_RatingDelta(t *testing.T) {
	tests := []struct {
		name      string
		oldRating int
		newRating int
		want      int
	}{
		{"positive change", 1500, 1600, 100},
		{"negative change", 1600, 1500, -100},
		{"no change", 1500, 1500, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := RatingChange{OldRating: tt.oldRating, NewRating: tt.newRating}
			if got := rc.RatingDelta(); got != tt.want {
				t.Errorf("RatingDelta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerdictConstants(t *testing.T) {
	verdicts := []string{
		VerdictOK, VerdictFailed, VerdictPartial,
		VerdictCompilationError, VerdictRuntimeError,
		VerdictWrongAnswer, VerdictPresentationError,
		VerdictTimeLimitExceeded, VerdictMemoryLimitExceeded,
		VerdictIdlenessLimitExc, VerdictSecurityViolated,
		VerdictCrashed, VerdictInputPrepCrashed,
		VerdictChallenged, VerdictSkipped,
		VerdictTesting, VerdictRejected,
	}

	for _, v := range verdicts {
		if v == "" {
			t.Error("Verdict constant should not be empty")
		}
	}
}

func TestContestPhaseConstants(t *testing.T) {
	phases := []string{
		PhaseBefore, PhaseCoding, PhasePendingTest,
		PhaseSystemTest, PhaseFinished,
	}

	for _, p := range phases {
		if p == "" {
			t.Error("Phase constant should not be empty")
		}
	}
}

func TestRankThresholds(t *testing.T) {
	// Test that all ranks have thresholds
	expectedRanks := []string{
		"newbie", "pupil", "specialist", "expert",
		"candidate master", "master", "international master",
		"grandmaster", "international grandmaster", "legendary grandmaster",
	}

	for _, rank := range expectedRanks {
		if _, ok := RankThresholds[rank]; !ok {
			t.Errorf("RankThresholds missing rank %s", rank)
		}
	}

	// Test that thresholds are in ascending order (by rank)
	if RankThresholds["newbie"] != 0 {
		t.Error("Newbie threshold should be 0")
	}
	if RankThresholds["pupil"] <= RankThresholds["newbie"] {
		t.Error("Pupil threshold should be greater than newbie")
	}
	if RankThresholds["legendary grandmaster"] != 3000 {
		t.Error("Legendary grandmaster threshold should be 3000")
	}
}

func TestProblemFields(t *testing.T) {
	problem := Problem{
		ContestID:      1325,
		ProblemsetName: "acmsguru",
		Index:          "A",
		Name:           "EhAb AnD gCd",
		Type:           "PROGRAMMING",
		Points:         500,
		Rating:         800,
		Tags:           []string{"math", "number theory"},
	}

	if problem.ContestID != 1325 {
		t.Errorf("ContestID = %v, want %v", problem.ContestID, 1325)
	}
	if problem.Name != "EhAb AnD gCd" {
		t.Errorf("Name = %v, want %v", problem.Name, "EhAb AnD gCd")
	}
	if len(problem.Tags) != 2 {
		t.Errorf("len(Tags) = %v, want %v", len(problem.Tags), 2)
	}
}

func TestUserFields(t *testing.T) {
	user := User{
		Handle:       "tourist",
		Rank:         "legendary grandmaster",
		Rating:       3800,
		MaxRank:      "legendary grandmaster",
		MaxRating:    3979,
		Contribution: 100,
		Avatar:       "https://example.com/avatar.jpg",
	}

	if user.Handle != "tourist" {
		t.Errorf("Handle = %v, want %v", user.Handle, "tourist")
	}
	if user.Rank != "legendary grandmaster" {
		t.Errorf("Rank = %v, want %v", user.Rank, "legendary grandmaster")
	}
	if user.Rating != 3800 {
		t.Errorf("Rating = %v, want %v", user.Rating, 3800)
	}
}
