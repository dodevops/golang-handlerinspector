// Package builder implements a tool to create http.Handler value using the builder pattern
package builder

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

// HandlerBuilder uses the builder pattern to generate http.Handler values
type HandlerBuilder struct {
	// Rules is a list of HandlerInspector Rules
	Rules []Rule
	// Called records which Rules have been Called for the Inspector
	Called map[string]int
	// Failed records if no non-matching rule was called for the Inspector
	Failed bool
}

// NewBuilder creates a new HandlerBuilder. Start here.
func NewBuilder() *HandlerBuilder {
	return &HandlerBuilder{}
}

// WithRule appends a Rule to the builder
func (b *HandlerBuilder) WithRule(r Rule) *HandlerBuilder {
	b.Rules = append(b.Rules, r)
	return b
}

// Build builds a http.Handler from the Rules in the builder
func (b *HandlerBuilder) Build() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Debugf("Checking rules for request %v", r)
		foundRule := false
		for _, rule := range b.Rules {
			logrus.Debugf("Checking rule %s", rule.Name)
			matches := true
			for _, c := range rule.conditions {
				logrus.Debugf("Checking condition %v", c)
				if !c.Matches(r) {
					matches = false
				}
			}
			if matches {
				foundRule = true
				logrus.Debugf("Carrying out matching rule %s", rule.Name)
				if b.Called == nil {
					b.Called = make(map[string]int)
				}
				if v, ok := b.Called[rule.Name]; ok {
					v++
				} else {
					b.Called[rule.Name] = 1
				}
				for key, value := range rule.headers {
					w.Header().Add(key, value)
				}
				w.WriteHeader(rule.code)
				if rule.useBodyFunc {
					_, _ = fmt.Fprint(w, rule.bodyFunc(r))
				} else {
					_, _ = fmt.Fprint(w, rule.body)
				}
				return
			}
		}

		if !foundRule {
			logrus.Errorf("Didn't find a rule for request %v with body %v", r, r.Body)
			b.Failed = true
		}
	})
}
