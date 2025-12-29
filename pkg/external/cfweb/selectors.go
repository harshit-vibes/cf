// Package cfweb provides web scraping and submission for Codeforces
package cfweb

// SelectorVersion represents a versioned set of CSS selectors
// When CF changes their HTML structure, we add a new version
type SelectorVersion struct {
	Version     string
	ValidFrom   string // Date when these selectors became valid
	ValidUntil  string // Empty if current
	Description string
}

// Selectors contains all CSS selectors for scraping CF pages
type Selectors struct {
	// Problem page selectors
	Problem ProblemSelectors

	// Login page selectors
	Login LoginSelectors

	// Submit page selectors
	Submit SubmitSelectors

	// Contest page selectors
	Contest ContestSelectors
}

// ProblemSelectors for problem page parsing
type ProblemSelectors struct {
	// Main content
	Title             string
	TimeLimit         string
	MemoryLimit       string
	Statement         string
	InputSpec         string
	OutputSpec        string
	Note              string

	// Sample tests
	SampleTests       string
	SampleInput       string
	SampleOutput      string

	// Tags and rating
	Tags              string
	Rating            string
}

// LoginSelectors for login page
type LoginSelectors struct {
	Form           string
	HandleOrEmail  string
	Password       string
	CSRFToken      string
	Remember       string
	SubmitButton   string
}

// SubmitSelectors for submit page
type SubmitSelectors struct {
	Form            string
	ProblemIndex    string
	LanguageSelect  string
	SourceCode      string
	CSRFToken       string
	SubmitButton    string
	FTAA            string
	BFAA            string
}

// ContestSelectors for contest page
type ContestSelectors struct {
	ProblemList     string
	ProblemRow      string
	ProblemLink     string
	ProblemName     string
	StandingsTable  string
}

// CurrentVersion returns the current selector version
var CurrentVersion = SelectorVersion{
	Version:     "2024.12",
	ValidFrom:   "2024-12-01",
	ValidUntil:  "",
	Description: "CF selectors as of December 2024",
}

// CurrentSelectors returns the current set of selectors
var CurrentSelectors = Selectors{
	Problem: ProblemSelectors{
		Title:        ".problem-statement .title",
		TimeLimit:    ".problem-statement .time-limit",
		MemoryLimit:  ".problem-statement .memory-limit",
		Statement:    ".problem-statement",
		InputSpec:    ".problem-statement .input-specification",
		OutputSpec:   ".problem-statement .output-specification",
		Note:         ".problem-statement .note",
		SampleTests:  ".sample-tests",
		SampleInput:  ".sample-tests .input pre",
		SampleOutput: ".sample-tests .output pre",
		Tags:        ".tag-box",
		Rating:      "span.tag-box[title='Difficulty']",
	},
	Login: LoginSelectors{
		Form:          "form#enterForm, form.enter-form",
		HandleOrEmail: "input[name='handleOrEmail']",
		Password:      "input[name='password']",
		CSRFToken:     "input[name='csrf_token'], meta[name='X-Csrf-Token']",
		Remember:      "input[name='remember']",
		SubmitButton:  "input[type='submit']",
	},
	Submit: SubmitSelectors{
		Form:           "form.submit-form",
		ProblemIndex:   "select[name='submittedProblemIndex'], input[name='submittedProblemIndex']",
		LanguageSelect: "select[name='programTypeId']",
		SourceCode:     "textarea[name='source'], #sourceCodeTextarea",
		CSRFToken:      "input[name='csrf_token']",
		SubmitButton:   "input[type='submit']",
		FTAA:           "input[name='ftaa']",
		BFAA:           "input[name='bfaa']",
	},
	Contest: ContestSelectors{
		ProblemList:    ".problems",
		ProblemRow:     ".problems tr",
		ProblemLink:    "td.id a",
		ProblemName:    "td a",
		StandingsTable: ".standings",
	},
}

// Language represents a programming language
type Language struct {
	ID          string
	Name        string
	Extension   string
	CompilerID  int
}

// SupportedLanguages lists all supported languages for submission
var SupportedLanguages = []Language{
	{ID: "cpp17", Name: "GNU G++17 7.3.0", Extension: ".cpp", CompilerID: 54},
	{ID: "cpp20", Name: "GNU G++20 11.2.0 (64 bit)", Extension: ".cpp", CompilerID: 89},
	{ID: "cpp23", Name: "GNU G++23 14.2 (64 bit)", Extension: ".cpp", CompilerID: 91},
	{ID: "python3", Name: "Python 3.8.10", Extension: ".py", CompilerID: 31},
	{ID: "pypy3", Name: "PyPy 3.10 (7.3.15)", Extension: ".py", CompilerID: 70},
	{ID: "java17", Name: "Java 17 64bit", Extension: ".java", CompilerID: 87},
	{ID: "java21", Name: "Java 21 64bit", Extension: ".java", CompilerID: 88},
	{ID: "go", Name: "Go 1.22.2", Extension: ".go", CompilerID: 32},
	{ID: "rust", Name: "Rust 1.75.0 (2021)", Extension: ".rs", CompilerID: 75},
	{ID: "kotlin", Name: "Kotlin 1.9.21", Extension: ".kt", CompilerID: 83},
	{ID: "csharp", Name: "C# 10, .NET SDK 6.0", Extension: ".cs", CompilerID: 79},
	{ID: "ruby", Name: "Ruby 3.2.2", Extension: ".rb", CompilerID: 67},
	{ID: "js", Name: "JavaScript V8 4.8.0", Extension: ".js", CompilerID: 34},
	{ID: "php", Name: "PHP 8.1.7", Extension: ".php", CompilerID: 6},
	{ID: "haskell", Name: "Haskell GHC 8.10.1", Extension: ".hs", CompilerID: 12},
	{ID: "scala", Name: "Scala 2.12.8", Extension: ".scala", CompilerID: 20},
}

// GetLanguageByExtension returns language info by file extension
func GetLanguageByExtension(ext string) *Language {
	for i := range SupportedLanguages {
		if SupportedLanguages[i].Extension == ext {
			return &SupportedLanguages[i]
		}
	}
	return nil
}

// GetLanguageByID returns language info by ID
func GetLanguageByID(id string) *Language {
	for i := range SupportedLanguages {
		if SupportedLanguages[i].ID == id {
			return &SupportedLanguages[i]
		}
	}
	return nil
}

// GetLanguageByCompilerID returns language info by compiler ID
func GetLanguageByCompilerID(compilerID int) *Language {
	for i := range SupportedLanguages {
		if SupportedLanguages[i].CompilerID == compilerID {
			return &SupportedLanguages[i]
		}
	}
	return nil
}
