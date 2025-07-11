package plugin

import (
	"github.com/maramkhaledn/protolint/internal/addon/plugin/shared"
	"github.com/maramkhaledn/protolint/linter/rule"

	"github.com/hashicorp/go-plugin"
)

// RegisterCustomRules registers custom rules.
func RegisterCustomRules(
	rules ...rule.Rule,
) {
	plugin.Serve(
		&plugin.ServeConfig{
			HandshakeConfig: shared.Handshake,
			Plugins: map[string]plugin.Plugin{
				"ruleSet": &shared.RuleSetGRPCPlugin{Impl: newRuleSet(rules)},
			},
			GRPCServer: plugin.DefaultGRPCServer,
		},
	)
}
