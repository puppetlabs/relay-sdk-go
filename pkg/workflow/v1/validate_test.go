package v1_test

import (
	"io/ioutil"
	"testing"

	v1 "github.com/puppetlabs/nebula-sdk/pkg/workflow/v1"
	"github.com/stretchr/testify/require"
)

func TestValidYaml(t *testing.T) {
	b, readErr := ioutil.ReadFile("testdata/valid.yaml")
	require.NoError(t, readErr)

	validationErr := v1.Validate(string(b))
	require.NoError(t, validationErr)
}

func TestBadSyntax(t *testing.T) {
	b, readErr := ioutil.ReadFile("testdata/badSyntax.yaml")
	require.NoError(t, readErr)

	validationErr := v1.Validate(string(b))
	require.Error(t, validationErr)
}

func TestInvalidSchema(t *testing.T) {
	b, readErr := ioutil.ReadFile("testdata/invalid.yaml")
	require.NoError(t, readErr)

	validationErr := v1.Validate(string(b))
	require.Error(t, validationErr)
}
