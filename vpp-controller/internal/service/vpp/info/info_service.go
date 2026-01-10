package info

import (
	"fmt"
	"time"

	"github.com/NikolayStepanov/RapidVPP/internal/domain"
	"github.com/NikolayStepanov/RapidVPP/internal/infrastructure/vpp"
	"go.fd.io/govpp/api"
	"go.fd.io/govpp/binapi/vpe"
)

type Service struct {
	client *vpp.Client
}

func (i *Service) GetVersion() (domain.Version, error) {
	var info domain.Version

	err := i.client.DoWithTimeout(5*time.Second, func(stream api.Stream) error {
		req := &vpe.ShowVersion{}

		if err := stream.SendMsg(req); err != nil {
			return fmt.Errorf("send request: %w", err)
		}

		msg, err := stream.RecvMsg()
		if err != nil {
			return fmt.Errorf("receive reply: %w", err)
		}

		reply, ok := msg.(*vpe.ShowVersionReply)
		if !ok {
			return fmt.Errorf("unexpected message type: %T, expected *vpe.ShowVersionReply", msg)
		}
		info = domain.Version{
			Version:   reply.Version,
			BuildDate: reply.BuildDate,
			BuildDir:  reply.BuildDirectory,
		}
		if err != nil {
			return fmt.Errorf("create version domain object: %w", err)
		}
		return nil
	})

	return info, err
}

func NewService(client *vpp.Client) *Service {
	return &Service{client: client}
}
