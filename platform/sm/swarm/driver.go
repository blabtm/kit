package swarm

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/blabtm/v2k/platform/sm/service"
)

func init() {
	service.Registry["swarm"] = New()
}

type Driver struct {
	mux map[string]*sync.Mutex
}

func New() *Driver {
	return &Driver{
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

	return nil, fmt.Errorf("unsupported")
}

func (prv *Driver) Up(ctx context.Context, spec *service.Spec, opts ...service.Option) error {
	mux := prv.getLock(spec.Name)
	mux.Lock()
	defer mux.Unlock()

	name := strings.ReplaceAll(spec.Name, ".", "_")
	path := filepath.Join(spec.DeployPath, "art", "compose.yaml")
	env := make([]string, 0, len(opts))

	for _, opt := range opts {
		env = append(env, fmt.Sprintf("%s:%v", strings.ToUpper(opt.Key), opt.Value))
	}

	return Deploy(name, path, env)
}

func (drv *Driver) Down(ctx context.Context, spec *service.Spec) error {
	mux := drv.getLock(spec.Name)
	mux.Lock()
	defer mux.Unlock()

	name := strings.ReplaceAll(spec.Name, ".", "_")

	return Rm(name)
}
