package rules_test

import (
    "reflect"
    "testing"

    "github.com/yoheimuta/go-protoparser/v4/parser"
    "github.com/yoheimuta/go-protoparser/v4/parser/meta"
    "github.com/yoheimuta/protolint/internal/addon/rules"
    "github.com/yoheimuta/protolint/linter/rule"
)

func TestRPCVersioningRule_Apply(t *testing.T) {
    tests := []struct {
        name     string
        inputProto *parser.Proto
        wantFailures []string // Expected failure messages
    }{
        {
            name: "no failures for valid RPC URLs with /v{num} prefix",
            inputProto: &parser.Proto{
                ProtoBody: []parser.Visitee{
                    &parser.Service{
                        ServiceBody: []parser.Visitee{
                            &parser.RPC{
                                RPCName: "GetUser",
                                Options: []*parser.Option{
                                    {
                                        OptionName: "(google.api.http)",
                                        Constant:   "/v1/users",
                                    },
                                },
                            },
                        },
                    },
                },
            },
            wantFailures: nil,
        },
        {
            name: "failures for RPC URLs without /v{num} prefix",
            inputProto: &parser.Proto{
                ProtoBody: []parser.Visitee{
                    &parser.Service{
                        ServiceBody: []parser.Visitee{
                            &parser.RPC{
                                RPCName: "GetUser",
                                Options: []*parser.Option{
                                    {
                                        OptionName: "(google.api.http)",
                                        Constant:   "/users",
                                    },
                                },
                                Meta: meta.Meta{
                                    Pos: meta.Position{
                                        Filename: "example.proto",
                                        Line:     10,
                                        Column:   5,
                                    },
                                },
                            },
                        },
                    },
                },
            },
            wantFailures: []string{
                `Option URL "/users" in RPC "GetUser" should have a prefix of the form "/v{num}"`,
            },
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            rule := rules.NewRPCVersioningRule(rule.SeverityError)

            failures, err := rule.Apply(test.inputProto)
            if err != nil {
                t.Errorf("unexpected error: %v", err)
                return
            }

            var gotFailures []string
            for _, failure := range failures {
                gotFailures = append(gotFailures, failure.Message())
            }

            if !reflect.DeepEqual(gotFailures, test.wantFailures) {
                t.Errorf("got %v, but want %v", gotFailures, test.wantFailures)
            }
        })
    }
}