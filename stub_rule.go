package wiremock

import (
	"encoding/json"
	"net/http"
)

type ParamMatcherInterface interface {
	Strategy() ParamMatchingStrategy
	Value() string
}

type URLMatcherInterface interface {
	Strategy() URLMatchingStrategy
	Value() string
}

type Request struct {
	urlMatcher   URLMatcherInterface
	method       string
	headers      map[string]ParamMatcherInterface
	queryParams  map[string]ParamMatcherInterface
	cookies      map[string]ParamMatcherInterface
	bodyPatterns []ParamMatcher
}

type Response struct {
	body    string
	headers map[string]string
	status  int64
}

type StubRule struct {
	request  Request
	response Response
}

// WithQueryParam adds query param and returns *StubRule
func (r *StubRule) WithQueryParam(param string, matcher ParamMatcherInterface) *StubRule {
	if r.request.queryParams == nil {
		r.request.queryParams = map[string]ParamMatcherInterface{
			param: matcher,
		}
		return r
	}

	r.request.queryParams[param] = matcher
	return r
}

// WithHeader adds header to Headers and returns *StubRule
func (r *StubRule) WithHeader(header string, matcher ParamMatcherInterface) *StubRule {
	if r.request.headers == nil {
		r.request.headers = map[string]ParamMatcherInterface{
			header: matcher,
		}
		return r
	}

	r.request.headers[header] = matcher
	return r
}

// WithCookie adds cookie and returns *StubRule
func (r *StubRule) WithCookie(cookie string, matcher ParamMatcherInterface) *StubRule {
	if r.request.cookies == nil {
		r.request.cookies = map[string]ParamMatcherInterface{
			cookie: matcher,
		}
		return r
	}

	r.request.cookies[cookie] = matcher
	return r
}

// WithBodyPattern adds body pattern and returns *StubRule
func (r *StubRule) WithBodyPattern(matcher ParamMatcher) *StubRule {
	r.request.bodyPatterns = append(r.request.bodyPatterns, matcher)
	return r
}

// WithBodyPattern sets response and returns *StubRule
func (r *StubRule) WillReturn(body string, headers map[string]string, status int64) *StubRule {
	r.response.body = body
	r.response.headers = headers
	r.response.status = status
	return r
}

// NewStubRule returns a new *StubRule.
func NewStubRule(method string, urlMatcher URLMatcher) *StubRule {
	return &StubRule{
		request: Request{
			urlMatcher: urlMatcher,
			method:     method,
		},
		response: Response{
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
