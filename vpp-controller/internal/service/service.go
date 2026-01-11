package service

import (
	"github.com/NikolayStepanov/RapidVPP/internal/domain"
	interfaces "go.fd.io/govpp/binapi/interface"
)

type Info interface {
	GetVersion() (domain.Version, error)
}
type Interface interface {
	List() ([]interfaces.SwInterfaceDetails, error)
}
type Services struct {
	Info      Info
	Interface Interface
}

func NewServices(info Info, inter Interface) *Services {
	return &Services{info, inter}
}
