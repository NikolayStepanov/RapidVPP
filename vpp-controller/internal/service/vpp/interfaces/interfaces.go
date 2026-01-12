package interfaces

import (
	"context"

	"github.com/NikolayStepanov/RapidVPP/internal/infrastructure/vpp"
	"go.fd.io/govpp/api"
	interfaces "go.fd.io/govpp/binapi/interface"
)

type Service struct {
	client *vpp.Client
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
