package grpc

import (
	"encoding/json"
	"fmt"
	"github.com/wiremock/go-wiremock"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const testDataDir = "testdata"

func TestStubRule_ToJson(t *testing.T) {
	testCases := []struct {
		Name             string
		StubRule         *wiremock.StubRule
		ExpectedFileName string
	}{
		{
			Name: "Error",
			StubRule: Method("greeting").
				WillReturn(Error(StatusUnavailable, "Unavailable")).
				Build("com.example.grpc.GreetingService"),
			ExpectedFileName: "error.json",
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
			err = json.Unmarshal([]byte(fmt.Sprintf(string(rawExpectedRequestBody), stubRule.UUID(), stubRule.UUID())), &expected)
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
