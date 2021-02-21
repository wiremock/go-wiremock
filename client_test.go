package wiremock

import (
	"encoding/json"
	"testing"
)

func TestStubRule_ToJson(t *testing.T) {
	newStubRule := NewStubRule("PATCH", URLMatching("/example"))
	expectedRequestBody := `{"request":{"method":"PATCH","urlPattern":"/example"},"response":{"status":200}}`
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
		WithBodyPattern(EqualToJson(`{"meta": "information"}`)).
		WithBodyPattern(Contains("information")).
		WithHeader("x-session", Matching("^\\S+@\\S+$")).
		WithCookie("session", EqualToXml("<xml>")).
		WillReturn(
			`{"code": 400, "detail": "detail"}`,
			map[string]string{"Content-Type": "application/json"},
			400,
		).
		AtPriority(1).
		InScenario("Scenario").
		WhenScenarioStateIs("Started").
		WillSetStateTo("Stopped")

	expectedRequestBody = `{"priority":1,"scenarioName":"Scenario","requiredScenarioState":"Started","newScenarioState":"Stopped",` +
		`"request":{"bodyPatterns":[{"equalToJson":"{\"meta\": \"information\"}"},{"contains":"information"}],` +
		`"cookies":{"session":{"equalToXml":"\u003cxml\u003e"}},` +
		`"headers":{"x-session":{"matches":"^\\S+@\\S+$"}},` +
		`"method":"POST","queryParameters":{"firstName":{"equalTo":"Jhon"},"lastName":{"doesNotMatch":"Black"}},"urlPath":"/example"},` +
		`"response":{"body":"{\"code\": 400, \"detail\": \"detail\"}","headers":{"Content-Type":"application/json"},"status":400}}`
	result, err = json.Marshal(postStubRule)

	if err != nil {
		t.Fatalf("StubRole json.Marshal error: %v", err)
	}
	if string(result) != expectedRequestBody {
		t.Errorf("expected requestBody %q; got %q", expectedRequestBody, string(result))
	}
}
