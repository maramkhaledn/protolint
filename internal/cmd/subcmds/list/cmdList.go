package list

import (
	"fmt"
	"io"

	"github.com/maramkhaledn/protolint/internal/addon/plugin/shared"

	"github.com/maramkhaledn/protolint/internal/linter/config"

	"github.com/maramkhaledn/protolint/internal/cmd/subcmds"
	"github.com/maramkhaledn/protolint/internal/osutil"
	"github.com/maramkhaledn/protolint/linter/autodisable"
	"github.com/maramkhaledn/protolint/linter/rule"
)

// CmdList is a rule list command.
type CmdList struct {
	stdout io.Writer
	stderr io.Writer
	flags  Flags
}

// NewCmdList creates a new CmdList.
func NewCmdList(
	flags Flags,
	stdout io.Writer,
	stderr io.Writer,
) *CmdList {
	return &CmdList{
		flags:  flags,
		stdout: stdout,
		stderr: stderr,
	}
}

// Run lists each rule description.
func (c *CmdList) Run() osutil.ExitCode {
	err := c.run()
	if err != nil {
		return osutil.ExitInternalFailure
	}
	return osutil.ExitSuccess
}

func (c *CmdList) run() error {
	rules, err := hasIDAndPurposes(c.flags.Plugins)
	if err != nil {
		return err
	}

	for _, r := range rules {
		_, err := fmt.Fprintf(
			c.stdout,
			"%s: %s\n",
			r.ID(),
			r.Purpose(),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

type hasIDAndPurpose interface {
	rule.HasID
	rule.HasPurpose
}

func hasIDAndPurposes(plugins []shared.RuleSet) ([]hasIDAndPurpose, error) {
	rs, err := subcmds.NewAllRules(config.RulesOption{}, false, autodisable.Noop, false, plugins)
	if err != nil {
		return nil, err
	}

	var rules []hasIDAndPurpose
	for _, r := range rs {
		rules = append(rules, r)
	}
	return rules, nil
}
