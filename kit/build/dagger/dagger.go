package dagger

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"dagger.io/dagger"
	"github.com/blabtm/kit/build"
)

const ZigVersion = "0.16.0"

type Arch string

const (
	Amd64 Arch = "x86_64"
	Arm64 Arch = "aarch64"
)

var arch Arch

func init() {
	switch runtime.GOARCH {
	case "arm64":
		arch = Arm64
	case "amd64":
		arch = Amd64
	default:
		panic("unsupported architecture")
	}
}

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
	ctx := context.Background()

	srcPath, _ := filepath.Abs(mod.ToolchainDir)
	buildPath := fmt.Sprintf("%s/build", srcPath)
	srcDir := dag.Host().Directory(srcPath)
	toolchain := getToolchain()
	builder := dag.Container().From("ubuntu:latest")

	builder = withPackges(builder, []string{
		"build-essential",
		"cmake",
		"curl",
		"git",
		"zlib1g-dev",
		"libabsl-dev",
		"libgsl-dev",
	})

	builder = withZig(builder, ZigVersion).
		WithDirectory("/src", srcDir).
		WithWorkdir("/src").
		WithExec([]string{"cmake", "--toolchain", toolchain, "-S", ".", "-B", "build"}).
		WithExec([]string{"cmake", "--build", "build", "--target", mod.Name})

	if _, err := builder.Sync(ctx); err != nil {
		if err, ok := errors.AsType[*dagger.ExecError](err); ok {
			fmt.Fprintln(os.Stderr, err.Stderr)
		}

		return fmt.Errorf("build: %w", err)
	}

	builder = builder.
		WithExec([]string{"CodeChecker"})

	if _, err :=
		getCommands(builder.File("/src/build/compile_commands.json"), srcPath).
			Export(ctx, buildPath+"/compile_commands.json"); err != nil {
		return fmt.Errorf("export: %w", err)
	}

	if _, err := builder.
		File(fmt.Sprintf("/src/build/%s/%s", mod.Name, mod.Name)).
		Export(ctx, "app"); err != nil {
		return fmt.Errorf("export: %w", err)
	}

	fmt.Println("Complete.")

	return nil
}

func withZig(cont *dagger.Container, version string) *dagger.Container {
	name := fmt.Sprintf("zig-%s-linux-%s", arch, version)
	tar := fmt.Sprintf("%s.tar.xz", name)
	url := fmt.Sprintf("https://ziglang.org/download/%s/%s", version, tar)
	dst := fmt.Sprintf("/opt/%s", name)

	return cont.
		WithExec([]string{"curl", "-LO", url}).
		WithExec([]string{"tar", "-xf", tar, "-C", "/opt"}).
		WithExec([]string{"mv", dst, "/opt/zig"}).
		WithExec([]string{"ln", "-s", "/opt/zig/zig", "/usr/local/bin/zig"})
}

func withPackges(cont *dagger.Container, pkgs []string) *dagger.Container {
	return cont.
		WithExec([]string{"apt-get", "update"}).
		WithExec(append([]string{"apt-get", "install", "-y"}, pkgs...))
}

func getToolchain() string {
	return fmt.Sprintf("zig-cmake-toolchains/zig-toolchain-%s-linux-gnu.cmake", arch)
}

func getCommands(file *dagger.File, path string) *dagger.File {
	return file.
		WithReplaced("\"/src", "\""+path, dagger.FileWithReplacedOpts{
			All: true,
		}).
		WithReplaced(" /src", " "+path, dagger.FileWithReplacedOpts{
			All: true,
		}).
		WithReplaced("-I/src", "-I"+path, dagger.FileWithReplacedOpts{
			All: true,
		})
}
