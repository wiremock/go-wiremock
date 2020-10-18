# go-wiremock
[![Actions Status](https://github.com/walkerus/go-wiremock/workflows/build/badge.svg)](https://github.com/walkerus/go-wiremock/actions?query=workflow%3Abuild)
[![Go Report Card](https://goreportcard.com/badge/github.com/walkerus/go-wiremock)](https://goreportcard.com/report/github.com/walkerus/go-wiremock)

A simple package for stub http resource using WireMock admin

## Install
```go get https://github.com/walkerus/go-wiremock```

## Usage
```
docker run -it --rm -p 8080:8080 rodolpheche/wiremock
```
```go
wiremockClient := NewClient("http://0.0.0.0:8080")
defer wiremockClient.Clear()
wiremockClient.StubFor(Post(URLPathEqualTo("/example")).
		WithQueryParam("firstName", EqualTo("Jhon")).
		WithQueryParam("lastName", NotMatching("Black")).
		WithBodyPattern(EqualToJson(`{"meta": "information"}`)).
		WithHeader("x-session", Matching("^\\S+fingerprint\\S+$")).
		WillReturn(
			`{"code": 400, "detail": "detail"}`,
			map[string]string{"Content-Type": "application/json"},
			400,
		))
```