package wiremock

import "encoding/json"

type MultiValueMatcher struct {
	strategy MultiValueMatchingStrategy
	matchers []BasicParamMatcher
}

// MarshalJSON returns the JSON encoding of the matcher.
func (m MultiValueMatcher) MarshalJSON() ([]byte, error) {
	jsonMap := map[string]interface{}{
		string(m.strategy): m.matchers,
	}

	return json.Marshal(jsonMap)
}

// HasExactly returns a matcher that matches when the parameter has exactly the specified values.
func HasExactly(matchers ...BasicParamMatcher) MultiValueMatcher {
	return MultiValueMatcher{
		strategy: ParamHasExactly,
		matchers: matchers,
	}
}

// Includes returns a matcher that matches when the parameter includes the specified values.
func Includes(matchers ...BasicParamMatcher) MultiValueMatcher {
	return MultiValueMatcher{
		strategy: ParamIncludes,
		matchers: matchers,
	}
}
