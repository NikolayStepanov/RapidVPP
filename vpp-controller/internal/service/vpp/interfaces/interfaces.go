package interfaces

import (
	"context"
	"fmt"

	"github.com/NikolayStepanov/RapidVPP/internal/infrastructure/vpp"
	"go.fd.io/govpp/api"
	interfaces "go.fd.io/govpp/binapi/interface"
)

type Service struct {
	client *vpp.Client
}

func (s *Service) CreateLoopback(ctx context.Context) (interfaces.CreateLoopbackReply, error) {
	req := &interfaces.CreateLoopback{}

	reply, err := vpp.DoRequest[*interfaces.CreateLoopback, *interfaces.CreateLoopbackReply](s.client, ctx, req)
	if err != nil {
		return interfaces.CreateLoopbackReply{}, fmt.Errorf("create loopback operation failed: %w", err)
	}

	return *reply, nil
}

func NewService(client *vpp.Client) *Service {
	return &Service{client: client}
}

func (s *Service) List(ctx context.Context) ([]interfaces.SwInterfaceDetails, error) {
	request := &interfaces.SwInterfaceDump{
		SwIfIndex:       0xFFFFFFFF,
		NameFilterValid: false,
	}

	converter := func(msg api.Message) (interfaces.SwInterfaceDetails, bool) {
		if details, ok := msg.(*interfaces.SwInterfaceDetails); ok {
			return *details, true
		}
		return interfaces.SwInterfaceDetails{}, false
	}

	return vpp.Dump(ctx, s.client, request, converter)
}
