package rules_test

import (
    "reflect"
    "testing"

    "github.com/yoheimuta/go-protoparser/v4/parser"
    "github.com/yoheimuta/go-protoparser/v4/parser/meta"
    "github.com/maramkhaledn/protolint/internal/addon/rules"
    "github.com/maramkhaledn/protolint/linter/rule"
)

func TestRPCSummaryPascalCaseRule_Apply(t *testing.T) {
    tests := []struct {
        name         string
        inputProto   *parser.Proto
        wantFailures []string
    }{
        {
            name: "no failures for valid PascalCase summary",
            inputProto: &parser.Proto{
                ProtoBody: []parser.Visitee{
                    &parser.Service{
                        ServiceBody: []parser.Visitee{
                            &parser.RPC{
                                RPCName: "CreateUser",
                                Options: []*parser.Option{
                                    {
                                        OptionName: "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation)",
                                        Constant:   "summary: \"CreateUser\"",
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
            name: "failure for empty summary",
            inputProto: &parser.Proto{
                ProtoBody: []parser.Visitee{
                    &parser.Service{
                        ServiceBody: []parser.Visitee{
                            &parser.RPC{
                                RPCName: "CreateUser",
                                Options: []*parser.Option{
                                    {
                                        OptionName: "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation)",
                                        Constant:   "summary: \"\"",
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
                },
            },
            wantFailures: []string{
                `Option Summary "" in RPC "CreateUser" should be non empty and in pascal case`,
            },
        },
        {
            name: "failure for non-PascalCase summary",
            inputProto: &parser.Proto{
                ProtoBody: []parser.Visitee{
                    &parser.Service{
                        ServiceBody: []parser.Visitee{
                            &parser.RPC{
                                RPCName: "DeleteUser",
                                Options: []*parser.Option{
                                    {
                                        OptionName: "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation)",
                                        Constant:   "summary: \"deleteUser\"",
                                        Meta: meta.Meta{
                                            Pos: meta.Position{
                                                Filename: "example.proto",
                                                Line:     15,
                                                Column:   5,
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                },
            },
            wantFailures: []string{
                `Option Summary "deleteUser" in RPC "DeleteUser" should be non empty and in pascal case`,
            },
        },
        {
            name: "failure for summary with underscore",
            inputProto: &parser.Proto{
                ProtoBody: []parser.Visitee{
                    &parser.Service{
                        ServiceBody: []parser.Visitee{
                            &parser.RPC{
                                RPCName: "UpdateUser",
                                Options: []*parser.Option{
                                    {
                                        OptionName: "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation)",
                                        Constant:   "summary: \"Update_User\"",
                                        Meta: meta.Meta{
                                            Pos: meta.Position{
                                                Filename: "example.proto",
                                                Line:     20,
                                                Column:   5,
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                },
            },
            wantFailures: []string{
                `Option Summary "Update_User" in RPC "UpdateUser" should be non empty and in pascal case`,
            },
        },
        {
            name: "failure for summary with space",
            inputProto: &parser.Proto{
                ProtoBody: []parser.Visitee{
                    &parser.Service{
                        ServiceBody: []parser.Visitee{
                            &parser.RPC{
                                RPCName: "GetUser",
                                Options: []*parser.Option{
                                    {
                                        OptionName: "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation)",
                                        Constant:   "summary: \"Get User\"",
                                        Meta: meta.Meta{
                                            Pos: meta.Position{
                                                Filename: "example.proto",
                                                Line:     25,
                                                Column:   5,
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                },
            },
            wantFailures: []string{
                `Option Summary "Get User" in RPC "GetUser" should be non empty and in pascal case`,
            },
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            rule := rules.NewRPCSummaryPascalCaseRule(rule.SeverityError)

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