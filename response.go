package wiremock

import (
	"net/http"
	"time"
)

type Fault string

const (
	FaultEmptyResponse          Fault = "EMPTY_RESPONSE"
	FaultMalformedResponseChunk Fault = "MALFORMED_RESPONSE_CHUNK"
	FaultRandomDataThenClose    Fault = "RANDOM_DATA_THEN_CLOSE"
	FaultConnectionResetByPeer  Fault = "CONNECTION_RESET_BY_PEER"
)

type ResponseInterface interface {
	ParseResponse() map[string]interface{}
}

type Response struct {
	body                *string
	base64Body          []byte
	bodyFileName        *string
	jsonBody            interface{}
	headers             map[string]string
	status              int64
	delayDistribution   DelayInterface
	chunkedDribbleDelay *chunkedDribbleDelay
	fault               *Fault
}

func NewResponse() Response {
	return Response{
		status: http.StatusOK,
	}
}

// WithLogNormalRandomDelay sets log normal random delay for response
func (r Response) WithLogNormalRandomDelay(mediana time.Duration, sigma float64) Response {
	r.delayDistribution = logNormalRandomDelay{
		median: mediana.Milliseconds(),
		sigma:  sigma,
	}

	return r
}

// WithUniformRandomDelay sets uniform random delay for response
func (r Response) WithUniformRandomDelay(lower, upper time.Duration) Response {
	r.delayDistribution = uniformRandomDelay{
		lower: lower.Milliseconds(),
		upper: upper.Milliseconds(),
	}

	return r
}

// WithFixedDelay sets fixed delay milliseconds for response
func (r Response) WithFixedDelay(time time.Duration) Response {
	r.delayDistribution = fixedDelay{
		milliseconds: time.Milliseconds(),
	}

	return r
}

// WithChunkedDribbleDelay sets chunked dribble delay for response
func (r Response) WithChunkedDribbleDelay(numberOfChunks int64, totalDuration time.Duration) Response {
	r.chunkedDribbleDelay = &chunkedDribbleDelay{
		numberOfChunks: numberOfChunks,
		totalDuration:  totalDuration.Milliseconds(),
	}

	return r
}

// WithStatus sets status for response
func (r Response) WithStatus(status int64) Response {
	r.status = status
	return r
}

// WithHeader sets header for response
func (r Response) WithHeader(key, value string) Response {
	if r.headers == nil {
		r.headers = make(map[string]string)
	}

	r.headers[key] = value

	return r
}

// WithHeaders sets headers for response
func (r Response) WithHeaders(headers map[string]string) Response {
	r.headers = headers
	return r
}

func (r Response) WithFault(fault Fault) Response {
	r.fault = &fault
	return r
}

// WithBody sets body for response
func (r Response) WithBody(body string) Response {
	r.body = &body
	return r
}

// WithBinaryBody sets binary body for response
func (r Response) WithBinaryBody(body []byte) Response {
	r.base64Body = body
	return r
}

// WithJSONBody sets json body for response
func (r Response) WithJSONBody(body interface{}) Response {
	r.jsonBody = body
	return r
}

// WithBodyFile sets body file name for response
func (r Response) WithBodyFile(fileName string) Response {
	r.bodyFileName = &fileName
	return r
}

func (r Response) ParseResponse() map[string]interface{} {
	jsonMap := map[string]interface{}{
		"status": r.status,
	}

	if r.body != nil {
		jsonMap["body"] = *r.body
	}

	if r.base64Body != nil {
		jsonMap["base64Body"] = r.base64Body
	}

	if r.bodyFileName != nil {
		jsonMap["bodyFileName"] = *r.bodyFileName
	}

	if r.jsonBody != nil {
		jsonMap["jsonBody"] = r.jsonBody
	}

	if r.headers != nil {
		jsonMap["headers"] = r.headers
	}

	if r.delayDistribution != nil {
		jsonMap["delayDistribution"] = r.delayDistribution.ParseDelay()
	}

	if r.chunkedDribbleDelay != nil {
		jsonMap["chunkedDribbleDelay"] = r.chunkedDribbleDelay
	}

	if r.fault != nil {
		jsonMap["fault"] = *r.fault
	}

	return jsonMap
}
