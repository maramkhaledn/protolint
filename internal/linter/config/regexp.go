package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/maramkhaledn/protolint/internal/filepathutil"
)

// compilePathRegexps compiles the exclude_regex patterns while the config is loaded and
// reports the first one that is not a valid regexp. ShouldSkipRule runs once per rule
// per file, so compiling the patterns on every match would repeat the same work for
// each of the dozens of rules that are checked against a single file.
func compilePathRegexps(configPath string, patterns []string) ([]*regexp.Regexp, error) {
	if len(patterns) == 0 {
		return nil, nil
	}

	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, pattern := range patterns {
		r, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("%s contains invalid exclude_regex pattern %q: %w", configPath, pattern, err)
		}
		compiled = append(compiled, r)
	}
	return compiled, nil
}

// matchesAnyPathRegexp reports whether the UNIX-style path matches any of the patterns.
//
// compiled holds the forms built by compilePathRegexps. It is empty when the config was
// assembled in code instead of being read from a file, and the patterns are compiled on
// the spot in that case. Patterns that do not compile are skipped there, because loading
// a config file is the only place left to report the error to.
func matchesAnyPathRegexp(unixPath string, patterns []string, compiled []*regexp.Regexp) bool {
	if len(compiled) == len(patterns) {
		for _, r := range compiled {
			if r.MatchString(unixPath) {
				return true
			}
		}
		return false
	}

	for _, pattern := range patterns {
		if matched, err := regexp.MatchString(pattern, unixPath); err == nil && matched {
			return true
		}
	}
	return false
}

// unixDirPath returns the directory part of a display path in the UNIX form, keeping
// the trailing separator so that a pattern naming a directory, like "^proto/generated/",
// matches it. The result is empty for a file sitting in the working directory itself.
func unixDirPath(displayPath string) string {
	unixPath := filepathutil.ToUnixPath(displayPath)
	return unixPath[:strings.LastIndex(unixPath, "/")+1]
}
