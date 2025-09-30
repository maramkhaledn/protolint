package rules

import (

    "github.com/yoheimuta/go-protoparser/v4/parser"
    "github.com/maramkhaledn/protolint/linter/report"
    "github.com/maramkhaledn/protolint/linter/rule"
    "github.com/maramkhaledn/protolint/linter/visitor"
	"unicode"
	"regexp"
)


// RPCSummaryPascalCase verifies that all RPCs have a non-empty pascal case summary.
type RPCSummaryPascalCaseRule struct {
    RuleWithSeverity
}

func NewRPCSummaryPascalCaseRule(severity rule.Severity) RPCSummaryPascalCaseRule {
    return RPCSummaryPascalCaseRule{
        RuleWithSeverity: RuleWithSeverity{severity: severity},
    }
}

// ID returns the ID of this rule.
func (r RPCSummaryPascalCaseRule) ID() string {
    return "RPC_SUMMARY_PASCAL_CASE"
}

// Purpose returns the purpose of this rule.
func (r RPCSummaryPascalCaseRule) Purpose() string {
    return "Verifies that all RPC Summaries are non-empty and in pascal case"
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r RPCSummaryPascalCaseRule) IsOfficial() bool {
    return true
}

// Apply applies the rule to the proto.
func (r RPCSummaryPascalCaseRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
    v := &rpcSummaryPascalCaseVisitor{
        BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID(), string(r.Severity())),
    }
    return visitor.RunVisitor(v, proto, r.ID())
}

type rpcSummaryPascalCaseVisitor struct {
    *visitor.BaseAddVisitor
}

func (v *rpcSummaryPascalCaseVisitor) VisitRPC(rpc *parser.RPC) bool {
    for _, option := range rpc.Options {
        if option.OptionName == "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation)" {
            optionSummary := extractSummaryFromOption(option.Constant)
            if isPascalCase(optionSummary) == false {
                v.AddFailuref(option.Meta.Pos, `Option Summary %q in RPC %q should be non empty and in pascal case`, optionSummary, rpc.RPCName)
            }
        }
    }
    return false
}

func extractSummaryFromOption(constant interface{}) string {
    
    if s, ok := constant.(string); ok {

        re := regexp.MustCompile(`summary\s*:\s*"([^"]*)"`)
        matches := re.FindStringSubmatch(s)
        if len(matches) > 1 {
            return matches[1]
        }

    }
    return ""
}

func isPascalCase(s string) bool {
	if len(s) == 0 {
		return false
	}
	capitalFound := 0
	for i, r := range s {
		if i == 0 && unicode.IsLower(r) {
			return false
		}
		if unicode.IsUpper(r) {
			capitalFound++
		}
		if !unicode.IsLetter(r) {
			return false 
		}
	}
	if capitalFound < 2{
		return false
	}
	
	return true
}


