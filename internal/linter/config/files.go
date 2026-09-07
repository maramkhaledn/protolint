package config

import (
	"regexp"

	"github.com/maramkhaledn/protolint/internal/filepathutil"
	"github.com/maramkhaledn/protolint/internal/stringsutil"
)

// Files represents the target files.
type Files struct {
	Exclude      []string `yaml:"exclude" json:"exclude" toml:"exclude"`
	ExcludeRegex []string `yaml:"exclude_regex" json:"exclude_regex" toml:"exclude_regex"`

	// excludeRegexps is ExcludeRegex compiled once by validate.
	excludeRegexps []*regexp.Regexp
}

func (d Files) shouldSkipRule(
	displayPath string,
) bool {
	return stringsutil.ContainsCrossPlatformPathInSlice(displayPath, d.Exclude) ||
		matchesAnyPathRegexp(filepathutil.ToUnixPath(displayPath), d.ExcludeRegex, d.excludeRegexps)
}

func (d *Files) validate() error {
	compiled, err := compilePathRegexps("lint.files", d.ExcludeRegex)
	if err != nil {
		return err
	}
	d.excludeRegexps = compiled
	return nil
}
