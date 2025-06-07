package grpc

import (
	"github.com/wiremock/go-wiremock"
)

type stubRuleBuilder interface {
	Build(serviceName string) *wiremock.StubRule
}

type Service struct {
	serviceName string
	wiremock    *wiremock.Client
}

// NewService creates a new instance of the Service
func NewService(serviceName string, wiremock *wiremock.Client) *Service {
	return &Service{
		serviceName: serviceName,
		wiremock:    wiremock,
	}
}

// StubFor creates a new stub mapping for grpc service.
func (s *Service) StubFor(builder stubRuleBuilder) error {
	return s.wiremock.StubFor(builder.Build(s.serviceName))
}
