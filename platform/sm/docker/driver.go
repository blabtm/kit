package docker

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"

	"github.com/blabtm/v2k/platform/sm/service"
	"github.com/compose-spec/compose-go/v2/cli"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/compose/v5/pkg/compose"
)

func init() {
	service.Registry["docker"] = New()
}

type Driver struct {
	api api.Compose
	mux map[string]*sync.Mutex
}

func New() *Driver {
	cli, err := command.NewDockerCli()

	if err != nil {
		log.Panicf("client: %v", err)
	}

	if err := cli.Initialize(&flags.ClientOptions{}); err != nil {
		log.Panicf("client: initialize: %v", err)
	}

	api, err := compose.NewComposeService(cli)

	if err != nil {
		log.Panicf("compose: %v", err)
	}

	return &Driver{
		api: api,
		mux: make(map[string]*sync.Mutex),
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

func (prv *Driver) Ps(ctx context.Context, spec *service.Spec) (*service.Status, error) {
	mux := prv.getLock(spec.Name)
	mux.Lock()
	defer mux.Unlock()

	name := strings.ReplaceAll(spec.Name, ".", "_")
	prj, err := prv.api.LoadProject(ctx, api.ProjectLoadOptions{
		ConfigPaths: []string{filepath.Join(spec.DeployPath, "art", "compose.yaml")},
		ProjectName: name,
	})

	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	sum, err := prv.api.Ps(ctx, name, api.PsOptions{
		Project: prj,
		All:     true,
	})

	if err != nil {
		return nil, fmt.Errorf("ps: %w", err)
	}

	if len(sum) == 0 {
		return &service.Status{State: service.Down}, nil
	}

	for _, con := range sum {
		if con.State != "running" {
			prv.api.Down(ctx, name, api.DownOptions{
				Project: prj,
			})

			return &service.Status{State: service.Down}, nil
		}
	}

	return &service.Status{State: service.Up}, nil
}

func (prv *Driver) Up(ctx context.Context, spec *service.Spec, opts ...service.Option) error {
	mux := prv.getLock(spec.Name)
	mux.Lock()
	defer mux.Unlock()

	env := make([]string, 0, len(opts))

	for _, opt := range opts {
		env = append(env, fmt.Sprintf("%s:%v", strings.ToUpper(opt.Key), opt.Value))
	}

	name := strings.ReplaceAll(spec.Name, ".", "_")
	prj, err := prv.api.LoadProject(ctx, api.ProjectLoadOptions{
		ConfigPaths: []string{filepath.Join(spec.DeployPath, "art", "compose.yaml")},
		ProjectName: name,
		ProjectOptionsFns: []cli.ProjectOptionsFn{
			cli.WithEnv(env),
		},
	})

	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	if err := prv.api.Up(ctx, prj, api.UpOptions{
		Create: api.CreateOptions{},
		Start:  api.StartOptions{},
	}); err != nil {
		return fmt.Errorf("up: %w", err)
	}

	return nil
}

func (drv *Driver) Down(ctx context.Context, spec *service.Spec) error {
	mux := drv.getLock(spec.Name)
	mux.Lock()
	defer mux.Unlock()

	name := strings.ReplaceAll(spec.Name, ".", "_")
	prj, err := drv.api.LoadProject(ctx, api.ProjectLoadOptions{
		ConfigPaths: []string{filepath.Join(spec.DeployPath, "art", "compose.yaml")},
		ProjectName: name,
	})

	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	if err := drv.api.Down(ctx, name, api.DownOptions{
		Project: prj,
	}); err != nil {
		return fmt.Errorf("down: %w", err)
	}

	return nil
}
