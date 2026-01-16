package ip

import (
	"fmt"
	"net"

	"github.com/NikolayStepanov/RapidVPP/internal/domain"
	"github.com/NikolayStepanov/RapidVPP/internal/infrastructure/vpp"
	"github.com/NikolayStepanov/RapidVPP/internal/mapper"
	"go.fd.io/govpp/binapi/fib_types"
	"go.fd.io/govpp/binapi/ip"
	"go.fd.io/govpp/binapi/ip_types"
	"golang.org/x/net/context"
)

type Service struct {
	client *vpp.Client
}

func (s *Service) AddRoute(ctx context.Context, route *domain.Route) error {
	paths, err := mapper.BuildFibPaths(route.NextHops)
	if err != nil {
		return fmt.Errorf("build fib paths: %w", err)
	}

	req := buildRouteRequest(route, paths)

	_, err = vpp.DoRequest[*ip.IPRouteAddDel, *ip.IPRouteAddDelReply](s.client, ctx, req)
	if err != nil {
		return fmt.Errorf("send route to VPP: %w", err)
	}

	return nil
}

func buildRouteRequest(route *domain.Route, paths []fib_types.FibPath) *ip.IPRouteAddDel {
	return &ip.IPRouteAddDel{
		IsAdd:       true,
		IsMultipath: len(paths) > 1,
		Route: ip.IPRoute{
			TableID: route.VRF,
			Prefix: ip_types.Prefix{
				Address: ip_types.NewAddress(net.ParseIP(route.Dst.Address)),
				Len:     route.Dst.Prefix,
			},
			NPaths: uint8(len(paths)),
			Paths:  paths,
		},
	}
}

func (s *Service) DeleteRoute(ctx context.Context, dst domain.IPWithPrefix, vrf uint32) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) ListRoutes(ctx context.Context, vrf uint32) ([]domain.Route, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetRoute(ctx context.Context, dst domain.IPWithPrefix, vrf uint32) (domain.Route, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) CreateVRF(ctx context.Context, id uint32, name string) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) DeleteVRF(ctx context.Context, id uint32) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) ListVRF(ctx context.Context) ([]domain.VRF, error) {
	//TODO implement me
	panic("implement me")
}

func NewService(client *vpp.Client) *Service {
	return &Service{client: client}
}
