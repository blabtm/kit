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
	"time"

	"dario.cat/mergo"
	"github.com/blabtm/v2k/platform/sm/cue"
	"github.com/blabtm/v2k/platform/sm/service/iid"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
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

	// ModelPath is an absolute path to service's models location.
	ModelPath string

	// ConfigPath is an absolute path to service's configurations location.
	ConfigPath string

	// DeployPath is an absolute path to service's deployment artifacts location.
	DeployPath string
}

// Load reads and unmarshals a service specification from a YAML file.
func Load(name string, id string) (*Spec, error) {
	spec := &Spec{
		Name:       name,
		ModelPath:  GetModelPath(name),
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

// GetModelPath returns the absolute path to the model for a service with the name.
// It translates the dot-notation service name (e.g., `em-es.sim`) into a directory path
// under storage root specified by the `MODULE_PATH` environment variable
// (e.g., `/etc/v2k/model/em_es/sim`).
func GetModelPath(name string) string {
	return filepath.Join(repoPath,
		"model", strings.ReplaceAll(name, ".", string(filepath.Separator)))
}

// GetConfigPath returns the absolute path to configuration for a service with the name.
// It translates the dot-notation service name (e.g., `em_es.sim`) into a directory path
// under storage root specified by the `CONFIG_PATH` environment variable
// (e.g., `/etc/v2k/config/em_es/sim`).
func GetConfigPath(name string) string {
	return filepath.Join(repoPath,
		"config", strings.ReplaceAll(name, ".", string(filepath.Separator)))
}

// GetDeployPath returns the absolute path to deployment artifacts for a service with
// the name. It translates the dot-notation service name (e.g., `em-es.sim`) into a
// directory path under storage root specified by the `DEPLOY_PATH` environment variable
// (e.g., `/etc/v2k/deploy/em-es/sim`).
func GetDeployPath(name string) string {
	return filepath.Join(repoPath,
		"deploy", strings.ReplaceAll(name, ".", string(filepath.Separator)))
}

// Config represents the desired service configuration.
type Config struct {
	// Active indicates the desired state for the service.
	Active bool `json:"active"`
}

// Ls returns the list of all registered services.
func Ls() ([]string, error) {
	dir := filepath.Join(repoPath, "deploy")
	res := make([]string, 0)

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.Type().IsRegular() && d.Name() == "service.yaml" {
			path, _ = strings.CutPrefix(path, dir+"/")
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
func Get(spec *Spec) ([]byte, error) {
	conf, err := os.ReadFile(filepath.Join(spec.ConfigPath, "config.json"))
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	return conf, nil
}

// Set updates service's configuration and enforces it's operational state.
func Set(spec *Spec, raw []byte) error {
	if spec.Type == Oneshot && spec.IID != nil {
		return fmt.Errorf("cannot update invocation configuration")
	}

	path := filepath.Join(spec.ConfigPath, "config.json")
	dst, err := os.ReadFile(path)
	if err != nil {
		slog.Error("read configuration", "err", err)
		return err
	}

	changed := false

	if len(raw) != 0 {
		dst, err = merge(dst, raw)
		if err != nil {
			slog.Error("merge", "err", err)
			return err
		}

		changed = true
	}

	err = cue.Validate(repoPath, filepath.Join(spec.ModelPath, "config.cue"), dst)
	if err != nil {
		slog.Error("validate", "err", err)
		return err
	}

	if changed {
		if err := os.WriteFile(path, dst, 0644); err != nil {
			return err
		}

		if err := commit(path, ""); err != nil {
			slog.Error("commit", "err", err)
		}
	}

	if spec.Type != Oneshot {
		if err := update(spec); err != nil {
			return err
		}
	}

	return nil
}

func Run(spec *Spec, raw []byte, msg string) (*Spec, error) {
	if spec.Type != Oneshot {
		return nil, fmt.Errorf("cannot fork from non-oneshot service")
	}

	if spec.IID != nil {
		return nil, fmt.Errorf("cannot fork from invocation")
	}

	id := iid.New(spec.Name)

	dst, err := os.ReadFile(filepath.Join(spec.ConfigPath, "config.json"))
	if err != nil {
		slog.Error("read configuration", "err", err)
		return nil, err
	}

	var conf Config
	if err := json.Unmarshal(dst, &conf); err != nil {
		return nil, err
	}

	if !conf.Active {
		return nil, fmt.Errorf("cannot fork from inactive service")
	}

	if len(raw) != 0 {
		dst, err = merge(dst, raw)
		if err != nil {
			slog.Error("merge", "err", err)
			return nil, err
		}
	}

	err = cue.Validate(repoPath, filepath.Join(spec.ModelPath, "config.cue"), dst)
	if err != nil {
		slog.Error("validate", "err", err)
		return nil, err
	}

	dir := filepath.Join(spec.ConfigPath, id.Short())
	if err := os.Mkdir(dir, 0755); err != nil {
		return nil, err
	}

	if err := os.WriteFile(filepath.Join(dir, "config.json"), dst, 0444); err != nil {
		return nil, err
	}

	if err := os.Chmod(dir, 0555); err != nil {
		return nil, err
	}

	if err := commit(dir, msg); err != nil {
		slog.Error("commit", "err", err)
	}

	newSpec := *spec
	newSpec.IID = id
	newSpec.ConfigPath = dir

	if err := update(&newSpec); err != nil {
		return nil, err
	}

	return &newSpec, nil
}

func merge(dst, src []byte) ([]byte, error) {
	var dstMap map[string]any
	var srcMap map[string]any

	dstDecoder := json.NewDecoder(bytes.NewReader(dst))
	dstDecoder.UseNumber()
	if err := dstDecoder.Decode(&dstMap); err != nil {
		return nil, err
	}

	srcDecoder := json.NewDecoder(bytes.NewReader(src))
	srcDecoder.UseNumber()
	if err := srcDecoder.Decode(&srcMap); err != nil {
		return nil, err
	}

	if err := mergo.Merge(&dstMap, srcMap, mergo.WithOverride); err != nil {
		return nil, err
	}

	merged, err := json.MarshalIndent(dstMap, "", "  ")
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

	// TODO: reload or notify service

	return nil
}

func commit(path, msg string) error {
	if msg == "" {
		msg = "configuration update"
	}

	repo, err := git.PlainOpen(configPath)
	if err != nil {
		return err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	relPath, _ := strings.CutPrefix(path, configPath+"/")
	if _, err := wt.Add(relPath); err != nil {
		return err
	}

	_, err = wt.Commit(fmt.Sprintf("sm: %s", msg), &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Service Manager",
			Email: "support@blabtm.org",
			When:  time.Now(),
		},
	})

	if err != nil {
		return err
	}

	return nil
}
