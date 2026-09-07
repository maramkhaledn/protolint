package config

import (
	"strings"

	"github.com/maramkhaledn/protolint/internal/filepathutil"
)

// Directories represents the target directories.
type Directories struct {
	Exclude      []string `yaml:"exclude" json:"exclude" toml:"exclude"`
	ExcludeRegex []string `yaml:"exclude_regex" json:"exclude_regex" toml:"exclude_regex"`
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
	return matchesAnyPathRegexp(displayPath, d.ExcludeRegex)
}

func (d Directories) validate() error {
	return validatePathRegexps("lint.directories", d.ExcludeRegex)
}
