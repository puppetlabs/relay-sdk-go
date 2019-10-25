package v1_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/puppetlabs/nebula-sdk/pkg/container/asset"
	v1 "github.com/puppetlabs/nebula-sdk/pkg/container/types/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStepContainerValid(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/container_valid.yaml")
	require.NoError(t, err)

	sct, err := v1.NewStepContainerFromString(string(b))
	assert.NoError(t, err)
	spew.Dump(sct)
}

func TestStepContainerTemplatesValid(t *testing.T) {
	for _, name := range []string{"bash.v1", "go.v1"} {
		s, err := asset.AssetString(fmt.Sprintf("templates/%s/container.yaml", name))
		require.NoError(t, err)

		sctt, err := v1.NewStepContainerTemplateFromString(s)
		assert.NoError(t, err)
		spew.Dump(sctt)
	}
}
