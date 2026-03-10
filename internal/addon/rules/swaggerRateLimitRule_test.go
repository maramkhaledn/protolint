package rules_test

import (
	"reflect"
	"testing"

	"github.com/maramkhaledn/protolint/internal/addon/rules"
	"github.com/maramkhaledn/protolint/linter/rule"
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
)

func TestSwaggerRateLimitRule_Apply(t *testing.T) {
	validConstant := `{
		info: { title: "My API" }
		rate_limit: {
			max_requests: 1500
			time_window: { seconds: 60 }
			burst: 1500
			burst_time_window: { seconds: 60 }
		}
	}`

	tests := []struct {
		name         string
		inputProto   *parser.Proto
		wantFailures []string
	}{
		{
			name: "no failures when rate_limit is fully specified",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Option{
						OptionName: "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger)",
						Constant:   validConstant,
					},
				},
			},
			wantFailures: nil,
		},
		{
			name: "failure when rate_limit block is missing entirely",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Option{
						OptionName: "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger)",
						Constant:   `{ info: { title: "My API" } }`,
						Meta: meta.Meta{
							Pos: meta.Position{Filename: "test.proto", Line: 3, Column: 1},
						},
					},
				},
			},
			wantFailures: []string{
				`Option "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger)" is missing required rate_limit sub-fields: [rate_limit]`,
			},
		},
		{
			name: "failure when max_requests is missing",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Option{
						OptionName: "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger)",
						Constant: `{
							rate_limit: {
								time_window: { seconds: 60 }
								burst: 1500
								burst_time_window: { seconds: 60 }
							}
						}`,
						Meta: meta.Meta{
							Pos: meta.Position{Filename: "test.proto", Line: 3, Column: 1},
						},
					},
				},
			},
			wantFailures: []string{
				`Option "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger)" is missing required rate_limit sub-fields: [rate_limit.max_requests]`,
			},
		},
		{
			name: "failure when time_window.seconds is missing",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Option{
						OptionName: "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger)",
						Constant: `{
							rate_limit: {
								max_requests: 1500
								time_window: {}
								burst: 1500
								burst_time_window: { seconds: 60 }
							}
						}`,
						Meta: meta.Meta{
							Pos: meta.Position{Filename: "test.proto", Line: 3, Column: 1},
						},
					},
				},
			},
			wantFailures: []string{
				`Option "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger)" is missing required rate_limit sub-fields: [rate_limit.time_window.seconds]`,
			},
		},
		{
			name: "failure when burst is missing",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Option{
						OptionName: "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger)",
						Constant: `{
							rate_limit: {
								max_requests: 1500
								time_window: { seconds: 60 }
								burst_time_window: { seconds: 60 }
							}
						}`,
						Meta: meta.Meta{
							Pos: meta.Position{Filename: "test.proto", Line: 3, Column: 1},
						},
					},
				},
			},
			wantFailures: []string{
				`Option "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger)" is missing required rate_limit sub-fields: [rate_limit.burst]`,
			},
		},
		{
			name: "failure when burst_time_window.seconds is missing",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Option{
						OptionName: "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger)",
						Constant: `{
							rate_limit: {
								max_requests: 1500
								time_window: { seconds: 60 }
								burst: 1500
								burst_time_window: {}
							}
						}`,
						Meta: meta.Meta{
							Pos: meta.Position{Filename: "test.proto", Line: 3, Column: 1},
						},
					},
				},
			},
			wantFailures: []string{
				`Option "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger)" is missing required rate_limit sub-fields: [rate_limit.burst_time_window.seconds]`,
			},
		},
		{
			name: "no failures for other option names",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Option{
						OptionName: "(google.api.http)",
						Constant:   `{ get: "/v1/users" }`,
					},
				},
			},
			wantFailures: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := rules.NewSwaggerRateLimitRule(rule.SeverityError)

			failures, err := r.Apply(test.inputProto)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			var gotFailures []string
			for _, f := range failures {
				gotFailures = append(gotFailures, f.Message())
			}

			if !reflect.DeepEqual(gotFailures, test.wantFailures) {
				t.Errorf("got %v, but want %v", gotFailures, test.wantFailures)
			}
		})
	}
}
