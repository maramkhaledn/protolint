package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/maramkhaledn/protolint/internal/addon/rules"
	"github.com/maramkhaledn/protolint/linter/report"
	"github.com/maramkhaledn/protolint/linter/rule"
)

func TestValidServiceNamesEndWithRule_Apply(t *testing.T) {
	validTestCase := struct {
		name         string
		inputProto   *parser.Proto
		wantFailures []report.Failure
	}{
		name: "no failures for proto with valid service names",
		inputProto: &parser.Proto{
			ProtoBody: []parser.Visitee{
				&parser.Service{
					ServiceName: "SomeServiceService",
				},
				&parser.Service{
					ServiceName: "AnotherService",
				},
			},
		},
	}

	t.Run(validTestCase.name, func(t *testing.T) {
		rule := rules.NewServiceNamesEndWithRule(rule.SeverityError, "Service")

		_, err := rule.Apply(validTestCase.inputProto)
		if err != nil {
			t.Errorf("got err %v, but want nil", err)
			return
		}
	})
}

func TestInvalidServiceNamesEndWithRule_Apply(t *testing.T) {
	invalidTestCase := struct {
		name         string
		inputProto   *parser.Proto
		wantFailures []report.Failure
	}{
		name: "failures for proto with invalid service names",
		inputProto: &parser.Proto{
			ProtoBody: []parser.Visitee{
				&parser.Service{
					ServiceName: "SomeThing",
				},
				&parser.Service{
					ServiceName: "AnotherThing",
				},
			},
		},
		wantFailures: []report.Failure{
			report.Failuref(meta.Position{}, "SERVICE_NAMES_END_WITH", string(rule.SeverityError), `Service name "SomeThing" must end with Service`),
			report.Failuref(meta.Position{}, "SERVICE_NAMES_END_WITH", string(rule.SeverityError), `Service name "AnotherThing" must end with Service`),
		},
	}

	t.Run(invalidTestCase.name, func(t *testing.T) {
		rule := rules.NewServiceNamesEndWithRule(rule.SeverityError, "Service")

		got, err := rule.Apply(invalidTestCase.inputProto)
		if err != nil {
			t.Errorf("got err %v, but want nil", err)
			return
		}
		if !reflect.DeepEqual(got, invalidTestCase.wantFailures) {
			t.Errorf("got %v, but want %v", got, invalidTestCase.wantFailures)
		}
	})
}
