package wiremock

import "encoding/json"

type BasicParamMatcher interface {
	json.Marshaler
	Or(stringMatcher BasicParamMatcher) BasicParamMatcher
	And(stringMatcher BasicParamMatcher) BasicParamMatcher
}

type StringValueMatcher struct {
	strategy ParamMatchingStrategy
	value    string
	flags    []string
}

func (m StringValueMatcher) MarshalJSON() ([]byte, error) {
	jsonMap := make(map[string]interface{}, 1+len(m.flags))
	if m.strategy != "" {
		jsonMap[string(m.strategy)] = m.value
	}
	for _, flag := range m.flags {
		jsonMap[flag] = true
	}

	return json.Marshal(jsonMap)
}

// Or returns a logical OR of the two matchers.
func (m StringValueMatcher) Or(matcher BasicParamMatcher) BasicParamMatcher {
	return Or(m, matcher)
}

// And returns a logical AND of the two matchers.
func (m StringValueMatcher) And(matcher BasicParamMatcher) BasicParamMatcher {
	return And(m, matcher)
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
