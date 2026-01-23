package acl

import (
	"context"
	"fmt"

	"github.com/NikolayStepanov/RapidVPP/internal/domain"
	"github.com/NikolayStepanov/RapidVPP/internal/infrastructure/vpp"
	"github.com/NikolayStepanov/RapidVPP/internal/mapper"
	"go.fd.io/govpp/api"
	"go.fd.io/govpp/binapi/acl"
)

type Service struct {
	client *vpp.Client
}

func (s *Service) Create(ctx context.Context, name string, rules []domain.ACLRule) (domain.AclID, error) {
	if len(rules) == 0 {
		return 0, fmt.Errorf("acl must contain at least one rule")
	}
	vppRules, err := mapper.ConvertACLRules(rules)
	if err != nil {
		return 0, err
	}
	req := &acl.ACLAddReplace{
		ACLIndex: 0xFFFFFFFF,
		Tag:      name,
		R:        vppRules,
	}

	reply, err := vpp.DoRequest[*acl.ACLAddReplace, *acl.ACLAddReplaceReply](s.client, ctx, req)
	if err != nil {
		return 0, fmt.Errorf("create acl operation failed: %w", err)
	}

	return domain.AclID(reply.ACLIndex), nil
}

func (s *Service) Update(ctx context.Context, id domain.AclID, rules []domain.ACLRule) error {
	if len(rules) == 0 {
		return fmt.Errorf("acl must contain at least one rule")
	}

	vppRules, err := mapper.ConvertACLRules(rules)
	if err != nil {
		return err
	}

	req := &acl.ACLAddReplace{
		ACLIndex: uint32(id),
		R:        vppRules,
	}

	_, err = vpp.DoRequest[*acl.ACLAddReplace, *acl.ACLAddReplaceReply](s.client, ctx, req)
	if err != nil {
		return fmt.Errorf("update acl operation failed: %w", err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id domain.AclID) error {
	req := &acl.ACLDel{
		ACLIndex: uint32(id),
	}

	_, err := vpp.DoRequest[*acl.ACLDel, *acl.ACLDelReply](s.client, ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete ACL %d: %w", id, err)
	}

	return nil
}

func (s *Service) List(ctx context.Context) ([]domain.ACLInfo, error) {
	request := &acl.ACLDump{}

	converter := func(msg api.Message) (domain.ACLInfo, bool) {
		details, ok := msg.(*acl.ACLDetails)
		if !ok {
			return domain.ACLInfo{}, false
		}

		rules, err := mapper.ConvertVPPACLRules(details.R)
		if err != nil {
			return domain.ACLInfo{}, false
		}
		return domain.ACLInfo{
			ID:    domain.AclID(details.ACLIndex),
			Name:  details.Tag,
			Rules: rules,
		}, true
	}

	return vpp.Dump(ctx, s.client, request, converter)
}

func NewService(client *vpp.Client) *Service {
	return &Service{client: client}

}
