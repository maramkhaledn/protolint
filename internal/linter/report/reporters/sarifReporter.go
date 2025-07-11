package reporters

import (
	"io"

	"github.com/chavacava/garif"
	"github.com/maramkhaledn/protolint/linter/report"
	"github.com/maramkhaledn/protolint/linter/rule"
)

// SarifReporter creates reports formatted as a JSON
// Document.
// The document format is used according to the SARIF
// Standard.
// Refer to http://docs.oasis-open.org/sarif/sarif/v2.1.0/sarif-v2.1.0.html
// for details to the format.
type SarifReporter struct{}

var allSeverities map[string]rule.Severity = map[string]rule.Severity{
	string(rule.SeverityError):   rule.SeverityError,
	string(rule.SeverityWarning): rule.SeverityWarning,
	string(rule.SeverityNote):    rule.SeverityNote,
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Report writes failures to w formatted as a SARIF document.
func (r SarifReporter) Report(w io.Writer, fs []report.Failure) error {
	rulesByID := make(map[string]*garif.ReportingDescriptor)
	allRules := []*garif.ReportingDescriptor{}
	artifactLocations := []string{}

	tool := garif.NewDriver("protolint").
		WithInformationUri("https://github.com/maramkhaledn/protolint")

	run := garif.NewRun(garif.NewTool(tool))

	for _, failure := range fs {
		_, ruleFound := rulesByID[failure.RuleID()]
		if !ruleFound {
			rule := garif.NewRule(
				failure.RuleID(),
			).
				WithHelpUri("https://github.com/maramkhaledn/protolint")

			rulesByID[failure.RuleID()] = rule
			allRules = append(allRules, rule)
		}

		if !(contains(artifactLocations, failure.Pos().Filename)) {
			artifactLocations = append(artifactLocations, failure.Pos().Filename)
		}

		run.WithResult(
			failure.RuleID(),
			failure.Message(),
			failure.Pos().Filename,
			failure.Pos().Line,
			failure.Pos().Column,
		)

		if len(run.Results) > 0 {

			recentResult := run.Results[len(run.Results)-1]
			recentResult.Kind = garif.ResultKind_Fail

			if lvl, ok := allSeverities[failure.Severity()]; ok {
				recentResult.Level = getResultLevel(lvl)
			}
		}
	}

	tool.WithRules(allRules...)
	run.WithArtifactsURIs(artifactLocations...)

	logFile := garif.NewLogFile([]*garif.Run{run}, garif.Version210)
	return logFile.PrettyWrite(w)
}

func getResultLevel(severity rule.Severity) garif.ResultLevel {
	switch severity {
	case rule.SeverityError:
		return garif.ResultLevel_Error
	case rule.SeverityWarning:
		return garif.ResultLevel_Warning
	case rule.SeverityNote:
		return garif.ResultLevel_None
	}

	return garif.ResultLevel_None
}
