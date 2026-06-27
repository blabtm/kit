package swarm

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"

	"github.com/blabtm/v2k/platform/sm/service"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
)

func init() {
	service.Registry["swarm"] = New()
}

type Driver struct {
	Swarm
	gmux sync.Mutex
	smux map[string]*sync.Mutex
}

func New() *Driver {
	cli, err := command.NewDockerCli()

	if err != nil {
		log.Panic(err)
	}

	if err := cli.Initialize(&flags.ClientOptions{}); err != nil {
		log.Panic(err)
	}

	return &Driver{
		Swarm: Swarm{cli},
		smux:  make(map[string]*sync.Mutex),
	}
}

func (drv *Driver) getLock(name string) *sync.Mutex {
	drv.gmux.Lock()
	defer drv.gmux.Unlock()

	mux, ok := drv.smux[name]

	if !ok {
		mux = &sync.Mutex{}
		drv.smux[name] = mux
	}

	return mux
}

func (drv *Driver) Ps(ctx context.Context, spec *service.Spec) (*service.Status, error) {
	mux := drv.getLock(spec.Name)
	mux.Lock()
	defer mux.Unlock()

	state, err := drv.Swarm.Ps(ctx, spec.String())

	if err != nil {
		return nil, err
	}

	return &service.Status{State: state}, nil
}

func (drv *Driver) Up(ctx context.Context, spec *service.Spec, opts ...service.Option) error {
	mux := drv.getLock(spec.Name)
	mux.Lock()
	defer mux.Unlock()

	path := filepath.Join(spec.DeployPath, "art", "compose.yaml")
	env := make([]string, 0, len(opts))

	for _, opt := range opts {
		env = append(env, fmt.Sprintf("%s=%v", strings.ToUpper(opt.Key), opt.Value))
	}

	if spec.Type == service.Oneshot {
		env = append(env, fmt.Sprintf("IID=%s", spec.IID.Short()))
	}

	return drv.Swarm.Deploy(spec.String(), path, env)
}

func (drv *Driver) Down(ctx context.Context, spec *service.Spec) error {
	mux := drv.getLock(spec.Name)
	mux.Lock()
	defer mux.Unlock()

	return drv.Swarm.Rm(spec.String())
}
