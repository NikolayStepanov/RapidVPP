package service

import (
	"context"

	"github.com/NikolayStepanov/RapidVPP/internal/domain"
	interfaces "go.fd.io/govpp/binapi/interface"
)

type Info interface {
	GetVersion(ctx context.Context) (domain.Version, error)
}
type Interface interface {
	List(ctx context.Context) ([]interfaces.SwInterfaceDetails, error)
	CreateLoopback(ctx context.Context) (interfaces.CreateLoopbackReply, error)
	DeleteLoopback(ctx context.Context, ifIndex uint32) error
	SetInterfaceAdminState(ctx context.Context, ifIndex uint32, up bool) error
	SetInterfaceIP(ctx context.Context, ifIndex uint32, IPPrefix domain.IPWithPrefix) error
	AttachACL(ctx context.Context, ifIndex uint32, aclID uint32, dir uint8) error
	DetachACL(ctx context.Context, ifIndex uint32, aclID uint32, dir uint8) error
}

type Route interface {
	AddRoute(ctx context.Context, route *domain.Route) error
	DeleteRoute(ctx context.Context, route *domain.Route) error
	ListRoutes(ctx context.Context, vrf uint32) ([]domain.Route, error)
	GetRoute(ctx context.Context, dst domain.IPWithPrefix, vrf uint32) (domain.Route, error)
}

type VRF interface {
	InitVRFCache(ctx context.Context) error
	CreateVRF(ctx context.Context, id uint32, name string) error
	DeleteVRF(ctx context.Context, id uint32) error
	ListVRF(ctx context.Context) ([]domain.VRF, error)
}

type ACL interface {
	Create(ctx context.Context, name string, rules []domain.ACLRule) (domain.AclID, error)
	Update(ctx context.Context, id domain.AclID, rules []domain.ACLRule) error
	Delete(ctx context.Context, id domain.AclID) error
	List(ctx context.Context) ([]domain.ACLInfo, error)
}

type IP interface {
	Route
	VRF
}

type Services struct {
	Info      Info
	Interface Interface
	IP        IP
	ACL       ACL
}

func NewServices(info Info, inter Interface, IPService IP, acl ACL) *Services {
	return &Services{
		info,
		inter,
		IPService,
		acl,
	}
}
