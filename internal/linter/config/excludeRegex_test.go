package config_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/maramkhaledn/protolint/internal/linter/config"
	yaml "gopkg.in/yaml.v2"
)

// eachConfigForm runs fn against the two ways an ExternalConfig comes into being: loaded
// from a config file, which compiles the exclude_regex patterns up front, and assembled
// in code, which leaves them to be compiled while matching. Both have to behave the same.
//
// The rules the tests below ask about are always enabled, so that the rules section never
// skips them on its own and the exclusion behaviour under test stays observable.
func eachConfigForm(
	t *testing.T,
	lintBody string,
	fn func(t *testing.T, externalConfig config.ExternalConfig),
) {
	t.Helper()

	source := "lint:\n" + lintBody + `
  rules:
    no_default: true
    add:
      - FIELD_NAMES_LOWER_SNAKE_CASE
      - MESSAGE_NAMES_UPPER_CAMEL_CASE
`

	t.Run("loaded from a config file", func(t *testing.T) {
		externalConfig, err := config.GetExternalConfig(writeConfigFile(t, ".protolint.yaml", source), "")
		if err != nil {
			t.Fatalf("got err %v while loading %s, but want nil", err, source)
		}
		fn(t, *externalConfig)
	})

	t.Run("assembled in code", func(t *testing.T) {
		var externalConfig config.ExternalConfig
		if err := yaml.UnmarshalStrict([]byte(source), &externalConfig); err != nil {
			t.Fatalf("got err %v while unmarshalling %s, but want nil", err, source)
		}
		fn(t, externalConfig)
	})
}

// writeConfigFile writes the content into a new file under a fresh temp dir and returns
// its path, so each supported config format can be exercised through GetExternalConfig.
func writeConfigFile(t *testing.T, fileName, content string) string {
	t.Helper()

	filePath := filepath.Join(t.TempDir(), fileName)
	if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
		t.Fatalf("got err %v, but want nil", err)
	}
	return filePath
}

func TestExternalConfig_ShouldSkipRuleWithExcludeRegexPatterns(t *testing.T) {
	for _, test := range []struct {
		name             string
		inputLintConfig  string
		inputDisplayPath string
		wantSkipRule     bool
	}{
		{
			name: "not exclude anything when exclude_regex is absent",
			inputLintConfig: `  files:
    exclude:
      - path/to/foo.proto
`,
			inputDisplayPath: "proto/search/search_generated.proto",
		},
		{
			name: "not exclude anything when exclude_regex is empty",
			inputLintConfig: `  files:
    exclude_regex: []
`,
			inputDisplayPath: "proto/search/search_generated.proto",
		},
		{
			name: "exclude the file matching the first of the multiple patterns",
			inputLintConfig: `  files:
    exclude_regex:
      - "^vendor/"
      - ".*_pb\\.proto$"
`,
			inputDisplayPath: "vendor/search/search.proto",
			wantSkipRule:     true,
		},
		{
			name: "exclude the file matching the last of the multiple patterns",
			inputLintConfig: `  files:
    exclude_regex:
      - "^vendor/"
      - ".*_pb\\.proto$"
`,
			inputDisplayPath: "proto/search/search_pb.proto",
			wantSkipRule:     true,
		},
		{
			name: "not exclude the file matching none of the multiple patterns",
			inputLintConfig: `  files:
    exclude_regex:
      - "^vendor/"
      - ".*_pb\\.proto$"
`,
			inputDisplayPath: "proto/search/search.proto",
		},
		{
			name: "match the pattern anywhere in the path because it is not anchored",
			inputLintConfig: `  files:
    exclude_regex:
      - "generated"
`,
			inputDisplayPath: "proto/deep/generated/search.proto",
			wantSkipRule:     true,
		},
		{
			name: "not match the anchored pattern in the middle of the path",
			inputLintConfig: `  files:
    exclude_regex:
      - "^generated"
`,
			inputDisplayPath: "proto/generated/search.proto",
		},
		{
			name: "exclude by the plain exclude even when exclude_regex does not match",
			inputLintConfig: `  files:
    exclude:
      - proto/search/search.proto
    exclude_regex:
      - "^vendor/"
`,
			inputDisplayPath: "proto/search/search.proto",
			wantSkipRule:     true,
		},
		{
			name: "exclude by exclude_regex even when the plain exclude does not match",
			inputLintConfig: `  files:
    exclude:
      - proto/other/other.proto
    exclude_regex:
      - "^proto/search/"
`,
			inputDisplayPath: "proto/search/search.proto",
			wantSkipRule:     true,
		},
		{
			name: "exclude every file under the excluded directory tree",
			inputLintConfig: `  directories:
    exclude_regex:
      - "^proto/generated/"
`,
			inputDisplayPath: "proto/generated/nested/deeply/search.proto",
			wantSkipRule:     true,
		},
		{
			name: "not match the directories pattern against the file name",
			inputLintConfig: `  directories:
    exclude_regex:
      - "search"
`,
			inputDisplayPath: "proto/other/search.proto",
		},
		{
			name: "match the directories pattern against a directory in the middle of the path",
			inputLintConfig: `  directories:
    exclude_regex:
      - "generated"
`,
			inputDisplayPath: "proto/generated/search.proto",
			wantSkipRule:     true,
		},
		{
			name: "not match the directories pattern against a file name that looks like the directory",
			inputLintConfig: `  directories:
    exclude_regex:
      - "generated"
`,
			inputDisplayPath: "proto/other/generated_thing.proto",
		},
		{
			name: "match the directories pattern ending with a separator against the directory it names",
			inputLintConfig: `  directories:
    exclude_regex:
      - "^path/to/generated/"
`,
			inputDisplayPath: "path/to/generated/search.proto",
			wantSkipRule:     true,
		},
		{
			name: "not match any directories pattern for a file sitting in the working directory",
			inputLintConfig: `  directories:
    exclude_regex:
      - "^proto/"
`,
			inputDisplayPath: "search.proto",
		},
		{
			name: "still exclude the file by the files pattern matching the file name",
			inputLintConfig: `  files:
    exclude_regex:
      - "search"
`,
			inputDisplayPath: "proto/other/search.proto",
			wantSkipRule:     true,
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			eachConfigForm(t, test.inputLintConfig, func(t *testing.T, externalConfig config.ExternalConfig) {
				got := externalConfig.ShouldSkipRule(
					"FIELD_NAMES_LOWER_SNAKE_CASE",
					test.inputDisplayPath,
					nil,
				)
				if got != test.wantSkipRule {
					t.Errorf("got %v, but want %v", got, test.wantSkipRule)
				}
			})
		})
	}
}

// TestExternalConfig_ShouldSkipRuleWithExcludeRegexAppliesToEveryRule pins down that
// exclude_regex silences the whole file, unlike ignores which is scoped to one rule.
func TestExternalConfig_ShouldSkipRuleWithExcludeRegexAppliesToEveryRule(t *testing.T) {
	eachConfigForm(t, `  files:
    exclude_regex:
      - ".*_generated\\.proto$"
`, func(t *testing.T, externalConfig config.ExternalConfig) {
		for _, ruleID := range []string{
			"FIELD_NAMES_LOWER_SNAKE_CASE",
			"MESSAGE_NAMES_UPPER_CAMEL_CASE",
		} {
			ruleID := ruleID
			t.Run(ruleID, func(t *testing.T) {
				if !externalConfig.ShouldSkipRule(ruleID, "proto/search/search_generated.proto", nil) {
					t.Errorf("got false, but want true for the excluded file")
				}
				if externalConfig.ShouldSkipRule(ruleID, "proto/search/search.proto", nil) {
					t.Errorf("got true, but want false for the file that is not excluded")
				}
			})
		}
	})
}

func TestGetExternalConfigLoadsExcludeRegex(t *testing.T) {
	wantFilesRegex := []string{`.*_generated\.proto$`, "^vendor/"}
	wantDirsRegex := []string{"^proto/generated/"}

	for _, test := range []struct {
		name          string
		inputFileName string
		inputContent  string
	}{
		{
			name:          "yaml",
			inputFileName: ".protolint.yaml",
			inputContent: `
lint:
  files:
    exclude_regex:
      - ".*_generated\\.proto$"
      - "^vendor/"
  directories:
    exclude_regex:
      - "^proto/generated/"
`,
		},
		{
			name:          "json",
			inputFileName: "package.json",
			inputContent: `{
  "name": "example",
  "protolint": {
    "files": {
      "exclude_regex": [".*_generated\\.proto$", "^vendor/"]
    },
    "directories": {
      "exclude_regex": ["^proto/generated/"]
    }
  }
}`,
		},
		{
			name:          "toml",
			inputFileName: "pyproject.toml",
			inputContent: `
[tools.protolint.files]
exclude_regex = [".*_generated\\.proto$", "^vendor/"]

[tools.protolint.directories]
exclude_regex = ["^proto/generated/"]
`,
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			filePath := writeConfigFile(t, test.inputFileName, test.inputContent)

			got, err := config.GetExternalConfig(filePath, "")
			if err != nil {
				t.Fatalf("got err %v, but want nil", err)
			}

			if !reflect.DeepEqual(got.Lint.Files.ExcludeRegex, wantFilesRegex) {
				t.Errorf("got files.exclude_regex %v, but want %v", got.Lint.Files.ExcludeRegex, wantFilesRegex)
			}
			if !reflect.DeepEqual(got.Lint.Directories.ExcludeRegex, wantDirsRegex) {
				t.Errorf("got directories.exclude_regex %v, but want %v", got.Lint.Directories.ExcludeRegex, wantDirsRegex)
			}
		})
	}
}

func TestGetExternalConfigReturnsErrorForInvalidExcludeRegex(t *testing.T) {
	for _, test := range []struct {
		name             string
		inputFileName    string
		inputContent     string
		wantErrSubstring string
	}{
		{
			name:          "invalid files.exclude_regex in yaml",
			inputFileName: ".protolint.yaml",
			inputContent: `
lint:
  files:
    exclude_regex:
      - "["
`,
			wantErrSubstring: `lint.files contains invalid exclude_regex pattern "["`,
		},
		{
			name:          "invalid directories.exclude_regex in yaml",
			inputFileName: ".protolint.yaml",
			inputContent: `
lint:
  directories:
    exclude_regex:
      - "*.proto"
`,
			wantErrSubstring: `lint.directories contains invalid exclude_regex pattern "*.proto"`,
		},
		{
			name:          "invalid files.exclude_regex in json",
			inputFileName: "package.json",
			inputContent: `{
  "protolint": {
    "files": {
      "exclude_regex": ["("]
    }
  }
}`,
			wantErrSubstring: `lint.files contains invalid exclude_regex pattern "("`,
		},
		{
			name:          "invalid directories.exclude_regex in json",
			inputFileName: "package.json",
			inputContent: `{
  "protolint": {
    "directories": {
      "exclude_regex": ["a{2,1}"]
    }
  }
}`,
			wantErrSubstring: `lint.directories contains invalid exclude_regex pattern "a{2,1}"`,
		},
		{
			name:          "invalid files.exclude_regex in toml",
			inputFileName: "pyproject.toml",
			inputContent: `
[tools.protolint.files]
exclude_regex = ["("]
`,
			wantErrSubstring: `lint.files contains invalid exclude_regex pattern "("`,
		},
		{
			name:          "invalid directories.exclude_regex in toml",
			inputFileName: "pyproject.toml",
			inputContent: `
[tools.protolint.directories]
exclude_regex = ["a{2,1}"]
`,
			wantErrSubstring: `lint.directories contains invalid exclude_regex pattern "a{2,1}"`,
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			filePath := writeConfigFile(t, test.inputFileName, test.inputContent)

			got, err := config.GetExternalConfig(filePath, "")
			if err == nil {
				t.Fatalf("got config %v and err nil, but want err", got)
			}
			if !strings.Contains(err.Error(), test.wantErrSubstring) {
				t.Errorf("got err %v, but want it to contain %s", err, test.wantErrSubstring)
			}
		})
	}
}

// TestGetExternalConfigReturnsErrorForInvalidExcludeRegexFoundByDirSearch covers the
// config file that is discovered by searching a directory instead of being named.
func TestGetExternalConfigReturnsErrorForInvalidExcludeRegexFoundByDirSearch(t *testing.T) {
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, ".protolint.yaml"), []byte(`
lint:
  files:
    exclude_regex:
      - "["
`), 0600)
	if err != nil {
		t.Fatalf("got err %v, but want nil", err)
	}

	got, err := config.GetExternalConfig("", dir)
	if err == nil {
		t.Fatalf("got config %v and err nil, but want err", got)
	}
}

func TestGetExternalConfigReturnsNoErrorForValidExcludeRegex(t *testing.T) {
	filePath := writeConfigFile(t, ".protolint.yaml", `
lint:
  files:
    exclude_regex:
      - ".*_generated\\.proto$"
  directories:
    exclude_regex:
      - "^proto/generated/"
`)

	if _, err := config.GetExternalConfig(filePath, ""); err != nil {
		t.Errorf("got err %v, but want nil", err)
	}
}
