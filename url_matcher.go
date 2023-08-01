package wiremock

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
