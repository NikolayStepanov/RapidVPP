package mapper

import (
	"fmt"
	"net"

	"github.com/NikolayStepanov/RapidVPP/internal/domain"
	"go.fd.io/govpp/binapi/acl_types"
	"go.fd.io/govpp/binapi/ip_types"
)

func ConvertACLRule(aclRule domain.ACLRule) (acl_types.ACLRule, error) {
	srcPrefix, err := IPWithPrefixToTypes(aclRule.Src)
	if err != nil {
		return acl_types.ACLRule{}, fmt.Errorf("invalid src prefix: %w", err)
	}

	dstPrefix, err := IPWithPrefixToTypes(aclRule.Dst)
	if err != nil {
		return acl_types.ACLRule{}, fmt.Errorf("invalid dst prefix: %w", err)
	}

	return acl_types.ACLRule{
		IsPermit: acl_types.ACLAction(aclRule.Action),

		SrcPrefix: srcPrefix,
		DstPrefix: dstPrefix,

		Proto: ip_types.IPProto(aclRule.Proto),

		SrcportOrIcmptypeFirst: aclRule.SrcPortLow,
		SrcportOrIcmptypeLast:  aclRule.SrcPortHigh,

		DstportOrIcmpcodeFirst: aclRule.DstPortLow,
		DstportOrIcmpcodeLast:  aclRule.DstPortHigh,

		TCPFlagsMask:  aclRule.TCPFlagsMask,
		TCPFlagsValue: aclRule.TCPFlagsValue,
	}, nil
}

func IPWithPrefixToTypes(prefix domain.IPWithPrefix) (ip_types.Prefix, error) {
	if prefix.Address == "" {
		return ip_types.Prefix{}, fmt.Errorf("empty ip address")
	}

	ip := net.ParseIP(prefix.Address)
	if ip == nil {
		return ip_types.Prefix{}, fmt.Errorf("invalid ip address: %s", prefix.Address)
	}

	var addr ip_types.Address

	if ip.To4() != nil {
		addr = ip_types.NewAddress(ip.To4())
	} else {
		addr = ip_types.NewAddress(ip)
	}

	return ip_types.Prefix{
		Address: addr,
		Len:     prefix.Prefix,
	}, nil
}

func IPWithPrefixFromTypes(prefix ip_types.Prefix) (domain.IPWithPrefix, error) {
	var ip net.IP
	ip = prefix.Address.ToIP()

	if ip == nil {
		return domain.IPWithPrefix{}, fmt.Errorf("invalid ip address")
	}

	return domain.IPWithPrefix{
		Address: ip.String(),
		Prefix:  prefix.Len,
	}, nil
}

func ConvertACLRules(rules []domain.ACLRule) ([]acl_types.ACLRule, error) {
	aclRules := make([]acl_types.ACLRule, 0, len(rules))
	for _, rule := range rules {
		ruleACL, err := ConvertACLRule(rule)
		if err != nil {
			return nil, err
		}
		aclRules = append(aclRules, ruleACL)
	}
	return aclRules, nil
}
