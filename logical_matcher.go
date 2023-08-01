package wiremock

import (
	"encoding/json"
)

type LogicalMatcher struct {
	operator string
	operands []BasicParamMatcher
}

func (m LogicalMatcher) MarshalJSON() ([]byte, error) {
	jsonMap := map[string]interface{}{
		m.operator: m.operands,
	}

	return json.Marshal(jsonMap)
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
