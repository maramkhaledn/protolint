package rules

import (
	"strings"
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/maramkhaledn/protolint/linter/report"
	"github.com/maramkhaledn/protolint/linter/rule"
	"github.com/maramkhaledn/protolint/linter/visitor"
)

// MessageNamesHaveResponseOrRequestSuffix verifies that all message names in RPC definitions end with Response or Request.
type MessageNamesHaveResponseOrRequestSuffix struct {
	RuleWithSeverity
}

// NewMessageNamesHaveResponseOrRequestSuffix creates a new MessageNamesHaveResponseOrRequestSuffix rule.
func NewMessageNamesHaveResponseOrRequestSuffix(severity rule.Severity) MessageNamesHaveResponseOrRequestSuffix {
	return MessageNamesHaveResponseOrRequestSuffix{
		RuleWithSeverity: RuleWithSeverity{severity: severity},
	}
}

// ID returns the ID of this rule.
func (r MessageNamesHaveResponseOrRequestSuffix) ID() string {
	return "MESSAGE_NAMES_HAVE_RESPONSE_OR_REQUEST_SUFFIX"
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r MessageNamesHaveResponseOrRequestSuffix) IsOfficial() bool {
	return true
}

// Purpose returns the purpose of this rule.
func (r MessageNamesHaveResponseOrRequestSuffix) Purpose() string {
	return `Verifies that all message names in RPC definitions end with "Response" or "Request".`
}

// Apply applies the rule to the proto.
func (r MessageNamesHaveResponseOrRequestSuffix) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &messageNamesHaveResponseOrRequestSuffixVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID(), string(r.Severity())),
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type messageNamesHaveResponseOrRequestSuffixVisitor struct {
	*visitor.BaseAddVisitor
}

// VisitRPC checks the request and response message names in RPC definitions.
func (v *messageNamesHaveResponseOrRequestSuffixVisitor) VisitRPC(rpc *parser.RPC) bool {
	if rpc == nil {
		return true
	}
	requestType := ""
	if rpc.RPCRequest != nil {
		requestType = rpc.RPCRequest.MessageType
	}

	responseType := ""
	if rpc.RPCResponse != nil {
		responseType = rpc.RPCResponse.MessageType
	}

	// Check request type
	if !hasRequestSuffix(requestType) && !isIgnoredType(requestType) {
		v.AddFailuref(rpc.Meta.Pos, "RPC request message name %q should end with 'Request'", requestType)
	}
	// Check response type
	if !hasResponseSuffix(responseType) && !isIgnoredType(responseType) {
		v.AddFailuref(rpc.Meta.Pos, "RPC response message name %q should end with 'Response'", responseType)
	}
	return true
}

var ignoredTypes = map[string]struct{}{
    "google.protobuf.Empty": {},
    "google.protobuf.Any":   {},
}

func isIgnoredType(name string) bool {
    _, ok := ignoredTypes[name]
    return ok
}

func hasResponseSuffix(name string) bool {
	return strings.HasSuffix(name, "Response")
}
func hasRequestSuffix(name string) bool {
	return strings.HasSuffix(name, "Request")
}


