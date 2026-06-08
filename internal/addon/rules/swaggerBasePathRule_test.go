package rules_test

import (
    "reflect"
    "testing"

    "github.com/maramkhaledn/protolint/internal/addon/rules"
    "github.com/maramkhaledn/protolint/linter/rule"
    "github.com/yoheimuta/go-protoparser/v4/parser"
    "github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

func TestSwaggerBasePathRule_Apply(t *testing.T) {
    tests := []struct {
        name         string
        inputProto   *parser.Proto
        wantFailures []string
    }{
        {
            name: "no failures for dash-separated base_path",
            inputProto: &parser.Proto{
                ProtoBody: []parser.Visitee{
                    &parser.Option{
                        OptionName: "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger)",
                        Constant:   `{ base_path: "/my-api-v1" }`,
                    },
                },
            },
            wantFailures: nil,
        },
        {
            name: "no failures when base_path is absent",
            inputProto: &parser.Proto{
                ProtoBody: []parser.Visitee{
                    &parser.Option{
                        OptionName: "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger)",
                        Constant:   `{ info: { title: "My API" } }`,
                    },
                },
            },
            wantFailures: nil,
        },
        {
            name: "ignores other options",
            inputProto: &parser.Proto{
                ProtoBody: []parser.Visitee{
                    &parser.Option{
                        OptionName: "(some.other.option)",
                        Constant:   `{ base_path: "/my_api" }`,
                    },
                },
            },
            wantFailures: nil,
        },
        {
            name: "failure for underscore in base_path",
            inputProto: &parser.Proto{
                ProtoBody: []parser.Visitee{
                    &parser.Option{
                        OptionName: "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger)",
                        Constant:   `{ base_path: "/my_api/v1" }`,
                        Meta: meta.Meta{
                            Pos: meta.Position{
                                Filename: "example.proto",
                                Line:     5,
                                Column:   1,
                            },
                        },
                    },
                },
            },
            wantFailures: []string{
                `base_path "/my_api/v1" in the openapiv2_swagger option must be dash-separated and never contain underscores`,
            },
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            r := rules.NewSwaggerBasePathRule(rule.SeverityError)

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
