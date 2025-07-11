package rules_test

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/maramkhaledn/protolint/internal/linter/file"
	"github.com/maramkhaledn/protolint/internal/setting_test"
	"github.com/maramkhaledn/protolint/internal/util_test"
	"github.com/maramkhaledn/protolint/linter/rule"
	"github.com/maramkhaledn/protolint/linter/strs"

	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/maramkhaledn/protolint/internal/addon/rules"
	"github.com/maramkhaledn/protolint/linter/report"
)

func TestFileNamesLowerSnakeCaseRule_Apply(t *testing.T) {
	tests := []struct {
		name          string
		inputProto    *parser.Proto
		inputExcluded []string
		wantFailures  []report.Failure
	}{
		{
			name: "no failures for proto with a valid file name",
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{
					Filename: "../proto/simple.proto",
				},
			},
		},
		{
			name: "no failures for proto with a valid lower snake case file name",
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{
					Filename: "../proto/lower_snake_case.proto",
				},
			},
		},
		{
			name: "no failures for excluded proto",
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{
					Filename: "proto/lowerSnakeCase.proto",
				},
			},
			inputExcluded: []string{
				"proto/lowerSnakeCase.proto",
			},
		},
		{
			name: "no failures for proto with disable directive",
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{
					Filename: "proto/lowerSnakeCase.proto",
				},
				ProtoBody: []parser.Visitee{
					&parser.Comment{
						Raw: "// protolint:disable FILE_NAMES_LOWER_SNAKE_CASE",
					},
				},
			},
		},
		{
			name: "a failure for proto with a camel case file name",
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{
					Filename: "proto/lowerSnakeCase.proto",
				},
			},
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "proto/lowerSnakeCase.proto",
						Offset:   0,
						Line:     1,
						Column:   1,
					},
					"FILE_NAMES_LOWER_SNAKE_CASE",
					string(rule.SeverityError),
					`File name "lowerSnakeCase.proto" should be lower_snake_case.proto like "lower_snake_case.proto".`,
				),
			},
		},
		{
			name: "a failure for proto with an invalid file extension",
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{
					Filename: "proto/lowerSnakeCase.txt",
				},
			},
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "proto/lowerSnakeCase.txt",
						Offset:   0,
						Line:     1,
						Column:   1,
					},
					"FILE_NAMES_LOWER_SNAKE_CASE",
					string(rule.SeverityError),
					`File name "lowerSnakeCase.txt" should be lower_snake_case.proto like "lower_snake_case.proto".`,
				),
			},
		},
		{
			name: "a failure for proto with an invalid separater",
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{
					Filename: "proto/dot.separated.proto",
				},
			},
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "proto/dot.separated.proto",
						Offset:   0,
						Line:     1,
						Column:   1,
					},
					"FILE_NAMES_LOWER_SNAKE_CASE",
					string(rule.SeverityError),
					`File name "dot.separated.proto" should be lower_snake_case.proto like "dot_separated.proto".`,
				),
			},
		},
		{
			name: "a failure for proto with a kebab case file name",
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{
					Filename: "proto/user-role.proto",
				},
			},
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "proto/user-role.proto",
						Offset:   0,
						Line:     1,
						Column:   1,
					},
					"FILE_NAMES_LOWER_SNAKE_CASE",
					string(rule.SeverityError),
					`File name "user-role.proto" should be lower_snake_case.proto like "user_role.proto".`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewFileNamesLowerSnakeCaseRule(rule.SeverityError, test.inputExcluded, false)

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

func TestFileNamesLowerSnakeCaseRule_Apply_fix(t *testing.T) {
	tests := []struct {
		name          string
		inputExcluded []string
		inputFilename string
		wantFilename  string
		wantAbort     bool
	}{
		{
			name:          "no fix for a correct proto",
			inputFilename: "lower_snake_case.proto",
			wantFilename:  "lower_snake_case.proto",
		},
		{
			name:          "abort to fix the proto because of alreadyExists",
			inputFilename: "lowerSnakeCase.proto",
			wantAbort:     true,
		},
		{
			name:          "fix for an incorrect proto",
			inputFilename: "UpperCamelCase.proto",
			wantFilename:  "upper_camel_case.proto",
		},
		{
			name:          "fix for a kebab case proto",
			inputFilename: "kebab-case.proto",
			wantFilename:  "kebab_case.proto",
		},
		{
			name:          "no fix for a proto with disable directive",
			inputFilename: "dot.separated.proto",
			wantFilename:  "dot.separated.proto",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			r := rules.NewFileNamesLowerSnakeCaseRule(rule.SeverityError, test.inputExcluded, true)

			dataDir := strs.ToLowerCamelCase(r.ID())
			input, err := util_test.NewTestData(setting_test.TestDataPath("rules", dataDir, test.inputFilename))
			if err != nil {
				t.Errorf("got err %v", err)
				return
			}

			proto, err := file.NewProtoFile(input.FilePath, input.FilePath).Parse(false)
			if err != nil {
				t.Errorf("%v", err.Error())
				return
			}
			fs, err := r.Apply(proto)
			if err != nil {
				t.Errorf("got err %v, but want nil", err)
				return
			}
			if test.wantAbort {
				if _, err := os.Stat(input.FilePath); os.IsNotExist(err) {
					t.Errorf("not found %q, but want to locate it", input.FilePath)
					return
				}
				for _, f := range fs {
					if strings.Contains(f.Message(), "Failed to rename") {
						return
					}
				}
				t.Error("not found failure message, but want to include it")
				return
			}

			wantPath := setting_test.TestDataPath("rules", dataDir, test.wantFilename)
			if _, err := os.Stat(wantPath); os.IsNotExist(err) {
				t.Errorf("not found %q, but want to locate it", wantPath)
				return
			}

			err = os.Rename(wantPath, input.FilePath)
			if err != nil {
				t.Errorf("got err %v", err)
			}
		})
	}
}
