package rules

import (
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/maramkhaledn/protolint/linter/rule"
)

func TestMessageNamesHaveResponseOrRequestSuffix_Apply(t *testing.T) {
	proto := &parser.Proto{
		ProtoBody: []parser.Visitee{
			&parser.Service{
				ServiceName: "TestService",
				ServiceBody: []parser.Visitee{
					&parser.RPC{
						RPCName: "GetBook",
						RPCRequest: &parser.RPCRequest{MessageType: "GetBookRequest"},
						RPCResponse: &parser.RPCResponse{MessageType: "GetBookResponse"},
					},
					&parser.RPC{
						RPCName: "BadExample",
						RPCRequest: &parser.RPCRequest{MessageType: "BadExampleReq"},
						RPCResponse: &parser.RPCResponse{MessageType: "BadExampleResp"},
					},
				},
			},
		},
	}

	rule := NewMessageNamesHaveResponseOrRequestSuffix(rule.SeverityError)
	failures, err := rule.Apply(proto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(failures) != 2 {
		t.Errorf("expected 2 failures, got %d", len(failures))
	}
	if len(failures) > 0 && failures[0].Message() != "RPC request message name \"BadExampleReq\" should end with 'Request'" {
		t.Errorf("unexpected failure message: %s", failures[0].Message())
	}
	if len(failures) > 1 && failures[1].Message() != "RPC response message name \"BadExampleResp\" should end with 'Response'" {
		t.Errorf("unexpected failure message: %s", failures[1].Message())
	}
}
