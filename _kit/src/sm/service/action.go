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
)

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

	err = cue.Validate(repoPath, filepath.Join(spec.SchemaPath, "config.cue"), dst)
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

func Fork(spec *Spec, raw []byte, msg string) (*Spec, error) {
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

	err = cue.Validate(repoPath, filepath.Join(spec.SchemaPath, "config.cue"), dst)
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
			Email: "support@blab.org",
			When:  time.Now(),
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func GetUI(spec *Spec) (string, string, error) {
	schema, err := cue.ToSchema(repoPath, filepath.Join(spec.SchemaPath, "config.cue"))
	if err != nil {
		return "", "", err
	}

	layout, err := cue.ToObject(repoPath, filepath.Join(spec.SchemaPath, "config.cue"))
	if err != nil {
		return "", "", err
	}

	return strings.ReplaceAll(schema, "%23", "#"), layout, nil
}
