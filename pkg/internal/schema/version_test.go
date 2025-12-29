package schema

import (
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Version
		wantErr bool
	}{
		{
			name:  "valid version",
			input: "1.2.3",
			want:  Version{Major: 1, Minor: 2, Patch: 3},
		},
		{
			name:  "version 1.0.0",
			input: "1.0.0",
			want:  Version{Major: 1, Minor: 0, Patch: 0},
		},
		{
			name:  "high numbers",
			input: "10.20.30",
			want:  Version{Major: 10, Minor: 20, Patch: 30},
		},
		{
			name:    "invalid format - two parts",
			input:   "1.2",
			wantErr: true,
		},
		{
			name:    "invalid format - four parts",
			input:   "1.2.3.4",
			wantErr: true,
		},
		{
			name:    "invalid format - non-numeric",
			input:   "a.b.c",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:  "negative number - parsed as is",
			input: "-1.0.0",
			want:  Version{Major: -1, Minor: 0, Patch: 0}, // strconv.Atoi accepts negative
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseVersion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_String(t *testing.T) {
	tests := []struct {
		name    string
		version Version
		want    string
	}{
		{
			name:    "version 1.0.0",
			version: Version{Major: 1, Minor: 0, Patch: 0},
			want:    "1.0.0",
		},
		{
			name:    "version 1.2.3",
			version: Version{Major: 1, Minor: 2, Patch: 3},
			want:    "1.2.3",
		},
		{
			name:    "zero version",
			version: Version{Major: 0, Minor: 0, Patch: 0},
			want:    "0.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.version.String(); got != tt.want {
				t.Errorf("Version.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_IsCompatible(t *testing.T) {
	tests := []struct {
		name    string
		current Version
		other   Version
		want    bool
	}{
		{
			name:    "same version",
			current: Version{Major: 1, Minor: 0, Patch: 0},
			other:   Version{Major: 1, Minor: 0, Patch: 0},
			want:    true,
		},
		{
			name:    "same major, different minor",
			current: Version{Major: 1, Minor: 2, Patch: 0},
			other:   Version{Major: 1, Minor: 1, Patch: 0},
			want:    true,
		},
		{
			name:    "same major, different patch",
			current: Version{Major: 1, Minor: 0, Patch: 5},
			other:   Version{Major: 1, Minor: 0, Patch: 3},
			want:    true,
		},
		{
			name:    "different major - incompatible",
			current: Version{Major: 2, Minor: 0, Patch: 0},
			other:   Version{Major: 1, Minor: 0, Patch: 0},
			want:    false,
		},
		{
			name:    "future major - incompatible",
			current: Version{Major: 1, Minor: 0, Patch: 0},
			other:   Version{Major: 2, Minor: 0, Patch: 0},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.current.IsCompatible(tt.other); got != tt.want {
				t.Errorf("Version.IsCompatible() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_NeedsMigration(t *testing.T) {
	tests := []struct {
		name    string
		current Version
		other   Version
		want    bool
	}{
		{
			name:    "same version - no migration",
			current: Version{Major: 1, Minor: 0, Patch: 0},
			other:   Version{Major: 1, Minor: 0, Patch: 0},
			want:    false,
		},
		{
			name:    "older minor - needs migration",
			current: Version{Major: 1, Minor: 2, Patch: 0},
			other:   Version{Major: 1, Minor: 1, Patch: 0},
			want:    true,
		},
		{
			name:    "same major/minor different patch - no migration",
			current: Version{Major: 1, Minor: 0, Patch: 5},
			other:   Version{Major: 1, Minor: 0, Patch: 3},
			want:    false, // Patch differences don't require migration
		},
		{
			name:    "different minor - needs migration",
			current: Version{Major: 1, Minor: 0, Patch: 0},
			other:   Version{Major: 1, Minor: 1, Patch: 0},
			want:    true, // Different minor versions need migration
		},
		{
			name:    "different major - needs migration",
			current: Version{Major: 1, Minor: 0, Patch: 0},
			other:   Version{Major: 2, Minor: 0, Patch: 0},
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.current.NeedsMigration(tt.other); got != tt.want {
				t.Errorf("Version.NeedsMigration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_Compare(t *testing.T) {
	tests := []struct {
		name string
		v1   Version
		v2   Version
		want int
	}{
		{
			name: "equal versions",
			v1:   Version{Major: 1, Minor: 0, Patch: 0},
			v2:   Version{Major: 1, Minor: 0, Patch: 0},
			want: 0,
		},
		{
			name: "v1 greater major",
			v1:   Version{Major: 2, Minor: 0, Patch: 0},
			v2:   Version{Major: 1, Minor: 0, Patch: 0},
			want: 1,
		},
		{
			name: "v1 lesser major",
			v1:   Version{Major: 1, Minor: 0, Patch: 0},
			v2:   Version{Major: 2, Minor: 0, Patch: 0},
			want: -1,
		},
		{
			name: "v1 greater minor",
			v1:   Version{Major: 1, Minor: 2, Patch: 0},
			v2:   Version{Major: 1, Minor: 1, Patch: 0},
			want: 1,
		},
		{
			name: "v1 greater patch",
			v1:   Version{Major: 1, Minor: 0, Patch: 2},
			v2:   Version{Major: 1, Minor: 0, Patch: 1},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v1.Compare(tt.v2); got != tt.want {
				t.Errorf("Version.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSchemaHeader(t *testing.T) {
	header := NewSchemaHeader(TypeProblem)

	if header.Version != CurrentVersion.String() {
		t.Errorf("NewSchemaHeader().Version = %v, want %v", header.Version, CurrentVersion.String())
	}
	if header.Type != TypeProblem {
		t.Errorf("NewSchemaHeader().Type = %v, want %v", header.Type, TypeProblem)
	}
}

func TestSchemaTypes(t *testing.T) {
	// Verify all schema types are defined
	types := []string{TypeWorkspace, TypeProblem, TypeSubmission, TypeProgress}
	for _, st := range types {
		if st == "" {
			t.Error("Schema type should not be empty")
		}
	}
}

func TestCurrentVersion(t *testing.T) {
	// Verify current version is valid
	if CurrentVersion.Major < 1 {
		t.Error("CurrentVersion.Major should be at least 1")
	}
	if CurrentVersion.Minor < 0 {
		t.Error("CurrentVersion.Minor should be non-negative")
	}
	if CurrentVersion.Patch < 0 {
		t.Error("CurrentVersion.Patch should be non-negative")
	}
}
