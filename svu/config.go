package svu

const (
	baseImage  = "ghcr.io/caarlos0/svu"
	defaultTag = "v1.9.0"
)

type config struct {
	repository string
	tag        string
	command    string
	prefix     string
}

func defaultConfig() config {
	return config{
		repository: baseImage,
		tag:        defaultTag,
		command:    "next",
		prefix:     "v",
	}
}

// Option is a function that configures the svu action.
type Option = func(config) config

// Major increases the major of the latest tag
func Major() Option {
	return func(c config) config {
		c.command = "major"

		return c
	}
}

// Minor increases the minor of the latest
func Minor() Option {
	return func(c config) config {
		c.command = "minor"

		return c
	}
}

// Patch increases the patch of the latest
func Patch() Option {
	return func(c config) config {
		c.command = "patch"

		return c
	}
}

// WithPrefix specifies a custom prefix
func WithPrefix(prefix string) Option {
	return func(c config) config {
		c.prefix = prefix

		return c
	}
}
