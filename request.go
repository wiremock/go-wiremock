package wiremock

import (
	"encoding/json"
)

// A Request is the part of StubRule describing the matching of the http request
type Request struct {
	urlMatcher           URLMatcherInterface
	method               string
	host                 BasicParamMatcher
	port                 *int64
	scheme               *string
	headers              map[string]MatcherInterface
	queryParams          map[string]MatcherInterface
	pathParams           map[string]MatcherInterface
	cookies              map[string]BasicParamMatcher
	formParameters       map[string]BasicParamMatcher
	bodyPatterns         []BasicParamMatcher
	multipartPatterns    []MultipartPatternInterface
	basicAuthCredentials *struct {
		username string
		password string
	}
}

// NewRequest constructs minimum possible Request
func NewRequest(method string, urlMatcher URLMatcherInterface) *Request {
	return &Request{
		method:     method,
		urlMatcher: urlMatcher,
	}
}

// WithPort is fluent-setter for port
func (r *Request) WithPort(port int64) *Request {
	r.port = &port
	return r
}

// WithScheme is fluent-setter for scheme
func (r *Request) WithScheme(scheme string) *Request {
	r.scheme = &scheme
	return r
}

// WithHost is fluent-setter for host
func (r *Request) WithHost(host BasicParamMatcher) *Request {
	r.host = host
	return r
}

// WithMethod is fluent-setter for http verb
func (r *Request) WithMethod(method string) *Request {
	r.method = method
	return r
}

// WithURLMatched is fluent-setter url matcher
func (r *Request) WithURLMatched(urlMatcher URLMatcherInterface) *Request {
	r.urlMatcher = urlMatcher
	return r
}

// WithBodyPattern adds body pattern to list
func (r *Request) WithBodyPattern(matcher BasicParamMatcher) *Request {
	r.bodyPatterns = append(r.bodyPatterns, matcher)
	return r
}

// WithFormParameter adds form parameter to list
func (r *Request) WithFormParameter(name string, matcher BasicParamMatcher) *Request {
	if r.formParameters == nil {
		r.formParameters = make(map[string]BasicParamMatcher, 1)
	}
	r.formParameters[name] = matcher
	return r
}

// WithMultipartPattern adds multipart pattern to list
func (r *Request) WithMultipartPattern(pattern *MultipartPattern) *Request {
	r.multipartPatterns = append(r.multipartPatterns, pattern)
	return r
}

// WithBasicAuth adds basic auth credentials to Request
func (r *Request) WithBasicAuth(username, password string) *Request {
	r.basicAuthCredentials = &struct {
		username string
		password string
	}{
		username: username,
		password: password,
	}
	return r
}

// WithQueryParam add param to query param list
func (r *Request) WithQueryParam(param string, matcher MatcherInterface) *Request {
	if r.queryParams == nil {
		r.queryParams = map[string]MatcherInterface{}
	}

	r.queryParams[param] = matcher
	return r
}

// WithPathParam add param to path param list
func (r *Request) WithPathParam(param string, matcher MatcherInterface) *Request {
	if r.pathParams == nil {
		r.pathParams = map[string]MatcherInterface{}
	}

	r.pathParams[param] = matcher
	return r
}

// WithHeader add header to header list
func (r *Request) WithHeader(header string, matcher MatcherInterface) *Request {
	if r.headers == nil {
		r.headers = map[string]MatcherInterface{}
	}

	r.headers[header] = matcher
	return r
}

// WithCookie is fluent-setter for cookie
func (r *Request) WithCookie(cookie string, matcher BasicParamMatcher) *Request {
	if r.cookies == nil {
		r.cookies = map[string]BasicParamMatcher{}
	}

	r.cookies[cookie] = matcher
	return r
}

// MarshalJSON gives valid JSON or error.
func (r *Request) MarshalJSON() ([]byte, error) {
	request := map[string]interface{}{
		"method":                        r.method,
		string(r.urlMatcher.Strategy()): r.urlMatcher.Value(),
	}

	if r.scheme != nil {
		request["scheme"] = r.scheme
	}

	if r.host != nil {
		request["host"] = r.host
	}

	if r.port != nil {
		request["port"] = r.port
	}

	if len(r.bodyPatterns) > 0 {
		request["bodyPatterns"] = r.bodyPatterns
	}
	if len(r.formParameters) > 0 {
		request["formParameters"] = r.formParameters
	}
	if len(r.multipartPatterns) > 0 {
		request["multipartPatterns"] = r.multipartPatterns
	}
	if len(r.headers) > 0 {
		request["headers"] = r.headers
	}
	if len(r.cookies) > 0 {
		request["cookies"] = r.cookies
	}
	if len(r.queryParams) > 0 {
		request["queryParameters"] = r.queryParams
	}
	if len(r.pathParams) > 0 {
		request["pathParameters"] = r.pathParams
	}

	if r.basicAuthCredentials != nil {
		request["basicAuthCredentials"] = map[string]string{
			"password": r.basicAuthCredentials.password,
			"username": r.basicAuthCredentials.username,
		}
	}

	return json.Marshal(request)
}
