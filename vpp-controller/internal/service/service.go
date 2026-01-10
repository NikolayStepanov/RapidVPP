package service

import "github.com/NikolayStepanov/RapidVPP/internal/domain"

type Info interface {
	GetVersion() (domain.Version, error)
}
type Services struct {
	Info Info
}

func NewServices(info Info) *Services {
	return &Services{info}
}
