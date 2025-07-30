package wiremock

// URLMatcherInterface is pair URLMatchingStrategy and string matched value
type URLMatcherInterface interface {
	Strategy() URLMatchingStrategy
	Value() string
}

// URLMatcher is structure for defining the type of url matching.
type URLMatcher struct {
	strategy URLMatchingStrategy
	value    string
}

// Strategy returns URLMatchingStrategy of URLMatcher.
func (m URLMatcher) Strategy() URLMatchingStrategy {
	return m.strategy
}

// Value returns value of URLMatcher.
func (m URLMatcher) Value() string {
	return m.value
}

// URLEqualTo returns URLMatcher with URLEqualToRule matching strategy.
func URLEqualTo(url string) URLMatcher {
	return URLMatcher{
		strategy: URLEqualToRule,
		value:    url,
	}
}

// URLPathEqualTo returns URLMatcher with URLPathEqualToRule matching strategy.
func URLPathEqualTo(url string) URLMatcher {
	return URLMatcher{
		strategy: URLPathEqualToRule,
		value:    url,
	}
}

// URLPathMatching returns URLMatcher with URLPathMatchingRule matching strategy.
func URLPathMatching(url string) URLMatcher {
	return URLMatcher{
		strategy: URLPathMatchingRule,
		value:    url,
	}
}

// URLMatching returns URLMatcher with URLMatchingRule matching strategy.
func URLMatching(url string) URLMatcher {
	return URLMatcher{
		strategy: URLMatchingRule,
		value:    url,
	}
}

// URLPathTemplate URL paths can be matched using URI templates, conforming to the same subset of the URI template standard as used in OpenAPI.
// Path variable matchers can also be used in the same manner as query and form parameters.
// Required wiremock >= 3.0.0
// Example: /contacts/{contactId}/addresses/{addressId}
func URLPathTemplate(url string) URLMatcher {
	return URLMatcher{
		strategy: URLPathTemplateRule,
		value:    url,
	}
}
