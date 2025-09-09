package config

// MessageNamesHaveResponseOrRequestSuffixOption represents the option for the MESSAGE_NAMES_HAVE_RESPONSE_OR_REQUEST_SUFFIX rule.
type MessageNamesHaveResponseOrRequestSuffixOption struct {
	CustomizableSeverityOption `yaml:",inline"`
	ShouldFollowGolangStyle    bool `yaml:"should_follow_golang_style" json:"should_follow_golang_style" toml:"should_follow_golang_style"`

}
