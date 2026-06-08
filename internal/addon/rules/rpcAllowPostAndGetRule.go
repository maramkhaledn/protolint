package rules

import (
    "regexp"

    "github.com/yoheimuta/go-protoparser/v4/parser"
    "github.com/maramkhaledn/protolint/linter/report"
    "github.com/maramkhaledn/protolint/linter/rule"
    "github.com/maramkhaledn/protolint/linter/visitor"
)

// httpMethodRegex extracts the HTTP method key from a stringified (google.api.http) option,
// e.g. `post:"/v1/foo" body:"*"` -> "post".
var httpMethodRegex = regexp.MustCompile(`(?m)\b(get|post|put|patch|delete|custom)\s*:`)

// allowedHTTPMethods is the set of HTTP methods permitted on an RPC's (google.api.http) option.
var allowedHTTPMethods = map[string]bool{
    "get":  true,
    "post": true,
}

// RPCAllowPostAndGetRule verifies that every RPC's (google.api.http) option uses only
// the GET or POST methods (never PUT, PATCH, DELETE, or custom).
type RPCAllowPostAndGetRule struct {
    RuleWithSeverity
}

// NewRPCAllowPostAndGetRule creates a new RPCAllowPostAndGetRule.
func NewRPCAllowPostAndGetRule(severity rule.Severity) RPCAllowPostAndGetRule {
    return RPCAllowPostAndGetRule{
        RuleWithSeverity: RuleWithSeverity{severity: severity},
    }
}

// ID returns the ID of this rule.
func (r RPCAllowPostAndGetRule) ID() string {
    return "RPC_ALLOW_POST_AND_GET"
}

// Purpose returns the purpose of this rule.
func (r RPCAllowPostAndGetRule) Purpose() string {
    return "Verifies that all RPCs only use the GET or POST methods in their (google.api.http) option."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r RPCAllowPostAndGetRule) IsOfficial() bool {
    return true
}

// Apply applies the rule to the proto.
func (r RPCAllowPostAndGetRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
    v := &rpcAllowPostAndGetVisitor{
        BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID(), string(r.Severity())),
    }
    return visitor.RunVisitor(v, proto, r.ID())
}

type rpcAllowPostAndGetVisitor struct {
    *visitor.BaseAddVisitor
}

// VisitRPC checks that the (google.api.http) option only uses GET or POST methods.
func (v *rpcAllowPostAndGetVisitor) VisitRPC(rpc *parser.RPC) bool {
    for _, option := range rpc.Options {
        if option.OptionName != "(google.api.http)" {
            continue
        }
        for _, method := range extractHTTPMethods(option.Constant) {
            if !allowedHTTPMethods[method] {
                v.AddFailuref(
                    option.Meta.Pos,
                    `RPC %q uses HTTP method %q in its (google.api.http) option, but only "get" and "post" are allowed`,
                    rpc.RPCName, method,
                )
            }
        }
    }
    return false
}

// extractHTTPMethods returns the HTTP method keys declared in a (google.api.http) option constant
// (e.g. `post:"/v1/foo" body:"*"` -> ["post"]).
func extractHTTPMethods(constant string) []string {
    var methods []string
    for _, match := range httpMethodRegex.FindAllStringSubmatch(constant, -1) {
        methods = append(methods, match[1])
    }
    return methods
}
