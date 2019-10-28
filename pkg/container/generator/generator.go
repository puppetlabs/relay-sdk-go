package generator

import (
	"os"

	"github.com/puppetlabs/nebula-sdk/pkg/container/def"
)

const (
	DefaultScriptFilename = "hack/build-container"
	DefaultRepoNameBase   = "sdk.nebula.localhost/generated"
)

type File struct {
	Ref     *def.FileRef
	Mode    os.FileMode
	Content string
}

type Generator struct {
	c *def.Container

	base           *def.FileRef
	scriptFilename string
	repoNameBase   string
}

func (g *Generator) Files() ([]*File, error) {
	var fs []*File

	// Generate build script
	f, err := g.generateScript()
	if err != nil {
		return nil, err
	}
	fs = append(fs, f)

	// Generate Dockerfiles
	ifs, err := g.generateImages()
	if err != nil {
		return nil, err
	}
	fs = append(fs, ifs...)

	return fs, nil
}

type Option func(g *Generator)

func WithScriptFilename(filename string) Option {
	return func(g *Generator) {
		g.scriptFilename = filename
	}
}

func WithRepoNameBase(base string) Option {
	return func(g *Generator) {
		g.repoNameBase = base
	}
}

func WithFilesRelativeTo(ref *def.FileRef) Option {
	return func(g *Generator) {
		g.base = ref.Dir()
	}
}

func New(container *def.Container, opts ...Option) *Generator {
	g := &Generator{
		c: container,

		base:           def.NewFileRef("."),
		scriptFilename: DefaultScriptFilename,
		repoNameBase:   DefaultRepoNameBase,
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}
