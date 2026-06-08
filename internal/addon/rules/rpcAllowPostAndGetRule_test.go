package rules_test

import (
    "reflect"
    "testing"

    "github.com/yoheimuta/go-protoparser/v4/parser"
    "github.com/yoheimuta/go-protoparser/v4/parser/meta"
    "github.com/maramkhaledn/protolint/internal/addon/rules"
    "github.com/maramkhaledn/protolint/linter/rule"
)

func TestRPCAllowPostAndGetRule_Apply(t *testing.T) {
    tests := []struct {
        name         string
        inputProto   *parser.Proto
        wantFailures []string
    }{
        {
            name: "no failures for GET method",
            inputProto: &parser.Proto{
                ProtoBody: []parser.Visitee{
                    &parser.Service{
                        ServiceBody: []parser.Visitee{
                            &parser.RPC{
                                RPCName: "GetUser",
                                Options: []*parser.Option{
                                    {
                                        OptionName: "(google.api.http)",
                                        Constant:   `get:"/v1/users"`,
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
            name: "no failures for POST method",
            inputProto: &parser.Proto{
                ProtoBody: []parser.Visitee{
                    &parser.Service{
                        ServiceBody: []parser.Visitee{
                            &parser.RPC{
                                RPCName: "CreateUser",
                                Options: []*parser.Option{
                                    {
                                        OptionName: "(google.api.http)",
                                        Constant:   `post:"/v1/users" body:"*"`,
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
            name: "failure for DELETE method",
            inputProto: &parser.Proto{
                ProtoBody: []parser.Visitee{
                    &parser.Service{
                        ServiceBody: []parser.Visitee{
                            &parser.RPC{
                                RPCName: "DeleteUser",
                                Options: []*parser.Option{
                                    {
                                        OptionName: "(google.api.http)",
                                        Constant:   `delete:"/v1/users/{id}"`,
                                    },
                                },
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
            wantFailures: []string{
                `RPC "DeleteUser" uses HTTP method "delete" in its (google.api.http) option, but only "get" and "post" are allowed`,
            },
        },
        {
            name: "failure for PUT method",
            inputProto: &parser.Proto{
                ProtoBody: []parser.Visitee{
                    &parser.Service{
                        ServiceBody: []parser.Visitee{
                            &parser.RPC{
                                RPCName: "ReplaceUser",
                                Options: []*parser.Option{
                                    {
                                        OptionName: "(google.api.http)",
                                        Constant:   `put:"/v1/users/{id}" body:"*"`,
                                    },
                                },
                            },
                        },
                    },
                },
            },
            wantFailures: []string{
                `RPC "ReplaceUser" uses HTTP method "put" in its (google.api.http) option, but only "get" and "post" are allowed`,
            },
        },
        {
            name: "failure for PATCH method",
            inputProto: &parser.Proto{
                ProtoBody: []parser.Visitee{
                    &parser.Service{
                        ServiceBody: []parser.Visitee{
                            &parser.RPC{
                                RPCName: "UpdateUser",
                                Options: []*parser.Option{
                                    {
                                        OptionName: "(google.api.http)",
                                        Constant:   `patch:"/v1/users/{id}" body:"*"`,
                                    },
                                },
                            },
                        },
                    },
                },
            },
            wantFailures: []string{
                `RPC "UpdateUser" uses HTTP method "patch" in its (google.api.http) option, but only "get" and "post" are allowed`,
            },
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            r := rules.NewRPCAllowPostAndGetRule(rule.SeverityError)

            failures, err := r.Apply(test.inputProto)
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
