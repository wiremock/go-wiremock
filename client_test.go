package wiremock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestStubRule_ToJson(t *testing.T) {
	newStubRule := NewStubRule("PATCH", URLMatching("/example"))
	expectedRequestBody := fmt.Sprintf(`{"uuid":"%s","id":"%s","request":{"method":"PATCH","urlPattern":"/example"},"response":{"status":200}}`, newStubRule.uuid, newStubRule.uuid)
	result, err := json.Marshal(newStubRule)

	if err != nil {
		t.Fatalf("StubRole json.Marshal error: %v", err)
	}
	if string(result) != expectedRequestBody {
		t.Errorf("expected requestBody %q; got %q", expectedRequestBody, string(result))
	}

	postStubRule := Post(URLPathEqualTo("/example")).
		WithHost(EqualTo("localhost")).
		WithScheme("http").
		WithPort(8080).
		WithQueryParam("firstName", EqualTo("John").Or(EqualTo("Jack"))).
		WithQueryParam("lastName", NotMatching("Black")).
		WithQueryParam("nickname", EqualToIgnoreCase("johnBlack")).
		WithQueryParam("address", Includes(EqualTo("1"), Contains("2"), NotContains("3"))).
		WithQueryParam("id", Contains("1").And(NotContains("2"))).
		WithBodyPattern(EqualToJson(`{"meta": "information"}`, IgnoreArrayOrder, IgnoreExtraElements)).
		WithBodyPattern(Contains("information")).
		WithMultipartPattern(
			NewMultipartPattern().
				WithName("info").
				WithHeader("Content-Type", Contains("charset")).
				WithBodyPattern(EqualToJson("{}", IgnoreExtraElements)),
		).
		WithBasicAuth("username", "password").
		WithHeader("x-absent", Absent()).
		WithCookie("absentcookie", Absent()).
		WithHeader("x-session", Matching("^\\S+@\\S+$")).
		WithCookie("session", EqualToXml("<xml>")).
		WillReturnResponse(
			NewResponse().
				WithStatus(http.StatusBadRequest).
				WithHeader("Content-Type", "application/json").
				WithBody(`{"code": 400, "detail": "detail"}`).
				WithFault(FaultConnectionResetByPeer).
				WithFixedDelay(time.Second*5),
		).
		WithPostServeAction("webhook", NewWebhook().
			WithMethod("POST").
			WithURL("http://my-target-host/callback").
			WithHeader("Content-Type", "application/json").
			WithBody(`{ "result": "SUCCESS" }`).
			WithFixedDelay(time.Second)).
		AtPriority(1).
		InScenario("Scenario").
		WhenScenarioStateIs("Started").
		WillSetStateTo("Stopped")

	rawExpectedRequestBody, err := os.ReadFile("expected-template.json")
	if err != nil {
		t.Fatalf("failed to read expected-template.json %v", err)
	}

	rawResult, err := json.Marshal(postStubRule)
	if err != nil {
		t.Fatalf("StubRole json.Marshal error: %v", err)
	}

	var expected map[string]interface{}
	err = json.Unmarshal([]byte(fmt.Sprintf(string(rawExpectedRequestBody), postStubRule.uuid, postStubRule.uuid)), &expected)
	if err != nil {
		t.Fatalf("StubRole json.Unmarshal error: %v", err)
	}

	var parsedResult map[string]interface{}
	err = json.Unmarshal(rawResult, &parsedResult)
	if err != nil {
		t.Fatalf("StubRole json.Unmarshal error: %v", err)
	}

	if !reflect.DeepEqual(parsedResult, expected) {
		t.Errorf("expected requestBody\n%v\n%v", parsedResult, expected)
	}
}
