package fn_test

import (
	"context"
	"testing"

	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/fn"
	"github.com/stretchr/testify/require"
)

func TestLibJSONUnmarshal(t *testing.T) {
	desc, err := fn.Library.Descriptor("jsonUnmarshal")
	require.NoError(t, err)

	invoker, err := desc.PositionalInvoker([]interface{}{`{"foo": "bar"}`})
	require.NoError(t, err)

	r, err := invoker.Invoke(context.Background())
	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{"foo": "bar"}, r)
}

func TestLibConcat(t *testing.T) {
	desc, err := fn.Library.Descriptor("concat")
	require.NoError(t, err)

	invoker, err := desc.PositionalInvoker([]interface{}{"Hello, ", "world!"})
	require.NoError(t, err)

	r, err := invoker.Invoke(context.Background())
	require.NoError(t, err)
	require.Equal(t, "Hello, world!", r)
}

func TestLibAppend(t *testing.T) {
	desc, err := fn.Library.Descriptor("append")
	require.NoError(t, err)

	invoker, err := desc.PositionalInvoker([]interface{}{
		[]interface{}{1, 2},
		3,
		4,
	})
	require.NoError(t, err)

	r, err := invoker.Invoke(context.Background())
	require.NoError(t, err)
	require.Equal(t, []interface{}{1, 2, 3, 4}, r)
}
