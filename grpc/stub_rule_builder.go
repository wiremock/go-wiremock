package grpc

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/wiremock/go-wiremock"
)

type StubRuleBuilder struct {
	method          string
	responseBuilder *ResponseBuilder
	bodyPatterns    []wiremock.BasicParamMatcher
}

// Method creates a new instance of the StubRuleBuilder with grpc method.
func Method(method string) *StubRuleBuilder {
	return &StubRuleBuilder{method: method}
}

// EqualToMessage grpc alias for wiremock.MustEqualToJson. May panic if there are problems with marshaling.
func EqualToMessage(message proto.Message) wiremock.BasicParamMatcher {
	data, err := protojson.Marshal(message)

	if err != nil {
		panic(fmt.Sprintf("failed to marshal proto message: %v", err))
	}

	return wiremock.EqualToJson(string(data))
}

// WithRequestMessage adds a request message matcher to the stub rule.
func (s *StubRuleBuilder) WithRequestMessage(matcher wiremock.BasicParamMatcher) *StubRuleBuilder {
	s.bodyPatterns = append(s.bodyPatterns, matcher)
	return s
}

// WillReturn sets the response for the stub rule.
func (s *StubRuleBuilder) WillReturn(responseBuilder *ResponseBuilder) *StubRuleBuilder {
	s.responseBuilder = responseBuilder
	return s
}

// Build builds a new instance of the StubRule.
func (s *StubRuleBuilder) Build(serviceName string) *wiremock.StubRule {
	stubRule := wiremock.Post(wiremock.URLPathEqualTo(fmt.Sprintf("/%s/%s", serviceName, s.method)))
	for _, bodyPattern := range s.bodyPatterns {
		stubRule = stubRule.WithBodyPattern(bodyPattern)
	}

	return stubRule.WillReturnResponse(s.responseBuilder.Build())
}
