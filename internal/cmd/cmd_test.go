package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/maramkhaledn/protolint/internal/cmd"
	"github.com/maramkhaledn/protolint/internal/osutil"
	"github.com/maramkhaledn/protolint/internal/setting_test"
)

// TestDo_LintHonorsExcludeRegex runs the whole lint command against a directory tree
// whose .protolint.yaml excludes files and directories by regexp, to confirm the
// patterns really reach the file walking done by the command.
func TestDo_LintHonorsExcludeRegex(t *testing.T) {
	defer chdir(t, setting_test.TestDataPath("excluderegex"))()

	var stdout, stderr bytes.Buffer
	got := cmd.Do([]string{"lint", "."}, &stdout, &stderr)

	if got != osutil.ExitLintFailure {
		t.Errorf("got exit code %v, but want %v. stderr: %s", got, osutil.ExitLintFailure, stderr.String())
	}

	output := stderr.String()
	for _, want := range []string{
		"proto/search/search.proto",
		"proto/generated_extra/extra.proto",
	} {
		if !strings.Contains(output, want) {
			t.Errorf("got output %s, but want it to report %s", output, want)
		}
	}
	for _, notWant := range []string{
		"proto/search/search_generated.proto",
		"proto/generated/gen.proto",
	} {
		if strings.Contains(output, notWant) {
			t.Errorf("got output %s, but want it to exclude %s", output, notWant)
		}
	}
}

// TestDo_LintHonorsExcludeRegexForAnExplicitlyPassedFile makes sure an excluded file
// stays excluded even when it is named on the command line rather than walked into.
func TestDo_LintHonorsExcludeRegexForAnExplicitlyPassedFile(t *testing.T) {
	defer chdir(t, setting_test.TestDataPath("excluderegex"))()

	var stdout, stderr bytes.Buffer
	got := cmd.Do([]string{"lint", "proto/search/search_generated.proto"}, &stdout, &stderr)

	if got != osutil.ExitSuccess {
		t.Errorf("got exit code %v, but want %v. stderr: %s", got, osutil.ExitSuccess, stderr.String())
	}
	if stderr.Len() > 0 {
		t.Errorf("got stderr %s, but want empty stderr", stderr.String())
	}
}

// TestDo_LintReportsAnInvalidExcludeRegex checks that an unparsable pattern aborts the
// command instead of being silently ignored while linting.
func TestDo_LintReportsAnInvalidExcludeRegex(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".protolint.yaml")
	err := os.WriteFile(configPath, []byte(`
lint:
  files:
    exclude_regex:
      - "["
`), 0600)
	if err != nil {
		t.Fatalf("got err %v, but want nil", err)
	}

	defer chdir(t, setting_test.TestDataPath("excluderegex"))()

	var stdout, stderr bytes.Buffer
	got := cmd.Do([]string{"lint", "-config_path=" + configPath, "."}, &stdout, &stderr)

	if got != osutil.ExitInternalFailure {
		t.Errorf("got exit code %v, but want %v. stderr: %s", got, osutil.ExitInternalFailure, stderr.String())
	}
	want := `lint.files contains invalid exclude_regex pattern "["`
	if !strings.Contains(stderr.String(), want) {
		t.Errorf("got stderr %s, but want it to contain %s", stderr.String(), want)
	}
}

// chdir moves into dir and returns the func to move back.
func chdir(t *testing.T, dir string) func() {
	t.Helper()

	prev, err := os.Getwd()
	if err != nil {
		t.Fatalf("got err %v, but want nil", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("got err %v, but want nil", err)
	}
	return func() {
		if err := os.Chdir(prev); err != nil {
			t.Fatalf("got err %v, but want nil", err)
		}
	}
}
