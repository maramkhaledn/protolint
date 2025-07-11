package rules

import (
	"github.com/yoheimuta/go-protoparser/v4/lexer"
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/maramkhaledn/protolint/linter/autodisable"
	"github.com/maramkhaledn/protolint/linter/fixer"
	"github.com/maramkhaledn/protolint/linter/rule"

	"github.com/maramkhaledn/protolint/linter/report"
	"github.com/maramkhaledn/protolint/linter/strs"
	"github.com/maramkhaledn/protolint/linter/visitor"
)

// RPCNamesUpperCamelCaseRule verifies that all rpc names are CamelCase (with an initial capital).
// See https://developers.google.com/protocol-buffers/docs/style#services.
type RPCNamesUpperCamelCaseRule struct {
	RuleWithSeverity
	fixMode         bool
	autoDisableType autodisable.PlacementType
}

// NewRPCNamesUpperCamelCaseRule creates a new RPCNamesUpperCamelCaseRule.
func NewRPCNamesUpperCamelCaseRule(
	severity rule.Severity,
	fixMode bool,
	autoDisableType autodisable.PlacementType,
) RPCNamesUpperCamelCaseRule {
	if autoDisableType != autodisable.Noop {
		fixMode = false
	}
	return RPCNamesUpperCamelCaseRule{
		RuleWithSeverity: RuleWithSeverity{severity: severity},
		fixMode:          fixMode,
		autoDisableType:  autoDisableType,
	}
}

// ID returns the ID of this rule.
func (r RPCNamesUpperCamelCaseRule) ID() string {
	return "RPC_NAMES_UPPER_CAMEL_CASE"
}

// Purpose returns the purpose of this rule.
func (r RPCNamesUpperCamelCaseRule) Purpose() string {
	return "Verifies that all rpc names are CamelCase (with an initial capital)."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r RPCNamesUpperCamelCaseRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r RPCNamesUpperCamelCaseRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base, err := visitor.NewBaseFixableVisitor(r.ID(), r.fixMode, proto, string(r.Severity()))
	if err != nil {
		return nil, err
	}

	v := &rpcNamesUpperCamelCaseVisitor{
		BaseFixableVisitor: base,
	}
	return visitor.RunVisitorAutoDisable(v, proto, r.ID(), r.autoDisableType)
}

type rpcNamesUpperCamelCaseVisitor struct {
	*visitor.BaseFixableVisitor
}

// VisitRPC checks the rpc.
func (v *rpcNamesUpperCamelCaseVisitor) VisitRPC(rpc *parser.RPC) bool {
	name := rpc.RPCName
	if !strs.IsUpperCamelCase(name) {
		expected := strs.ToUpperCamelCase(name)
		v.AddFailuref(rpc.Meta.Pos, "RPC name %q must be UpperCamelCase like %q", name, expected)

		err := v.Fixer.SearchAndReplace(rpc.Meta.Pos, func(lex *lexer.Lexer) fixer.TextEdit {
			lex.NextKeyword()
			lex.Next()
			return fixer.TextEdit{
				Pos:     lex.Pos.Offset,
				End:     lex.Pos.Offset + len(lex.Text) - 1,
				NewText: []byte(expected),
			}
		})
		if err != nil {
			panic(err)
		}
	}
	return false
}
