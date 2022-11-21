// Package svu contains dagger helpers to manage semantic version at ease.
// It is base on Semantic Version Util [https://github.com/caarlos0/svu] is a tool to manage semantic versions
package svu

//go:generate gomarkdoc --config ../.gomarkdoc.yml --output {{.Dir}}/README.md ./...

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	"github.com/pkg/errors"
)

// Exec run svu command to get the version
func Exec(ctx context.Context, client *dagger.Client, workdir *dagger.Directory, opts ...Option) (string, error) {
	config := defaultConfig()

	for _, opt := range opts {
		config = opt(config)
	}

	container := client.Container().
		From(fmt.Sprintf("%s:%s", config.repository, config.tag)).
		WithMountedDirectory("/src", workdir).
		WithWorkdir("/src")

	container = container.Exec(
		dagger.ContainerExecOpts{
			Args: []string{
				config.command,
				"--prefix",
				config.prefix,
			},
		},
	)

	if _, err := container.ExitCode(ctx); err != nil {
		return "", errors.Errorf("Couldn't run SVU: %s", err)
	}

	version, err := container.Stdout().Contents(ctx)
	if err != nil {
		return "", errors.Errorf("Couldn't get the ouptut of SVU: %s", err)
	}

	return version, nil
}
