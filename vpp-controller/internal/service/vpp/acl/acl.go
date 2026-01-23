package acl

import (
	"context"
	"fmt"

	"github.com/NikolayStepanov/RapidVPP/internal/domain"
	"github.com/NikolayStepanov/RapidVPP/internal/infrastructure/vpp"
	"github.com/NikolayStepanov/RapidVPP/internal/mapper"
	"go.fd.io/govpp/binapi/acl"
)

type Service struct {
	client *vpp.Client
}

func (s Service) Create(ctx context.Context, name string, rules []domain.ACLRule) (domain.AclID, error) {
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

func (s Service) Update(ctx context.Context, id domain.AclID, rules []domain.ACLRule) error {
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

func (s Service) Delete(ctx context.Context, id domain.AclID) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) Get(ctx context.Context, id domain.AclID) (domain.ACLInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) List(ctx context.Context) ([]domain.ACLInfo, error) {
	//TODO implement me
	panic("implement me")
}

func NewService(client *vpp.Client) *Service {
	return &Service{client: client}

}
