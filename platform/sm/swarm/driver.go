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
	mux map[string]*sync.Mutex
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
		mux:   make(map[string]*sync.Mutex),
	}
}

func (prv *Driver) getLock(name string) *sync.Mutex {
	mux, ok := prv.mux[name]

	if !ok {
		mux = &sync.Mutex{}
		prv.mux[name] = mux
	}

	return mux
}

func (drv *Driver) Ps(ctx context.Context, spec *service.Spec) (*service.Status, error) {
	mux := drv.getLock(spec.Name)
	mux.Lock()
	defer mux.Unlock()

	name := strings.ReplaceAll(spec.Name, ".", "_")
	state, err := drv.Swarm.Ps(ctx, name)

	if err != nil {
		return nil, err
	}

	return &service.Status{State: state}, nil
}

func (drv *Driver) Up(ctx context.Context, spec *service.Spec, opts ...service.Option) error {
	mux := drv.getLock(spec.Name)
	mux.Lock()
	defer mux.Unlock()

	name := strings.ReplaceAll(spec.Name, ".", "_")
	path := filepath.Join(spec.DeployPath, "art", "compose.yaml")
	env := make([]string, 0, len(opts))

	for _, opt := range opts {
		env = append(env, fmt.Sprintf("%s:%v", strings.ToUpper(opt.Key), opt.Value))
	}

	return drv.Swarm.Deploy(name, path, env)
}

func (drv *Driver) Down(ctx context.Context, spec *service.Spec) error {
	mux := drv.getLock(spec.Name)
	mux.Lock()
	defer mux.Unlock()

	name := strings.ReplaceAll(spec.Name, ".", "_")

	return drv.Swarm.Rm(name)
}
