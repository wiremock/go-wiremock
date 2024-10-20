package wiremock

import (
	"encoding/json"
)

type LogicalMatcher struct {
	operator string
	operands []BasicParamMatcher
}

// MarshalJSON returns the JSON encoding of the matcher.
func (m LogicalMatcher) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.ParseMatcher())
}

// ParseMatcher returns the map representation of the structure.
func (m LogicalMatcher) ParseMatcher() map[string]interface{} {
	if m.operator == "not" {
		return map[string]interface{}{
			m.operator: m.operands[0],
		}
	}

	return map[string]interface{}{
		m.operator: m.operands,
	}
}

// Or returns a logical OR of the current matcher and the given matcher.
func (m LogicalMatcher) Or(matcher BasicParamMatcher) BasicParamMatcher {
	if m.operator == "or" {
		m.operands = append(m.operands, matcher)
		return m
	}

	return Or(m, matcher)
}

// And returns a logical AND of the current matcher and the given matcher.
func (m LogicalMatcher) And(matcher BasicParamMatcher) BasicParamMatcher {
	if m.operator == "and" {
		m.operands = append(m.operands, matcher)
		return m
	}

	return And(m, matcher)
}

// Or returns a logical OR of the list of matchers.
func Or(matchers ...BasicParamMatcher) LogicalMatcher {
	return LogicalMatcher{
		operator: "or",
		operands: matchers,
	}
}

// And returns a logical AND of the list of matchers.
func And(matchers ...BasicParamMatcher) LogicalMatcher {
	return LogicalMatcher{
		operator: "and",
		operands: matchers,
	}
}

// Not returns a logical NOT of the given matcher. Required wiremock version >= 3.0.0
func Not(matcher BasicParamMatcher) LogicalMatcher {
	return LogicalMatcher{
		operator: "not",
		operands: []BasicParamMatcher{matcher},
	}
}
