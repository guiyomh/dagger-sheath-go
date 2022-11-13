package markdown

const (
	baseImage  = "tmknom/markdownlint"
	defaultTag = "0.31.1"
)

type config struct {
	repository string
	tag        string
	files      []string
}

func defaultConfig() config {
	return config{
		repository: baseImage,
		tag:        defaultTag,
	}
}

// Option is a function that configures the markdown action.
type Option = func(config) config

// WithFiles specifies the list of file, directory or Glob pattern to analyze
func WithFiles(files ...string) Option {
	return func(c config) config {
		c.files = files

		return c
	}
}
