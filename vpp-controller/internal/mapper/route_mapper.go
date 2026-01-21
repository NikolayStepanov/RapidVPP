package mapper

import (
	"errors"
	"fmt"
	"net"

	"github.com/NikolayStepanov/RapidVPP/internal/domain"
	"go.fd.io/govpp/binapi/fib_types"
	"go.fd.io/govpp/binapi/ip"
	"go.fd.io/govpp/binapi/ip_types"
)

func BuildFibPaths(nextHops []domain.NextHop) ([]fib_types.FibPath, error) {
	paths := make([]fib_types.FibPath, 0, len(nextHops))

	for _, nh := range nextHops {
		path, err := BuildFibPath(nh)
		if err != nil {
			return nil, fmt.Errorf("next-hop %v: %w", nh, err)
		}
		paths = append(paths, path)
	}

	return paths, nil
}

func BuildFibPath(nh domain.NextHop) (fib_types.FibPath, error) {
	if nh.Drop {
		return fib_types.FibPath{
			Type:   fib_types.FIB_API_PATH_TYPE_DROP,
			Weight: nh.Weight,
		}, nil
	}

	if nh.IP == nil {
		return fib_types.FibPath{}, errors.New("next-hop ip is nil")
	}

	proto := fib_types.FIB_API_PATH_NH_PROTO_IP4
	if nh.IP.To4() == nil {
		proto = fib_types.FIB_API_PATH_NH_PROTO_IP6
	}

	return fib_types.FibPath{
		SwIfIndex: nh.IfIndex,
		Weight:    nh.Weight,
		Type:      fib_types.FIB_API_PATH_TYPE_NORMAL,
		Proto:     proto,
		Nh: fib_types.FibPathNh{
			Address: ip_types.NewAddress(nh.IP).Un,
		},
	}, nil
}

func ConvertRouteDetails(details *ip.IPRouteDetails) (domain.Route, error) {
	route := domain.Route{
		VRF: details.Route.TableID,
	}

	prefix := details.Route.Prefix
	addr := prefix.Address.ToIP()

	route.Dst = domain.IPWithPrefix{
		Address: addr.String(),
		Prefix:  prefix.Len,
	}

	paths := details.Route.Paths
	route.NextHops = make([]domain.NextHop, len(paths))

	for i, path := range paths {
		nextHop, err := ConvertFibPathToDomainNextHop(path)
		if err != nil {
			return domain.Route{}, fmt.Errorf("failed to convert FibPath to domain.NextHop: %w", err)
		}
		route.NextHops[i] = nextHop
	}

	return route, nil
}

func ConvertFibPathToDomainNextHop(path fib_types.FibPath) (domain.NextHop, error) {
	isDrop := path.Type == fib_types.FIB_API_PATH_TYPE_DROP

	nextHop := domain.NextHop{
		IfIndex: path.SwIfIndex,
		Weight:  path.Weight,
		Drop:    isDrop,
	}

	if isDrop {
		return nextHop, nil
	}

	if path.Type != fib_types.FIB_API_PATH_TYPE_NORMAL {
		return domain.NextHop{}, fmt.Errorf("unsupported path type: %v", path.Type)
	}

	var ip net.IP

	switch path.Proto {
	case fib_types.FIB_API_PATH_NH_PROTO_IP4:
		ip4 := path.Nh.Address.GetIP4()
		if ip4 != [4]uint8{} {
			ip = ip4[:]
		}
	case fib_types.FIB_API_PATH_NH_PROTO_IP6:
		ip6 := path.Nh.Address.GetIP6()
		if ip6 != [16]uint8{} {
			ip = ip6[:]
		}
	default:
		return domain.NextHop{}, fmt.Errorf("unsupported protocol: %v", path.Proto)
	}

	nextHop.IP = ip
	return nextHop, nil
}
