// Package cue provides utilities to manipulate configuration data
// in the CUE file.
package cue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/cue/parser"
	"cuelang.org/go/cue/token"
)

func Validate(root string, conf string, raw []byte) error {
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

// Merge merges JSON into CUE.
func Merge(cue string, raw []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()

	var data map[string]any
	if err := decoder.Decode(&data); err != nil {
		return err
	}

	mod, err := parser.ParseFile(cue, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, dec := range mod.Decls {
		if emb, ok := dec.(*ast.EmbedDecl); ok {
			str, err := searchExport(emb.Expr)
			if err != nil {
				return err
			}

			if err := inject(str, data); err != nil {
				return err
			}

			break
		}
	}

	out, err := format.Node(mod)
	if err != nil {
		return err
	}

	if err := os.WriteFile(cue, out, 0644); err != nil {
		return err
	}

	return nil
}

func searchExport(src ast.Expr) (*ast.StructLit, error) {
	switch e := src.(type) {
	case *ast.BinaryExpr:
		return searchExport(e.Y)
	case *ast.StructLit:
		return e, nil
	default:
		return nil, fmt.Errorf("malformed configuration")
	}
}

func inject(dst *ast.StructLit, src map[string]any) error {
	for _, decl := range dst.Elts {
		field, ok := decl.(*ast.Field)

		if !ok {
			return fmt.Errorf("malformed configuration")
		}

		name := field.Label.(*ast.Ident).Name
		val, ok := src[name]

		if !ok {
			continue
		}

		switch v := field.Value.(type) {
		case *ast.StructLit:
			structVal, ok := val.(map[string]any)
			if !ok {
				return fmt.Errorf("incompatible configuration")
			}

			if err := inject(v, structVal); err != nil {
				return err
			}
		case *ast.BasicLit:
			v.Value = fmt.Sprint(val)
		}

		delete(src, name)
	}

	dst.Elts = append(dst.Elts, toDecl(src)...)

	return nil
}

func toDecl(src map[string]any) []ast.Decl {
	elts := make([]ast.Decl, 0, len(src))

	for k, v := range src {
		var rhs ast.Expr

		switch val := v.(type) {
		case map[string]any:
			rhs = &ast.StructLit{
				Elts: toDecl(val),
			}
		case string:
			rhs = ast.NewLit(token.FLOAT, val)
		case bool:
			rhs = ast.NewBool(val)
		case json.Number:
			if strings.Contains(val.String(), ".") {
				rhs = ast.NewLit(token.FLOAT, val.String())
			} else {
				rhs = ast.NewLit(token.INT, val.String())
			}
		}

		if rhs == nil {
			return nil
		}

		elts = append(elts, &ast.Field{
			Label: ast.NewIdent(k),
			Value: rhs,
		})
	}

	return elts
}
