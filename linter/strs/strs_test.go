package strs_test

import (
	"reflect"
	"testing"

	"github.com/maramkhaledn/protolint/linter/strs"
)

func TestIsUpperCamelCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "the first letter is not an uppercase character",
			input: "hello",
		},
		{
			name:  "_ is included",
			input: "Hello_world",
		},
		{
			name:  ". is included",
			input: "Hello.world",
		},
		{
			name:  "the first letter is an uppercase character",
			input: "Hello",
			want:  true,
		},
		{
			name:  "the first letter is an uppercase character and rest is a camel case",
			input: "HelloWorld",
			want:  true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got := strs.IsUpperCamelCase(test.input)
			if got != test.want {
				t.Errorf("got %v, but want %v", got, test.want)
			}
		})
	}
}

func TestIsLowerCamelCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "the first letter is an uppercase character",
			input: "Hello",
		},
		{
			name:  "_ is included",
			input: "hello_world",
		},
		{
			name:  ". is included",
			input: "hello.world",
		},
		{
			name:  "the first letter is a lowercase character",
			input: "hello",
			want:  true,
		},
		{
			name:  "the first letter is a lowercase character and rest is a camel case",
			input: "helloWorld",
			want:  true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got := strs.IsLowerCamelCase(test.input)
			if got != test.want {
				t.Errorf("got %v, but want %v", got, test.want)
			}
		})
	}
}
func TestIsUpperSnakeCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name: "empty is not uppercase",
		},
		{
			name:  "includes lowercase characters",
			input: "hello",
		},
		{
			name:  "includes a lowercase character",
			input: "hELLO",
		},
		{
			name:  "all uppercase",
			input: "HELLO",
			want:  true,
		},
		{
			name:  "all uppercase with underscore",
			input: "FIRST_VALUE",
			want:  true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got := strs.IsUpperSnakeCase(test.input)
			if got != test.want {
				t.Errorf("got %v, but want %v", got, test.want)
			}
		})
	}
}

func TestIsLowerSnakeCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name: "empty is not lowercase",
		},
		{
			name:  "includes uppercase characters",
			input: "HELLO",
		},
		{
			name:  "includes a uppercase character",
			input: "Hello",
		},
		{
			name:  "all lowercase",
			input: "hello",
			want:  true,
		},
		{
			name:  "all lowercase with underscore",
			input: "song_name",
			want:  true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got := strs.IsLowerSnakeCase(test.input)
			if got != test.want {
				t.Errorf("got %v, but want %v", got, test.want)
			}
		})
	}
}

func TestSplitCamelCaseWord(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name: "if s is empty, returns nil",
		},
		{
			name:  "if s is not camel_case, returns nil",
			input: "not_camel",
		},
		{
			name:  "input consists of one word",
			input: "Account",
			want: []string{
				"Account",
			},
		},
		{
			name:  "input consists of words with an initial capital",
			input: "AccountStatus",
			want: []string{
				"Account",
				"Status",
			},
		},
		{
			name:  "input consists of words without an initial capital",
			input: "accountStatus",
			want: []string{
				"account",
				"Status",
			},
		},
		{
			name:  "input consists of words with continuous upper letters",
			input: "ACCOUNTStatusException",
			want: []string{
				"ACCOUNT",
				"Status",
				"Exception",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got := strs.SplitCamelCaseWord(test.input)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("got %v, but want %v", got, test.want)
			}
		})
	}
}

func TestToUpperSnakeCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "s is not camel_case",
			input: "not_camel",
			want:  "NOT_CAMEL",
		},
		{
			name:  "input consists of one word",
			input: "Account",
			want:  "ACCOUNT",
		},
		{
			name:  "input consists of words with an initial capital",
			input: "AccountStatus",
			want:  "ACCOUNT_STATUS",
		},
		{
			name:  "input consists of words without an initial capital",
			input: "accountStatus",
			want:  "ACCOUNT_STATUS",
		},
		{
			name:  "convert from camel-case strings starting with a 2-letter abbreviation #341",
			input: "ITDepartmentRegion",
			want:  "IT_DEPARTMENT_REGION",
		},
		{
			name:  "convert from camel-case strings includes OAuth #351",
			input: "ListAccountPipedriveOAuthsEnabledFilter",
			want:  "LIST_ACCOUNT_PIPEDRIVE_OAUTHS_ENABLED_FILTER",
		},
		{
			name:  "convert from camel-case strings starting with OAuth #351",
			input: "OAuthsEnabledFilter",
			want:  "OAUTHS_ENABLED_FILTER",
		},
		{
			name:  "convert from camel-case strings includes with a 2-letter abbreviation #341",
			input: "ListITDepartmentRegion",
			want:  "LIST_IT_DEPARTMENT_REGION",
		},
		{
			name:  "convert from camel-case strings starting with a 3-letter abbreviation",
			input: "ITCDepartmentRegion",
			want:  "ITC_DEPARTMENT_REGION",
		},
		{
			name:  "convert from camel-case strings includes a 3-letter abbreviation",
			input: "ListITCDepartmentRegion",
			want:  "LIST_ITC_DEPARTMENT_REGION",
		},
		{
			name:  "input consists of kebab case",
			input: "account-status",
			want:  "ACCOUNT_STATUS",
		},
		{
			name:  "input consists of .",
			input: "account.status",
			want:  "ACCOUNT_STATUS",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got := strs.ToUpperSnakeCase(test.input)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("got %v, but want %v", got, test.want)
			}
		})
	}
}

func TestToLowerSnakeCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "input consists of one word",
			input: "Account",
			want:  "account",
		},
		{
			name:  "input consists of words with an initial capital",
			input: "AccountStatus",
			want:  "account_status",
		},
		{
			name:  "input consists of words without an initial capital",
			input: "accountStatus",
			want:  "account_status",
		},
		{
			name:  "input consists of kebab case",
			input: "account-status",
			want:  "account_status",
		},
		{
			name:  "input consists of .",
			input: "account.status",
			want:  "account_status",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got := strs.ToLowerSnakeCase(test.input)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("got %v, but want %v", got, test.want)
			}
		})
	}
}

func TestToUpperCamelCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "input consists of one word",
			input: "account",
			want:  "Account",
		},
		{
			name:  "input consists of words with an initial capital",
			input: "AccountStatus",
			want:  "AccountStatus",
		},
		{
			name:  "input consists of words without an initial capital",
			input: "accountStatus",
			want:  "AccountStatus",
		},
		{
			name:  "input consists of words without capital letters",
			input: "accountstatus",
			want:  "Accountstatus",
		},
		{
			name:  "input lower_snake_case",
			input: "account_status",
			want:  "AccountStatus",
		},
		{
			name:  "input UPPER_SNAKE_CASE",
			input: "ACCOUNT_STATUS",
			want:  "AccountStatus",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got := strs.ToUpperCamelCase(test.input)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("got %v, but want %v", got, test.want)
			}
		})
	}
}

func TestToLowerCamelCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "input consists of one word",
			input: "account",
			want:  "account",
		},
		{
			name:  "input consists of words with an initial capital",
			input: "AccountStatus",
			want:  "accountStatus",
		},
		{
			name:  "input consists of words without an initial capital",
			input: "accountStatus",
			want:  "accountStatus",
		},
		{
			name:  "input consists of words without capital letters",
			input: "accountstatus",
			want:  "accountstatus",
		},
		{
			name:  "input lower_snake_case",
			input: "account_status",
			want:  "accountStatus",
		},
		{
			name:  "input UPPER_SNAKE_CASE",
			input: "ACCOUNT_STATUS",
			want:  "accountStatus",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got := strs.ToLowerCamelCase(test.input)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("got %v, but want %v", got, test.want)
			}
		})
	}
}

func TestSplitSnakeCaseWord(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name: "if s is empty, returns nil",
		},
		{
			name:  "if s is not snake_case, returns nil",
			input: "_not_snake",
		},
		{
			name:  "input consists of one word",
			input: "HELLO",
			want: []string{
				"HELLO",
			},
		},
		{
			name:  "input consists of multiple upper case words",
			input: "REASON_FOR_ERROR",
			want: []string{
				"REASON",
				"FOR",
				"ERROR",
			},
		},
		{
			name:  "input consists of multiple lower case words",
			input: "reason_for_error",
			want: []string{
				"reason",
				"for",
				"error",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got := strs.SplitSnakeCaseWord(test.input)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("got %v, but want %v", got, test.want)
			}
		})
	}
}
