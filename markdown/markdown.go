// Package markdown contains dagger helpers to manipulate markdown file
package markdown

//go:generate gomarkdoc --config ../.gomarkdoc.yml --output {{.Dir}}/README.md ./...

import (
	"context"
	"fmt"
	"strings"

	"dagger.io/dagger"
	"github.com/hashicorp/go-multierror"
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

	commandLine := fmt.Sprintf("/home/nonroot/entrypoint.sh --output /tmp/errors --quiet %s || true", strings.Join(config.files, " "))
	contents, err := container.Exec(dagger.ContainerExecOpts{
		Args: []string{"sh", "-c", commandLine},
	}).File("/tmp/errors").Contents(ctx)

	if err != nil {
		// file does not exist, the command exited properly
		return nil
	}
	if len(contents) > 0 {
		var errs error
		lines := strings.Split(strings.ReplaceAll(contents, "\r\n", "\n"), "\n")
		for _, line := range lines {
			errs = multierror.Append(errs, errors.New(line))
		}
		return errs
	}

	return nil
}
