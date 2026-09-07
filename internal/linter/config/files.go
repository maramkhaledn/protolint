package config

import "github.com/maramkhaledn/protolint/internal/stringsutil"

// Files represents the target files.
type Files struct {
	Exclude      []string `yaml:"exclude" json:"exclude" toml:"exclude"`
	ExcludeRegex []string `yaml:"exclude_regex" json:"exclude_regex" toml:"exclude_regex"`
}

func (d Files) shouldSkipRule(
	displayPath string,
) bool {
	return stringsutil.ContainsCrossPlatformPathInSlice(displayPath, d.Exclude) ||
		matchesAnyPathRegexp(displayPath, d.ExcludeRegex)
}

func (d Files) validate() error {
	return validatePathRegexps("lint.files", d.ExcludeRegex)
}
