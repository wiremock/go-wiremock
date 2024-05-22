package wiremock

type ServeEvents struct {
	ServeEvents []ServeEvent `json:"serveEvents,omitempty"`
}

type ServeEvent struct {
	RequestJournalDisabled bool             `json:"requestJournalDisabled,omitempty"`
	Meta                   Meta             `json:"meta,omitempty"`
	Requests               []RequestElement `json:"requests,omitempty"`

	Request            ServeEventRequest       `json:"request,omitempty"`
	StubMapping        StubMapping             `json:"stubMapping,omitempty"`
	WasMatched         bool                    `json:"wasMatched,omitempty"`
	Response           PurpleResponse          `json:"response,omitempty"`
	Timing             Timing                  `json:"timing,omitempty"`
	ID                 string                  `json:"id,omitempty"`
	ResponseDefinition ResponseDefinitionClass `json:"responseDefinition,omitempty"`
}

type Meta struct {
	Total int64 `json:"total,omitempty"`
}

type RequestElement struct {
	Request            ServeEventRequest       `json:"request,omitempty"`
	StubMapping        StubMapping             `json:"stubMapping,omitempty"`
	WasMatched         bool                    `json:"wasMatched,omitempty"`
	Response           PurpleResponse          `json:"response,omitempty"`
	Timing             Timing                  `json:"timing,omitempty"`
	ID                 string                  `json:"id,omitempty"`
	ResponseDefinition ResponseDefinitionClass `json:"responseDefinition,omitempty"`
}

type ServeEventRequest struct {
	Headers             RequestHeaders `json:"headers,omitempty"`
	Method              string         `json:"method,omitempty"`
	Scheme              string         `json:"scheme,omitempty"`
	BrowserProxyRequest bool           `json:"browserProxyRequest,omitempty"`
	QueryParams         Cookies        `json:"queryParams,omitempty"`
	AbsoluteURL         string         `json:"absoluteUrl,omitempty"`
	Body                string         `json:"body,omitempty"`
	LoggedDate          int64          `json:"loggedDate,omitempty"`
	URL                 string         `json:"url,omitempty"`
	Cookies             Cookies        `json:"cookies,omitempty"`
	Protocol            string         `json:"protocol,omitempty"`
	Port                int64          `json:"port,omitempty"`
	LoggedDateString    string         `json:"loggedDateString,omitempty"`
	ClientIP            string         `json:"clientIp,omitempty"`
	Host                string         `json:"host,omitempty"`
	BodyAsBase64        string         `json:"bodyAsBase64,omitempty"`
}

type Cookies struct {
}

type RequestHeaders struct {
	B3                 string `json:"B3,omitempty"`
	GrpcAcceptEncoding string `json:"Grpc-Accept-Encoding,omitempty"`
	Connection         string `json:"Connection,omitempty"`
	UserAgent          string `json:"User-Agent,omitempty"`
	Host               string `json:"Host,omitempty"`
	AcceptEncoding     string `json:"Accept-Encoding,omitempty"`
	ContentLength      string `json:"Content-Length,omitempty"`
	XRealIP            string `json:"X-Real-IP,omitempty"`
	Traceparent        string `json:"Traceparent,omitempty"`
	ContentType        string `json:"Content-Type,omitempty"`
}

type PurpleResponse struct {
	Headers      PurpleHeaders `json:"headers,omitempty"`
	BodyAsBase64 string        `json:"bodyAsBase64,omitempty"`
	Body         string        `json:"body,omitempty"`
	Status       int64         `json:"status,omitempty"`
}

type PurpleHeaders struct {
	MatchedStubID   string `json:"Matched-Stub-Id,omitempty"`
	MatchedStubName string `json:"Matched-Stub-Name,omitempty"`
	ContentType     string `json:"Content-Type,omitempty"`
}

type ResponseDefinitionClass struct {
	Headers ResponseDefinitionHeaders `json:"headers,omitempty"`
	Body    string                    `json:"body,omitempty"`
	Status  int64                     `json:"status,omitempty"`
}

type ResponseDefinitionHeaders struct {
	ContentType string `json:"Content-Type,omitempty"`
}

type StubMapping struct {
	Request  StubMappingRequest      `json:"request,omitempty"`
	Metadata Metadata                `json:"metadata,omitempty"`
	Response ResponseDefinitionClass `json:"response,omitempty"`
	Name     string                  `json:"name,omitempty"`
	ID       string                  `json:"id,omitempty"`
	UUID     string                  `json:"uuid,omitempty"`
}

type Metadata struct {
	Description string `json:"description,omitempty"`
}

type StubMappingRequest struct {
	Method  string `json:"method,omitempty"`
	URLPath string `json:"urlPath,omitempty"`
}

type Timing struct {
	ServeTime        int64 `json:"serveTime,omitempty"`
	TotalTime        int64 `json:"totalTime,omitempty"`
	ProcessTime      int64 `json:"processTime,omitempty"`
	ResponseSendTime int64 `json:"responseSendTime,omitempty"`
	AddedDelay       int64 `json:"addedDelay,omitempty"`
}
