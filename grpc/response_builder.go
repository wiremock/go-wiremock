package grpc

import (
	"regexp"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/wiremock/go-wiremock"
)

const (
	responseStatusName   = "grpc-status-name"
	responseStatusReason = "grpc-status-reason"
)

type ResponseBuilder struct {
	grpcResponseStatus codes.Code
	grpcStatusReason   *string
	fault              *wiremock.Fault
	body               string
	delay              wiremock.DelayInterface
}

// Error creates a response builder with a gRPC error status and reason.
func Error(grpcResponseStatus codes.Code, grpcResponseReason string) *ResponseBuilder {
	return &ResponseBuilder{
		grpcResponseStatus: grpcResponseStatus,
		grpcStatusReason:   &grpcResponseReason,
	}
}

// JSON creates a response builder with a JSON body.
func JSON(json string) *ResponseBuilder {
	return &ResponseBuilder{
		body: json,
	}
}

// Message creates a response builder with a serialized proto message as the body. Can panic if there are problems with marshaling.
func Message(message proto.Message) *ResponseBuilder {
	bytes, err := protojson.Marshal(message)
	if err != nil {
		panic("failed to marshal proto message: " + err.Error())
	}

	return &ResponseBuilder{
		body: string(bytes),
	}
}

func Fault(fault wiremock.Fault) *ResponseBuilder {
	return &ResponseBuilder{
		fault: &fault,
	}
}

func (b *ResponseBuilder) WithDelay(delay wiremock.DelayInterface) *ResponseBuilder {
	b.delay = delay
	return b
}

func (b *ResponseBuilder) Build() wiremock.Response {
	response := wiremock.OK().WithHeader(responseStatusName, grpcCodeToWireMockCode(b.grpcResponseStatus))

	if b.grpcStatusReason != nil {
		return response.WithHeader(responseStatusReason, *b.grpcStatusReason)
	}

	if b.fault != nil {
		return response.WithFault(*b.fault)
	}

	if b.delay != nil {
		response = response.WithDelay(b.delay)
	}

	return response.WithBody(b.body)
}

func grpcCodeToWireMockCode(code codes.Code) string {
	if code == codes.Canceled {
		return "CANCELLED"
	}

	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(code.String(), "${1}_${2}")
	return strings.ToUpper(snake)
}
