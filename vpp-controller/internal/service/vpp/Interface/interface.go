package Interface

import (
	"time"

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

func (s *Service) List() ([]interfaces.SwInterfaceDetails, error) {
	return vpp.DumpWithTimeout(
		s.client,
		30*time.Second,
		&interfaces.SwInterfaceDump{
			SwIfIndex:       0xFFFFFFFF,
			NameFilterValid: false,
		},
		func(msg api.Message) (interfaces.SwInterfaceDetails, bool) {
			if details, ok := msg.(*interfaces.SwInterfaceDetails); ok {
				return *details, true
			}
			return interfaces.SwInterfaceDetails{}, false
		},
	)
}
