package info

import (
	"context"
	"fmt"

	"github.com/NikolayStepanov/RapidVPP/internal/domain"
	"github.com/NikolayStepanov/RapidVPP/internal/infrastructure/vpp"
	"go.fd.io/govpp/binapi/vpe"
)

type Service struct {
	client *vpp.Client
}

func NewService(client *vpp.Client) *Service {
	return &Service{client: client}
}

func (s *Service) GetVersion(ctx context.Context) (domain.Version, error) {
	req := &vpe.ShowVersion{}

	reply, err := vpp.DoRequest[*vpe.ShowVersion, *vpe.ShowVersionReply](s.client, ctx, req)
	if err != nil {
		return domain.Version{}, fmt.Errorf("get version operation failed: %w", err)
	}

	return domain.Version{
		Version:   reply.Version,
		BuildDate: reply.BuildDate,
		BuildDir:  reply.BuildDirectory,
	}, nil
}
