package wiremock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

const testDataDir = "testdata"

func TestStubRule_ToJson(t *testing.T) {
	testCases := []struct {
		Name             string
		StubRule         *StubRule
		ExpectedFileName string
	}{
		{
			Name:             "BasicStubRule",
			StubRule:         NewStubRule("PATCH", URLMatching("/example")),
			ExpectedFileName: "expected-template-basic.json",
		},
		{
			Name: "StubRuleWithScenario",
			StubRule: Post(URLPathEqualTo("/example")).
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
				WithFormParameter("form1", EqualTo("value1")).
				WithFormParameter("form2", Matching("value2")).
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
				WillSetStateTo("Stopped"),
			ExpectedFileName: "expected-template-scenario.json",
		},
		{
			Name: "MustEqualToJson",
			StubRule: NewStubRule("PATCH", URLMatching("/example")).
				WithBodyPattern(MustEqualToJson(map[string]interface{}{"meta": "information"}, IgnoreArrayOrder, IgnoreExtraElements)),
			ExpectedFileName: "must-equal-to-json.json",
		},
		{
			Name: "StubRuleWithBearerToken_StartsWithMatcher",
			StubRule: Post(URLPathEqualTo("/example")).
				WithHost(EqualTo("localhost")).
				WithScheme("http").
				WithPort(8080).
				WithBearerToken(StartsWith("token")).
				WillReturnResponse(OK()),
			ExpectedFileName: "expected-template-bearer-auth-startsWith.json",
		},
		{
			Name: "StubRuleWithBearerToken_EqualToMatcher",
			StubRule: Post(URLPathEqualTo("/example")).
				WithHost(EqualTo("localhost")).
				WithScheme("http").
				WithPort(8080).
				WithBearerToken(EqualTo("token")).
				WillReturnResponse(OK()),
			ExpectedFileName: "expected-template-bearer-auth-equalTo.json",
		},
		{
			Name: "StubRuleWithBearerToken_ContainsMatcher",
			StubRule: Post(URLPathEqualTo("/example")).
				WithHost(EqualTo("localhost")).
				WithScheme("http").
				WithPort(8080).
				WithBearerToken(Contains("token")).
				WillReturnResponse(OK()),
			ExpectedFileName: "expected-template-bearer-auth-contains.json",
		},
		{
			Name: "StubRuleWithBearerToken_LogicalMatcher",
			StubRule: Post(URLPathEqualTo("/example")).
				WithHost(EqualTo("localhost")).
				WithScheme("http").
				WithPort(8080).
				WithBearerToken(EqualTo("token123").And(StartsWith("token"))).
				WillReturnResponse(OK()),
			ExpectedFileName: "expected-template-bearer-auth-logicalMatcher.json",
		},
		{
			Name: "NotLogicalMatcher",
			StubRule: Post(URLPathEqualTo("/example")).
				WithQueryParam("firstName", Not(EqualTo("John").Or(EqualTo("Jack")))).
				WillReturnResponse(OK()),
			ExpectedFileName: "not-logical-expression.json",
		},
		{
			Name: "JsonSchemaMatcher",
			StubRule: Post(URLPathEqualTo("/example")).
				WithQueryParam("firstName", MatchesJsonSchema(
					`{
  "type": "object",
  "required": [
    "name"
  ],
  "properties": {
    "name": {
      "type": "string"
    },
    "tag": {
      "type": "string"
    }
  }
}`,
					"V202012",
				)).
				WillReturnResponse(OK()),
			ExpectedFileName: "matches-Json-schema.json",
		},
		{
			Name: "URLPathTemplateMatcher",
			StubRule: Get(URLPathTemplate("/contacts/{contactId}/addresses/{addressId}")).
				WithPathParam("contactId", EqualTo("12345")).
				WithPathParam("addressId", EqualTo("99876")).
				WillReturnResponse(OK()),
			ExpectedFileName: "url-path-templating.json",
		},
		{
			Name: "StubRuleWithScenarioWithTransformerParameters",
			StubRule: Get(URLPathTemplate("/templated")).
				WillReturnResponse(
					NewResponse().
						WithStatus(http.StatusOK).
						WithTransformers("response-template").
						WithTransformerParameter("MyCustomParameter", "Parameter Value")),
			ExpectedFileName: "expected-template-transformerParameters.json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			stubRule := tc.StubRule

			rawExpectedRequestBody, err := os.ReadFile(filepath.Join(testDataDir, tc.ExpectedFileName))
			if err != nil {
				t.Fatalf("failed to read expected JSON file %s: %v", tc.ExpectedFileName, err)
			}

			var expected map[string]interface{}
			err = json.Unmarshal([]byte(fmt.Sprintf(string(rawExpectedRequestBody), stubRule.uuid, stubRule.uuid)), &expected)
			if err != nil {
				t.Fatalf("StubRule json.Unmarshal error: %v", err)
			}

			rawResult, err := json.Marshal(stubRule)
			if err != nil {
				t.Fatalf("StubRule json.Marshal error: %v", err)
			}

			var parsedResult map[string]interface{}
			err = json.Unmarshal(rawResult, &parsedResult)
			if err != nil {
				t.Fatalf("StubRule json.Unmarshal error: %v", err)
			}

			if !reflect.DeepEqual(parsedResult, expected) {
				t.Errorf("expected JSON:\n%v\nactual JSON:\n%v", parsedResult, expected)
			}
		})
	}
}
