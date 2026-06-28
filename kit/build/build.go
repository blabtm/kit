package main

import (
	"context"
	"fmt"
)

func (m *Kit) Build(ctx context.Context, name string) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"echo", fmt.Sprintf("\"Hi, building %s...\"", name)}).
		Stdout(ctx)
}

func (m *Kit) BuildNative() (string, error) {
	return "", nil
}
