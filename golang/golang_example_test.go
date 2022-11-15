package golang_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
	"github.com/guiyomh/dagger-sheath-go/golang"
)

func createDirectory(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func listDirectory(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		fmt.Println(path)

		return nil
	})
}

func ExampleBuild() {
	ctx := context.Background()
	client, err := dagger.Connect(ctx /*, dagger.WithLogOutput(os.Stdout)*/)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	source := client.Host().Directory("./test")

	output := golang.Build(ctx, client, source)

	if err = createDirectory("dist/"); err != nil {
		panic(err)
	}

	_, err = output.Export(ctx, "dist/")
	if err != nil {
		panic(err)
	}

	listDirectory("dist/")
	//Output:
	// dist/example_test
}

func ExampleBuild_multiarch() {
	ctx := context.Background()
	client, err := dagger.Connect(ctx /*, dagger.WithLogOutput(os.Stdout)*/)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	source := client.Host().Directory("./test")

	output := golang.Build(
		ctx,
		client,
		source,
		golang.WithArch(golang.ArchAmd64),
		golang.WithOses(golang.OsLinux, golang.OsDarwin),
		golang.WithBinaryName("my_binary"),
	)

	if err = createDirectory("multi_dist/"); err != nil {
		panic(err)
	}

	_, err = output.Export(ctx, "multi_dist/")
	if err != nil {
		panic(err)
	}

	listDirectory("multi_dist/")
	//Output:
	// multi_dist/darwin/amd64/my_binary_darwin_amd64
	// multi_dist/linux/amd64/my_binary_linux_amd64
}
