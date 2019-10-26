package generator

import (
	"os"

	"github.com/puppetlabs/nebula-sdk/pkg/container/def"
)

type File struct {
	Ref     *def.FileRef
	Mode    os.FileMode
	Content string
}

type Generator struct {
	name string
	c    *def.Container

	base           *def.FileRef
	scriptFilename string
}

func (g *Generator) Files() ([]*File, error) {
	var fs []*File

	// Generate build script
	f, err := generateScript(g.name, g.c, g.base.Join(g.scriptFilename))
	if err != nil {
		return nil, err
	}
	fs = append(fs, f)

	// Generate Dockerfiles
	ifs, err := generateImages(g.name, g.c, g.base)
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

func WithFilesRelativeTo(ref *def.FileRef) Option {
	return func(g *Generator) {
		g.base = ref.Dir()
	}
}

func New(name string, container *def.Container, opts ...Option) *Generator {
	g := &Generator{
		name: name,
		c:    container,

		base:           def.NewFileRef("."),
		scriptFilename: DefaultScriptFilename,
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}
