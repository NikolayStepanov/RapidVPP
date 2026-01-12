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
}
type Services struct {
	Info      Info
	Interface Interface
}

func NewServices(info Info, inter Interface) *Services {
	return &Services{info, inter}
}
