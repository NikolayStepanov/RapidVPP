package acl

import (
	"fmt"

	"github.com/NikolayStepanov/RapidVPP/internal/domain"
)

func ConvertRulesRequestToDomain(rules []RulesRequest) ([]domain.ACLRule, error) {
	if len(rules) == 0 {
		return nil, fmt.Errorf("rules list is empty")
	}

	out := make([]domain.ACLRule, 0, len(rules))

	for i, r := range rules {
		if r.Src.Address == "" {
			return nil, fmt.Errorf("rule %d: source address is empty", i)
		}
		if r.Dst.Address == "" {
			return nil, fmt.Errorf("rule %d: destination address is empty", i)
		}

		out = append(out, domain.ACLRule{
			Action:        domain.ACLAction(r.Action), // 0=deny,1=permit
			Proto:         r.Proto,
			Src:           domain.IPWithPrefix{Address: r.Src.Address, Prefix: r.Src.Prefix},
			Dst:           domain.IPWithPrefix{Address: r.Dst.Address, Prefix: r.Dst.Prefix},
			SrcPortLow:    r.SrcPortLow,
			SrcPortHigh:   r.SrcPortHigh,
			DstPortLow:    r.DstPortLow,
			DstPortHigh:   r.DstPortHigh,
			TCPFlagsMask:  r.TCPFlagsMask,
			TCPFlagsValue: r.TCPFlagsValue,
		})
	}

	return out, nil
}

func AclInfoToResponse(info domain.ACLInfo) ACLResponse {
	rules := make([]RulesResponse, len(info.Rules))

	for i, r := range info.Rules {
		rules[i] = RulesResponse{
			Action:        uint8(r.Action),
			Proto:         r.Proto,
			Src:           IPWithPrefix{Address: r.Src.Address, Prefix: r.Src.Prefix},
			Dst:           IPWithPrefix{Address: r.Dst.Address, Prefix: r.Dst.Prefix},
			SrcPortLow:    r.SrcPortLow,
			SrcPortHigh:   r.SrcPortHigh,
			DstPortLow:    r.DstPortLow,
			DstPortHigh:   r.DstPortHigh,
			TCPFlagsMask:  r.TCPFlagsMask,
			TCPFlagsValue: r.TCPFlagsValue,
		}
	}

	return ACLResponse{
		ID:    uint32(info.ID),
		Name:  info.Name,
		Rules: rules,
	}
}

func InfosToResponse(acls []domain.ACLInfo) []ACLResponse {
	responses := make([]ACLResponse, len(acls))
	for i, acl := range acls {
		responses[i] = AclInfoToResponse(acl)
	}
	return responses
}
