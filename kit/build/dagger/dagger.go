package dagger

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/blabtm/kit/build"
)

type Dagger struct {
	*dagger.Client
}

func New() (*Dagger, error) {
	dag, err := dagger.Connect(context.Background(),
		dagger.WithLogOutput(os.Stdout))

	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	return &Dagger{dag}, nil
}

func (dag *Dagger) Build(mod *build.Module) error {
	switch mod.Toolchain {
	case build.Native:
		return dag.BuildNative(mod)
	}

	return fmt.Errorf("unsupported toolchain")
}

func (dag *Dagger) BuildNative(mod *build.Module) error {
	src := dag.Host().Directory(mod.ToolchainDir)
	pkgs := []string{
		"build-essential",
		"cmake",
	}

	builder := dag.Container().
		From("ubuntu:latest").
		WithExec([]string{"apt-get", "update"}).
		WithExec(append([]string{"apt-get", "install", "-y"}, pkgs...)).
		WithDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"mkdir", "-p", "build"}).
		WithExec([]string{"cmake", "-S", ".", "-B", "build"}).
		WithExec([]string{"cmake", "--build", "build", "--target", mod.Name})

	image := dag.Container().
		From("alpine:latest").
		WithFile("/bin/run", builder.File(fmt.Sprintf("/src/build/%s/%s", mod.Name, mod.Name))).
		WithEntrypoint([]string{"/bin/run"})

	addr, err := image.Export(context.Background(), "./" + mod.Name + ".tar")
	if err != nil {
		return fmt.Errorf("export: %w", err)
	}

	fmt.Printf("Exported at %s\n", addr)

	return nil
}
