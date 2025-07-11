package config

// RPCVersioningOption represents the option for the RPC_VERSIONING rule.
type RPCVersioningOption struct {
    CustomizableSeverityOption `yaml:",inline"`
    RPCVersioning      bool `yaml:"require_version_prefix" json:"require_version_prefix" toml:"require_version_prefix"`
}