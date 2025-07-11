package strs_test

import (
	"testing"

	"github.com/maramkhaledn/protolint/linter/strs"
)

func TestPluralizeClient_ToPlural(t *testing.T) {
	tests := []struct {
		name           string
		word           string
		pluralizedWord string
	}{
		{"PluralizeNormalSingularWord", "car", "cars"},
		{"PluralizeSingularWord", "person", "people"},
		{"PluralizePluralWord", "people", "people"},
		{"PluralizeNonstandardPluralWord", "persons", "people"},
		{"PluralizeNoPluralFormWord", "moose", "moose"},
		{"PluralizePluralLatinWord", "cacti", "cacti"},
		{"PluralizeNonstandardPluralLatinWord", "cactuses", "cacti"},
		{"PluralizePluralCamelCaseWord", "office_chairs", "office_chairs"},
		{"PluralizeSingularCamelCaseWord", "office_chair", "office_chairs"},
		{"special case #312: appends s to uri", "uri", "uris"},
		{"special case #312: appends s to module_uri", "module_uri", "module_uris"},
		{"special case #312: not change uris", "uris", "uris"},
		{"special case #312: not change module_uris", "module_uris", "module_uris"},
		{"case #447: PluralizeUncountableWord", "info", "info"},
		{"case #447: PluralizeUncountableWord2", "additional_info", "additional_info"},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			c := strs.NewPluralizeClient()
			if got := c.ToPlural(test.word); got != test.pluralizedWord {
				t.Errorf("Plural(%s) got %s, but want %s", test.word, got, test.pluralizedWord)
			}
		})
	}
}

func TestPluralizeClient_AddPluralRule(t *testing.T) {
	tests := []struct {
		name           string
		word           string
		pluralizedWord string
		rule           string
		replacement    string
	}{
		{
			name:           "normal conversion",
			word:           "regex",
			pluralizedWord: "regexes",
		},
		{
			name:           "special conversion after adding manually",
			word:           "regex",
			pluralizedWord: "regexii",
			rule:           "(?i)gex$",
			replacement:    "gexii",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			c := strs.NewPluralizeClient()
			if 0 < len(test.rule) {
				c.AddPluralRule(test.rule, test.replacement)
			}
			if got := c.ToPlural(test.word); got != test.pluralizedWord {
				t.Errorf("Plural(%s) got %s, but want %s", test.word, got, test.pluralizedWord)
			}
		})
	}
}

func TestPluralizeClient_AddSingularRule(t *testing.T) {
	tests := []struct {
		name           string
		word           string
		pluralizedWord string
		rule           string
		replacement    string
	}{
		{
			name:           "normal conversion",
			word:           "singles",
			pluralizedWord: "singles",
		},
		{
			name:           "special conversion after adding manually",
			word:           "singles",
			pluralizedWord: "singulars",
			rule:           "(?i)singles$",
			replacement:    "singular",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			c := strs.NewPluralizeClient()
			if 0 < len(test.rule) {
				c.AddSingularRule(test.rule, test.replacement)
			}
			if got := c.ToPlural(test.word); got != test.pluralizedWord {
				t.Errorf("Singular(%s) got %s, but want %s", test.word, got, test.pluralizedWord)
			}
		})
	}
}
