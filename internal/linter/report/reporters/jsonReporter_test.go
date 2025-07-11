package reporters_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/maramkhaledn/protolint/internal/linter/report/reporters"
	"github.com/maramkhaledn/protolint/linter/report"
	"github.com/maramkhaledn/protolint/linter/rule"
)

func TestJSONReporter_Report(t *testing.T) {
	tests := []struct {
		name          string
		inputFailures []report.Failure
		wantOutput    func(basedir string) string
	}{
		{
			name: "Prints failures in JSON format",
			inputFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   100,
						Line:     5,
						Column:   10,
					},
					"ENUM_NAMES_UPPER_CAMEL_CASE",
					string(rule.SeverityError),
					`EnumField name "fIRST_VALUE" must be CAPITALS_WITH_UNDERSCORES`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     10,
						Column:   20,
					},
					"ENUM_NAMES_UPPER_CAMEL_CASE",
					string(rule.SeverityError),
					`EnumField name "SECOND.VALUE" must be CAPITALS_WITH_UNDERSCORES`,
				),
			},
			wantOutput: func(basedir string) string {
				return `{
  "basedir": "` + basedir + `",
  "lints": [
    {
      "filename": "example.proto",
      "line": 5,
      "column": 10,
      "message": "EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES",
      "rule": "ENUM_NAMES_UPPER_CAMEL_CASE",
      "severity": "error"
    },
    {
      "filename": "example.proto",
      "line": 10,
      "column": 20,
      "message": "EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES",
      "rule": "ENUM_NAMES_UPPER_CAMEL_CASE",
      "severity": "error"
    }
  ]
}
`
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := reporters.JSONReporter{}.Report(buf, test.inputFailures)
			if err != nil {
				t.Errorf("got err %v, but want nil", err)
				return
			}

			cwd, err := os.Getwd()
			if err != nil {
				t.Errorf("Failed to get current working directory: %v", err)
				return
			}

			wantedOutput := test.wantOutput(cwd)
			if buf.String() != wantedOutput {
				t.Errorf("got %s, but want %s", buf.String(), wantedOutput)
			}
		})
	}
}
