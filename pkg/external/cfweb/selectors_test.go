package cfweb

import (
	"testing"
)

func TestCurrentVersion(t *testing.T) {
	if CurrentVersion.Version == "" {
		t.Error("CurrentVersion.Version should not be empty")
	}
	if CurrentVersion.ValidFrom == "" {
		t.Error("CurrentVersion.ValidFrom should not be empty")
	}
	if CurrentVersion.Description == "" {
		t.Error("CurrentVersion.Description should not be empty")
	}
}

func TestCurrentSelectors_Problem(t *testing.T) {
	sel := CurrentSelectors.Problem

	if sel.Title == "" {
		t.Error("Problem.Title selector should not be empty")
	}
	if sel.TimeLimit == "" {
		t.Error("Problem.TimeLimit selector should not be empty")
	}
	if sel.MemoryLimit == "" {
		t.Error("Problem.MemoryLimit selector should not be empty")
	}
	if sel.Statement == "" {
		t.Error("Problem.Statement selector should not be empty")
	}
	if sel.SampleTests == "" {
		t.Error("Problem.SampleTests selector should not be empty")
	}
	if sel.Tags == "" {
		t.Error("Problem.Tags selector should not be empty")
	}
}

func TestCurrentSelectors_Login(t *testing.T) {
	sel := CurrentSelectors.Login

	if sel.Form == "" {
		t.Error("Login.Form selector should not be empty")
	}
	if sel.HandleOrEmail == "" {
		t.Error("Login.HandleOrEmail selector should not be empty")
	}
	if sel.Password == "" {
		t.Error("Login.Password selector should not be empty")
	}
	if sel.CSRFToken == "" {
		t.Error("Login.CSRFToken selector should not be empty")
	}
}

func TestCurrentSelectors_Submit(t *testing.T) {
	sel := CurrentSelectors.Submit

	if sel.Form == "" {
		t.Error("Submit.Form selector should not be empty")
	}
	if sel.ProblemIndex == "" {
		t.Error("Submit.ProblemIndex selector should not be empty")
	}
	if sel.LanguageSelect == "" {
		t.Error("Submit.LanguageSelect selector should not be empty")
	}
	if sel.SourceCode == "" {
		t.Error("Submit.SourceCode selector should not be empty")
	}
	if sel.CSRFToken == "" {
		t.Error("Submit.CSRFToken selector should not be empty")
	}
}

func TestCurrentSelectors_Contest(t *testing.T) {
	sel := CurrentSelectors.Contest

	if sel.ProblemList == "" {
		t.Error("Contest.ProblemList selector should not be empty")
	}
	if sel.ProblemRow == "" {
		t.Error("Contest.ProblemRow selector should not be empty")
	}
	if sel.ProblemLink == "" {
		t.Error("Contest.ProblemLink selector should not be empty")
	}
}

func TestSupportedLanguages(t *testing.T) {
	if len(SupportedLanguages) == 0 {
		t.Fatal("SupportedLanguages should not be empty")
	}

	for i, lang := range SupportedLanguages {
		if lang.ID == "" {
			t.Errorf("Language %d: ID should not be empty", i)
		}
		if lang.Name == "" {
			t.Errorf("Language %d: Name should not be empty", i)
		}
		if lang.Extension == "" {
			t.Errorf("Language %d: Extension should not be empty", i)
		}
		if lang.CompilerID == 0 {
			t.Errorf("Language %d: CompilerID should not be zero", i)
		}
	}
}

func TestGetLanguageByExtension(t *testing.T) {
	tests := []struct {
		ext      string
		wantID   string
		wantNil  bool
	}{
		{".cpp", "cpp17", false},
		{".py", "python3", false},
		{".java", "java17", false},
		{".go", "go", false},
		{".rs", "rust", false},
		{".kt", "kotlin", false},
		{".cs", "csharp", false},
		{".rb", "ruby", false},
		{".js", "js", false},
		{".xyz", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			lang := GetLanguageByExtension(tt.ext)
			if tt.wantNil {
				if lang != nil {
					t.Errorf("GetLanguageByExtension(%s) = %v, want nil", tt.ext, lang)
				}
			} else {
				if lang == nil {
					t.Errorf("GetLanguageByExtension(%s) = nil, want %s", tt.ext, tt.wantID)
				} else if lang.ID != tt.wantID {
					t.Errorf("GetLanguageByExtension(%s).ID = %v, want %v", tt.ext, lang.ID, tt.wantID)
				}
			}
		})
	}
}

func TestGetLanguageByID(t *testing.T) {
	tests := []struct {
		id       string
		wantName string
		wantNil  bool
	}{
		{"cpp17", "GNU G++17 7.3.0", false},
		{"cpp20", "GNU G++20 11.2.0 (64 bit)", false},
		{"python3", "Python 3.8.10", false},
		{"pypy3", "PyPy 3.10 (7.3.15)", false},
		{"go", "Go 1.22.2", false},
		{"rust", "Rust 1.75.0 (2021)", false},
		{"nonexistent", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			lang := GetLanguageByID(tt.id)
			if tt.wantNil {
				if lang != nil {
					t.Errorf("GetLanguageByID(%s) = %v, want nil", tt.id, lang)
				}
			} else {
				if lang == nil {
					t.Errorf("GetLanguageByID(%s) = nil", tt.id)
				} else if lang.Name != tt.wantName {
					t.Errorf("GetLanguageByID(%s).Name = %v, want %v", tt.id, lang.Name, tt.wantName)
				}
			}
		})
	}
}

func TestGetLanguageByCompilerID(t *testing.T) {
	tests := []struct {
		compilerID int
		wantID     string
		wantNil    bool
	}{
		{54, "cpp17", false},
		{89, "cpp20", false},
		{31, "python3", false},
		{32, "go", false},
		{75, "rust", false},
		{9999, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.wantID, func(t *testing.T) {
			lang := GetLanguageByCompilerID(tt.compilerID)
			if tt.wantNil {
				if lang != nil {
					t.Errorf("GetLanguageByCompilerID(%d) = %v, want nil", tt.compilerID, lang)
				}
			} else {
				if lang == nil {
					t.Errorf("GetLanguageByCompilerID(%d) = nil", tt.compilerID)
				} else if lang.ID != tt.wantID {
					t.Errorf("GetLanguageByCompilerID(%d).ID = %v, want %v", tt.compilerID, lang.ID, tt.wantID)
				}
			}
		})
	}
}

func TestSelectorVersion_Fields(t *testing.T) {
	sv := SelectorVersion{
		Version:     "2024.12",
		ValidFrom:   "2024-12-01",
		ValidUntil:  "",
		Description: "Test selectors",
	}

	if sv.Version != "2024.12" {
		t.Errorf("Version = %v, want 2024.12", sv.Version)
	}
	if sv.ValidFrom != "2024-12-01" {
		t.Errorf("ValidFrom = %v, want 2024-12-01", sv.ValidFrom)
	}
}

func TestLanguage_Fields(t *testing.T) {
	lang := Language{
		ID:         "cpp17",
		Name:       "GNU G++17 7.3.0",
		Extension:  ".cpp",
		CompilerID: 54,
	}

	if lang.ID != "cpp17" {
		t.Errorf("ID = %v, want cpp17", lang.ID)
	}
	if lang.Extension != ".cpp" {
		t.Errorf("Extension = %v, want .cpp", lang.Extension)
	}
	if lang.CompilerID != 54 {
		t.Errorf("CompilerID = %v, want 54", lang.CompilerID)
	}
}

func TestUniqueLanguageIDs(t *testing.T) {
	seen := make(map[string]bool)
	for _, lang := range SupportedLanguages {
		if seen[lang.ID] {
			t.Errorf("Duplicate language ID: %s", lang.ID)
		}
		seen[lang.ID] = true
	}
}

func TestUniqueCompilerIDs(t *testing.T) {
	seen := make(map[int]bool)
	for _, lang := range SupportedLanguages {
		if seen[lang.CompilerID] {
			t.Errorf("Duplicate compiler ID: %d (%s)", lang.CompilerID, lang.ID)
		}
		seen[lang.CompilerID] = true
	}
}

func TestCppLanguagesExist(t *testing.T) {
	cppVersions := []string{"cpp17", "cpp20", "cpp23"}
	for _, v := range cppVersions {
		if GetLanguageByID(v) == nil {
			t.Errorf("Expected C++ version %s to be supported", v)
		}
	}
}

func TestPythonLanguagesExist(t *testing.T) {
	pythonVersions := []string{"python3", "pypy3"}
	for _, v := range pythonVersions {
		if GetLanguageByID(v) == nil {
			t.Errorf("Expected Python version %s to be supported", v)
		}
	}
}
