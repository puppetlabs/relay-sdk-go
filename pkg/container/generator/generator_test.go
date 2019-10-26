package generator_test

import (
	"fmt"
	"testing"

	"github.com/puppetlabs/nebula-sdk/pkg/container/def"
	"github.com/puppetlabs/nebula-sdk/pkg/container/generator"
	"github.com/stretchr/testify/require"
)

func TestGeneratorFiles(t *testing.T) {
	tpl, err := def.NewTemplateFromFileRef(def.NewFileRef("bash.v1", def.WithFileRefResolver(def.SDKResolver)))
	require.NoError(t, err)

	c := &def.Container{
		Common:      tpl.Common,
		Title:       "Test",
		Description: "The test task does the best testing.",
	}
	c.Settings["AdditionalPackages"].Value = []string{"xmlstarlet"}
	c.Settings["AdditionalCommands"].Value = []string{"do\nmy\nbidding"}

	g := generator.New("test", c)

	m, err := g.Files()
	require.NoError(t, err)

	for _, f := range m {
		fmt.Println("====", f.Name, "====")
		fmt.Println(f.Content)
	}
}
