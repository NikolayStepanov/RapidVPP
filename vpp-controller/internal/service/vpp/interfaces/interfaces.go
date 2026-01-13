package interfaces

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/NikolayStepanov/RapidVPP/internal/domain"
	"github.com/NikolayStepanov/RapidVPP/internal/infrastructure/vpp"
	"github.com/NikolayStepanov/RapidVPP/pkg/logger"
	"go.fd.io/govpp/api"
	interfaces "go.fd.io/govpp/binapi/interface"
	"go.fd.io/govpp/binapi/interface_types"
	"go.fd.io/govpp/binapi/ip_types"
	"go.uber.org/zap"
)

var (
	ErrNotFound      = errors.New("interface not found")
	ErrAlreadyExists = errors.New("resource already exists")
)

type (
	SwIfFlagsReq        = interfaces.SwInterfaceSetFlags
	SwIfFlagsReply      = interfaces.SwInterfaceSetFlagsReply
	SwIfAddDelAddrReq   = interfaces.SwInterfaceAddDelAddress
	SwIfAddDelAddrReply = interfaces.SwInterfaceAddDelAddressReply
)

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

func (s *Service) DeleteLoopback(ctx context.Context, ifIndex uint32) error {
	req := &interfaces.DeleteLoopback{
		SwIfIndex: interface_types.InterfaceIndex(ifIndex),
	}

	_, err := vpp.DoRequest[*interfaces.DeleteLoopback, *interfaces.DeleteLoopbackReply](s.client, ctx, req)
	if err != nil {
		return fmt.Errorf("delete loopback failed: %w", err)
	}
	return nil
}

func (s *Service) SetInterfaceIP(ctx context.Context, ifIndex uint32, IPPrefix domain.IPWithPrefix) error {
	req := &interfaces.SwInterfaceAddDelAddress{
		SwIfIndex: interface_types.InterfaceIndex(ifIndex),
		IsAdd:     true,
		Prefix: ip_types.AddressWithPrefix{
			Address: ip_types.NewAddress(net.ParseIP(IPPrefix.Address)),
			Len:     IPPrefix.Prefix,
		},
	}
	logger.Debug("req ", zap.Any("req", req))
	_, err := vpp.DoRequest[*SwIfAddDelAddrReq, *SwIfAddDelAddrReply](s.client, ctx, req)
	if err != nil {
		logger.Error("set interface ip address failed", zap.Error(err))
		return mapVppError(err)
	}
	return nil
}

func mapVppError(err error) error {
	switch {
	case errors.Is(err, api.NO_SUCH_ENTRY), errors.Is(err, api.INVALID_SW_IF_INDEX):
		return ErrNotFound
	case errors.Is(err, api.DUPLICATE_IF_ADDRESS), errors.Is(err, api.IF_ALREADY_EXISTS), errors.Is(err, api.ADDRESS_IN_USE):
		return ErrAlreadyExists
	}
	return err
}
