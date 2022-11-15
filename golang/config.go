package golang

type goos string

const (
	// OsAix is the value for GOOS
	OsAix goos = "aix"
	// OsDarwin is the value for GOOS
	OsDarwin goos = "darwin"
	// OsDragonfly is the value for GOOS
	OsDragonfly goos = "dragonfly"
	// OsFreeBSD is the value for GOOS
	OsFreeBSD goos = "freebsd"
	// OsIllumos is the value for GOOS
	OsIllumos goos = "illumos"
	// OsJs is the value for GOOS
	OsJs goos = "js"
	// OsLinux is the value for GOOS
	OsLinux goos = "linux"
	// OsNetBSD is the value for GOOS
	OsNetBSD goos = "netbsd"
	// OsOpenBSD is the value for GOOS
	OsOpenBSD goos = "openbsd"
	// OsPlan9 is the value for GOOS
	OsPlan9 goos = "plan9"
	// OsSolaris is the value for GOOS
	OsSolaris goos = "solaris"
	// OsWindows is the value for GOOS
	OsWindows goos = "windows"
)

type goarch string

const (
	// Arch386 is the value for GOARCH
	Arch386 goarch = "386"
	// ArchAmd64 is the value for GOARCH
	ArchAmd64 goarch = "amd64"
	// ArchArm is the value for GOARCH
	ArchArm goarch = "arm"
	// ArchArm64 is the value for GOARCH
	ArchArm64 goarch = "arm64"
	// ArchPpc64 is the value for GOARCH
	ArchPpc64 goarch = "ppc64"
	// ArchWasm is the value for GOARCH
	ArchWasm goarch = "wasm"
)

type config struct {
	repository string
	tag        string
	os         []goos
	arch       []goarch
	path       string
	binaryName string
}

func (c config) Matrix() bool {
	return len(c.os) > 0 && len(c.arch) > 0
}

func defaultConfig() config {
	return config{
		repository: "golang",
		tag:        "latest",
		os:         []goos{},
		arch:       []goarch{},
		path:       "",
		binaryName: "",
	}
}

// BuildOption is a function that configures the build action.
type BuildOption = func(config) config

// WithOses defines the list of OS used to build the binary
func WithOses(os goos, oses ...goos) BuildOption {
	return func(c config) config {
		c.os = append(c.os, append([]goos{os}, oses...)...)
		return c
	}
}

// WithArch defines the list of ARCH used to build the binary
func WithArch(arch goarch, arches ...goarch) BuildOption {
	return func(c config) config {
		c.arch = append(c.arch, append([]goarch{arch}, arches...)...)
		return c
	}
}

// WithRepository allows to specify a custom repository and tag
func WithRepository(repository, tag string) BuildOption {
	return func(c config) config {
		c.repository = repository
		c.tag = tag
		return c
	}
}

// WithPath specifies the path of the executable
func WithPath(path string) BuildOption {
	return func(c config) config {
		c.path = path
		return c
	}
}

// WithBinaryName specifies the name a the builded binary
//
// If you build multiple os and arch, this name will be used as a prefix
func WithBinaryName(name string) BuildOption {
	return func(c config) config {
		c.binaryName = name
		return c
	}
}
