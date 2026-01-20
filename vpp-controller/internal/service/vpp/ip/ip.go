package ip

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/NikolayStepanov/RapidVPP/internal/domain"
	"github.com/NikolayStepanov/RapidVPP/internal/infrastructure/vpp"
	"github.com/NikolayStepanov/RapidVPP/internal/mapper"
	"go.fd.io/govpp/api"
	"go.fd.io/govpp/binapi/fib_types"
	"go.fd.io/govpp/binapi/ip"
	"go.fd.io/govpp/binapi/ip_types"
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

func (s *Service) DeleteRoute(ctx context.Context, route *domain.Route) error {
	paths, err := mapper.BuildFibPaths(route.NextHops)
	if err != nil {
		return fmt.Errorf("build fib paths: %w", err)
	}

	req := &ip.IPRouteAddDel{
		IsAdd:       false,
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

	_, err = vpp.DoRequest[*ip.IPRouteAddDel, *ip.IPRouteAddDelReply](s.client, ctx, req)
	if err != nil {
		return fmt.Errorf("delete route from VPP: %w", err)
	}

	return nil
}

func (s *Service) ListRoutes(ctx context.Context, vrf uint32) ([]domain.Route, error) {
	request := &ip.IPRouteDump{
		Table: ip.IPTable{
			TableID: vrf,
			IsIP6:   false,
			Name:    "",
		},
	}

	converter := func(msg api.Message) (domain.Route, bool) {
		details, ok := msg.(*ip.IPRouteDetails)
		if !ok {
			return domain.Route{}, false
		}

		route, err := mapper.ConvertRouteDetails(details)
		if err != nil {
			log.Printf("Failed to convert route: %v", err)
			return domain.Route{}, false
		}

		return route, true
	}

	routesIPv4, err := vpp.Dump(ctx, s.client, request, converter)
	if err != nil {
		return nil, fmt.Errorf("failed to dump IPv4 routes: %w", err)
	}

	request.Table.IsIP6 = true
	routesIPv6, err := vpp.Dump(ctx, s.client, request, converter)
	if err != nil {
		return nil, fmt.Errorf("failed to dump IPv6 routes: %w", err)
	}

	allRoutes := append(routesIPv4, routesIPv6...)

	return allRoutes, nil
}

func (s *Service) GetRoute(ctx context.Context, dst domain.IPWithPrefix, vrf uint32) (domain.Route, error) {
	dstIP := net.ParseIP(dst.Address)
	if dstIP == nil {
		return domain.Route{}, fmt.Errorf("invalid IP address: %s", dst.Address)
	}

	req := &ip.IPRouteLookup{
		TableID: vrf,
		Prefix: ip_types.Prefix{
			Address: ip_types.NewAddress(dstIP),
			Len:     dst.Prefix,
		},
		Exact: 1,
	}

	reply, err := vpp.DoRequest[*ip.IPRouteLookup, *ip.IPRouteLookupReply](
		s.client, ctx, req)
	if err != nil {
		return domain.Route{}, fmt.Errorf("IPRouteLookup failed: %w", err)
	}
	routeDetails := &ip.IPRouteDetails{Route: reply.Route}
	route, err := mapper.ConvertRouteDetails(routeDetails)
	if err != nil {
		return domain.Route{}, fmt.Errorf("failed to convert route details: %w", err)
	}

	return route, nil
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
