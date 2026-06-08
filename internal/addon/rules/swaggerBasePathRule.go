package rules

import (
    "regexp"
    "strings"

    "github.com/maramkhaledn/protolint/linter/report"
    "github.com/maramkhaledn/protolint/linter/rule"
    "github.com/maramkhaledn/protolint/linter/visitor"
    "github.com/yoheimuta/go-protoparser/v4/parser"
)

// basePathRegex extracts the base_path value from an openapiv2_swagger option, e.g.
// `base_path: "/my-api"` -> "/my-api".
var basePathRegex = regexp.MustCompile(`base_path\s*:\s*"([^"]*)"`)

// SwaggerBasePathRule verifies that the base_path field in the openapiv2_swagger option is
// dash-separated and never contains underscores.
type SwaggerBasePathRule struct {
    RuleWithSeverity
}

// NewSwaggerBasePathRule creates a new SwaggerBasePathRule.
func NewSwaggerBasePathRule(severity rule.Severity) SwaggerBasePathRule {
    return SwaggerBasePathRule{
        RuleWithSeverity: RuleWithSeverity{severity: severity},
    }
}

// ID returns the ID of this rule.
func (r SwaggerBasePathRule) ID() string {
    return "SWAGGER_BASE_PATH"
}

// Purpose returns the purpose of this rule.
func (r SwaggerBasePathRule) Purpose() string {
    return "Verifies that the base_path in the openapiv2_swagger option is dash-separated and never contains underscores."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r SwaggerBasePathRule) IsOfficial() bool {
    return true
}

// Apply applies the rule to the proto.
func (r SwaggerBasePathRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
    v := &swaggerBasePathVisitor{
        BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID(), string(r.Severity())),
    }
    return visitor.RunVisitor(v, proto, r.ID())
}

type swaggerBasePathVisitor struct {
    *visitor.BaseAddVisitor
}

// VisitOption checks that the openapiv2_swagger option's base_path uses dashes, not underscores.
func (v *swaggerBasePathVisitor) VisitOption(option *parser.Option) bool {
    if option.OptionName != "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger)" {
        return false
    }

    m := basePathRegex.FindStringSubmatch(option.Constant)
    if len(m) < 2 {
        return false
    }

    basePath := m[1]
    if strings.Contains(basePath, "_") {
        v.AddFailuref(
            option.Meta.Pos,
            `base_path %q in the openapiv2_swagger option must be dash-separated and never contain underscores`,
            basePath,
        )
    }

    return false
}
