package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/maramkhaledn/protolint/internal/setting_test"
	"github.com/maramkhaledn/protolint/internal/util_test"

	"github.com/maramkhaledn/protolint/internal/linter/file"

	"github.com/maramkhaledn/protolint/internal/addon/rules"
	"github.com/maramkhaledn/protolint/linter/report"
	"github.com/maramkhaledn/protolint/linter/rule"
)

func testImportSortedProtoPath(name string) string {
	return setting_test.TestDataPath("rules", "importsSorted", name)
}

func TestImportsSortedRule_Apply(t *testing.T) {
	tests := []struct {
		name          string
		inputFilename string
		wantFailures  []report.Failure
		wantExistErr  bool
	}{
		{
			name:          "no failures for proto with sorted imports",
			inputFilename: "sorted.proto",
		},
		{
			name:          "no failures for proto with sorted imports separated by a newline",
			inputFilename: "sortedWithNewline.proto",
		},
		{
			name:          "failures for proto with not sorted imports",
			inputFilename: "notSorted.proto",
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: testImportSortedProtoPath("notSorted.proto"),
						Offset:   20,
						Line:     3,
						Column:   1,
					},
					"IMPORTS_SORTED",
					string(rule.SeverityError),
					`Imports are not sorted.`,
				),
				report.Failuref(
					meta.Position{
						Filename: testImportSortedProtoPath("notSorted.proto"),
						Offset:   47,
						Line:     4,
						Column:   1,
					},
					"IMPORTS_SORTED",
					string(rule.SeverityError),
					`Imports are not sorted.`,
				),
			},
		},
		{
			name:          "failures for proto with not sorted imports separated by a newline",
			inputFilename: "notSortedWithNewline.proto",
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: testImportSortedProtoPath("notSortedWithNewline.proto"),
						Offset:   20,
						Line:     3,
						Column:   1,
					},
					"IMPORTS_SORTED",
					string(rule.SeverityError),
					`Imports are not sorted.`,
				),
				report.Failuref(
					meta.Position{
						Filename: testImportSortedProtoPath("notSortedWithNewline.proto"),
						Offset:   42,
						Line:     4,
						Column:   1,
					},
					"IMPORTS_SORTED",
					string(rule.SeverityError),
					`Imports are not sorted.`,
				),
				report.Failuref(
					meta.Position{
						Filename: testImportSortedProtoPath("notSortedWithNewline.proto"),
						Offset:   151,
						Line:     9,
						Column:   1,
					},
					"IMPORTS_SORTED",
					string(rule.SeverityError),
					`Imports are not sorted.`,
				),
				report.Failuref(
					meta.Position{
						Filename: testImportSortedProtoPath("notSortedWithNewline.proto"),
						Offset:   190,
						Line:     10,
						Column:   1,
					},
					"IMPORTS_SORTED",
					string(rule.SeverityError),
					`Imports are not sorted.`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewImportsSortedRule(
				rule.SeverityError,
				false,
			)

			protoPath := testImportSortedProtoPath(test.inputFilename)
			proto, err := file.NewProtoFile(protoPath, protoPath).Parse(false)
			if err != nil {
				t.Errorf("%v", err)
				return
			}

			got, err := rule.Apply(proto)
			if test.wantExistErr {
				if err == nil {
					t.Errorf("got err nil, but want err")
				}
				return
			}
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

func newTestImportsSortedData(
	fileName string,
) (util_test.TestData, error) {
	return util_test.NewTestData(testImportSortedProtoPath(fileName))
}

func TestImportsSortedRule_Apply_fix(t *testing.T) {
	tests := []struct {
		name          string
		inputFilename string
		wantFilename  string
	}{
		{
			name:          "no fix for proto with sorted imports",
			inputFilename: "sorted.proto",
			wantFilename:  "sorted.proto",
		},
		{
			name:          "no fix for proto with sorted imports separated by a newline",
			inputFilename: "sortedWithNewline.proto",
			wantFilename:  "sortedWithNewline.proto",
		},
		{
			name:          "fix for proto with not sorted imports",
			inputFilename: "notSorted.proto",
			wantFilename:  "sorted.proto",
		},
		{
			name:          "fix for proto with sorted imports separated by a newline",
			inputFilename: "notSortedWithNewline.proto",
			wantFilename:  "sortedWithNewline.proto",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewImportsSortedRule(
				rule.SeverityError,
				true,
			)

			input, err := newTestImportsSortedData(test.inputFilename)
			if err != nil {
				t.Errorf("got err %v", err)
				return
			}

			want, err := newTestImportsSortedData(test.wantFilename)
			if err != nil {
				t.Errorf("got err %v", err)
				return
			}

			proto, err := file.NewProtoFile(input.FilePath, input.FilePath).Parse(false)
			if err != nil {
				t.Errorf("%v", err)
				return
			}

			_, err = rule.Apply(proto)
			if err != nil {
				t.Errorf("got err %v, but want nil", err)
				return
			}

			got, err := input.Data()
			if err != nil {
				t.Errorf("got err %v, but want nil", err)
				return
			}
			if !reflect.DeepEqual(got, want.OriginData) {
				t.Errorf(
					"got %s(%v), but want %s(%v)",
					string(got), got,
					string(want.OriginData), want.OriginData,
				)
			}

			err = input.Restore()
			if err != nil {
				t.Errorf("got err %v", err)
			}
		})
	}
}
