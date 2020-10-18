package wiremock

import (
	"encoding/json"
	"net/http"
)

// ParamMatcherInterface is pair ParamMatchingStrategy key and string value
type ParamMatcherInterface interface {
	Strategy() ParamMatchingStrategy
	Value() string
}

// URLMatcherInterface is pair URLMatchingStrategy key and string value
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
	body    string
	headers map[string]string
	status  int64
}

// StubRule is struct of http request body to WireMock
type StubRule struct {
	request  request
	response response
}

// WithQueryParam adds query param and returns *StubRule
func (s *StubRule) WithQueryParam(param string, matcher ParamMatcherInterface) *StubRule {
	if s.request.queryParams == nil {
		s.request.queryParams = map[string]ParamMatcherInterface{
			param: matcher,
		}
		return s
	}

	s.request.queryParams[param] = matcher
	return s
}

// WithHeader adds header to Headers and returns *StubRule
func (s *StubRule) WithHeader(header string, matcher ParamMatcherInterface) *StubRule {
	if s.request.headers == nil {
		s.request.headers = map[string]ParamMatcherInterface{
			header: matcher,
		}
		return s
	}

	s.request.headers[header] = matcher
	return s
}

// WithCookie adds cookie and returns *StubRule
func (s *StubRule) WithCookie(cookie string, matcher ParamMatcherInterface) *StubRule {
	if s.request.cookies == nil {
		s.request.cookies = map[string]ParamMatcherInterface{
			cookie: matcher,
		}
		return s
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

func (s *StubRule) MarshalJSON() ([]byte, error) {
	jsonStubRule := struct {
		Request  map[string]interface{} `json:"request"`
		Response struct {
			Body    string            `json:"body,omitempty"`
			Headers map[string]string `json:"headers,omitempty"`
			Status  int64             `json:"status,omitempty"`
		} `json:"response"`
	}{}
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
