package wiremock

// Types of params matching.
const (
	ParamEqualTo           ParamMatchingStrategy = "equalTo"
	ParamMatches           ParamMatchingStrategy = "matches"
	ParamContains          ParamMatchingStrategy = "contains"
	ParamEqualToXml        ParamMatchingStrategy = "equalToXml"
	ParamEqualToJson       ParamMatchingStrategy = "equalToJson"
	ParamMatchesXPath      ParamMatchingStrategy = "matchesXPath"
	ParamMatchesJsonPath   ParamMatchingStrategy = "matchesJsonPath"
	ParamAbsent            ParamMatchingStrategy = "absent"
	ParamDoesNotMatch      ParamMatchingStrategy = "doesNotMatch"
	ParamDoesNotContains   ParamMatchingStrategy = "doesNotContain"
	ParamMatchesJsonSchema ParamMatchingStrategy = "matchesJsonSchema"
)

// Types of url matching.
const (
	URLEqualToRule      URLMatchingStrategy = "url"
	URLPathEqualToRule  URLMatchingStrategy = "urlPath"
	URLPathMatchingRule URLMatchingStrategy = "urlPathPattern"
	URLMatchingRule     URLMatchingStrategy = "urlPattern"
	URLPathTemplateRule URLMatchingStrategy = "urlPathTemplate"
)

// Type of less strict matching flags.
const (
	IgnoreArrayOrder    EqualFlag = "ignoreArrayOrder"
	IgnoreExtraElements EqualFlag = "ignoreExtraElements"
)

const (
	ParamHasExactly MultiValueMatchingStrategy = "hasExactly"
	ParamIncludes   MultiValueMatchingStrategy = "includes"
)

// EqualFlag is enum of less strict matching flag.
type EqualFlag string

// URLMatchingStrategy is enum url matching type.
type URLMatchingStrategy string

// ParamMatchingStrategy is enum params matching type.
type ParamMatchingStrategy string

// MultiValueMatchingStrategy is enum multi value matching type.
type MultiValueMatchingStrategy string
