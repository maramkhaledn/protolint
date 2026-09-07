package config

import (
	"regexp"
	"strings"

	"github.com/maramkhaledn/protolint/internal/filepathutil"
)

// Directories represents the target directories.
type Directories struct {
	Exclude      []string `yaml:"exclude" json:"exclude" toml:"exclude"`
	ExcludeRegex []string `yaml:"exclude_regex" json:"exclude_regex" toml:"exclude_regex"`

	// excludeRegexps is ExcludeRegex compiled once by validate.
	excludeRegexps []*regexp.Regexp
}

func (d Directories) shouldSkipRule(
	displayPath string,
) bool {
	for _, exclude := range d.Exclude {
		if !strings.HasSuffix(exclude, string(filepathutil.OSPathSeparator)) {
			exclude += string(filepathutil.OSPathSeparator)
		}
		if filepathutil.HasUnixPathPrefix(displayPath, exclude) {
			return true
		}
	}
	return matchesAnyPathRegexp(unixDirPath(displayPath), d.ExcludeRegex, d.excludeRegexps)
}

func (d *Directories) validate() error {
	compiled, err := compilePathRegexps("lint.directories", d.ExcludeRegex)
	if err != nil {
		return err
	}
	d.excludeRegexps = compiled
	return nil
}
