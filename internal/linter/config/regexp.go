package config

import (
	"fmt"
	"regexp"

	"github.com/maramkhaledn/protolint/internal/filepathutil"
)

func validatePathRegexps(configPath string, patterns []string) error {
	for _, pattern := range patterns {
		if _, err := regexp.Compile(pattern); err != nil {
			return fmt.Errorf("%s contains invalid exclude_regex pattern %q: %w", configPath, pattern, err)
		}
	}
	return nil
}

func matchesAnyPathRegexp(displayPath string, patterns []string) bool {
	unixDisplayPath := filepathutil.ToUnixPath(displayPath)
	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, unixDisplayPath)
		if err == nil && matched {
			return true
		}
	}
	return false
}
