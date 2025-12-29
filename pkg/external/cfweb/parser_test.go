package cfweb

import (
	"net/http"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestNewParser(t *testing.T) {
	session := &Session{}
	parser := NewParser(session)

	if parser == nil {
		t.Fatal("NewParser() returned nil")
	}
	if parser.session != session {
		t.Error("Parser.session not set correctly")
	}
}

func TestNewParserWithClient(t *testing.T) {
	client := &http.Client{}
	parser := NewParserWithClient(client)

	if parser == nil {
		t.Fatal("NewParserWithClient() returned nil")
	}
	if parser.session == nil {
		t.Error("Parser.session should not be nil")
	}
	if parser.session.client != client {
		t.Error("Parser.session.client not set correctly")
	}
}

func TestParsedProblem_ToSchemaProblem(t *testing.T) {
	parsed := &ParsedProblem{
		ContestID:   1325,
		Index:       "A",
		Name:        "EhAb AnD gCd",
		TimeLimit:   "1 second",
		MemoryLimit: "256 megabytes",
		Statement:   "Test statement",
		Rating:      800,
		Tags:        []string{"math", "number theory"},
		URL:         "https://codeforces.com/contest/1325/problem/A",
		Samples: []Sample{
			{Index: 1, Input: "1\n", Output: "2\n"},
			{Index: 2, Input: "2\n", Output: "3 1\n"},
		},
	}

	problem := parsed.ToSchemaProblem()

	if problem == nil {
		t.Fatal("ToSchemaProblem() returned nil")
	}
	if problem.ID != "1325A" {
		t.Errorf("ID = %v, want 1325A", problem.ID)
	}
	if problem.Platform != "codeforces" {
		t.Errorf("Platform = %v, want codeforces", problem.Platform)
	}
	if problem.ContestID != 1325 {
		t.Errorf("ContestID = %v, want 1325", problem.ContestID)
	}
	if problem.Index != "A" {
		t.Errorf("Index = %v, want A", problem.Index)
	}
	if problem.Name != "EhAb AnD gCd" {
		t.Errorf("Name = %v, want EhAb AnD gCd", problem.Name)
	}
	if problem.URL != parsed.URL {
		t.Errorf("URL = %v, want %v", problem.URL, parsed.URL)
	}
	if problem.Limits.TimeLimit != "1 second" {
		t.Errorf("TimeLimit = %v, want '1 second'", problem.Limits.TimeLimit)
	}
	if problem.Limits.MemoryLimit != "256 megabytes" {
		t.Errorf("MemoryLimit = %v, want '256 megabytes'", problem.Limits.MemoryLimit)
	}
	if problem.Metadata.Rating != 800 {
		t.Errorf("Rating = %v, want 800", problem.Metadata.Rating)
	}
	if len(problem.Metadata.Tags) != 2 {
		t.Errorf("len(Tags) = %v, want 2", len(problem.Metadata.Tags))
	}
	if len(problem.Samples) != 2 {
		t.Errorf("len(Samples) = %v, want 2", len(problem.Samples))
	}
	if problem.Samples[0].Index != 1 {
		t.Errorf("Samples[0].Index = %v, want 1", problem.Samples[0].Index)
	}
}

func TestCleanTitle(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"A. EhAb AnD gCd", "EhAb AnD gCd"},
		{"B. Problem Name", "Problem Name"},
		{"B1. Subproblem", "Subproblem"},
		{"C2. Another One", "Another One"},
		{"  D. Spaces  ", "Spaces"},
		{"Problem Without Index", "Problem Without Index"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := cleanTitle(tt.input)
			if got != tt.want {
				t.Errorf("cleanTitle(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestExtractLimit(t *testing.T) {
	tests := []struct {
		text   string
		prefix string
		want   string
	}{
		{"time limit per test1 second", "time limit per test", "1 second"},
		{"memory limit per test256 megabytes", "memory limit per test", "256 megabytes"},
		{"  time limit per test 2 seconds  ", "time limit per test", "2 seconds"},
		{"something else", "time limit per test", "something else"},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			got := extractLimit(tt.text, tt.prefix)
			if got != tt.want {
				t.Errorf("extractLimit(%q, %q) = %q, want %q", tt.text, tt.prefix, got, tt.want)
			}
		})
	}
}

func TestCleanHTML(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"  hello world  ", "hello world"},
		{"multiple   spaces", "multiple spaces"},
		{"\n\t text \n\t", "text"},
		{"a  b   c    d", "a b c d"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := cleanHTML(tt.input)
			if got != tt.want {
				t.Errorf("cleanHTML(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestExtractProblemIndex(t *testing.T) {
	tests := []struct {
		href string
		want string
	}{
		{"/contest/123/problem/A", "A"},
		{"/contest/1/problem/B", "B"},
		{"/contest/999/problem/C1", "C1"},
		{"/contest/1325/problem/B2", "B2"},
		{"/contest/1/problems", ""},
		{"/problemset/problem/1/A", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.href, func(t *testing.T) {
			got := extractProblemIndex(tt.href)
			if got != tt.want {
				t.Errorf("extractProblemIndex(%q) = %q, want %q", tt.href, got, tt.want)
			}
		})
	}
}

func TestParseRating(t *testing.T) {
	tests := []struct {
		text string
		want int
	}{
		{"*800", 800},
		{"*1400", 1400},
		{"*2100", 2100},
		{"*3500", 3500},
		{"no rating", 0},
		{"", 0},
		{"* 800", 0}, // space after asterisk
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			got := parseRating(tt.text)
			if got != tt.want {
				t.Errorf("parseRating(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}

func TestExtractPreContent(t *testing.T) {
	html := `<pre>1
2
3</pre>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	pre := doc.Find("pre").First()
	got := extractPreContent(pre)

	if !strings.Contains(got, "1") || !strings.Contains(got, "2") || !strings.Contains(got, "3") {
		t.Errorf("extractPreContent() = %q, expected to contain 1, 2, 3", got)
	}
}

func TestExtractPreContent_WithBR(t *testing.T) {
	html := `<pre>1<br>2<br/>3</pre>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	pre := doc.Find("pre").First()
	got := extractPreContent(pre)

	lines := strings.Split(got, "\n")
	if len(lines) < 3 {
		t.Errorf("extractPreContent() should convert <br> to newlines, got: %q", got)
	}
}

func TestBuildStatement(t *testing.T) {
	html := `<div class="problem-statement">
		This is the problem text.
		<div class="input-specification">Input</div>
		<div class="output-specification">Output</div>
	</div>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	statement := doc.Find(".problem-statement").First()
	got := buildStatement(statement)

	// Should contain some text
	if got == "" {
		t.Error("buildStatement() returned empty string")
	}
}

func TestBuildStatement_Nil(t *testing.T) {
	got := buildStatement(nil)
	if got != "" {
		t.Errorf("buildStatement(nil) = %q, want empty", got)
	}
}

func TestParsedProblem_Fields(t *testing.T) {
	p := ParsedProblem{
		ContestID:   1325,
		Index:       "A",
		Name:        "Test",
		TimeLimit:   "1s",
		MemoryLimit: "256MB",
		Statement:   "Do something",
		InputSpec:   "Input spec",
		OutputSpec:  "Output spec",
		Note:        "Some note",
		Rating:      800,
		Tags:        []string{"math"},
		URL:         "https://codeforces.com/problem/1325/A",
		Samples: []Sample{
			{Index: 1, Input: "1", Output: "2"},
		},
	}

	if p.ContestID != 1325 {
		t.Errorf("ContestID = %v, want 1325", p.ContestID)
	}
	if p.Index != "A" {
		t.Errorf("Index = %v, want A", p.Index)
	}
	if len(p.Samples) != 1 {
		t.Errorf("len(Samples) = %v, want 1", len(p.Samples))
	}
}

func TestSample_Fields(t *testing.T) {
	s := Sample{
		Index:  1,
		Input:  "input data",
		Output: "output data",
	}

	if s.Index != 1 {
		t.Errorf("Index = %v, want 1", s.Index)
	}
	if s.Input != "input data" {
		t.Errorf("Input = %v, want 'input data'", s.Input)
	}
	if s.Output != "output data" {
		t.Errorf("Output = %v, want 'output data'", s.Output)
	}
}

func TestParser_Selectors(t *testing.T) {
	parser := NewParser(nil)

	// Verify parser uses current selectors
	if parser.selectors.Problem.Title != CurrentSelectors.Problem.Title {
		t.Error("Parser should use CurrentSelectors")
	}
}

func TestParsedProblem_ToSchemaProblem_EmptySamples(t *testing.T) {
	parsed := &ParsedProblem{
		ContestID: 1,
		Index:     "A",
		Name:      "Test",
		Samples:   nil,
	}

	problem := parsed.ToSchemaProblem()

	if len(problem.Samples) != 0 {
		t.Errorf("len(Samples) = %v, want 0", len(problem.Samples))
	}
}

func TestParsedProblem_ToSchemaProblem_EmptyTags(t *testing.T) {
	parsed := &ParsedProblem{
		ContestID: 1,
		Index:     "A",
		Name:      "Test",
		Tags:      nil,
	}

	problem := parsed.ToSchemaProblem()

	if problem.Metadata.Tags != nil && len(problem.Metadata.Tags) > 0 {
		t.Errorf("Tags should be empty, got %v", problem.Metadata.Tags)
	}
}

func TestParser_ParseSamples(t *testing.T) {
	html := `<div class="sample-tests">
		<div class="sample-test">
			<div class="input"><div class="title">Input</div><pre>1 2</pre></div>
			<div class="output"><div class="title">Output</div><pre>3</pre></div>
		</div>
		<div class="sample-test">
			<div class="input"><div class="title">Input</div><pre>4 5</pre></div>
			<div class="output"><div class="title">Output</div><pre>9</pre></div>
		</div>
	</div>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	samples := parseSamples(doc.Find(".sample-tests"), CurrentSelectors.Problem)
	if len(samples) != 2 {
		t.Errorf("len(samples) = %d, want 2", len(samples))
	}

	if len(samples) > 0 {
		if samples[0].Input != "1 2" {
			t.Errorf("samples[0].Input = %q, want '1 2'", samples[0].Input)
		}
		if samples[0].Output != "3" {
			t.Errorf("samples[0].Output = %q, want '3'", samples[0].Output)
		}
	}
}

func TestParser_ParseProblemHTML(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head><title>Mock Problem - Codeforces</title></head>
<body>
<div class="problem-statement">
	<div class="header">
		<div class="title">B. Mock Problem</div>
		<div class="time-limit"><div class="property-title">time limit per test</div>2 seconds</div>
		<div class="memory-limit"><div class="property-title">memory limit per test</div>512 megabytes</div>
	</div>
	<div class="problem-statement-text">
		<p>This is the problem statement.</p>
	</div>
	<div class="input-specification">
		<div class="section-title">Input</div>
		<p>The input format.</p>
	</div>
	<div class="output-specification">
		<div class="section-title">Output</div>
		<p>The output format.</p>
	</div>
	<div class="sample-tests">
		<div class="sample-test">
			<div class="input"><div class="title">Input</div><pre>1 2</pre></div>
			<div class="output"><div class="title">Output</div><pre>3</pre></div>
		</div>
	</div>
</div>
<span class="tag-box" title="Difficulty">*1200</span>
<span class="tag-box" style="font-size:1.2rem;">math</span>
</body>
</html>`

	parser := NewParser(nil)
	problem, err := parser.parseProblemHTML(strings.NewReader(html), 123, "B", "https://codeforces.com/contest/123/problem/B")
	if err != nil {
		t.Fatalf("parseProblemHTML() error = %v", err)
	}

	if problem.Name != "Mock Problem" {
		t.Errorf("Name = %v, want Mock Problem", problem.Name)
	}
	if problem.ContestID != 123 {
		t.Errorf("ContestID = %v, want 123", problem.ContestID)
	}
	if problem.Index != "B" {
		t.Errorf("Index = %v, want B", problem.Index)
	}
	if problem.Rating != 1200 {
		t.Errorf("Rating = %v, want 1200", problem.Rating)
	}
}
