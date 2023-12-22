package builder

import "net/http"

// A Rule is a collection of conditions with a name that will apply headers, body and return code if all conditions match
type Rule struct {
	// The Name of this rule
	Name string
	// The conditions that must be met for this rule to be applied
	conditions []Condition
	// headers that should be set for the response
	headers map[string]string
	// the body to set in the response
	body string
	// the bodyFunc can return the body dynamically from a request
	bodyFunc func(r *http.Request) string
	// useBodyFunc tells the builder to use the bodyFunc instead of the body
	useBodyFunc bool
	// the code to set for the response code
	code int
}

// The RuleBuilder is a builder interface to generate a new Rule
type RuleBuilder struct {
	rule Rule
}

// NewRule creates a new named rule
func NewRule(name string) *RuleBuilder {
	return (&RuleBuilder{}).Named(name)
}

// Named sets the Name for the rule
func (r *RuleBuilder) Named(n string) *RuleBuilder {
	r.rule.Name = n
	return r
}

// WithCondition adds a new Condition to the rule
func (r *RuleBuilder) WithCondition(c Condition) *RuleBuilder {
	r.rule.conditions = append(r.rule.conditions, c)
	return r
}

// ReturnHeader adds a new header to set on the response
func (r *RuleBuilder) ReturnHeader(key string, value string) *RuleBuilder {
	if r.rule.headers == nil {
		r.rule.headers = make(map[string]string)
	}
	r.rule.headers[key] = value
	return r
}

// ReturnBody sets the body to set on the response
func (r *RuleBuilder) ReturnBody(body string) *RuleBuilder {
	r.rule.body = body
	return r
}

func (r *RuleBuilder) ReturnBodyFromFunction(f func(r *http.Request) string) *RuleBuilder {
	r.rule.bodyFunc = f
	r.rule.useBodyFunc = true
	return r
}

// ReturnCode sets the response code in the response
func (r *RuleBuilder) ReturnCode(code int) *RuleBuilder {
	r.rule.code = code
	return r
}

// Build creates a new rule
func (r *RuleBuilder) Build() Rule {
	if r.rule.code == 0 {
		r.rule.code = 200
	}
	return r.rule
}
