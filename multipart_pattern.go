package wiremock

import (
	"encoding/json"
	"fmt"
)

const (
	MultipartMatchingTypeAny = "ANY"
	MultipartMatchingTypeAll = "ALL"
)

type MultipartMatchingType string

type MultipartPatternInterface interface {
	json.Marshaler
	ParseMultipartPattern() map[string]interface{}
}

type MultipartPattern struct {
	matchingType MultipartMatchingType
	headers      map[string]MatcherInterface
	bodyPatterns []BasicParamMatcher
}

func NewMultipartPattern() *MultipartPattern {
	return &MultipartPattern{
		matchingType: MultipartMatchingTypeAny,
	}
}

func (m *MultipartPattern) WithName(name string) *MultipartPattern {
	if m.headers == nil {
		m.headers = map[string]MatcherInterface{}
	}

	m.headers["Content-Disposition"] = Contains(fmt.Sprintf(`name="%s"`, name))
	return m
}

func (m *MultipartPattern) WithMatchingType(matchingType MultipartMatchingType) *MultipartPattern {
	m.matchingType = matchingType
	return m
}

func (m *MultipartPattern) WithAllMatchingType() *MultipartPattern {
	m.matchingType = MultipartMatchingTypeAll
	return m
}

func (m *MultipartPattern) WithAnyMatchingType() *MultipartPattern {
	m.matchingType = MultipartMatchingTypeAny
	return m
}

func (m *MultipartPattern) WithBodyPattern(matcher BasicParamMatcher) *MultipartPattern {
	m.bodyPatterns = append(m.bodyPatterns, matcher)
	return m
}

func (m *MultipartPattern) WithHeader(header string, matcher MatcherInterface) *MultipartPattern {
	if m.headers == nil {
		m.headers = map[string]MatcherInterface{}
	}

	m.headers[header] = matcher
	return m
}

// MarshalJSON gives valid JSON or error.
func (m *MultipartPattern) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.ParseMultipartPattern())
}

func (m *MultipartPattern) ParseMultipartPattern() map[string]interface{} {
	multipart := map[string]interface{}{
		"matchingType": m.matchingType,
	}

	if len(m.bodyPatterns) > 0 {
		multipart["bodyPatterns"] = m.bodyPatterns
	}

	if len(m.headers) > 0 {
		multipart["headers"] = m.headers
	}

	return multipart
}
