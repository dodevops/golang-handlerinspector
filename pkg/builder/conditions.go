package builder

import (
	"io"
	"net/http"
)

// A Condition to match the incoming http.Request again
type Condition interface {
	Matches(*http.Request) bool
}

// HasPathCondition checks if the given path was called
type HasPathCondition struct {
	path string
}

// HasPath creates a new HasPathCondition
func HasPath(path string) HasPathCondition {
	return HasPathCondition{
		path: path,
	}
}

func (h HasPathCondition) Matches(request *http.Request) bool {
	return request.URL.Path == h.path
}

// HasMethodCondition checks if the given method was used
type HasMethodCondition struct {
	method string
}

// HasMethod creates a new HasMethodCondition
func HasMethod(method string) HasMethodCondition {
	return HasMethodCondition{
		method: method,
	}
}

func (h HasMethodCondition) Matches(request *http.Request) bool {
	return request.Method == h.method
}

// HasHeaderCondition checks if the given header (key/value) exist in the request
type HasHeaderCondition struct {
	key   string
	value string
}

// HasHeader creates a new HasHeaderCondition
func HasHeader(key string, value string) *HasHeaderCondition {
	return &HasHeaderCondition{
		key:   key,
		value: value,
	}
}

func (h HasHeaderCondition) Matches(request *http.Request) bool {
	return request.Header.Get(h.key) == h.value
}

// HasQueryParamCondition checks if the given query param (key/value) was specified
type HasQueryParamCondition struct {
	key   string
	value string
}

// HasQueryParam creates a new HasHeaderCondition
func HasQueryParam(key string, value string) *HasQueryParamCondition {
	return &HasQueryParamCondition{
		key:   key,
		value: value,
	}
}

func (h HasQueryParamCondition) Matches(request *http.Request) bool {
	return request.URL.Query().Get(h.key) == h.value
}

// HasBodyCondition checks if the request's body is equal to the given string
type HasBodyCondition struct {
	body string
}

// HasBody creates a new HasBodyCondition
func HasBody(body string) *HasBodyCondition {
	return &HasBodyCondition{
		body: body,
	}
}

func (h HasBodyCondition) Matches(request *http.Request) bool {
	if b, err := io.ReadAll(request.Body); err != nil {
		return false
	} else {
		return string(b) == h.body
	}
}
