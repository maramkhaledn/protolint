package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/maramkhaledn/protolint/internal/addon/rules"
	"github.com/maramkhaledn/protolint/linter/report"
	"github.com/maramkhaledn/protolint/linter/rule"
)

func TestOrderRule_Apply(t *testing.T) {
	tests := []struct {
		name         string
		inputProto   *parser.Proto
		wantFailures []report.Failure
	}{
		{
			name: "no failures for proto including all elements in order",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Syntax{},
					&parser.Package{},
					&parser.Import{},
					&parser.Import{},
					&parser.Option{},
					&parser.Option{},
					&parser.Message{},
					&parser.Enum{},
					&parser.Service{},
					&parser.Extend{},
				},
			},
		},
		{
			name: "no failures with editions for proto including all elements in order",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Edition{},
					&parser.Package{},
					&parser.Import{},
					&parser.Import{},
					&parser.Option{},
					&parser.Option{},
					&parser.Message{},
					&parser.Enum{},
					&parser.Service{},
					&parser.Extend{},
				},
			},
		},
		{
			name: "no failures for proto omitting the syntax",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Package{},
					&parser.Import{},
					&parser.Import{},
					&parser.Option{},
					&parser.Option{},
					&parser.Message{},
					&parser.Enum{},
					&parser.Service{},
					&parser.Extend{},
				},
			},
		},
		{
			name: "no failures for proto omitting the package",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Syntax{},
					&parser.Import{},
					&parser.Import{},
					&parser.Option{},
					&parser.Option{},
					&parser.Message{},
					&parser.Enum{},
					&parser.Service{},
					&parser.Extend{},
				},
			},
		},
		{
			name: "no failures for proto omitting the imports",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Syntax{},
					&parser.Package{},
					&parser.Option{},
					&parser.Option{},
					&parser.Message{},
					&parser.Enum{},
					&parser.Service{},
					&parser.Extend{},
				},
			},
		},
		{
			name: "no failures for proto omitting the options",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Syntax{},
					&parser.Package{},
					&parser.Import{},
					&parser.Import{},
					&parser.Message{},
					&parser.Enum{},
					&parser.Service{},
					&parser.Extend{},
				},
			},
		},
		{
			name: "no failures for proto omitting the everything else",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Syntax{},
					&parser.Package{},
					&parser.Import{},
					&parser.Import{},
					&parser.Option{},
					&parser.Option{},
				},
			},
		},
		{
			name: "no failures for proto including only the everything else",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{},
					&parser.Enum{},
					&parser.Service{},
					&parser.Extend{},
				},
			},
		},
		{
			name: "no failures for proto omitting the syntax and the package",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Import{},
					&parser.Import{},
					&parser.Option{},
					&parser.Option{},
					&parser.Message{},
					&parser.Enum{},
					&parser.Service{},
					&parser.Extend{},
				},
			},
		},
		{
			name: "no failures for proto omitting the syntax, the package and the imports",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Option{},
					&parser.Option{},
					&parser.Message{},
					&parser.Enum{},
					&parser.Service{},
					&parser.Extend{},
				},
			},
		},
		{
			name: "no failures for proto omitting the syntax, the package and the options",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Import{},
					&parser.Import{},
					&parser.Message{},
					&parser.Enum{},
					&parser.Service{},
					&parser.Extend{},
				},
			},
		},
		{
			name: "failures for proto in which the order of the syntax is invalid",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Package{},
					&parser.Syntax{
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "example.proto",
								Offset:   100,
								Line:     5,
								Column:   10,
							},
						},
					},
				},
			},
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   100,
						Line:     5,
						Column:   10,
					},
					"ORDER",
					string(rule.SeverityError),
					`Syntax should be located at the top. Check if the file is ordered in the correct manner.`,
				),
			},
		},
		{
			name: "failures for proto in which the order of the package is invalid",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Syntax{},
					&parser.Import{},
					&parser.Package{
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "example.proto",
								Offset:   100,
								Line:     5,
								Column:   10,
							},
						},
					},
				},
			},
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   100,
						Line:     5,
						Column:   10,
					},
					"ORDER",
					string(rule.SeverityError),
					`The order of Package is invalid. Check if the file is ordered in the correct manner.`,
				),
			},
		},
		{
			name: "failures for proto in which the order of the imports is invalid",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Syntax{},
					&parser.Package{},
					&parser.Message{},
					&parser.Import{
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "example.proto",
								Offset:   100,
								Line:     5,
								Column:   10,
							},
						},
					},
					&parser.Import{},
					&parser.Option{},
					&parser.Import{
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "example.proto",
								Offset:   200,
								Line:     10,
								Column:   20,
							},
						},
					},
				},
			},
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   100,
						Line:     5,
						Column:   10,
					},
					"ORDER",
					string(rule.SeverityError),
					`The order of Import is invalid. Check if the file is ordered in the correct manner.`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     10,
						Column:   20,
					},
					"ORDER",
					string(rule.SeverityError),
					`The order of Import is invalid. Check if the file is ordered in the correct manner.`,
				),
			},
		},
		{
			name: "failures for proto in which the order of the options is invalid",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Option{},
					&parser.Extend{},
					&parser.Option{
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "example.proto",
								Offset:   100,
								Line:     5,
								Column:   10,
							},
						},
					},
					&parser.Option{},
					&parser.Service{},
					&parser.Option{
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "example.proto",
								Offset:   200,
								Line:     10,
								Column:   20,
							},
						},
					},
				},
			},
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   100,
						Line:     5,
						Column:   10,
					},
					"ORDER",
					string(rule.SeverityError),
					`The order of Option is invalid. Check if the file is ordered in the correct manner.`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     10,
						Column:   20,
					},
					"ORDER",
					string(rule.SeverityError),
					`The order of Option is invalid. Check if the file is ordered in the correct manner.`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewOrderRule(rule.SeverityError, false)

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

func TestOrderRule_Apply_fix(t *testing.T) {
	tests := []struct {
		name          string
		inputFilename string
		wantFilename  string
	}{
		{
			name:          "no fix for a correct proto",
			inputFilename: "order.proto",
			wantFilename:  "order.proto",
		},
		{
			name:          "no fix for a correct proto with editions",
			inputFilename: "orderEditions.proto",
			wantFilename:  "orderEditions.proto",
		},
		{
			name:          "fix for an incorrect proto",
			inputFilename: "invalid.proto",
			wantFilename:  "order.proto",
		},
		{
			name:          "fix for an incorrect proto with editions",
			inputFilename: "invalidEditions.proto",
			wantFilename:  "orderEditions.proto",
		},
		{
			name:          "fix for an incorrect proto while keeping contiguous misc elements",
			inputFilename: "invalidContiguousMisc.proto",
			wantFilename:  "orderContiguousMisc.proto",
		},
		{
			name:          "fix for an incorrect proto filled with many elements",
			inputFilename: "invalidMany.proto",
			wantFilename:  "orderMany.proto",
		},
		{
			name:          "fix for an incorrect proto filled with many comments",
			inputFilename: "invalidWithComments.proto",
			wantFilename:  "orderWithComments.proto",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			r := rules.NewOrderRule(rule.SeverityError, true)
			testApplyFix(t, r, test.inputFilename, test.wantFilename)
		})
	}
}
