// Package service provides the core interfaces and types for managing
// heterogeneous services within the platform. It defines the abstraction
// for service providers and the specification for service deployment.
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"dario.cat/mergo"
	"github.com/blabtm/v2k/platform/sm/cue"
	"sigs.k8s.io/yaml"
)

var modulePath string
var configPath string
var deployPath string

func init() {
	modulePath = os.Getenv("MODULE_PATH")
	configPath = os.Getenv("CONFIG_PATH")
	deployPath = os.Getenv("DEPLOY_PATH")
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

// GetModelPath returns the absolute path to the model for a service with the name.
// It translates the dot-notation service name (e.g., `em-es.sim`) into a directory path
// under storage root specified by the `MODULE_PATH` environment variable
// (e.g., `/etc/v2k/model/em_es/sim`).
func GetModelPath(name string) string {
	return filepath.Join(modulePath,
		"model", strings.ReplaceAll(name, ".", string(filepath.Separator)))
}

// GetConfigPath returns the absolute path to configuration for a service with the name.
// It translates the dot-notation service name (e.g., `em_es.sim`) into a directory path
// under storage root specified by the `CONFIG_PATH` environment variable
// (e.g., `/etc/v2k/config/em_es/sim`).
func GetConfigPath(name string) string {
	return filepath.Join(configPath,
		strings.ReplaceAll(name, ".", string(filepath.Separator)))
}

// GetDeployPath returns the absolute path to deployment artifacts for a service with
// the name. It translates the dot-notation service name (e.g., `em-es.sim`) into a
// directory path under storage root specified by the `DEPLOY_PATH` environment variable
// (e.g., `/etc/v2k/deploy/em-es/sim`).
func GetDeployPath(name string) string {
	return filepath.Join(deployPath,
		strings.ReplaceAll(name, ".", string(filepath.Separator)))
}

// Config represents the desired service configuration.
type Config struct {
	// Active indicates the desired state for the service.
	Active bool `json:"active"`
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

// Get returns runtime configuration for the service.
func Get(name string) ([]byte, error) {
	spec, err := Load(name)
	if err != nil {
		return nil, err
	}

	conf, err := os.ReadFile(filepath.Join(spec.ConfigPath, "config.json"))
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	return conf, nil
}

// Set updates service's configuration and enforces it's operational state.
func Set(name string, raw []byte) error {
	spec, err := Load(name)

	if err != nil {
		return err
	}

	if len(raw) != 0 {
		dst := filepath.Join(spec.ConfigPath, "config.json")

		merged, err := merge(dst, raw)
		if err != nil {
			slog.Error("merge", "err", err)
			return err
		}

		err = cue.Validate(modulePath, filepath.Join(GetModelPath(name), "config.cue"), merged)
		if err != nil {
			slog.Error("validate", "err", err)
			return err
		}

		if err := os.WriteFile(dst, merged, 0644); err != nil {
			return err
		}
	}

	if err := update(spec); err != nil {
		return err
	}

	return nil
}

func merge(path string, raw []byte) ([]byte, error) {
	base, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var baseMap map[string]any
	var overrideMap map[string]any

	baseDecoder := json.NewDecoder(bytes.NewReader(base))
	baseDecoder.UseNumber()
	if err := baseDecoder.Decode(&baseMap); err != nil {
		return nil, err
	}

	overrideDecoder := json.NewDecoder(bytes.NewReader(raw))
	overrideDecoder.UseNumber()
	if err := overrideDecoder.Decode(&overrideMap); err != nil {
		return nil, err
	}

	if err := mergo.Merge(&baseMap, overrideMap, mergo.WithOverride); err != nil {
		return nil, err
	}

	merged, err := json.MarshalIndent(baseMap, "", "  ")
	if err != nil {
		return nil, err
	}

	return merged, nil
}

func update(spec *Spec) error {
	ctx := context.Background()

	raw, err := os.ReadFile(filepath.Join(spec.ConfigPath, "config.json"))
	if err != nil {
		return err
	}

	var conf Config
	if err := json.Unmarshal(raw, &conf); err != nil {
		return err
	}

	status, err := Registry[spec.Driver].Ps(ctx, spec)
	if err != nil {
		return err
	}

	// Underlying engine will handle reconciliation loop.
	// We only need to deploy or withdraw the job here.

	if status.State != Stopped && !conf.Active {
		return Registry[spec.Driver].Down(ctx, spec)
	}

	if status.State == Stopped && conf.Active {
		return Registry[spec.Driver].Up(ctx, spec)
	}

	return nil
}
