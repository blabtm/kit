package systemd

import (
	"context"

	"github.com/blabtm/v2k/platform/sm/service"
)

func init() {
	service.Registry["systemd"] = New()
}

type Driver struct{}

func New() *Driver {
	return &Driver{}
}

func (prv *Driver) Ps(ctx context.Context, spec *service.Spec) (*service.Status, error) {
	return &service.Status{State: service.Stopped}, nil
}

func (prv *Driver) Up(ctx context.Context, spec *service.Spec, opts ...service.Option) error {
	return nil
}

func (drv *Driver) Down(ctx context.Context, spec *service.Spec) error {
	return nil
}
