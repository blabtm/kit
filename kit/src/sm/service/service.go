// Package service provides the core interfaces and types for managing
// heterogeneous services within the platform. It defines the abstraction
// for service providers and the specification for service deployment.
package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/blabtm/v2k/platform/sm/service/iid"
	"sigs.k8s.io/yaml"
)

var repoPath string
var configPath string

func init() {
	repoPath = os.Getenv("REPO_PATH")
	configPath = filepath.Join(repoPath, "config")
}

// Registry is a global driver storage. Each driver must register itself within the
// registry (e.g. in the init function) so it can be resolved later.
var Registry = make(map[string]Driver)

// Option is an additional argument to pass when the service starts up.
type Option struct {
	Key   string
	Value any
}

// State represents the operational state of a service.
type State string

const (
	// Stopped indicates that the service was explicitly turned off or not started yet.
	// For multi-task services, it indicates that no tasks have been spawned yet.
	Stopped State = "STOPPED"

	// Pending indicates that the service is transitioning between two states.
	// For multi-task services, it indicates that at least one task is pending.
	// Finished tasks are considered PENDING while the overall service state is unresolved.
	Pending State = "PENDING"

	// Running indicates that the service is operational and executing.
	// For multi-task services, it indicates that all tasks are in an operational state.
	Running State = "RUNNING"

	// Completed indicates that the oneshot service finished flawlessly.
	// For multi-task services, it indicates that all tasks completed successfully.
	Completed State = "COMPLETED"

	// Failed indicates that the oneshot service finished with an error.
	// For multi-task services, it indicates that all tasks are finished and
	// at least one failed.
	Failed State = "FAILED"
)

// Status provides a detailed report of a service's current operational state.
type Status struct {
	// State is the current operational state of the service.
	State State `json:"state"`
}

// Driver defines the interface for managing the lifecycle and status of a service
// within a specific execution environment (e.g., swarm or systemd).
type Driver interface {
	// Ps reports the current operational status of the service specified by the Spec.
	Ps(ctx context.Context, spec *Spec) (*Status, error)

	// Up brings the service specified by the Spec online, starting its execution.
	Up(ctx context.Context, spec *Spec, opts ...Option) error

	// Down takes the service specified by the Spec offline, stopping its execution.
	Down(ctx context.Context, spec *Spec) error
}

// Config represents the desired service configuration.
type Config struct {
	// Active indicates the desired state for the service.
	Active bool `json:"active"`
}

type Type string

const (
	Service Type = "service"
	Oneshot Type = "oneshot"
)

type Capability string

const (
	// Live indicates the service's ability to reload configuration at runtime.
	Live Capability = "live"
)

// Spec represents the declarative specification for a service.
type Spec struct {
	// Driver is a driver name.
	Driver string `yaml:"driver"`

	// Type is service type.
	Type Type `yaml:"type"`

	// Capabilities holds service optional functionality.
	Capabilities map[string]bool `yaml:"capabilities"`

	// Options holds driver-specific configuration parameters.
	Options any `yaml:"options"`

	// IID is a unique oneshot invocation identifier.
	// It's ignored for non-oneshot services.
	IID *iid.IID

	// Name is a fully-qualified unique service identifier in dot-notation.
	// For example, `em_es.sim`.
	Name string

	// SchemaPath is an absolute path to service's models location.
	SchemaPath string

	// ConfigPath is an absolute path to service's configurations location.
	ConfigPath string

	// DeployPath is an absolute path to service's deployment artifacts location.
	DeployPath string
}

// Load reads and unmarshals a service specification from a YAML file.
func Load(name string, id string) (*Spec, error) {
	spec := &Spec{
		Name:       name,
		SchemaPath: GetSchemaPath(name),
		ConfigPath: GetConfigPath(name),
		DeployPath: GetDeployPath(name),
	}

	body, err := os.ReadFile(filepath.Join(spec.DeployPath, "service.yaml"))
	if err != nil {
		return nil, fmt.Errorf("spec: %w", err)
	}

	if err := yaml.Unmarshal(body, spec); err != nil {
		return nil, fmt.Errorf("spec: %w", err)
	}

	if id != "" {
		spec.IID = iid.Parse(name, id)
		spec.ConfigPath = filepath.Join(spec.ConfigPath, spec.IID.Short())

		if _, err := os.Stat(spec.ConfigPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("no such invocation")
		}
	}

	return spec, nil
}

func (s *Spec) String() string {
	if s.IID != nil {
		return s.IID.Full()
	}

	return strings.ReplaceAll(s.Name, ".", "_")
}

// GetSchemaPath returns the absolute path to the schema for a
// service with the name.
func GetSchemaPath(name string) string {
	path := strings.ReplaceAll(name, ".", string(filepath.Separator))
	root, _ := os.ReadDir(repoPath)

	for _, ent := range root {
		if ent.IsDir() &&
			ent.Name() != "dom" &&
			ent.Name() != "config" &&
			ent.Name() != "blob" &&
			ent.Name() != "deploy" {
			pred := filepath.Join(repoPath, ent.Name(), path)

			if _, err := os.Stat(pred); err == nil {
				return pred
			}
		}
	}

	return ""
}

// GetConfigPath returns the absolute path to configuration for a
// service with the name.
func GetConfigPath(name string) string {
	return filepath.Join(repoPath,
		"config", strings.ReplaceAll(name, ".", string(filepath.Separator)))
}

// GetDeployPath returns the absolute path to deployment artifacts for
// a service with the name.
func GetDeployPath(name string) string {
	return filepath.Join(repoPath,
		"deploy", strings.ReplaceAll(name, ".", string(filepath.Separator)))
}
