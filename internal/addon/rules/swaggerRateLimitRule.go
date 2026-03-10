package rules

import (
	"regexp"

	"github.com/maramkhaledn/protolint/linter/report"
	"github.com/maramkhaledn/protolint/linter/rule"
	"github.com/maramkhaledn/protolint/linter/visitor"
	"github.com/yoheimuta/go-protoparser/v4/parser"
)

var (
	rateLimitBlockRegex       = regexp.MustCompile(`rate_limit\s*:\s*\{`)
	rateLimitMaxRequestsRegex = regexp.MustCompile(`max_requests\s*:\s*\d+`)
	rateLimitTimeWindowRegex  = regexp.MustCompile(`\btime_window\s*:\s*\{[^}]*seconds\s*:\s*\d+`)
	rateLimitBurstRegex       = regexp.MustCompile(`\bburst\s*:\s*\d+`)
	rateLimitBurstTWRegex     = regexp.MustCompile(`burst_time_window\s*:\s*\{[^}]*seconds\s*:\s*\d+`)
)

// SwaggerRateLimitRule verifies that the openapiv2_swagger option always contains a rate_limit field
// with all required sub-fields: max_requests, time_window (with seconds), burst, and burst_time_window (with seconds).
type SwaggerRateLimitRule struct {
	RuleWithSeverity
}

// NewSwaggerRateLimitRule creates a new SwaggerRateLimitRule.
func NewSwaggerRateLimitRule(severity rule.Severity) SwaggerRateLimitRule {
	return SwaggerRateLimitRule{
		RuleWithSeverity: RuleWithSeverity{severity: severity},
	}
}

// ID returns the ID of this rule.
func (r SwaggerRateLimitRule) ID() string {
	return "SWAGGER_RATE_LIMIT"
}

// Purpose returns the purpose of this rule.
func (r SwaggerRateLimitRule) Purpose() string {
	return "Verifies that the openapiv2_swagger option contains a rate_limit field with max_requests, time_window, burst, and burst_time_window."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r SwaggerRateLimitRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r SwaggerRateLimitRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &swaggerRateLimitVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID(), string(r.Severity())),
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type swaggerRateLimitVisitor struct {
	*visitor.BaseAddVisitor
}

// VisitOption checks that the openapiv2_swagger option has a valid rate_limit field.
func (v *swaggerRateLimitVisitor) VisitOption(option *parser.Option) bool {
	if option.OptionName != "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger)" {
		return false
	}

	missing := checkRateLimitFields(option.Constant)
	if len(missing) > 0 {
		v.AddFailuref(option.Meta.Pos, `Option "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger)" is missing required rate_limit sub-fields: %v`, missing)
	}

	return false
}

// checkRateLimitFields returns a list of missing required fields within the rate_limit block.
func checkRateLimitFields(constant string) []string {
	if !rateLimitBlockRegex.MatchString(constant) {
		return []string{"rate_limit"}
	}

	var missing []string
	if !rateLimitMaxRequestsRegex.MatchString(constant) {
		missing = append(missing, "rate_limit.max_requests")
	}
	if !rateLimitTimeWindowRegex.MatchString(constant) {
		missing = append(missing, "rate_limit.time_window.seconds")
	}
	if !rateLimitBurstRegex.MatchString(constant) {
		missing = append(missing, "rate_limit.burst")
	}
	if !rateLimitBurstTWRegex.MatchString(constant) {
		missing = append(missing, "rate_limit.burst_time_window.seconds")
	}
	return missing
}
