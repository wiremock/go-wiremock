# go-wiremock

[![GoDoc](https://godoc.org/github.com/wiremock/go-wiremock?status.svg)](http://godoc.org/github.com/wiremock/go-wiremock)
[![Actions Status](https://github.com/wiremock/go-wiremock/workflows/build/badge.svg)](https://github.com/wiremock/go-wiremock/actions?query=workflow%3Abuild)
[![Slack](https://img.shields.io/badge/slack.wiremock.org-%23wiremock—go-brightgreen?style=flat&logo=slack)](https://slack.wiremock.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/wiremock/go-wiremock)](https://goreportcard.com/report/github.com/wiremock/go-wiremock)

<a href="https://go.wiremock.org" target="_blank">
    <img width="128px" align="right" src="docs/images/logo/logo.png" alt="Go WireMock Logo"/>
</a>

---

<table>
<tr>
<td>
<img src="https://wiremock.org/images/wiremock-cloud/wiremock_cloud_logo.png" alt="WireMock Cloud Logo" height="20" align="left">
<strong>WireMock open source is supported by <a href="https://www.wiremock.io/cloud-overview?utm_source=github.com&utm_campaign=go-wiremock-README.md-banner">WireMock Cloud</a>. Please consider trying it out if your team needs advanced capabilities such as OpenAPI, dynamic state, data sources and more.</strong>
</td>
</tr>
</table>

---

The Golang client library to stub API resources in [WireMock](https://wiremock.org) using its
[REST API](https://wiremock.org/docs/api/).
The project connects to the instance and allows
setting up stubs and response templating,
or using administrative API to extract observability data.

Learn more: [Golang & WireMock Solutions page]( https://wiremock.org/docs/solutions/golang/)

## Documentation

[![GoDoc](https://godoc.org/github.com/wiremock/go-wiremock?status.svg)](http://godoc.org/github.com/wiremock/go-wiremock)

## Compatibility

The library was tested with the following distributions
of WireMock:

- WireMock 2.x - standalone deployments, including but not limited to official Docker images, Helm charts and the Java executable
- WireMock 3.x Beta - partial support, some features are
  yet to be implemented. Contributions are welcome!
- [WireMock Cloud](https://www.wiremock.io/product) -
  proprietary SaaS edition by WireMock Inc.

Note that the CI pipelines run only against the official community distributions of WireMock.
It may work for custom builds and other distributions.
Should there be any issues, contact their vendors/maintainers.

## Usage

Launch a standalone Docker instance:

```shell
docker run -it --rm -p 8080:8080 wiremock/wiremock
```

Connect to it using the client library:

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
        WithBearerToken(wiremock.StartsWith("token")).
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
## gRPC
You can mock grpc services using the library as well.

### Prerequisites

```bash
# create descriptor file for wiremock
protoc --include_imports --descriptor_set_out ./proto/services_nebius.dsc ./proto/greeting_service.proto

# download the wiremock grpc extension
curl -L https://github-registry-files.githubusercontent.com/694300421/0a0c0b80-089b-11f0-8fd1-f7e53ce5fce0?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAVCODYLSA53PQK4ZA%2F20250609%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20250609T202456Z&X-Amz-Expires=300&X-Amz-Signature=85f6bc85c1b86f2416ceb5c5031fde515f9f3003d612e2e375801e3dadde0d1c&X-Amz-SignedHeaders=host&response-content-disposition=filename%3Dwiremock-grpc-extension-standalone-0.10.0.jar&response-content-type=application%2Foctet-stream -o ./extensions/wiremock-grpc-0.10.0.jar

# run wiremock with grpc extension
docker run -it --rm -p 8080:8080 -v $(pwd)/extensions:/var/wiremock/extensions -v $(pwd)/proto:/home/wiremock/grpc wiremock/wiremock --verbose 
```

### Example of usage

```go
package main

import (
	"google.golang.org/grpc/codes"
	
	"github.com/wiremock/go-wiremock"
	wiremockGRPC "github.com/wiremock/go-wiremock/grpc"
	proto "github.com/wiremock/go-wiremock/grpc/testdata"
)

func main() {
	wiremockClient := wiremock.NewClient("http://0.0.0.0:8080")
	service := wiremockGRPC.NewService("com.example.greeting.v1.GreetingService", wiremockClient)

	service.StubFor(
		wiremockGRPC.Method("Greet").
			WithRequestMessage(wiremockGRPC.EqualToMessage(&proto.GreetingRequest{
				Name: "Tom",
			})).
			WillReturn(wiremockGRPC.Message(&proto.GreetingResponse{
				Greeting: "Hello, Tom!",
			})),
	)

	service.StubFor(
		wiremockGRPC.Method("Greet").
			WithRequestMessage(wiremockGRPC.EqualToMessage(&proto.GreetingRequest{
				Name: "Robert",
			})).
			WillReturn(wiremockGRPC.Error(codes.NotFound, "Robert not found")),
	)

	service.StubFor(
		wiremockGRPC.Method("Greet").
			WithRequestMessage(wiremockGRPC.EqualToMessage(&proto.GreetingRequest{
				Name: "Richard",
			})).
			WillReturn(wiremockGRPC.Fault(wiremock.FaultConnectionResetByPeer)),
	)
}

```

## Recording Stubs

Alternatively, you can use `wiremock` to record stubs and play them back:

```go
wiremockClient.StartRecording("https://my.saas.endpoint.com")
defer wiremockClient.StopRecording()
//… do some requests to Wiremock
//… do some assertions using your Saas' SDK
```

## Support for Authentication Schemes

The library provides support for common authentication schemes, i.e.: Basic Authentication, API Token Authentication, Bearer Authentication, Digest Access Authentication.
All of them are equivalent to manually specifying the "Authorization" header value with the appropriate prefix.
E.g. `WithBearerToken(wiremock.EqualTo("token123")).` works the same as `WithHeader("Authorization", wiremock.EqualTo("Bearer token123")).`.

### Example of usage

```go

basicAuthStub := wiremock.Get(wiremock.URLPathEqualTo("/basic")).
    WithBasicAuth("username", "password"). // same as: WithHeader("Authorization", wiremock.EqualTo("Basic dXNlcm5hbWU6cGFzc3dvcmQ=")).
    WillReturnResponse(wiremock.NewResponse().WithStatus(http.StatusOK))

bearerTokenStub := wiremock.Get(wiremock.URLPathEqualTo("/bearer")).
    WithBearerToken(wiremock.Matching("^\\S+abc\\S+$")). // same as: WithHeader("Authorization", wiremock.Matching("^Bearer \\S+abc\\S+$")).
    WillReturnResponse(wiremock.NewResponse().WithStatus(http.StatusOK))

apiTokenStub := wiremock.Get(wiremock.URLPathEqualTo("/token")).
    WithAuthToken(wiremock.StartsWith("myToken123")). // same as: WithHeader("Authorization", wiremock.StartsWith("Token myToken123")).
    WillReturnResponse(wiremock.NewResponse().WithStatus(http.StatusOK))

digestAuthStub := wiremock.Get(wiremock.URLPathEqualTo("/digest")).
    WithDigestAuth(wiremock.Contains("realm")). // same as: WithHeader("Authorization", wiremock.StartsWith("Digest ").And(Contains("realm"))).
    WillReturnResponse(wiremock.NewResponse().WithStatus(http.StatusOK))

```

## License

[MIT License](./LICENSE)

## See also

- [Golang & WireMock Solutions page]( https://wiremock.org/docs/solutions/golang/)
- [WireMock module for Testcontainers Go](https://wiremock.org/docs/solutions/testcontainers/)
