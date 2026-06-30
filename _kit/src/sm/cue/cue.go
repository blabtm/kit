// Package cue provides utilities to manipulate configuration data
// in the CUE file.
package cue

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/jsonschema"
)

func Validate(root, conf string, raw []byte) error {
	ctx := cuecontext.New()
	ins := load.Instances([]string{conf}, &load.Config{
		ModuleRoot: root,
	})

	if len(ins) == 0 {
		return fmt.Errorf("no instances to load")
	}

	res, err := ctx.BuildInstances(ins)
	if err != nil {
		return err
	}

	cfg := res[0].LookupPath(cue.ParsePath("#Config"))

	rawCue := ctx.CompileBytes(raw, cue.Filename("input.json"))
	if err := rawCue.Err(); err != nil {
		return err
	}

	return cfg.Unify(rawCue).Validate(cue.Final())
}

func ToSchema(root, conf string) (string, error) {
	ctx := cuecontext.New()
	ins := load.Instances([]string{conf}, &load.Config{
		ModuleRoot: root,
	})

	if len(ins) == 0 {
		return "", fmt.Errorf("no instances to load")
	}

	res, err := ctx.BuildInstances(ins)
	if err != nil {
		return "", err
	}

	cfg := res[0].LookupPath(cue.ParsePath("#Config"))

	ast, err := jsonschema.Generate(cfg, &jsonschema.GenerateConfig{
		Version:      jsonschema.VersionDraft2020_12,
		ExplicitOpen: false,
	})
	if err != nil {
		return "", err
	}

	jsonValue := ctx.BuildExpr(ast)

	jsonBytes, err := jsonValue.MarshalJSON()
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func ToObject(root, conf string) (string, error) {
	ctx := cuecontext.New()
	ins := load.Instances([]string{conf}, &load.Config{
		ModuleRoot: root,
	})

	if len(ins) == 0 {
		return "", fmt.Errorf("no instances to load")
	}

	res, err := ctx.BuildInstances(ins)
	if err != nil {
		return "", err
	}

	cfg := res[0]

	jsonBytes, err := cfg.MarshalJSON()
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
