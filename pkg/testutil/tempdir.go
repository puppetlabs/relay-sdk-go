package testutil

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func WithTemporaryDirectory(t *testing.T, prefix string, fn func(name string)) {
	d, err := ioutil.TempDir("", prefix)
	require.NoError(t, err)
	defer os.RemoveAll(d)

	t.Logf("Allocated temporary directory: %s", d)

	fn(d)
}
