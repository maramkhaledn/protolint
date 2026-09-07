package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGetExternalConfigCompilesExcludeRegexOnce guards the reason compilePathRegexps
// exists: ShouldSkipRule runs for every rule of every file, so the patterns have to be
// compiled while the config is loaded rather than on each match.
func TestGetExternalConfigCompilesExcludeRegexOnce(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), ".protolint.yaml")
	err := os.WriteFile(filePath, []byte(`
lint:
  files:
    exclude_regex:
      - .*_generated\.proto$
      - ^vendor/
  directories:
    exclude_regex:
      - ^proto/generated/
`), 0600)
	if err != nil {
		t.Fatalf("got err %v, but want nil", err)
	}

	got, err := GetExternalConfig(filePath, "")
	if err != nil {
		t.Fatalf("got err %v, but want nil", err)
	}

	if len(got.Lint.Files.excludeRegexps) != len(got.Lint.Files.ExcludeRegex) {
		t.Errorf(
			"got %d compiled files patterns, but want %d",
			len(got.Lint.Files.excludeRegexps), len(got.Lint.Files.ExcludeRegex),
		)
	}
	if len(got.Lint.Directories.excludeRegexps) != len(got.Lint.Directories.ExcludeRegex) {
		t.Errorf(
			"got %d compiled directories patterns, but want %d",
			len(got.Lint.Directories.excludeRegexps), len(got.Lint.Directories.ExcludeRegex),
		)
	}
}

func TestCompilePathRegexps(t *testing.T) {
	for _, test := range []struct {
		name          string
		inputPatterns []string
		wantLen       int
		wantErr       bool
	}{
		{
			name:          "compiles nothing for no patterns",
			inputPatterns: nil,
		},
		{
			name:          "compiles one regexp per pattern",
			inputPatterns: []string{`.*_generated\.proto$`, "^vendor/"},
			wantLen:       2,
		},
		{
			name:          "reports the first pattern that does not compile",
			inputPatterns: []string{"^vendor/", "["},
			wantErr:       true,
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got, err := compilePathRegexps("lint.files", test.inputPatterns)

			if test.wantErr {
				if err == nil {
					t.Fatalf("got %v and err nil, but want err", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("got err %v, but want nil", err)
			}
			if len(got) != test.wantLen {
				t.Errorf("got %d compiled patterns, but want %d", len(got), test.wantLen)
			}
		})
	}
}

func TestUnixDirPath(t *testing.T) {
	for _, test := range []struct {
		name             string
		inputDisplayPath string
		want             string
	}{
		{
			name:             "keeps the trailing separator of the directory",
			inputDisplayPath: "proto/generated/search.proto",
			want:             "proto/generated/",
		},
		{
			name:             "returns an empty path for a file in the working directory",
			inputDisplayPath: "search.proto",
			want:             "",
		},
		{
			name:             "keeps every directory of a deep path",
			inputDisplayPath: "proto/a/b/c/search.proto",
			want:             "proto/a/b/c/",
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got := unixDirPath(test.inputDisplayPath)
			if got != test.want {
				t.Errorf("got %q, but want %q", got, test.want)
			}
		})
	}
}
