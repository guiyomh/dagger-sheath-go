package svu_test

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/guiyomh/dagger-sheath-go/svu"
)

func ExampleExec() {
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithWorkdir(".."))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	source := client.Host().Workdir()

	version, err := svu.Exec(ctx, client, source, svu.WithPrefix("foo/v"))
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stdout, "Next version: %s", version)

	// Output:
	// Next version: foo/v0.1.0
}
