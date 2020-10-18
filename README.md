# go-swagger

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