package build

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	ErrNotExist = errors.New("module does not exist")
)

type Builder interface {
	Build(proj *Module) error
}

type Toolchain int

const (
	Native Toolchain = iota
)

var Toolchains = map[Toolchain]string{
	Native: "native",
}

type Module struct {
	Name         string
	Dir          string
	Toolchain    Toolchain
	ToolchainDir string
}

func OpenModule(name string) (*Module, error) {
	mod := &Module{
		Name: name,
	}

	for toolchain, dir := range Toolchains {
		tcDir := filepath.Join(".", "src", dir)
		dir := filepath.Join(tcDir, name)
		info, err := os.Stat(dir)

		if err != nil || !info.IsDir() {
			continue
		}

		mod.Dir = dir
		mod.Toolchain = toolchain
		mod.ToolchainDir = tcDir
	}

	if mod.Dir == "" {
		return nil, ErrNotExist
	}

	return mod, nil
}
