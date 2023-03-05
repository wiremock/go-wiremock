package wiremock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		WithQueryParam("firstName", EqualTo("Jhon")).
		WithQueryParam("lastName", NotMatching("Black")).
		WithQueryParam("nickname", EqualToIgnoreCase("jhonBlack")).
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
		WillReturn(
			`{"code": 400, "detail": "detail"}`,
			map[string]string{"Content-Type": "application/json"},
			400,
		).
		WithFixedDelayMilliseconds(time.Second * 5).
		AtPriority(1).
		InScenario("Scenario").
		WhenScenarioStateIs("Started").
		WillSetStateTo("Stopped")

	rawExpectedRequestBody, err := ioutil.ReadFile("expected-template.json")
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
