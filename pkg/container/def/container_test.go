package def_test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/puppetlabs/nebula-sdk/pkg/container/def"
	"github.com/stretchr/testify/require"
)

func TestContainerTemplateResolution(t *testing.T) {
	c, err := def.NewFromFilePath("testdata/container_valid.yaml")
	require.NoError(t, err)

	spew.Dump(c)
}
