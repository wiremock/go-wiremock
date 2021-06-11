package wiremock

import (
	"encoding/json"
	"net/http"
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

type request struct {
	urlMatcher   URLMatcherInterface
	method       string
	headers      map[string]ParamMatcherInterface
	queryParams  map[string]ParamMatcherInterface
	cookies      map[string]ParamMatcherInterface
	bodyPatterns []ParamMatcher
}

type response struct {
	body                   string
	headers                map[string]string
	status                 int64
	fixedDelayMilliseconds time.Duration
}

// StubRule is struct of http request body to WireMock
type StubRule struct {
	request               request
	response              response
	priority              *int64
	scenarioName          *string
	requiredScenarioState *string
	newScenarioState      *string
}

// NewStubRule returns a new *StubRule.
func NewStubRule(method string, urlMatcher URLMatcher) *StubRule {
	return &StubRule{
		request: request{
			urlMatcher: urlMatcher,
			method:     method,
		},
		response: response{
			status: http.StatusOK,
		},
	}
}

// WithQueryParam adds query param and returns *StubRule
func (s *StubRule) WithQueryParam(param string, matcher ParamMatcherInterface) *StubRule {
	if s.request.queryParams == nil {
		s.request.queryParams = map[string]ParamMatcherInterface{}
	}

	s.request.queryParams[param] = matcher
	return s
}

// WithHeader adds header to Headers and returns *StubRule
func (s *StubRule) WithHeader(header string, matcher ParamMatcherInterface) *StubRule {
	if s.request.headers == nil {
		s.request.headers = map[string]ParamMatcherInterface{}
	}

	s.request.headers[header] = matcher
	return s
}

// WithCookie adds cookie and returns *StubRule
func (s *StubRule) WithCookie(cookie string, matcher ParamMatcherInterface) *StubRule {
	if s.request.cookies == nil {
		s.request.cookies = map[string]ParamMatcherInterface{}
	}

	s.request.cookies[cookie] = matcher
	return s
}

// WithBodyPattern adds body pattern and returns *StubRule
func (s *StubRule) WithBodyPattern(matcher ParamMatcher) *StubRule {
	s.request.bodyPatterns = append(s.request.bodyPatterns, matcher)
	return s
}

// WillReturn sets response and returns *StubRule
func (s *StubRule) WillReturn(body string, headers map[string]string, status int64) *StubRule {
	s.response.body = body
	s.response.headers = headers
	s.response.status = status
	return s
}

func (s *StubRule) WithFixedDelayMilliseconds(time time.Duration) *StubRule {
	s.response.fixedDelayMilliseconds = time
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

//MarshalJSON makes json body for http request
func (s *StubRule) MarshalJSON() ([]byte, error) {
	jsonStubRule := struct {
		Priority                      *int64                 `json:"priority,omitempty"`
		ScenarioName                  *string                `json:"scenarioName,omitempty"`
		RequiredScenarioScenarioState *string                `json:"requiredScenarioState,omitempty"`
		NewScenarioState              *string                `json:"newScenarioState,omitempty"`
		Request                       map[string]interface{} `json:"request"`
		Response                      struct {
			Body                   string            `json:"body,omitempty"`
			Headers                map[string]string `json:"headers,omitempty"`
			Status                 int64             `json:"status,omitempty"`
			FixedDelayMilliseconds int               `json:"fixedDelayMilliseconds,omitempty"`
		} `json:"response"`
	}{}
	jsonStubRule.Priority = s.priority
	jsonStubRule.ScenarioName = s.scenarioName
	jsonStubRule.RequiredScenarioScenarioState = s.requiredScenarioState
	jsonStubRule.NewScenarioState = s.newScenarioState
	jsonStubRule.Response.Body = s.response.body
	jsonStubRule.Response.Headers = s.response.headers
	jsonStubRule.Response.Status = s.response.status
	jsonStubRule.Request = map[string]interface{}{
		"method":                                s.request.method,
		string(s.request.urlMatcher.Strategy()): s.request.urlMatcher.Value(),
	}
	if len(s.request.bodyPatterns) > 0 {
		bodyPatterns := make([]map[ParamMatchingStrategy]string, len(s.request.bodyPatterns))
		for i, bodyPattern := range s.request.bodyPatterns {
			bodyPatterns[i] = map[ParamMatchingStrategy]string{
				bodyPattern.Strategy(): bodyPattern.Value(),
			}
		}
		jsonStubRule.Request["bodyPatterns"] = bodyPatterns
	}
	if len(s.request.headers) > 0 {
		headers := make(map[string]map[ParamMatchingStrategy]string, len(s.request.bodyPatterns))
		for key, header := range s.request.headers {
			headers[key] = map[ParamMatchingStrategy]string{
				header.Strategy(): header.Value(),
			}
		}
		jsonStubRule.Request["headers"] = headers
	}
	if len(s.request.cookies) > 0 {
		cookies := make(map[string]map[ParamMatchingStrategy]string, len(s.request.cookies))
		for key, cookie := range s.request.cookies {
			cookies[key] = map[ParamMatchingStrategy]string{
				cookie.Strategy(): cookie.Value(),
			}
		}
		jsonStubRule.Request["cookies"] = cookies
	}
	if len(s.request.queryParams) > 0 {
		params := make(map[string]map[ParamMatchingStrategy]string, len(s.request.queryParams))
		for key, param := range s.request.queryParams {
			params[key] = map[ParamMatchingStrategy]string{
				param.Strategy(): param.Value(),
			}
		}
		jsonStubRule.Request["queryParameters"] = params
	}

	return json.Marshal(jsonStubRule)
}
