package interfaces

import (
	"context"
	"fmt"

	"github.com/NikolayStepanov/RapidVPP/internal/infrastructure/vpp"
	"go.fd.io/govpp/api"
	interfaces "go.fd.io/govpp/binapi/interface"
	"go.fd.io/govpp/binapi/interface_types"
)

type SwIfFlagsReq = interfaces.SwInterfaceSetFlags
type SwIfFlagsReply = interfaces.SwInterfaceSetFlagsReply

type Service struct {
	client *vpp.Client
}

func NewService(client *vpp.Client) *Service {
	return &Service{client: client}
}

func (s *Service) CreateLoopback(ctx context.Context) (interfaces.CreateLoopbackReply, error) {
	req := &interfaces.CreateLoopback{}

	reply, err := vpp.DoRequest[*interfaces.CreateLoopback, *interfaces.CreateLoopbackReply](s.client, ctx, req)
	if err != nil {
		return interfaces.CreateLoopbackReply{}, fmt.Errorf("create loopback operation failed: %w", err)
	}

	return *reply, nil
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

func (s *Service) SetInterfaceAdminState(ctx context.Context, ifIndex uint32, up bool) error {
	var flags interface_types.IfStatusFlags

	if up {
		flags = interface_types.IF_STATUS_API_FLAG_ADMIN_UP
	}

	req := &SwIfFlagsReq{
		SwIfIndex: interface_types.InterfaceIndex(ifIndex),
		Flags:     flags,
	}

	_, err := vpp.DoRequest[*SwIfFlagsReq, *SwIfFlagsReply](s.client, ctx, req)
	if err != nil {
		return fmt.Errorf("set interface admin state failed: %w", err)
	}

	return nil
}
