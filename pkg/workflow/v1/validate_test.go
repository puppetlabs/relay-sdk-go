package v1_test

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	v1 "github.com/puppetlabs/nebula-sdk/pkg/workflow/v1"
	"github.com/stretchr/testify/require"
)

func TestFixtureValidation(t *testing.T) {
	fs, err := filepath.Glob("testdata/*.yaml")
	require.NoError(t, err)

	for _, file := range fs {
		t.Run(filepath.Base(file), func(t *testing.T) {
			b, err := ioutil.ReadFile(file)
			require.NoError(t, err)

			err = v1.ValidateYAML(string(b))
			if strings.HasSuffix(file[:len(file)-len(filepath.Ext(file))], "_invalid") {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
