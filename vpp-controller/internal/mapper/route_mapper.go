package mapper

import (
	"errors"
	"fmt"

	"github.com/NikolayStepanov/RapidVPP/internal/domain"
	"go.fd.io/govpp/binapi/fib_types"
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
