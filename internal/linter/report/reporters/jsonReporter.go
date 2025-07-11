package reporters

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/maramkhaledn/protolint/linter/report"
)

// JSONReporter prints failures as a single JSON struct, allowing
// for simple machine-readable output.
//
// The format is:
//
//	 {
//			"lints":
//				[
//					{"filename": FILENAME, "line": LINE, "column": COL, "message": MESSAGE, "rule": RULE}
//				],
//	 }
type JSONReporter struct{}

type lintJSON struct {
	Filename string `json:"filename"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Message  string `json:"message"`
	Rule     string `json:"rule"`
	Severity string `json:"severity"`
}

type outJSON struct {
	Basedir string     `json:"basedir"`
	Lints   []lintJSON `json:"lints"`
}

// Report writes failures to w.
func (r JSONReporter) Report(w io.Writer, fs []report.Failure) error {
	out := outJSON{}
	// Write base dir
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	out.Basedir = pwd
	// Write failures
	for _, failure := range fs {
		out.Lints = append(out.Lints, lintJSON{
			Filename: failure.Pos().Filename,
			Line:     failure.Pos().Line,
			Column:   failure.Pos().Column,
			Message:  failure.Message(),
			Rule:     failure.RuleID(),
			Severity: failure.Severity(),
		})
	}

	bs, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, string(bs))
	if err != nil {
		return err
	}

	return nil
}
