package dagger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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
	srcPath, _ := filepath.Abs(mod.ToolchainDir)
	buildPath := fmt.Sprintf("%s/build", srcPath)
	srcDir := dag.Host().Directory(srcPath)

	pkgs := []string{
		"build-essential",
		"cmake",
		"zlib1g-dev",
		"libabsl-dev",
	}

	builder := dag.Container().
		From("ubuntu:latest").
		WithExec([]string{"apt-get", "update"}).
		WithExec(append([]string{"apt-get", "install", "-y"}, pkgs...)).
		WithDirectory(srcPath, srcDir).
		WithWorkdir(srcPath).
		WithExec([]string{"mkdir", "-p", "build"}).
		WithExec([]string{"cmake", "-S", ".", "-B", "build"}).
		WithExec([]string{"cmake", "--build", "build", "--target", mod.Name})

	if _, err := builder.
		Directory(buildPath).
		Export(context.Background(), buildPath); err != nil {
		return fmt.Errorf("export: %w", err)
	}

	fmt.Println("Ok.")

	return nil
}

func (dag *Dagger) PackageNative(mod *build.Module, exec *dagger.File) error {
	image := dag.Container().
		From("alpine:latest").
		WithFile("/bin/run", exec).
		WithEntrypoint([]string{"/bin/run"})

	if _, err := image.Export(context.Background(), "./"+mod.Name+".tar"); err != nil {
		return fmt.Errorf("export: %w", err)
	}

	return nil
}
