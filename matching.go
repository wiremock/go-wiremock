package wiremock

// Types of params matching.
const (
	ParamEqualTo         ParamMatchingStrategy = "equalTo"
	ParamMatches         ParamMatchingStrategy = "matches"
	ParamContains        ParamMatchingStrategy = "contains"
	ParamEqualToXml      ParamMatchingStrategy = "equalToXml"
	ParamEqualToJson     ParamMatchingStrategy = "equalToJson"
	ParamMatchesXPath    ParamMatchingStrategy = "matchesXPath"
	ParamMatchesJsonPath ParamMatchingStrategy = "matchesJsonPath"
	ParamAbsent          ParamMatchingStrategy = "absent"
	ParamDoesNotMatch    ParamMatchingStrategy = "doesNotMatch"
)

// Types of url matching.
const (
	URLEqualToRule      URLMatchingStrategy = "url"
	URLPathEqualToRule  URLMatchingStrategy = "urlPath"
	URLPathMatchingRule URLMatchingStrategy = "urlPathPattern"
	URLMatchingRule     URLMatchingStrategy = "urlPattern"
)

// Type of less strict matching flags.
const (
	IgnoreArrayOrder    EqualFlag = "ignoreArrayOrder"
	IgnoreExtraElements EqualFlag = "ignoreExtraElements"
)

// EqualFlag is enum of less strict matching flag.
type EqualFlag string

// URLMatchingStrategy is enum url matching type.
type URLMatchingStrategy string

// ParamMatchingStrategy is enum params matching type.
type ParamMatchingStrategy string

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

// ParamMatcher is structure for defining the type of params.
type ParamMatcher struct {
	strategy ParamMatchingStrategy
	value    string
	flags    map[string]bool
}

// Strategy returns ParamMatchingStrategy of ParamMatcher.
func (m ParamMatcher) Strategy() ParamMatchingStrategy {
	return m.strategy
}

// Value returns value of ParamMatcher.
func (m ParamMatcher) Value() string {
	return m.value
}

// Flags return value of ParamMatcher.
func (m ParamMatcher) Flags() map[string]bool {
	return m.flags
}

// EqualTo returns ParamMatcher with ParamEqualTo matching strategy.
func EqualTo(param string) ParamMatcher {
	return ParamMatcher{
		strategy: ParamEqualTo,
		value:    param,
	}
}

// EqualToIgnoreCase returns ParamMatcher with ParamEqualToIgnoreCase matching strategy
func EqualToIgnoreCase(param string) ParamMatcher {
	return ParamMatcher{
		strategy: ParamEqualTo,
		value:    param,
		flags: map[string]bool{
			"caseInsensitive": true,
		},
	}
}

// Matching returns ParamMatcher with ParamMatches matching strategy.
func Matching(param string) ParamMatcher {
	return ParamMatcher{
		strategy: ParamMatches,
		value:    param,
	}
}

// Contains returns ParamMatcher with ParamContains matching strategy.
func Contains(param string) ParamMatcher {
	return ParamMatcher{
		strategy: ParamContains,
		value:    param,
	}
}

// EqualToXml returns ParamMatcher with ParamEqualToXml matching strategy.
func EqualToXml(param string) ParamMatcher {
	return ParamMatcher{
		strategy: ParamEqualToXml,
		value:    param,
	}
}

// EqualToJson returns ParamMatcher with ParamEquignoreArrayOrderalToJson matching strategy.
func EqualToJson(param string, flags ...EqualFlag) ParamMatcher {
	mflags := make(map[string]bool, len(flags))
	for _, flag := range flags {
		mflags[string(flag)] = true
	}

	return ParamMatcher{
		strategy: ParamEqualToJson,
		value:    param,
		flags:    mflags,
	}
}

// MatchingXPath returns ParamMatcher with ParamMatchesXPath matching strategy.
func MatchingXPath(param string) ParamMatcher {
	return ParamMatcher{
		strategy: ParamMatchesXPath,
		value:    param,
	}
}

// MatchingJsonPath returns ParamMatcher with ParamMatchesJsonPath matching strategy.
func MatchingJsonPath(param string) ParamMatcher {
	return ParamMatcher{
		strategy: ParamMatchesJsonPath,
		value:    param,
	}
}

// NotMatching returns ParamMatcher with ParamDoesNotMatch matching strategy.
func NotMatching(param string) ParamMatcher {
	return ParamMatcher{
		strategy: ParamDoesNotMatch,
		value:    param,
	}
}

func Absent() ParamMatcher {
	return ParamMatcher{
		strategy: ParamAbsent,
		value:    "",
		flags: map[string]bool{
			string(ParamAbsent): true,
		},
	}
}
