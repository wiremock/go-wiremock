# go-wiremock

[![Actions Status](https://github.com/wiremock/go-wiremock/workflows/build/badge.svg)](https://github.com/wiremock/go-wiremock/actions?query=workflow%3Abuild)
[![Go Report Card](https://goreportcard.com/badge/github.com/wiremock/go-wiremock)](https://goreportcard.com/report/github.com/wiremock/go-wiremock)

The simple package to stub HTTP resource using [WireMock admin](http://wiremock.org/docs/api/)

## Documentation
### NPM Documentation

[![GoDoc](https://godoc.org/github.com/wiremock/go-wiremock?status.svg)](http://godoc.org/github.com/wiremock/go-wiremock)

### Wiremock Solutions for GoLang
[![Golang Wiremock Solutions](https://wiremock.org/images/logos/wiremock/logo_square.svg)]( https://wiremock.org/docs/solutions/golang/)

## Usage

```shell
docker run -it --rm -p 8080:8080 wiremock/wiremock
```

```go
package main

import (
    "net/http"
    "testing"

    "github.com/wiremock/go-wiremock"
)

func TestSome(t *testing.T) {
    wiremockClient := wiremock.NewClient("http://0.0.0.0:8080")
    defer wiremockClient.Reset()

    // stubbing POST http://0.0.0.0:8080/example
    wiremockClient.StubFor(wiremock.Post(wiremock.URLPathEqualTo("/example")).
        WithQueryParam("firstName", wiremock.EqualTo("John")).
        WithQueryParam("lastName", wiremock.NotMatching("Black")).
        WithBodyPattern(wiremock.EqualToJson(`{"meta": "information"}`)).
        WithHeader("x-session", wiremock.Matching("^\\S+fingerprint\\S+$")).
        WillReturnResponse(
            wiremock.NewResponse().
                WithJSONBody(map[string]interface{}{
                    "code":   400,
                    "detail": "detail",
                }).
                WithHeader("Content-Type", "application/json").
                WithStatus(http.StatusBadRequest),
        ).
        AtPriority(1))

    // scenario
    defer wiremockClient.ResetAllScenarios()
    wiremockClient.StubFor(wiremock.Get(wiremock.URLPathEqualTo("/status")).
        WillReturnResponse(
            wiremock.NewResponse().
                WithJSONBody(map[string]interface{}{
                    "status": nil,
                }).
                WithHeader("Content-Type", "application/json").
                WithStatus(http.StatusOK),
        ).
        InScenario("Set status").
        WhenScenarioStateIs(wiremock.ScenarioStateStarted))

    wiremockClient.StubFor(wiremock.Post(wiremock.URLPathEqualTo("/state")).
        WithBodyPattern(wiremock.EqualToJson(`{"status": "started"}`)).
        InScenario("Set status").
        WillSetStateTo("Status started"))

    statusStub := wiremock.Get(wiremock.URLPathEqualTo("/status")).
        WillReturnResponse(
            wiremock.NewResponse().
                WithJSONBody(map[string]interface{}{
                    "status": "started",
                }).
                WithHeader("Content-Type", "application/json").
                WithStatus(http.StatusOK),
        ).
        InScenario("Set status").
        WhenScenarioStateIs("Status started")
    wiremockClient.StubFor(statusStub)

    //testing code...

    verifyResult, _ := wiremockClient.Verify(statusStub.Request(), 1)
    if !verifyResult {
        //...
    }

    wiremockClient.DeleteStub(statusStub)
}
```