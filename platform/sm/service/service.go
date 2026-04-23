// Package service provides the core interfaces and types for managing
// heterogeneous services within the platform. It defines the abstraction
// for service providers and the specification for service deployment.
package service

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/yaml"
)

var configPath string
var deployPath string

func init() {
	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		configPath = "/etc/v2k/config"
	}

	deployPath = os.Getenv("DEPLOY_PATH")

	if deployPath == "" {
		deployPath = "/etc/v2k/deploy"
	}
}

// Registry is a global driver storage. Each driver must register itself within the registry
// (e.g. in the `init` function) so it can be resolved later.
var Registry = make(map[string]Driver)

// Option is an additional argument to pass when service starts up.
type Option struct {
	Key   string
	Value any
}

// State represents the operational state of a service.
type State string

const (
	Up      State = "Up"      // Up indicates the service is operational and running.
	Down    State = "Down"    // Down indicates the service is transitioning, not operational or stopped.
	Partial State = "Partial" // Partial indicates the service is not fully functional.
)

// Status provides a detailed report of a service's current operational state.
type Status struct {
	// State is the current operational state of the service.
	State State `json:"state"`
}

// Driver defines the interface for managing the lifecycle and status of a service
// within a specific execution environment (e.g., Docker, Flink).
type Driver interface {
	// Ps reports the current operational status of the service specified by the Spec.
	Ps(ctx context.Context, spec *Spec) (*Status, error)

	// Up brings the service specified by the Spec online, starting its execution.
	Up(ctx context.Context, spec *Spec, opts ...Option) error

	// Down takes the service specified by the Spec offline, stopping its execution.
	Down(ctx context.Context, spec *Spec) error
}

// Spec represents the declarative specification for a service.
type Spec struct {
	// Driver is a driver name.
	Driver string `yaml:"driver"`

	// Config holds driver-specific configuration parameters.
	Config map[string]any `yaml:"config"`

	// Name is a fully-qualified unique service identifier in dot-notation.
	// For example, `em-es.sim`.
	Name string

	// ConfigPath is an absolute path to service's configurations location.
	ConfigPath string

	// DeployPath is an absolute path to service's deployment artifacts location.
	DeployPath string
}

// Load reads and unmarshals a service specification from a YAML file.
func Load(name string) (*Spec, error) {
	def := &Spec{
		Name:       name,
		ConfigPath: GetConfigPath(name),
		DeployPath: GetDeployPath(name),
	}

	body, err := os.ReadFile(filepath.Join(def.DeployPath, "service.yaml"))

	if err != nil {
		return nil, fmt.Errorf("spec: %w", err)
	}

	if err := yaml.Unmarshal(body, def); err != nil {
		return nil, fmt.Errorf("spec: %w", err)
	}

	return def, nil
}

// GetConfigPath returns the absolute path to configuration for a service with the name.
// It translates the dot-notation service name (e.g., `em-es.sim`)
// into a directory path under storage root specified by the `CONFIG_PATH` environment variable
// (e.g., `/etc/v2k/config/em-es/sim`).
func GetConfigPath(name string) string {
	return filepath.Join(configPath, strings.ReplaceAll(name, ".", string(filepath.Separator)))
}

// GetDeployPath returns the absolute path to deployment artifacts for a service with the name.
// It translates the dot-notation service name (e.g., `em-es.sim`)
// into a directory path under storage root specified by the `DEPLOY_PATH` environment variable
// (e.g., `/etc/v2k/deploy/em-es/sim`).
func GetDeployPath(name string) string {
	return filepath.Join(deployPath, strings.ReplaceAll(name, ".", string(filepath.Separator)))
}

// Ls returns the list of all registered services.
func Ls() ([]string, error) {
	res := make([]string, 0)

	err := filepath.WalkDir(deployPath, func(path string, d fs.DirEntry, err error) error {
		if d.Type().IsRegular() && d.Name() == "service.yaml" {
			path, _ = strings.CutPrefix(path, deployPath+"/")
			path, _ = strings.CutSuffix(path, "/service.yaml")
			path = strings.ReplaceAll(path, string(filepath.Separator), ".")

			res = append(res, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetConfig returns the configurations for the service in YAML.
func GetConfig(name string) ([]byte, error) {
	conf, err := os.ReadFile(filepath.Join(
		GetConfigPath(name),
		"config.yaml",
	))

	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	return conf, nil
}

// PutConfig stores the configurations for the service in YAML.
func PutConfig(name string, conf []byte) error {
	if err := os.WriteFile(filepath.Join(GetConfigPath(name), "config.yaml"), conf, 0666); err != nil {
		return fmt.Errorf("write: %v", err)
	}

	return nil
}

