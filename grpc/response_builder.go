package grpc

import (
	"strconv"

	"github.com/wiremock/go-wiremock"
)

const (
	responseStatusName   = "grpc-status-name"
	responseStatusReason = "grpc-status-reason"
)

type GRPCResponseStatus int

const (
	StatusOK = GRPCResponseStatus(iota)
	StatusCanceled
	StatusUnknown
	StatusInvalidArgument
	StatusDeadlineExceeded
	StatusNotFound
	StatusAlreadyExists
	StatusPermissionDenied
	StatusResourceExhausted
	StatusFailedPrecondition
	StatusAborted
	StatusOutOfRange
	StatusUnimplemented
	StatusInternal
	StatusUnavailable
	StatusDataLoss
	StatusUnauthenticated
)

type ResponseBuilder struct {
	grpcResponseStatus GRPCResponseStatus
	grpcStatusReason   *string
	fault              *wiremock.Fault
	json               any
	body               string
	delay              wiremock.DelayInterface
}

func Error(grpcResponseStatus GRPCResponseStatus, grpcResponseReason string) *ResponseBuilder {
	return &ResponseBuilder{
		grpcResponseStatus: grpcResponseStatus,
		grpcStatusReason:   &grpcResponseReason,
	}
}

func JSON(json string) *ResponseBuilder {
	return &ResponseBuilder{
		body: json,
	}
}

func (b *ResponseBuilder) Message(json any) *ResponseBuilder {
	return &ResponseBuilder{
		json: json,
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
	if b.fault != nil {
		return wiremock.NewResponse().WithFault(*b.fault)
	}

	response := wiremock.OK().WithHeader(responseStatusName, strconv.Itoa(int(b.grpcResponseStatus)))

	if b.grpcStatusReason != nil {
		return response.WithHeader(responseStatusReason, *b.grpcStatusReason)
	}

	if b.fault != nil {
		return response.WithFault(*b.fault)
	}

	if b.delay != nil {
		response = response.WithDelay(b.delay)
	}

	if b.json != nil {
		return response.WithJSONBody(b.json)
	}

	return response.WithBody(b.body)
}
