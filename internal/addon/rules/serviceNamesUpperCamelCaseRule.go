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

// ServiceNamesUpperCamelCaseRule verifies that all service names are CamelCase (with an initial capital).
// See https://developers.google.com/protocol-buffers/docs/style#services.
type ServiceNamesUpperCamelCaseRule struct {
	RuleWithSeverity
	fixMode         bool
	autoDisableType autodisable.PlacementType
}

// NewServiceNamesUpperCamelCaseRule creates a new ServiceNamesUpperCamelCaseRule.
func NewServiceNamesUpperCamelCaseRule(
	severity rule.Severity,
	fixMode bool,
	autoDisableType autodisable.PlacementType,
) ServiceNamesUpperCamelCaseRule {
	if autoDisableType != autodisable.Noop {
		fixMode = false
	}
	return ServiceNamesUpperCamelCaseRule{
		RuleWithSeverity: RuleWithSeverity{severity: severity},
		fixMode:          fixMode,
		autoDisableType:  autoDisableType,
	}
}

// ID returns the ID of this rule.
func (r ServiceNamesUpperCamelCaseRule) ID() string {
	return "SERVICE_NAMES_UPPER_CAMEL_CASE"
}

// Purpose returns the purpose of this rule.
func (r ServiceNamesUpperCamelCaseRule) Purpose() string {
	return "Verifies that all service names are CamelCase (with an initial capital)."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r ServiceNamesUpperCamelCaseRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r ServiceNamesUpperCamelCaseRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base, err := visitor.NewBaseFixableVisitor(r.ID(), r.fixMode, proto, string(r.Severity()))
	if err != nil {
		return nil, err
	}

	v := &serviceNamesUpperCamelCaseVisitor{
		BaseFixableVisitor: base,
	}
	return visitor.RunVisitorAutoDisable(v, proto, r.ID(), r.autoDisableType)
}

type serviceNamesUpperCamelCaseVisitor struct {
	*visitor.BaseFixableVisitor
}

// VisitService checks the service.
func (v *serviceNamesUpperCamelCaseVisitor) VisitService(service *parser.Service) bool {
	name := service.ServiceName
	if !strs.IsUpperCamelCase(name) {
		expected := strs.ToUpperCamelCase(name)
		v.AddFailuref(service.Meta.Pos, "Service name %q must be UpperCamelCase like %q", name, expected)

		err := v.Fixer.SearchAndReplace(service.Meta.Pos, func(lex *lexer.Lexer) fixer.TextEdit {
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
