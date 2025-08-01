package journal

type Cookies map[string]string

type Headers map[string]string

type Params map[string]Param

type Param struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

type GetAllRequestsResponse struct {
	Requests               []GetRequestResponse `json:"requests,omitempty"`
	Meta                   Meta                 `json:"meta,omitempty"`
	RequestJournalDisabled bool                 `json:"requestJournalDisabled,omitempty"`
}

type GetRequestResponse struct {
	ID                 string             `json:"id,omitempty"`
	Request            Request            `json:"request,omitempty"`
	ResponseDefinition ResponseDefinition `json:"responseDefinition,omitempty"`
	Response           Response           `json:"response,omitempty"`
	WasMatched         bool               `json:"wasMatched,omitempty"`
	Timing             Timing             `json:"timing,omitempty"`
	StubMapping        StubMapping        `json:"stubMapping,omitempty"`
}

type FindRequestsByCriteriaResponse struct {
	Requests []Request `json:"requests,omitempty"`
}

type FindUnmatchedRequestsResponse struct {
	Requests []Request `json:"requests,omitempty"`
}

type DeleteRequestByCriteriaResponse struct {
	Requests []GetRequestResponse `json:"serveEvents,omitempty"`
}

type Request struct {
	URL                 string  `json:"url,omitempty"`
	AbsoluteURL         string  `json:"absoluteUrl,omitempty"`
	Method              string  `json:"method,omitempty"`
	ClientIP            string  `json:"clientIp,omitempty"`
	Headers             Headers `json:"headers,omitempty"`
	Cookies             Cookies `json:"cookies,omitempty"`
	BrowserProxyRequest bool    `json:"browserProxyRequest,omitempty"`
	LoggedDate          int64   `json:"loggedDate,omitempty"`
	BodyAsBase64        string  `json:"bodyAsBase64,omitempty"`
	Body                string  `json:"body,omitempty"`
	Protocol            string  `json:"protocol,omitempty"`
	Scheme              string  `json:"scheme,omitempty"`
	LoggedDateString    string  `json:"loggedDateString,omitempty"`
	Host                string  `json:"host,omitempty"`
	Port                int64   `json:"port,omitempty"`
	QueryParams         Params  `json:"queryParams,omitempty"`
	FormParams          Params  `json:"formParams,omitempty"`
}

type ResponseDefinition struct {
	Headers            Headers `json:"headers,omitempty"`
	Body               string  `json:"body,omitempty"`
	Status             int64   `json:"status,omitempty"`
	FromConfiguredStub bool    `json:"fromConfiguredStub,omitempty"`
}

type Response struct {
	Headers      Headers `json:"headers,omitempty"`
	BodyAsBase64 string  `json:"bodyAsBase64,omitempty"`
	Body         string  `json:"body,omitempty"`
	Status       int64   `json:"status,omitempty"`
}

type Timing struct {
	ServeTime        int64 `json:"serveTime,omitempty"`
	TotalTime        int64 `json:"totalTime,omitempty"`
	ProcessTime      int64 `json:"processTime,omitempty"`
	ResponseSendTime int64 `json:"responseSendTime,omitempty"`
	AddedDelay       int64 `json:"addedDelay,omitempty"`
}

type StubMapping struct {
	ID       string             `json:"id,omitempty"`
	Request  StubMappingRequest `json:"request,omitempty"`
	Response ResponseDefinition `json:"response,omitempty"`
	UUID     string             `json:"uuid,omitempty"`
}

type StubMappingRequest struct {
	Method          string `json:"method,omitempty"`
	URL             string `json:"url,omitempty"`
	URLPattern      string `json:"urlPattern,omitempty"`
	URLPath         string `json:"urlPath,omitempty"`
	URLPathPattern  string `json:"urlPathPattern,omitempty"`
	URLPathTemplate string `json:"urlPathTemplate,omitempty"`
}

type Meta struct {
	Total int64 `json:"total,omitempty"`
}

type ServeEvent struct {
	Request            Request            `json:"request,omitempty"`
	StubMapping        StubMapping        `json:"stubMapping,omitempty"`
	WasMatched         bool               `json:"wasMatched,omitempty"`
	Response           Response           `json:"response,omitempty"`
	Timing             Timing             `json:"timing,omitempty"`
	ID                 string             `json:"id,omitempty"`
	ResponseDefinition ResponseDefinition `json:"responseDefinition,omitempty"`
}
