package wiremock

import (
	"testing"
)

func TestStringValueMatcher_AddPrefixToMatcher(t *testing.T) {
	testCases := []struct {
		name     string
		strategy ParamMatchingStrategy
		value    string
		prefix   string
		expected string
	}{
		{
			name:     "EqualTo - Add prefix to the value",
			strategy: ParamEqualTo,
			value:    "abc",
			prefix:   "pre_",
			expected: "pre_abc",
		},
		{
			name:     "Matches - Add prefix to regex value without start anchor",
			strategy: ParamMatches,
			value:    "abc",
			prefix:   "pre_",
			expected: "^pre_abc",
		},
		{
			name:     "Matches - Add prefix to regex value with start anchor",
			strategy: ParamMatches,
			value:    "^abc",
			prefix:   "pre_",
			expected: "^pre_abc",
		},
		{
			name:     "Matches - Add prefix to regex value with end anchor",
			strategy: ParamMatches,
			value:    "t?o?ken$",
			prefix:   "pre_",
			expected: "^pre_t?o?ken$",
		},
		{
			name:     "Matches - Should add prefix to wildcard regex",
			strategy: ParamMatches,
			value:    ".*",
			prefix:   "pre_",
			expected: "^pre_.*",
		},
		{
			name:     "Matches - Should add prefix to empty regex",
			strategy: ParamMatches,
			value:    "",
			prefix:   "pre_",
			expected: "^pre_",
		},
		{
			name:     "DoesNotMatch - Add prefix to regex value without start anchor",
			strategy: ParamDoesNotMatch,
			value:    "abc",
			prefix:   "pre_",
			expected: "^pre_abc",
		},
		{
			name:     "DoesNotMatch - Add prefix to regex value with start anchor",
			strategy: ParamDoesNotMatch,
			value:    "^abc",
			prefix:   "pre_",
			expected: "^pre_abc",
		},
		{
			name:     "DoesNotMatch - wildcard regex",
			strategy: ParamDoesNotMatch,
			value:    ".*",
			prefix:   "pre_",
			expected: "^pre_.*",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			matcher := StringValueMatcher{
				strategy: tc.strategy,
				value:    tc.value,
			}

			modifiedMatcher := matcher.addPrefixToMatcher(tc.prefix)

			if modifiedMatcher.(StringValueMatcher).value != tc.expected {
				t.Errorf("Expected: %s, Got: %s", tc.expected, modifiedMatcher.(StringValueMatcher).value)
			}
		})
	}
}
