# Handler Inspector

A tool to create http.Handler values using the builder pattern with an included inspector interface for http tests.

## Installation

Use

    go get github.com/dodevops/golang-handlerinspector@latest

to install and require the module.

## Usage

### Building http.Handler objects

This module includes a tool to create http.Handler objects using the 
[builder pattern](https://refactoring.guru/design-patterns/builder). The returned object will handle http requests
using a set of rules. Each rule is tested against a list of conditions. If all these conditions match, the response
is generated.

Example:

```go
builder.NewBuilder().
    WithRule(
        handlerinspector.NewRule("test-endpoint").
            WithCondition(handlerinspector.HasPath("/api/endpoint")).
            WithCondition(handlerinspector.HasMethod("GET")).
            WithCondition(handlerinspector.HasHeader("Authorization", "Bearer TESTTOKEN")).
            ReturnBody("[]").
            ReturnHeader("Content-Type", "application/json").
            Build(),
    )
```

### Inspecting the generated handler

In a test context using [httptest's server type](https://pkg.go.dev/net/http/httptest#Server), the supplied Inspector
type can be used to inspect the calls to the previously generated http.Handler and use asserts on it.

Example:

```go
h := builder.NewBuilder().
    WithRule(
        handlerinspector.NewRule("test-endpoint").
            WithCondition(handlerinspector.HasPath("/api/endpoint")).
            WithCondition(handlerinspector.HasMethod("GET")).
            WithCondition(handlerinspector.HasHeader("Authorization", "Bearer TESTTOKEN")).
            ReturnBody("[]").
            ReturnHeader("Content-Type", "application/json").
            Build(),
    )
s := httptest.NewServer(h.Build())
defer s.Close()

// (...)

i := inspector.NewInspector(h)
assert.Equal(t, i.Failed(), false)
assert.Equal(t, i.AllWereCalled(), true)
```

## Motivation

While creating a test suite for our [vmware-rest-proxy](https://github.com/dodevops/vmware-rest-proxy) we were looking
for a nice mock server implementation only to find it in a Golang core package, which was a pleasant suprise.

However, creating http.Handlers for it and especially testing them was a bit of a hassle, and we didn't find any easy
tool to achieve it.

Thus, here's Handler Inspector. It's currently designed to meet our needs for the vmware-rest-proxy, but can easily
be extended if other needs arise.