# go-wiremock

## 📢 Important Notice:
This repository is a fork maintained for backward compatibility purposes.
The primary development and latest updates now take place in [the main repository](https://github.com/wiremock/go-wiremock)
under [the wiremock organization](https://github.com/wiremock).
We highly recommend referring to the new repository for the most recent version and to submit issues or pull requests.

[![Actions Status](https://github.com/walkerus/go-wiremock/workflows/build/badge.svg)](https://github.com/walkerus/go-wiremock/actions?query=workflow%3Abuild)
[![Go Report Card](https://goreportcard.com/badge/github.com/walkerus/go-wiremock)](https://goreportcard.com/report/github.com/walkerus/go-wiremock)

The simple package to stub HTTP resource using [WireMock admin](http://wiremock.org/docs/api/)

## Documentation

[![GoDoc](https://godoc.org/github.com/walkerus/go-wiremock?status.svg)](http://godoc.org/github.com/walkerus/go-wiremock)

## Usage

```shell
docker run -it --rm -p 8080:8080 wiremock/wiremock
```

```go
package main

import (
    "net/http"
    "testing"

    "github.com/walkerus/go-wiremock"
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