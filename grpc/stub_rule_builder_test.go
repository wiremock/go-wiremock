package grpc

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"

	"github.com/wiremock/go-wiremock"
	"github.com/wiremock/go-wiremock/grpc/testdata"
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
				WillReturn(Error(codes.Unavailable, "Unavailable")).
				Build("com.example.grpc.GreetingService"),
			ExpectedFileName: "error.json",
		},
		{
			Name: "OK Message",
			StubRule: Method("greeting").
				WithRequestMessage(EqualToMessage(&testdata.GreetingRequest{Name: "Tom"})).
				WillReturn(Message(&testdata.GreetingResponse{Greeting: "Hello Tom"})).
				Build("com.example.grpc.GreetingService"),
			ExpectedFileName: "ok.json",
		},
		{
			Name: "OK JSON",
			StubRule: Method("greeting").
				WithRequestMessage(wiremock.EqualToJson(`{"name":"Tom"}`)).
				WillReturn(JSON(`{"greeting":"Hello Tom"}`)).
				Build("com.example.grpc.GreetingService"),
			ExpectedFileName: "ok.json",
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
