package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/maramkhaledn/protolint/internal/linter/config"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/maramkhaledn/protolint/internal/addon/rules"
	"github.com/maramkhaledn/protolint/linter/report"
	"github.com/maramkhaledn/protolint/linter/rule"
)

func TestQuoteConsistentRule_Apply(t *testing.T) {
	tests := []struct {
		name         string
		inputProto   *parser.Proto
		inputQuote   config.QuoteType
		wantFailures []report.Failure
	}{
		{
			name: "no failures for proto with consistent double-quoted strings",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Syntax{
						ProtobufVersionQuote: `"proto3"`,
					},
					&parser.Import{
						Location: `"google/protobuf/empty.proto"`,
					},
					&parser.Option{
						Constant: `"com.example.foo"`,
					},
					&parser.Enum{
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								EnumValueOptions: []*parser.EnumValueOption{
									{
										Constant: `"custom option"`,
									},
								},
							},
						},
					},
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldOptions: []*parser.FieldOption{
									{
										Constant: `"field option"`,
									},
								},
							},
						},
					},
				},
			},
			inputQuote: config.DoubleQuote,
		},
		{
			name: "no failures for proto with consistent single-quoted strings",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Syntax{
						ProtobufVersionQuote: `'proto3'`,
					},
					&parser.Import{
						Location: `'google/protobuf/empty.proto'`,
					},
					&parser.Option{
						Constant: `'com.example.foo'`,
					},
					&parser.Enum{
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								EnumValueOptions: []*parser.EnumValueOption{
									{
										Constant: `'custom option'`,
									},
								},
							},
						},
					},
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldOptions: []*parser.FieldOption{
									{
										Constant: `'field option'`,
									},
								},
							},
						},
					},
				},
			},
			inputQuote: config.SingleQuote,
		},
		{
			name: "failures for proto with an inconsistent double-quoted strings",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Syntax{
						ProtobufVersionQuote: `"proto3"`,
					},
					&parser.Import{
						Location: `"google/protobuf/empty.proto"`,
					},
					&parser.Option{
						Constant: `"com.example.foo"`,
					},
					&parser.Enum{
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								EnumValueOptions: []*parser.EnumValueOption{
									{
										Constant: `"custom option"`,
									},
								},
							},
						},
					},
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldOptions: []*parser.FieldOption{
									{
										Constant: `"field option"`,
									},
								},
							},
						},
					},
				},
			},
			inputQuote: config.SingleQuote,
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{},
					"QUOTE_CONSISTENT",
					string(rule.SeverityError),
					`Quoted string should be 'proto3' but was "proto3".`,
				),
				report.Failuref(
					meta.Position{},
					"QUOTE_CONSISTENT",
					string(rule.SeverityError),
					`Quoted string should be 'google/protobuf/empty.proto' but was "google/protobuf/empty.proto".`,
				),
				report.Failuref(
					meta.Position{},
					"QUOTE_CONSISTENT",
					string(rule.SeverityError),
					`Quoted string should be 'com.example.foo' but was "com.example.foo".`,
				),
				report.Failuref(
					meta.Position{},
					"QUOTE_CONSISTENT",
					string(rule.SeverityError),
					`Quoted string should be 'custom option' but was "custom option".`,
				),
				report.Failuref(
					meta.Position{},
					"QUOTE_CONSISTENT",
					string(rule.SeverityError),
					`Quoted string should be 'field option' but was "field option".`,
				),
			},
		},
		{
			name: "failures for proto with an inconsistent single-quoted strings",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Syntax{
						ProtobufVersionQuote: `'proto3'`,
					},
					&parser.Import{
						Location: `'google/protobuf/empty.proto'`,
					},
					&parser.Option{
						Constant: `'com.example.foo'`,
					},
					&parser.Enum{
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								EnumValueOptions: []*parser.EnumValueOption{
									{
										Constant: `'custom option'`,
									},
								},
							},
						},
					},
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldOptions: []*parser.FieldOption{
									{
										Constant: `'field option'`,
									},
								},
							},
						},
					},
				},
			},
			inputQuote: config.DoubleQuote,
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{},
					"QUOTE_CONSISTENT",
					string(rule.SeverityError),
					`Quoted string should be "proto3" but was 'proto3'.`,
				),
				report.Failuref(
					meta.Position{},
					"QUOTE_CONSISTENT",
					string(rule.SeverityError),
					`Quoted string should be "google/protobuf/empty.proto" but was 'google/protobuf/empty.proto'.`,
				),
				report.Failuref(
					meta.Position{},
					"QUOTE_CONSISTENT",
					string(rule.SeverityError),
					`Quoted string should be "com.example.foo" but was 'com.example.foo'.`,
				),
				report.Failuref(
					meta.Position{},
					"QUOTE_CONSISTENT",
					string(rule.SeverityError),
					`Quoted string should be "custom option" but was 'custom option'.`,
				),
				report.Failuref(
					meta.Position{},
					"QUOTE_CONSISTENT",
					string(rule.SeverityError),
					`Quoted string should be "field option" but was 'field option'.`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewQuoteConsistentRule(rule.SeverityError, test.inputQuote, false)

			got, err := rule.Apply(test.inputProto)
			if err != nil {
				t.Errorf("got err %v, but want nil", err)
				return
			}
			if !reflect.DeepEqual(got, test.wantFailures) {
				t.Errorf("got %v, but want %v", got, test.wantFailures)
			}
		})
	}
}

func TestQuoteConsistentRule_Apply_fix(t *testing.T) {
	tests := []struct {
		name          string
		inputQuote    config.QuoteType
		inputFilename string
		wantFilename  string
	}{
		{
			name:          "no fix for a double-quoted proto",
			inputQuote:    config.DoubleQuote,
			inputFilename: "double-quoted.proto",
			wantFilename:  "double-quoted.proto",
		},
		{
			name:          "fix for an inconsistent proto with double-quoted consistency",
			inputQuote:    config.DoubleQuote,
			inputFilename: "inconsistent.proto",
			wantFilename:  "double-quoted.proto",
		},
		{
			name:          "fix for an inconsistent proto with single-quoted consistency",
			inputQuote:    config.SingleQuote,
			inputFilename: "inconsistent.proto",
			wantFilename:  "single-quoted.proto",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			r := rules.NewQuoteConsistentRule(
				rule.SeverityError,
				test.inputQuote,
				true,
			)
			testApplyFix(t, r, test.inputFilename, test.wantFilename)
		})
	}
}
