// Package markdown contains dagger helpers to manipulate markdown file
package markdown

//go:generate gomarkdoc --output {{.Dir}}/README.md ./...

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	"github.com/pkg/errors"
)

// Lint run markdown lint to check markdown files and flag style issues.
func Lint(ctx context.Context, client *dagger.Client, workdir *dagger.Directory, opts ...Option) error {
	config := defaultConfig()
	for _, opt := range opts {
		config = opt(config)
	}

	if len(config.files) == 0 {
		return errors.New("You must specify files/Directory. Use WithFiles")
	}

	container := client.Container().
		From(fmt.Sprintf("%s:%s", config.repository, config.tag)).
		WithMountedDirectory("/src", workdir).
		WithWorkdir("/src")

	container = container.Exec(dagger.ContainerExecOpts{
		Args: append([]string{"markdownlint"}, config.files...),
	})

	_, err := container.ExitCode(ctx)
	if err != nil {
		return err
	}

	return nil
}
