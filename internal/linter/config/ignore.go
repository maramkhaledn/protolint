package config

import (
	"fmt"
	"regexp"

	"github.com/maramkhaledn/protolint/internal/filepathutil"
	"github.com/maramkhaledn/protolint/internal/stringsutil"
)

// Ignore represents files ignoring the specific rule.
type Ignore struct {
	ID           string   `yaml:"id" json:"id" toml:"id"`
	Files        []string `yaml:"files" json:"files" toml:"files"`
	ExcludeRegex []string `yaml:"exclude_regex" json:"exclude_regex" toml:"exclude_regex"`

	// excludeRegexps is ExcludeRegex compiled once by validate.
	excludeRegexps []*regexp.Regexp
}

func (i Ignore) shouldSkipRule(
	ruleID string,
	displayPath string,
) bool {
	if i.ID != ruleID {
		return false
	}
	return stringsutil.ContainsCrossPlatformPathInSlice(displayPath, i.Files) ||
		matchesAnyPathRegexp(filepathutil.ToUnixPath(displayPath), i.ExcludeRegex, i.excludeRegexps)
}

func (i *Ignore) validate(index int) error {
	compiled, err := compilePathRegexps(fmt.Sprintf("lint.ignores[%d]", index), i.ExcludeRegex)
	if err != nil {
		return err
	}
	i.excludeRegexps = compiled
	return nil
}
