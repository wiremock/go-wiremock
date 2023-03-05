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

type MultipartPattern struct {
	matchingType MultipartMatchingType
	headers      map[string]ParamMatcherInterface
	bodyPatterns []ParamMatcher
}

func NewMultipartPattern() *MultipartPattern {
	return &MultipartPattern{
		matchingType: MultipartMatchingTypeAny,
	}
}

func (m *MultipartPattern) WithName(name string) *MultipartPattern {
	if m.headers == nil {
		m.headers = map[string]ParamMatcherInterface{}
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

func (m *MultipartPattern) WithBodyPattern(matcher ParamMatcher) *MultipartPattern {
	m.bodyPatterns = append(m.bodyPatterns, matcher)
	return m
}

func (m *MultipartPattern) WithHeader(header string, matcher ParamMatcherInterface) *MultipartPattern {
	if m.headers == nil {
		m.headers = map[string]ParamMatcherInterface{}
	}

	m.headers[header] = matcher
	return m
}

// MarshalJSON gives valid JSON or error.
func (m *MultipartPattern) MarshalJSON() ([]byte, error) {
	multipart := map[string]interface{}{
		"matchingType": m.matchingType,
	}

	if len(m.bodyPatterns) > 0 {
		bodyPatterns := make([]map[string]interface{}, len(m.bodyPatterns))
		for i, bodyPattern := range m.bodyPatterns {
			bodyPatterns[i] = map[string]interface{}{
				string(bodyPattern.Strategy()): bodyPattern.Value(),
			}

			for flag, value := range bodyPattern.flags {
				bodyPatterns[i][flag] = value
			}
		}
		multipart["bodyPatterns"] = bodyPatterns
	}

	if len(m.headers) > 0 {
		headers := make(map[string]map[string]interface{}, len(m.headers))
		for key, header := range m.headers {
			headers[key] = map[string]interface{}{
				string(header.Strategy()): header.Value(),
			}

			for flag, value := range header.Flags() {
				headers[key][flag] = value
			}
		}
		multipart["headers"] = headers
	}

	return json.Marshal(multipart)
}
