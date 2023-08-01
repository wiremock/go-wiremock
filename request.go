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
	headers              map[string]json.Marshaler
	queryParams          map[string]json.Marshaler
	cookies              map[string]BasicParamMatcher
	bodyPatterns         []BasicParamMatcher
	multipartPatterns    []*MultipartPattern
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
func (r *Request) WithQueryParam(param string, matcher json.Marshaler) *Request {
	if r.queryParams == nil {
		r.queryParams = map[string]json.Marshaler{}
	}

	r.queryParams[param] = matcher
	return r
}

// WithHeader add header to header list
func (r *Request) WithHeader(header string, matcher json.Marshaler) *Request {
	if r.headers == nil {
		r.headers = map[string]json.Marshaler{}
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

	if r.basicAuthCredentials != nil {
		request["basicAuthCredentials"] = map[string]string{
			"password": r.basicAuthCredentials.password,
			"username": r.basicAuthCredentials.username,
		}
	}

	return json.Marshal(request)
}
