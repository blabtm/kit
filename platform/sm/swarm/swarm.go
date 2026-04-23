package swarm

import (
	"os/exec"
)

func Deploy(name, path string, env []string) error {
	cmd := exec.Command("docker", "stack", "deploy", "-c", path, name)
	cmd.Env = append(cmd.Env, env...)
	return cmd.Run()
}

func Rm(name string) error {
	cmd := exec.Command("docker", "stack", "rm", name)
	return cmd.Run()
}
