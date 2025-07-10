package rules

import (
    "regexp"

    "github.com/yoheimuta/go-protoparser/v4/parser"
    "github.com/yoheimuta/protolint/linter/report"
    "github.com/yoheimuta/protolint/linter/rule"
    "github.com/yoheimuta/protolint/linter/visitor"
)

// Define the regex for versioning
var versioningRegex = regexp.MustCompile(`^/v\d+`)

// RPCVersioningRule verifies that all RPC URLs have a prefix /v{num}.
type RPCVersioningRule struct {
    RuleWithSeverity
}

// NewRPCVersioningRule creates a new RPCVersioningRule.
func NewRPCVersioningRule(severity rule.Severity) RPCVersioningRule {
    return RPCVersioningRule{
        RuleWithSeverity: RuleWithSeverity{severity: severity},
    }
}

// ID returns the ID of this rule.
func (r RPCVersioningRule) ID() string {
    return "RPC_VERSIONING"
}

// Purpose returns the purpose of this rule.
func (r RPCVersioningRule) Purpose() string {
    return "Verifies that all RPC URLs have a prefix /v{num}."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r RPCVersioningRule) IsOfficial() bool {
    return false
}

// Apply applies the rule to the proto.
func (r RPCVersioningRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
    v := &rpcVersioningVisitor{
        BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID(), string(r.Severity())),
    }
    return visitor.RunVisitor(v, proto, r.ID())
}

type rpcVersioningVisitor struct {
    *visitor.BaseAddVisitor
}

// VisitRPC checks the rpc for URL prefix /v{num} in its options.
func (v *rpcVersioningVisitor) VisitRPC(rpc *parser.RPC) bool {
    // Check the options for google.api.http
    for _, option := range rpc.Options {
        if option.OptionName == "(google.api.http)" { // Use the correct field name
            optionURL := extractURLFromOption(option.Constant)
            if optionURL != "" && !versioningRegex.MatchString(optionURL) {
                v.AddFailuref(option.Meta.Pos, `Option URL %q in RPC %q should have a prefix of the form "/v{num}"`, optionURL, rpc.RPCName)
            }
        }
    }
    return false
}

// extractURLFromOption extracts the URL from the option constant.
func extractURLFromOption(constant interface{}) string {
    // Assuming the constant is a string or structured data containing the URL.
    // You may need to parse the constant based on its actual structure.
    if str, ok := constant.(string); ok {
        return str
    }
    return ""
}





