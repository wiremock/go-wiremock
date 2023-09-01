package wiremock

import (
	"encoding/json"
	"net/http"
	"time"

	uuidPkg "github.com/google/uuid"
)

const ScenarioStateStarted = "Started"
const authorizationHeader = "Authorization"

// StubRule is struct of http Request body to WireMock
type StubRule struct {
	uuid                   string
	request                *Request
	response               ResponseInterface
	fixedDelayMilliseconds *int64
	priority               *int64
	scenarioName           *string
	requiredScenarioState  *string
	newScenarioState       *string
	postServeActions       []WebhookInterface
}

// NewStubRule returns a new *StubRule.
func NewStubRule(method string, urlMatcher URLMatcher) *StubRule {
	uuid, _ := uuidPkg.NewRandom()
	return &StubRule{
		uuid:     uuid.String(),
		request:  NewRequest(method, urlMatcher),
		response: NewResponse(),
	}
}

// Request is getter for Request
func (s *StubRule) Request() *Request {
	return s.request
}

// WithQueryParam adds query param and returns *StubRule
func (s *StubRule) WithQueryParam(param string, matcher MatcherInterface) *StubRule {
	s.request.WithQueryParam(param, matcher)
	return s
}

// WithPort adds port and returns *StubRule
func (s *StubRule) WithPort(port int64) *StubRule {
	s.request.WithPort(port)
	return s
}

// WithScheme adds scheme and returns *StubRule
func (s *StubRule) WithScheme(scheme string) *StubRule {
	s.request.WithScheme(scheme)
	return s
}

// WithHost adds host and returns *StubRule
func (s *StubRule) WithHost(host BasicParamMatcher) *StubRule {
	s.request.WithHost(host)
	return s
}

// WithHeader adds header to Headers and returns *StubRule
func (s *StubRule) WithHeader(header string, matcher MatcherInterface) *StubRule {
	s.request.WithHeader(header, matcher)
	return s
}

// WithCookie adds cookie and returns *StubRule
func (s *StubRule) WithCookie(cookie string, matcher BasicParamMatcher) *StubRule {
	s.request.WithCookie(cookie, matcher)
	return s
}

// WithBodyPattern adds body pattern and returns *StubRule
func (s *StubRule) WithBodyPattern(matcher BasicParamMatcher) *StubRule {
	s.request.WithBodyPattern(matcher)
	return s
}

// WithMultipartPattern adds multipart body pattern and returns *StubRule
func (s *StubRule) WithMultipartPattern(pattern *MultipartPattern) *StubRule {
	s.request.WithMultipartPattern(pattern)
	return s
}

// WithAuthToken adds Authorization header with Token auth method *StubRule
func (s *StubRule) WithAuthToken(tokenMatcher BasicParamMatcher) *StubRule {
	methodPrefix := "Token "
	m := addAuthMethodToMatcher(tokenMatcher, methodPrefix)
	s.WithHeader(authorizationHeader, HasExactly(StartsWith(methodPrefix), m))
	return s
}

// WithBearerToken adds Authorization header with Bearer auth method *StubRule
func (s *StubRule) WithBearerToken(tokenMatcher BasicParamMatcher) *StubRule {
	methodPrefix := "Bearer "
	m := addAuthMethodToMatcher(tokenMatcher, methodPrefix)
	s.WithHeader(authorizationHeader, StartsWith(methodPrefix).And(m))
	return s
}

// WithDigestAuth adds Authorization header with Digest auth method *StubRule
func (s *StubRule) WithDigestAuth(matcher BasicParamMatcher) *StubRule {
	methodPrefix := "Digest "
	m := addAuthMethodToMatcher(matcher, methodPrefix)
	s.WithHeader(authorizationHeader, HasExactly(StartsWith(methodPrefix), m))
	return s
}

// Deprecated: Use WillReturnResponse(NewResponse().WithBody(body).WithHeaders(headers).WithStatus(status)) instead
// WillReturn sets response and returns *StubRule
func (s *StubRule) WillReturn(body string, headers map[string]string, status int64) *StubRule {
	s.response = NewResponse().WithBody(body).WithStatus(status).WithHeaders(headers)
	return s
}

// Deprecated: Use WillReturnResponse(NewResponse().WithBinaryBody(body).WithHeaders(headers).WithStatus(status)) instead
// WillReturnBinary sets response with binary body and returns *StubRule
func (s *StubRule) WillReturnBinary(body []byte, headers map[string]string, status int64) *StubRule {
	s.response = NewResponse().WithBinaryBody(body).WithStatus(status).WithHeaders(headers)
	return s
}

// Deprecated: Use WillReturnResponse(NewResponse().WithBodyFile(file).WithHeaders(headers).WithStatus(status)) instead
// WillReturnFileContent sets response with some file content and returns *StubRule
func (s *StubRule) WillReturnFileContent(bodyFileName string, headers map[string]string, status int64) *StubRule {
	s.response = NewResponse().WithBodyFile(bodyFileName).WithStatus(status).WithHeaders(headers)
	return s
}

// Deprecated: Use WillReturnResponse(NewResponse().WithJsonBody(json).WithHeaders(headers).WithStatus(status)) instead
// WillReturnJSON sets response with json body and returns *StubRule
func (s *StubRule) WillReturnJSON(json interface{}, headers map[string]string, status int64) *StubRule {
	s.response = NewResponse().WithJSONBody(json).WithStatus(status).WithHeaders(headers)
	return s
}

// Deprecated: Use WillReturnResponse(NewResponse().WithFixedDelay(time.Second)) instead
// WithFixedDelayMilliseconds adds delay to response and returns *StubRule
func (s *StubRule) WithFixedDelayMilliseconds(delay time.Duration) *StubRule {
	milliseconds := delay.Milliseconds()
	s.fixedDelayMilliseconds = &milliseconds
	return s
}

// WillReturnResponse sets response and returns *StubRule
func (s *StubRule) WillReturnResponse(response ResponseInterface) *StubRule {
	s.response = response
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

func (s *StubRule) WithPostServeAction(extensionName string, webhook WebhookInterface) *StubRule {
	s.postServeActions = append(s.postServeActions, webhook.WithName(extensionName))
	return s
}

// MarshalJSON makes json body for http Request
func (s *StubRule) MarshalJSON() ([]byte, error) {
	jsonStubRule := struct {
		UUID                          string                 `json:"uuid,omitempty"`
		ID                            string                 `json:"id,omitempty"`
		Priority                      *int64                 `json:"priority,omitempty"`
		ScenarioName                  *string                `json:"scenarioName,omitempty"`
		RequiredScenarioScenarioState *string                `json:"requiredScenarioState,omitempty"`
		NewScenarioState              *string                `json:"newScenarioState,omitempty"`
		Request                       *Request               `json:"request"`
		Response                      map[string]interface{} `json:"response"`
		PostServeActions              []WebhookInterface     `json:"postServeActions,omitempty"`
	}{}

	jsonStubRule.Priority = s.priority
	jsonStubRule.ScenarioName = s.scenarioName
	jsonStubRule.RequiredScenarioScenarioState = s.requiredScenarioState
	jsonStubRule.NewScenarioState = s.newScenarioState
	jsonStubRule.Response = s.response.ParseResponse()
	jsonStubRule.PostServeActions = s.postServeActions

	if s.fixedDelayMilliseconds != nil {
		jsonStubRule.Response["fixedDelayMilliseconds"] = *s.fixedDelayMilliseconds
	}

	jsonStubRule.Request = s.request
	jsonStubRule.ID = s.uuid
	jsonStubRule.UUID = s.uuid

	return json.Marshal(jsonStubRule)
}

func addAuthMethodToMatcher(matcher BasicParamMatcher, methodPrefix string) BasicParamMatcher {
	switch m := matcher.(type) {
	case StringValueMatcher:
		return m.addPrefixToMatcher(methodPrefix)
	case LogicalMatcher:
		for i, operand := range m.operands {
			m.operands[i] = addAuthMethodToMatcher(operand, methodPrefix)
		}
		return m
	default:
		return matcher
	}
}
