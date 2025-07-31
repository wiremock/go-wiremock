package wiremock_test

import (
	"context"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/wiremock/go-wiremock"
)

type WiremockTestService struct {
	container tc.Container
	client    *wiremock.Client
	baseURL   string
}

func getWiremockTestService(ctx context.Context, t *testing.T) *WiremockTestService {
	req := tc.ContainerRequest{
		Name:         "go-wiremock",
		Image:        "wiremock/wiremock:latest",
		ExposedPorts: []string{"8080/tcp"},
		Cmd:          []string{"--verbose"},
		WaitingFor:   wait.ForHealthCheck(),
	}

	c, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            true,
	})
	requireNoError(t, err)

	baseURL, err := c.PortEndpoint(ctx, "8080/tcp", "http")
	requireNoError(t, err)

	client := wiremock.NewClient(baseURL)

	return &WiremockTestService{
		container: c,
		client:    client,
		baseURL:   baseURL,
	}
}

// Reset is a shortcut for client.Reset
func (s *WiremockTestService) Reset() error {
	return s.client.Reset()
}

func TestClient_GetAllRequests(t *testing.T) {
	ctx := context.Background()
	svc := getWiremockTestService(ctx, t)

	t.Run("empty request journal", func(t *testing.T) {
		err := svc.Reset()
		requireNoError(t, err)

		events, err := svc.client.GetAllRequests()
		requireNoError(t, err)

		assertEqual(t, 0, events.Meta.Total)
		assertEqual(t, 0, len(events.Requests))
	})

	t.Run("with requests", func(t *testing.T) {
		err := svc.Reset()
		requireNoError(t, err)

		stub := wiremock.NewStubRule("GET", wiremock.URLMatching("/test")).WillReturnResponse(wiremock.OK())
		err = svc.client.StubFor(stub)
		requireNoError(t, err)

		_, err = http.Get(svc.baseURL + "/test")
		requireNoError(t, err)

		_, err = http.Get(svc.baseURL + "/not-a-stub?param=1234")
		requireNoError(t, err)

		events, err := svc.client.GetAllRequests()
		requireNoError(t, err)

		assertEqual(t, 2, events.Meta.Total)
		assertEqual(t, 2, len(events.Requests))
		assertEqual(t, "/not-a-stub?param=1234", events.Requests[0].Request.URL)
		assertEqual(t, "/test", events.Requests[1].Request.URL)
	})
}

func TestClient_GetRequestsByID(t *testing.T) {
	ctx := context.Background()
	svc := getWiremockTestService(ctx, t)

	t.Run("invalid request id", func(t *testing.T) {
		err := svc.Reset()
		requireNoError(t, err)

		id := uuid.New()

		request, err := svc.client.GetRequestByID(id.String())
		if err == nil {
			t.Fatal("expected error, got none")
		}
		if !strings.Contains(err.Error(), "bad response status: 404") {
			t.Errorf("expected error message to contain 'bad response status: 404', got %s", err.Error())
		}

		assertNil(t, request)
	})

	t.Run("valid request id", func(t *testing.T) {
		err := svc.Reset()
		requireNoError(t, err)

		_, err = http.Get(svc.baseURL + "/test-1")
		requireNoError(t, err)

		events, err := svc.client.GetAllRequests()
		requireNoError(t, err)

		// Add another request to the journal to ensure we're really getting the one we asked for
		_, err = http.Get(svc.baseURL + "/test-2")
		requireNoError(t, err)

		reqID := events.Requests[0].ID

		req, err := svc.client.GetRequestByID(reqID)
		requireNoError(t, err)

		assertEqual(t, "/test-1", req.Request.URL)
	})
}

func TestClient_FindRequestsByCriteria(t *testing.T) {
	ctx := context.Background()
	svc := getWiremockTestService(ctx, t)

	t.Run("no requests to find", func(t *testing.T) {
		err := svc.Reset()
		requireNoError(t, err)

		stub := wiremock.NewStubRule("GET", wiremock.URLMatching("/test")).WillReturnResponse(wiremock.OK())
		err = svc.client.StubFor(stub)
		requireNoError(t, err)

		_, err = http.Get(svc.baseURL + "/other")
		requireNoError(t, err)

		// Find requests for /test
		resp, err := svc.client.FindRequestsByCriteria(stub.Request())
		requireNoError(t, err)
		assertEqual(t, 0, len(resp.Requests))
	})

	t.Run("with requests to find", func(t *testing.T) {
		err := svc.Reset()
		requireNoError(t, err)

		stub := wiremock.NewStubRule("GET", wiremock.URLMatching("/test")).WillReturnResponse(wiremock.OK())
		err = svc.client.StubFor(stub)
		requireNoError(t, err)

		_, err = http.Get(svc.baseURL + "/test")
		requireNoError(t, err)

		_, err = http.Get(svc.baseURL + "/test")
		requireNoError(t, err)

		_, err = http.Get(svc.baseURL + "/other")
		requireNoError(t, err)

		// Find requests for /test
		resp, err := svc.client.FindRequestsByCriteria(stub.Request())
		requireNoError(t, err)

		assertEqual(t, 2, len(resp.Requests))
		for _, r := range resp.Requests {
			assertEqual(t, "/test", r.URL)
		}

		events, err := svc.client.GetAllRequests()
		requireNoError(t, err)

		// All requests should be in the events
		assertEqual(t, 3, len(events.Requests))
	})
}

func TestClient_FindUnmatchedRequests(t *testing.T) {
	ctx := context.Background()
	svc := getWiremockTestService(ctx, t)

	t.Run("no requests to find", func(t *testing.T) {
		err := svc.Reset()
		requireNoError(t, err)

		stub := wiremock.NewStubRule("GET", wiremock.URLMatching("/test")).WillReturnResponse(wiremock.OK())
		err = svc.client.StubFor(stub)
		requireNoError(t, err)

		_, err = http.Get(svc.baseURL + "/test")
		requireNoError(t, err)

		resp, err := svc.client.FindUnmatchedRequests()
		requireNoError(t, err)
		assertEqual(t, 0, len(resp.Requests))
	})

	t.Run("with requests to find", func(t *testing.T) {
		err := svc.Reset()
		requireNoError(t, err)

		stub := wiremock.NewStubRule("GET", wiremock.URLMatching("/test")).WillReturnResponse(wiremock.OK())
		err = svc.client.StubFor(stub)
		requireNoError(t, err)

		_, err = http.Get(svc.baseURL + "/test")
		requireNoError(t, err)

		_, err = http.Get(svc.baseURL + "/unmatched")
		requireNoError(t, err)

		_, err = http.Get(svc.baseURL + "/unmatched")
		requireNoError(t, err)

		// Find requests for /test
		resp, err := svc.client.FindUnmatchedRequests()
		requireNoError(t, err)

		assertEqual(t, 2, len(resp.Requests))
		for _, r := range resp.Requests {
			assertEqual(t, "/unmatched", r.URL)
		}
	})
}

func TestClient_DeleteAllRequests(t *testing.T) {
	ctx := context.Background()
	svc := getWiremockTestService(ctx, t)

	err := svc.Reset()
	requireNoError(t, err)

	_, err = http.Get(svc.baseURL + "/test")
	requireNoError(t, err)

	_, err = http.Get(svc.baseURL + "/test2")
	requireNoError(t, err)

	events, err := svc.client.GetAllRequests()
	requireNoError(t, err)

	assertEqual(t, 2, events.Meta.Total)

	err = svc.client.DeleteAllRequests()
	requireNoError(t, err)

	events, err = svc.client.GetAllRequests()
	requireNoError(t, err)

	assertEqual(t, 0, events.Meta.Total)
}

func TestClient_DeleteRequestByID(t *testing.T) {
	ctx := context.Background()
	svc := getWiremockTestService(ctx, t)

	t.Run("invalid request id", func(t *testing.T) {
		err := svc.Reset()
		requireNoError(t, err)

		id := uuid.New()

		err = svc.client.DeleteRequestByID(id.String())
		// Deleting a non-existing request does nothing and returns a 200.
		requireNoError(t, err)
	})

	t.Run("valid request id", func(t *testing.T) {
		err := svc.Reset()
		requireNoError(t, err)

		_, err = http.Get(svc.baseURL + "/test-1")
		requireNoError(t, err)

		events, err := svc.client.GetAllRequests()
		requireNoError(t, err)

		// Add another request to the journal, it should still be there after the deletion
		_, err = http.Get(svc.baseURL + "/test-2")
		requireNoError(t, err)

		reqID := events.Requests[0].ID

		// Delete request for /test-1
		err = svc.client.DeleteRequestByID(reqID)
		requireNoError(t, err)

		events, err = svc.client.GetAllRequests()
		requireNoError(t, err)

		assertEqual(t, 1, len(events.Requests))
		assertEqual(t, "/test-2", events.Requests[0].Request.URL)
	})
}

func TestClient_DeleteRequestsByCriteria(t *testing.T) {
	ctx := context.Background()
	svc := getWiremockTestService(ctx, t)

	t.Run("no requests to delete", func(t *testing.T) {
		err := svc.Reset()
		requireNoError(t, err)

		stub := wiremock.NewStubRule("GET", wiremock.URLMatching("/test")).WillReturnResponse(wiremock.OK())
		err = svc.client.StubFor(stub)
		requireNoError(t, err)

		_, err = http.Get(svc.baseURL + "/other")
		requireNoError(t, err)

		// Delete requests for /test
		resp, err := svc.client.DeleteRequestsByCriteria(stub.Request())
		requireNoError(t, err)
		assertEqual(t, 0, len(resp.Requests))

		events, err := svc.client.GetAllRequests()
		requireNoError(t, err)

		// The other request should still be present
		assertEqual(t, 1, len(events.Requests))
		assertEqual(t, "/other", events.Requests[0].Request.URL)
	})

	t.Run("with requests to delete", func(t *testing.T) {
		err := svc.Reset()
		requireNoError(t, err)

		stub := wiremock.NewStubRule("GET", wiremock.URLMatching("/test")).WillReturnResponse(wiremock.OK())
		err = svc.client.StubFor(stub)
		requireNoError(t, err)

		_, err = http.Get(svc.baseURL + "/test")
		requireNoError(t, err)

		_, err = http.Get(svc.baseURL + "/test")
		requireNoError(t, err)

		_, err = http.Get(svc.baseURL + "/other")
		requireNoError(t, err)

		// Delete requests for /test
		resp, err := svc.client.DeleteRequestsByCriteria(stub.Request())
		requireNoError(t, err)

		assertEqual(t, 2, len(resp.Requests))
		for _, evt := range resp.Requests {
			assertEqual(t, "/test", evt.Request.URL)
		}

		events, err := svc.client.GetAllRequests()
		requireNoError(t, err)

		// The other request should be the only one remaining
		assertEqual(t, 1, len(events.Requests))
		assertEqual(t, "/other", events.Requests[0].Request.URL)
	})
}

func requireNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func assertNil[T any](t *testing.T, value *T) {
	if value != nil {
		t.Errorf("expected nil, got %v", value)
	}
}

func assertEqual[T any](t *testing.T, expected, actual T) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %T(%v), got %T(%v)", expected, expected, actual, actual)
	}
}
