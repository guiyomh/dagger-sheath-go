package golang

//go:generate gomarkdoc --output {{.Dir}}/README.md ./...

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"dagger.io/dagger"
)

const (
	buildPath = "build/"
)

// Build compiles the binary
func Build(ctx context.Context, client *dagger.Client, workdir *dagger.Directory, opts ...BuildOption) *dagger.Directory {
	cfg := defaultConfig()

	for _, opt := range opts {
		cfg = opt(cfg)
	}

	// create the golang container and mount the source
	container := client.Container().From(fmt.Sprintf("%s:%s", cfg.repository, cfg.tag))
	container = container.
		WithMountedDirectory("/src", workdir).
		WithWorkdir("/src")

	// identify the golang module and prefix the path with it
	moduleName := moduleName(ctx, container)
	if moduleName != "" {
		cfg.path = fmt.Sprintf("%s/%s", moduleName, cfg.path)
	}
	if cfg.binaryName == "" {
		cfg.binaryName = filepath.Base(moduleName)
	}

	// prepare the output
	var output *dagger.Directory

	if cfg.Matrix() {
		output = multiarch(client, container, cfg)
	} else {
		output = build(container, filepath.Join(buildPath, cfg.binaryName), cfg)
	}

	return output
}

func moduleName(ctx context.Context, container *dagger.Container) string {
	container = container.Exec(dagger.ContainerExecOpts{Args: []string{
		"go", "list", "-m",
	}})

	if _, err := container.ExitCode(ctx); err != nil {
		return ""
	}
	moduleName, err := container.Stdout().Contents(ctx)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(moduleName)
}

func multiarch(client *dagger.Client, container *dagger.Container, cfg config) *dagger.Directory {
	outputs := client.Directory()

	for _, os := range cfg.os {
		for _, arch := range cfg.arch {
			binaryName := fmt.Sprintf("%s_%s_%s", cfg.binaryName, os, arch)
			outputPath := filepath.Join(string(os), string(arch), binaryName)

			osarch := container.WithEnvVariable("GOOS", string(os))
			osarch = osarch.WithEnvVariable("GOARCH", string(arch))

			outputs = outputs.WithDirectory(filepath.Dir(outputPath), build(osarch, outputPath, cfg))
		}
	}

	return outputs
}

func build(container *dagger.Container, outputPath string, cfg config) *dagger.Directory {
	container = container.Exec(
		dagger.ContainerExecOpts{Args: []string{"go", "build", "-o", outputPath, cfg.path}},
	)

	return container.Directory(filepath.Dir(outputPath))
}
