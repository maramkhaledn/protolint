package rules

import (
	"regexp"

	"github.com/maramkhaledn/protolint/linter/report"
	"github.com/maramkhaledn/protolint/linter/rule"
	"github.com/maramkhaledn/protolint/linter/visitor"
	"github.com/yoheimuta/go-protoparser/v4/parser"
)

const (
	defaultMaxRpcSummaryChars = 30
)

type MaxRpcSummaryLengthRule struct {
	RuleWithSeverity
	maxChars int
}

// MaxRpcSummaryLengthRule creates a new MaxRpcSummaryLengthRule.
func NewMaxRpcSummaryLengthRule(
	severity rule.Severity,
	maxChars int,
) MaxRpcSummaryLengthRule {
	if maxChars == 0 {
		maxChars = defaultMaxRpcSummaryChars
	}
	return MaxRpcSummaryLengthRule{
		RuleWithSeverity: RuleWithSeverity{severity: severity},
		maxChars:         maxChars,
	}
}

// ID returns the ID of this rule.
func (r MaxRpcSummaryLengthRule) ID() string {
	return "MAX_RPC_SUMMARY_LENGTH"
}

// Purpose returns the purpose of this rule.
func (r MaxRpcSummaryLengthRule) Purpose() string {
	return "Verifies that all RPC Summaries are less than or equal the maximum allowed length."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r MaxRpcSummaryLengthRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r MaxRpcSummaryLengthRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &MaxRpcSummaryLengthVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID(), string(r.Severity())),
		maxChars:       r.maxChars,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type MaxRpcSummaryLengthVisitor struct {
	*visitor.BaseAddVisitor
	maxChars int
}

func (v *MaxRpcSummaryLengthVisitor) VisitRPC(rpc *parser.RPC) bool {
	for _, option := range rpc.Options {
		if option.OptionName == "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation)" {
			optionSummary := extractRpcSummaryFromOption(option.Constant)
			if len(optionSummary) > v.maxChars {
				v.AddFailuref(option.Meta.Pos, `Option Summary %q in RPC %q should be less than or equal the maximum length (%d)`, optionSummary, rpc.RPCName, v.maxChars)
			}
		}
	}
	return false
}

func extractRpcSummaryFromOption(constant interface{}) string {
	if s, ok := constant.(string); ok {
		re := regexp.MustCompile(`summary\s*:\s*"([^"]*)"`)
		matches := re.FindStringSubmatch(s)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}
