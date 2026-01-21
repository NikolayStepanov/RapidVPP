package ip

import (
	"fmt"
	"net"

	"github.com/NikolayStepanov/RapidVPP/internal/domain"
)

func (r *AddRouteRequest) ToDomain() (*domain.Route, error) {
	ip, ipnet, err := net.ParseCIDR(r.Destination)
	if err != nil {
		return nil, fmt.Errorf("parse destination: %w", err)
	}

	prefix, _ := ipnet.Mask.Size()

	nextHops := make([]domain.NextHop, 0, len(r.NextHops))
	for _, nh := range r.NextHops {
		var nextHopIP net.IP
		if nh.IP != "" {
			nextHopIP = net.ParseIP(nh.IP)
			if nextHopIP == nil {
				return nil, fmt.Errorf("invalid next-hop IP: %s", nh.IP)
			}
		}

		nextHops = append(nextHops, domain.NextHop{
			IP:      nextHopIP,
			IfIndex: nh.IfIndex,
			Weight:  nh.Weight,
			Drop:    nh.Drop,
		})
	}

	return &domain.Route{
		Dst: domain.IPWithPrefix{
			Address: ip.String(),
			Prefix:  uint8(prefix),
		},
		VRF:      r.VRF,
		NextHops: nextHops,
	}, nil
}

func VRFToResponse(domains []domain.VRF) []VRFResponse {
	res := make([]VRFResponse, 0, len(domains))

	for _, v := range domains {
		res = append(res, VRFResponse{
			ID:         v.ID,
			Name:       v.Name,
			IPv4:       v.IPv4,
			IPv6:       v.IPv6,
			RouteCount: v.RouteCount,
		})
	}

	return res
}
