package config

// MaxRpcSummaryLengthOption represents the option for the MAX_RPC_SUMMARY_LENGTH rule.
type MaxRpcSummaryLengthOption struct {
	CustomizableSeverityOption `yaml:",inline"`
	MaxChars                   int `yaml:"max_chars" json:"max_chars" toml:"max_chars"`
}
