package ip

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/netip"
	"sort"
	"sync"

	"github.com/NikolayStepanov/RapidVPP/internal/domain"
	"github.com/NikolayStepanov/RapidVPP/internal/infrastructure/vpp"
	"github.com/NikolayStepanov/RapidVPP/internal/mapper"
	"go.fd.io/govpp/api"
	"go.fd.io/govpp/binapi/fib_types"
	"go.fd.io/govpp/binapi/ip"
	"go.fd.io/govpp/binapi/ip_types"
)

var systemPrefixes = []netip.Prefix{
	netip.MustParsePrefix("fe80::/10"), // Link-local
	netip.MustParsePrefix("ff00::/8"),  // Multicast
	netip.MustParsePrefix("::1/128"),   // Loopback
	netip.MustParsePrefix("100::/64"),  // Discard-only
}

type Service struct {
	client   *vpp.Client
	vrfCache *VRFCache
}

type VRFCache struct {
	mu   sync.RWMutex
	data map[uint32]*domain.VRFEntry
}

func NewService(client *vpp.Client) *Service {
	return &Service{
		client: client,
		vrfCache: &VRFCache{
			data: make(map[uint32]*domain.VRFEntry),
		},
	}

}

func (s *Service) InitVRFCache(ctx context.Context) error {
	vrfEntries, err := s.dumpAllVPPVRF(ctx)
	if err != nil {
		return err
	}
	for _, vrf := range vrfEntries {
		s.addOrUpdateVRFCache(vrf.ID, vrf.Name)
	}
	return nil
}

func (s *Service) dumpAllVPPVRF(ctx context.Context) ([]domain.VRF, error) {
	var vrfList []domain.VRF

	converter := func(msg api.Message) (interface{}, bool) {
		tbl, ok := msg.(*ip.IPTableDetails)
		if !ok {
			return nil, false
		}

		vrfList = append(vrfList, domain.VRF{
			ID:   tbl.Table.TableID,
			Name: tbl.Table.Name,
		})

		return nil, false
	}

	req := &ip.IPTableDump{}
	_, err := vpp.Dump(ctx, s.client, req, converter)
	if err != nil {
		return nil, fmt.Errorf("failed to dump VPP VRFs: %w", err)
	}

	return vrfList, nil
}

func (s *Service) addVRFsFromCache(vrfMap map[uint32]*domain.VRF) {
	s.vrfCache.mu.RLock()
	defer s.vrfCache.mu.RUnlock()

	for id, entry := range s.vrfCache.data {
		if _, exists := vrfMap[id]; !exists {
			vrfMap[id] = &domain.VRF{
				ID:         id,
				Name:       entry.Name,
				IPv4:       false,
				IPv6:       false,
				RouteCount: 0,
			}
		}
	}
}

func (s *Service) getEntryVRFCache(id uint32) (*domain.VRFEntry, error) {
	s.vrfCache.mu.RLock()
	defer s.vrfCache.mu.RUnlock()
	entry, ok := s.vrfCache.data[id]
	if !ok {
		return nil, fmt.Errorf("VRF with ID %d not found in cache", id)
	}
	return entry, nil
}
func (s *Service) deleteVRFCache(id uint32) {
	s.vrfCache.mu.Lock()
	defer s.vrfCache.mu.Unlock()
	delete(s.vrfCache.data, id)
}

func IsSystemRoute(details *ip.IPRouteDetails) bool {
	for _, path := range details.Route.Paths {
		switch path.Type {
		case fib_types.FIB_API_PATH_TYPE_DROP,
			fib_types.FIB_API_PATH_TYPE_LOCAL,
			fib_types.FIB_API_PATH_TYPE_INTERFACE_RX:
			return true
		}
	}
	if IsIPv6SystemRoute(details.Route.Prefix) {
		return true
	}

	return false
}

func IsIPv6SystemRoute(pfx ip_types.Prefix) bool {
	if pfx.Address.Af != ip_types.ADDRESS_IP6 {
		return false
	}

	ip6 := pfx.Address.Un.GetIP6().ToIP()
	addr := netip.AddrFrom16([16]byte(ip6))
	prefixLen := pfx.Len

	for _, sys := range systemPrefixes {
		if sys.Contains(addr) && prefixLen >= uint8(sys.Bits()) {
			return true
		}
	}
	return false
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
	ipv4Req := &ip.IPTableAddDel{
		IsAdd: true,
		Table: ip.IPTable{
			TableID: id,
			IsIP6:   false,
			Name:    name,
		},
	}

	_, err := vpp.DoRequest[*ip.IPTableAddDel, *ip.IPTableAddDelReply](s.client, ctx, ipv4Req)
	if err != nil {
		return fmt.Errorf("failed to create IPv4 VRF table: %w", err)
	}

	ipv6Req := &ip.IPTableAddDel{
		IsAdd: true,
		Table: ip.IPTable{
			TableID: id,
			IsIP6:   true,
			Name:    name,
		},
	}

	_, err = vpp.DoRequest[*ip.IPTableAddDel, *ip.IPTableAddDelReply](s.client, ctx, ipv6Req)
	if err != nil {
		defer s.deleteVRF(ctx, id, name, false)
		return fmt.Errorf("failed to create IPv6 VRF table: %w", err)
	}
	s.addOrUpdateVRFCache(id, name)

	return nil
}

func (s *Service) deleteVRF(ctx context.Context, id uint32, name string, isIPv6 bool) error {
	req := &ip.IPTableAddDel{
		IsAdd: false,
		Table: ip.IPTable{
			TableID: id,
			IsIP6:   isIPv6,
			Name:    name,
		},
	}

	_, err := vpp.DoRequest[*ip.IPTableAddDel, *ip.IPTableAddDelReply](s.client, ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete VRF table: %w", err)
	}
	return nil
}

func (s *Service) enrichVRFWithNames(vrfMap map[uint32]*domain.VRF) {
	s.vrfCache.mu.RLock()
	defer s.vrfCache.mu.RUnlock()

	for id, vrf := range vrfMap {
		if entry, exists := s.vrfCache.data[id]; exists {
			vrf.Name = entry.Name
		}
	}
}
func (s *Service) addOrUpdateVRFCache(id uint32, name string) {
	s.vrfCache.mu.Lock()
	defer s.vrfCache.mu.Unlock()

	s.vrfCache.data[id] = &domain.VRFEntry{
		Name: name,
	}
}

func (s *Service) DeleteVRF(ctx context.Context, id uint32) error {
	vfrEntry, err := s.getEntryVRFCache(id)
	if err != nil {
		return err
	}
	var errs []error

	errIPv4 := s.deleteVRF(ctx, id, vfrEntry.Name, false)
	if errIPv4 != nil {
		errs = append(errs, fmt.Errorf("IPv4: %w", errIPv4))
	}

	errIPv6 := s.deleteVRF(ctx, id, vfrEntry.Name, true)
	if errIPv6 != nil {
		errs = append(errs, fmt.Errorf("IPv6: %w", errIPv6))
	}

	if err = errors.Join(errs...); err != nil {
		s.deleteVRFCache(id)
		return fmt.Errorf("failed to delete VRF tables: %w", err)
	}

	s.deleteVRFCache(id)

	return nil
}

func (s *Service) ListVRF(ctx context.Context) ([]domain.VRF, error) {
	vrfMap := make(map[uint32]*domain.VRF)

	if err := s.collectRouteStats(ctx, vrfMap); err != nil {
		return nil, fmt.Errorf("route stats: %w", err)
	}

	s.addVRFsFromCache(vrfMap)
	s.enrichVRFWithNames(vrfMap)

	result := make([]domain.VRF, 0, len(vrfMap))
	for _, vrf := range vrfMap {
		result = append(result, *vrf)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result, nil
}

func (s *Service) collectRouteStats(ctx context.Context, vrfMap map[uint32]*domain.VRF) error {
	s.vrfCache.mu.RLock()
	vrfIDs := make([]uint32, 0, len(s.vrfCache.data))
	for id := range s.vrfCache.data {
		vrfIDs = append(vrfIDs, id)
	}
	s.vrfCache.mu.RUnlock()

	for _, id := range vrfIDs {
		if err := s.collectIPRoutes(ctx, vrfMap, id, false); err != nil {
			return fmt.Errorf("IPv4: %w", err)
		}

		if err := s.collectIPRoutes(ctx, vrfMap, id, true); err != nil {
			return fmt.Errorf("IPv6: %w", err)
		}
	}

	return nil
}

func (s *Service) collectIPRoutes(ctx context.Context, vrfMap map[uint32]*domain.VRF, vrfID uint32, isIPv6 bool) error {
	request := &ip.IPRouteDump{
		Table: ip.IPTable{
			TableID: vrfID,
			IsIP6:   isIPv6,
		},
	}

	converter := func(msg api.Message) (interface{}, bool) {
		details, ok := msg.(*ip.IPRouteDetails)
		if !ok {
			return nil, false
		}
		if IsSystemRoute(details) {
			return nil, false
		}
		vrfID := details.Route.TableID

		vrf, exists := vrfMap[vrfID]
		if !exists {
			vrf = &domain.VRF{
				ID:   vrfID,
				IPv4: !isIPv6,
				IPv6: isIPv6,
			}
			vrfMap[vrfID] = vrf
		} else {
			if isIPv6 {
				vrf.IPv6 = true
			} else {
				vrf.IPv4 = true
			}
		}

		vrf.RouteCount++
		return nil, false
	}

	_, err := vpp.Dump(ctx, s.client, request, converter)
	return err
}
