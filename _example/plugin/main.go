package main

import (
	"flag"

	"github.com/maramkhaledn/protolint/_example/plugin/customrules"
	"github.com/maramkhaledn/protolint/internal/addon/rules"
	"github.com/maramkhaledn/protolint/linter/rule"
	"github.com/maramkhaledn/protolint/plugin"
)

var (
	goStyle = flag.Bool("go_style", true, "the comments should follow a golang style")
)

func main() {
	flag.Parse()

	plugin.RegisterCustomRules(
		// The purpose of this line just illustrates that you can implement the same as internal linter rules.
		rules.NewEnumsHaveCommentRule(rule.SeverityWarning, *goStyle),

		// A common custom rule example. It's simple.
		customrules.NewEnumNamesLowerSnakeCaseRule(),

		// Wrapping with RuleGen allows referring to command-line flags.
		plugin.RuleGen(func(
			verbose bool,
			fixMode bool,
		) rule.Rule {
			return customrules.NewSimpleRule(verbose, fixMode, rule.SeverityError)
		}),
	)
}
