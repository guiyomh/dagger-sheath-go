package markdown_test

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	"github.com/guiyomh/dagger-sheath-go/markdown"
)

func ExampleLint() {
	ctx := context.Background()
	client, err := dagger.Connect(ctx /*, dagger.WithLogOutput(os.Stdout)*/)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	source := client.Host().Workdir()

	err = markdown.Lint(ctx, client, source, markdown.WithFiles("README.md"))
	if err != nil {
		panic(err)
	}
	fmt.Println("no-error")

	//Output:
	// no-error
}
