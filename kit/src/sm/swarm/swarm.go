package swarm

import (
	"context"
	"os/exec"

	"github.com/blabtm/v2k/platform/sm/service"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/compose/convert"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
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

// Ps loads swarm stack status and maps it to the SM service state.
//
// Swarm stack composed of multiple swarm services, while each swarm service
// contains multiple swarm tasks.
//
// Swarm stack is mapped to multi-task SM service, where each swarm service is
// an SM task. Rules about state transitioning apply recursively to swarm services:
// swarm service becomes multi-task SM service, where each swarm task is an SM task.
func (sw *Swarm) Ps(ctx context.Context, name string) (service.State, error) {
	api := sw.Cli.Client()

	res, err := api.ServiceList(ctx, swarm.ServiceListOptions{
		Filters: getStackPsArgs(name),
		Status:  true,
	})

	if err != nil {
		return service.Stopped, err
	}

	if len(res) == 0 {
		return service.Stopped, nil
	}

	total := len(res)
	pending := 0
	running := 0
	completed := 0
	failed := 0

	for _, svc := range res {
		state, err := getServiceState(api, &svc)

		if err != nil {
			return service.Stopped, err
		}

		switch state {
		case service.Pending:
			pending += 1
		case service.Running:
			running += 1
		case service.Completed:
			completed += 1
		case service.Failed:
			failed += 1
		}
	}

	if completed+failed == total {
		if failed != 0 {
			return service.Failed, nil
		}

		return service.Completed, nil
	}

	if running == total {
		return service.Running, nil
	}

	return service.Pending, nil
}

func getStackPsArgs(name string) filters.Args {
	return filters.NewArgs(
		filters.Arg("label", convert.LabelNamespace+"="+name),
	)
}

func getServiceState(cli client.APIClient, svc *swarm.Service) (service.State, error) {
	tasks, err := cli.TaskList(context.Background(), swarm.TaskListOptions{
		Filters: getServicePsArgs(svc.Spec.Name),
	})

	if err != nil {
		return service.Stopped, err
	}

	total := len(tasks)
	pending := 0
	running := 0
	completed := 0
	failed := 0

	for _, task := range tasks {
		switch task.Status.State {
		case swarm.TaskStateNew,
			swarm.TaskStateAllocated,
			swarm.TaskStatePending,
			swarm.TaskStateAssigned,
			swarm.TaskStateAccepted,
			swarm.TaskStatePreparing,
			swarm.TaskStateReady,
			swarm.TaskStateStarting,
			swarm.TaskStateRejected,
			swarm.TaskStateRemove,
			swarm.TaskStateOrphaned:
			pending += 1
		case swarm.TaskStateRunning:
			running += 1
		case swarm.TaskStateComplete:
			completed += 1
		case swarm.TaskStateFailed:
			failed += 1
		}
	}

	if completed+failed == total {
		if failed != 0 {
			return service.Failed, nil
		}

		return service.Completed, nil
	}

	if running == total {
		return service.Running, nil
	}

	return service.Pending, nil
}

func getServicePsArgs(name string) filters.Args {
	return filters.NewArgs(
		filters.Arg("service", name),
	)
}
