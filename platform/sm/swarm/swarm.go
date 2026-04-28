package swarm

import (
	"context"
	"os/exec"

	"github.com/blabtm/v2k/platform/sm/service"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/compose/convert"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
)

type Swarm struct {
	Cli command.Cli
}

func (*Swarm) Deploy(name, path string, env []string) error {
	cmd := exec.Command("docker", "stack", "deploy", "-c", path, name)
	cmd.Env = append(cmd.Env, env...)
	return cmd.Run()
}

func (*Swarm) Rm(name string) error {
	cmd := exec.Command("docker", "stack", "rm", name)
	return cmd.Run()
}

func getPsArgs(name string) filters.Args {
	return filters.NewArgs(
		filters.Arg("label", convert.LabelNamespace+"="+name),
	)
}

func (sw *Swarm) Ps(ctx context.Context, name string) (service.State, error) {
	api := sw.Cli.Client()

	res, err := api.ServiceList(ctx, swarm.ServiceListOptions{
		Filters: getPsArgs(name),
		Status:  true,
	})

	if err != nil {
		return service.Down, err
	}

	total := 0
	up := 0

	for _, service := range res {
		total += 1

		if service.ServiceStatus.RunningTasks == service.ServiceStatus.DesiredTasks {
			up += 1
		}
	}

	if up == 0 {
		return service.Down, nil
	}

	if up == total {
		return service.Up, nil
	}

	return service.Partial, nil
}
