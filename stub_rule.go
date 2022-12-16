package wiremock

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	uuidPkg "github.com/google/uuid"
)

const ScenarioStateStarted = "Started"

// ParamMatcherInterface is pair ParamMatchingStrategy and string matched value
type ParamMatcherInterface interface {
	Strategy() ParamMatchingStrategy
	Value() string
}

// URLMatcherInterface is pair URLMatchingStrategy and string matched value
type URLMatcherInterface interface {
	Strategy() URLMatchingStrategy
	Value() string
}

type response struct {
	body                   *string
	base64Body             []byte
	bodyFileName           *string
	jsonBody               interface{}
	headers                map[string]string
	status                 int64
	fixedDelayMilliseconds time.Duration
}

// StubRule is struct of http Request body to WireMock
type StubRule struct {
	uuid                  string
	request               *Request
	response              response
	priority              *int64
	scenarioName          *string
	requiredScenarioState *string
	newScenarioState      *string
}

// NewStubRule returns a new *StubRule.
func NewStubRule(method string, urlMatcher URLMatcher) *StubRule {
	uuid, _ := uuidPkg.NewRandom()
	return &StubRule{
		uuid:    uuid.String(),
		request: NewRequest(method, urlMatcher),
		response: response{
			status: http.StatusOK,
		},
	}
}

// Request is getter for Request
func (s *StubRule) Request() *Request {
	return s.request
}

// WithQueryParam adds query param and returns *StubRule
func (s *StubRule) WithQueryParam(param string, matcher ParamMatcherInterface) *StubRule {
	s.request.WithQueryParam(param, matcher)
	return s
}

// WithHeader adds header to Headers and returns *StubRule
func (s *StubRule) WithHeader(header string, matcher ParamMatcherInterface) *StubRule {
	s.request.WithHeader(header, matcher)
	return s
}

// WithCookie adds cookie and returns *StubRule
func (s *StubRule) WithCookie(cookie string, matcher ParamMatcherInterface) *StubRule {
	s.request.WithCookie(cookie, matcher)
	return s
}

// WithBodyPattern adds body pattern and returns *StubRule
func (s *StubRule) WithBodyPattern(matcher ParamMatcher) *StubRule {
	s.request.WithBodyPattern(matcher)
	return s
}

// WithMultipartPattern adds multipart body pattern and returns *StubRule
func (s *StubRule) WithMultipartPattern(pattern *MultipartPattern) *StubRule {
	s.request.WithMultipartPattern(pattern)
	return s
}

// WillReturn sets response and returns *StubRule
func (s *StubRule) WillReturn(body string, headers map[string]string, status int64) *StubRule {
	s.response.body = &body
	s.response.headers = headers
	s.response.status = status
	return s
}

// WillReturnBinary sets response with binary body and returns *StubRule
func (s *StubRule) WillReturnBinary(body []byte, headers map[string]string, status int64) *StubRule {
	s.response.base64Body = body
	s.response.headers = headers
	s.response.status = status
	return s
}

// WillReturnFileContent sets response with some file content and returns *StubRule
func (s *StubRule) WillReturnFileContent(bodyFileName string, headers map[string]string, status int64) *StubRule {
	s.response.bodyFileName = &bodyFileName
	s.response.headers = headers
	s.response.status = status
	return s
}

// WillReturnJSON sets response with json body and returns *StubRule
func (s *StubRule) WillReturnJSON(json interface{}, headers map[string]string, status int64) *StubRule {
	s.response.jsonBody = json
	s.response.headers = headers
	s.response.status = status
	return s
}

// WithFixedDelayMilliseconds sets fixed delay milliseconds for response
func (s *StubRule) WithFixedDelayMilliseconds(time time.Duration) *StubRule {
	s.response.fixedDelayMilliseconds = time
	return s
}

// WithBasicAuth adds basic auth credentials
func (s *StubRule) WithBasicAuth(username, password string) *StubRule {
	s.request.WithBasicAuth(username, password)
	return s
}

// AtPriority sets priority and returns *StubRule
func (s *StubRule) AtPriority(priority int64) *StubRule {
	s.priority = &priority
	return s
}

// InScenario sets scenarioName and returns *StubRule
func (s *StubRule) InScenario(scenarioName string) *StubRule {
	s.scenarioName = &scenarioName
	return s
}

// WhenScenarioStateIs sets requiredScenarioState and returns *StubRule
func (s *StubRule) WhenScenarioStateIs(scenarioState string) *StubRule {
	s.requiredScenarioState = &scenarioState
	return s
}

// WillSetStateTo sets newScenarioState and returns *StubRule
func (s *StubRule) WillSetStateTo(scenarioState string) *StubRule {
	s.newScenarioState = &scenarioState
	return s
}

// UUID is getter for uuid
func (s *StubRule) UUID() string {
	return s.uuid
}

// Post returns *StubRule for POST method.
func Post(urlMatchingPair URLMatcher) *StubRule {
	return NewStubRule(http.MethodPost, urlMatchingPair)
}

// Get returns *StubRule for GET method.
func Get(urlMatchingPair URLMatcher) *StubRule {
	return NewStubRule(http.MethodGet, urlMatchingPair)
}

// Delete returns *StubRule for DELETE method.
func Delete(urlMatchingPair URLMatcher) *StubRule {
	return NewStubRule(http.MethodDelete, urlMatchingPair)
}

// Put returns *StubRule for PUT method.
func Put(urlMatchingPair URLMatcher) *StubRule {
	return NewStubRule(http.MethodPut, urlMatchingPair)
}

// Patch returns *StubRule for PATCH method.
func Patch(urlMatchingPair URLMatcher) *StubRule {
	return NewStubRule(http.MethodPatch, urlMatchingPair)
}

//MarshalJSON makes json body for http Request
func (s *StubRule) MarshalJSON() ([]byte, error) {
	jsonStubRule := struct {
		UUID                          string   `json:"uuid,omitempty"`
		ID                            string   `json:"id,omitempty"`
		Priority                      *int64   `json:"priority,omitempty"`
		ScenarioName                  *string  `json:"scenarioName,omitempty"`
		RequiredScenarioScenarioState *string  `json:"requiredScenarioState,omitempty"`
		NewScenarioState              *string  `json:"newScenarioState,omitempty"`
		Request                       *Request `json:"request"`
		Response                      struct {
			Body                   string            `json:"body,omitempty"`
			Base64Body             string            `json:"base64Body,omitempty"`
			BodyFileName           string            `json:"bodyFileName,omitempty"`
			JSONBody               interface{}       `json:"jsonBody,omitempty"`
			Headers                map[string]string `json:"headers,omitempty"`
			Status                 int64             `json:"status,omitempty"`
			FixedDelayMilliseconds int               `json:"fixedDelayMilliseconds,omitempty"`
		} `json:"response"`
	}{}
	jsonStubRule.Priority = s.priority
	jsonStubRule.ScenarioName = s.scenarioName
	jsonStubRule.RequiredScenarioScenarioState = s.requiredScenarioState
	jsonStubRule.NewScenarioState = s.newScenarioState

	if s.response.body != nil {
		jsonStubRule.Response.Body = *s.response.body
	} else if len(s.response.base64Body) > 0 {
		jsonStubRule.Response.Base64Body = base64.StdEncoding.EncodeToString(s.response.base64Body)
	} else if s.response.bodyFileName != nil {
		jsonStubRule.Response.BodyFileName = *s.response.bodyFileName
	} else if s.response.jsonBody != nil {
		jsonStubRule.Response.JSONBody = s.response.jsonBody
	}

	jsonStubRule.Response.Headers = s.response.headers
	jsonStubRule.Response.Status = s.response.status
	jsonStubRule.Response.FixedDelayMilliseconds = int(s.response.fixedDelayMilliseconds.Milliseconds())
	jsonStubRule.Request = s.request
	jsonStubRule.ID = s.uuid
	jsonStubRule.UUID = s.uuid

	return json.Marshal(jsonStubRule)
}
