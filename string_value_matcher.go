package wiremock

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type MatcherInterface interface {
	json.Marshaler
	ParseMatcher() map[string]interface{}
}

type BasicParamMatcher interface {
	json.Marshaler
	ParseMatcher() map[string]interface{}
	Or(stringMatcher BasicParamMatcher) BasicParamMatcher
	And(stringMatcher BasicParamMatcher) BasicParamMatcher
}

type StringValueMatcher struct {
	strategy ParamMatchingStrategy
	value    string
	flags    []string
}

// MarshalJSON returns the JSON encoding of the matcher.
func (m StringValueMatcher) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.ParseMatcher())
}

// ParseMatcher returns the map representation of the structure.
func (m StringValueMatcher) ParseMatcher() map[string]interface{} {
	jsonMap := make(map[string]interface{}, 1+len(m.flags))
	if m.strategy != "" {
		jsonMap[string(m.strategy)] = m.value
	}

	for _, flag := range m.flags {
		jsonMap[flag] = true
	}

	return jsonMap
}

// Or returns a logical OR of the two matchers.
func (m StringValueMatcher) Or(matcher BasicParamMatcher) BasicParamMatcher {
	return Or(m, matcher)
}

// And returns a logical AND of the two matchers.
func (m StringValueMatcher) And(matcher BasicParamMatcher) BasicParamMatcher {
	return And(m, matcher)
}

// addPrefixToMatcher adds prefix to matcher.
// In case of "contains", "absent", "doesNotContain" prefix is not added as it doesn't affect the match result
func (m StringValueMatcher) addPrefixToMatcher(prefix string) BasicParamMatcher {
	switch m.strategy {
	case ParamEqualTo, ParamEqualToJson, ParamEqualToXml, ParamMatchesJsonPath, ParamMatchesXPath:
		m.value = prefix + m.value
	case ParamMatches, ParamDoesNotMatch:
		if regexContainsStartAnchor(m.value) {
			m.value = m.value[1:]
		}
		m.value = fmt.Sprintf("^%s", prefix) + m.value
	}
	return m
}

// NewStringValueMatcher creates a new StringValueMatcher.
func NewStringValueMatcher(strategy ParamMatchingStrategy, value string, flags ...string) StringValueMatcher {
	return StringValueMatcher{
		strategy: strategy,
		value:    value,
		flags:    flags,
	}
}

// EqualTo returns a matcher that matches when the parameter equals the specified value.
func EqualTo(value string) BasicParamMatcher {
	return NewStringValueMatcher(ParamEqualTo, value)
}

// EqualToIgnoreCase returns a matcher that matches when the parameter equals the specified value, ignoring case.
func EqualToIgnoreCase(value string) BasicParamMatcher {
	return NewStringValueMatcher(ParamEqualTo, value, "caseInsensitive")
}

// Matching returns a matcher that matches when the parameter matches the specified regular expression.
func Matching(param string) BasicParamMatcher {
	return NewStringValueMatcher(ParamMatches, param)
}

// EqualToXml returns a matcher that matches when the parameter is equal to the specified XML.
func EqualToXml(param string) BasicParamMatcher {
	return NewStringValueMatcher(ParamEqualToXml, param)
}

// EqualToJson returns a matcher that matches when the parameter is equal to the specified JSON.
func EqualToJson(param string, equalJsonFlags ...EqualFlag) BasicParamMatcher {
	flags := make([]string, len(equalJsonFlags))
	for i, flag := range equalJsonFlags {
		flags[i] = string(flag)
	}

	return NewStringValueMatcher(ParamEqualToJson, param, flags...)
}

// MustEqualToJson returns a matcher that matches when the parameter is equal to the specified JSON.
// This method panics if param cannot be marshaled to JSON.
func MustEqualToJson(param any, equalJsonFlags ...EqualFlag) BasicParamMatcher {
	if str, ok := param.(string); ok {
		return EqualToJson(str, equalJsonFlags...)
	}

	if jsonParam, err := json.Marshal(param); err != nil {
		panic(fmt.Sprintf("Unable to marshal parameter to JSON: %v", err))
	} else {
		return EqualToJson(string(jsonParam), equalJsonFlags...)
	}
}

// MatchingXPath returns a matcher that matches when the parameter matches the specified XPath.
func MatchingXPath(param string) BasicParamMatcher {
	return NewStringValueMatcher(ParamMatchesXPath, param)
}

// MatchingJsonPath returns a matcher that matches when the parameter matches the specified JSON path.
func MatchingJsonPath(param string) BasicParamMatcher {
	return NewStringValueMatcher(ParamMatchesJsonPath, param)
}

// NotMatching returns a matcher that matches when the parameter does not match the specified regular expression.
func NotMatching(param string) BasicParamMatcher {
	return NewStringValueMatcher(ParamDoesNotMatch, param)
}

// Absent returns a matcher that matches when the parameter is absent.
func Absent() BasicParamMatcher {
	return StringValueMatcher{
		flags: []string{string(ParamAbsent)},
	}
}

// Contains returns a matcher that matches when the parameter contains the specified value.
func Contains(param string) BasicParamMatcher {
	return NewStringValueMatcher(ParamContains, param)
}

// NotContains returns a matcher that matches when the parameter does not contain the specified value.
func NotContains(param string) BasicParamMatcher {
	return NewStringValueMatcher(ParamDoesNotContains, param)
}

// StartsWith returns a matcher that matches when the parameter starts with the specified prefix.
// Matches also when prefix alone is the whole expression
func StartsWith(prefix string) BasicParamMatcher {
	regex := fmt.Sprintf(`^%s\s*\S*`, regexp.QuoteMeta(prefix))
	return NewStringValueMatcher(ParamMatches, regex)
}

type JSONSchemaMatcher struct {
	StringValueMatcher
	schemaVersion string
}

// MatchesJsonSchema returns a matcher that matches when the parameter matches the specified JSON schema.
// Required wiremock version >= 3.0.0
func MatchesJsonSchema(schema string, schemaVersion string) BasicParamMatcher {
	return JSONSchemaMatcher{
		StringValueMatcher: NewStringValueMatcher(ParamMatchesJsonSchema, schema),
		schemaVersion:      schemaVersion,
	}
}

// MarshalJSON returns the JSON encoding of the matcher.
func (m JSONSchemaMatcher) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		string(m.strategy): m.value,
		"schemaVersion":    m.schemaVersion,
	})
}

func regexContainsStartAnchor(regex string) bool {
	return len(regex) > 0 && regex[0] == '^'
}
