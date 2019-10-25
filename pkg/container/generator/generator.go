package generator

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/puppetlabs/nebula-sdk/pkg/container/def"
)

type Generator struct {
}

func New(container *def.Container) *Generator {
	spew.Dump(container)
	return nil
}
