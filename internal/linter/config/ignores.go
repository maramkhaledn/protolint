package config

// Ignores represents list about files ignoring the specific rule.
type Ignores []Ignore

func (is Ignores) shouldSkipRule(
	ruleID string,
	displayPath string,
) bool {
	for i := range is {
		if is[i].shouldSkipRule(ruleID, displayPath) {
			return true
		}
	}
	return false
}

func (is *Ignores) validate() error {
	for i := range *is {
		if err := (*is)[i].validate(i); err != nil {
			return err
		}
	}
	return nil
}
