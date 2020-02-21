package fnlib_test

import (
	"context"
	"testing"

	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/fnlib"
	"github.com/stretchr/testify/require"
)

func TestEquals(t *testing.T) {
	desc, err := fnlib.Library().Descriptor("equals")
	require.NoError(t, err)

	cases := [][]interface{}{
		[]interface{}{"foobar", "foobar"},
		[]interface{}{10, 10},
	}

	for _, c := range cases {
		invoker, err := desc.PositionalInvoker(c)
		require.NoError(t, err)
		r, err := invoker.Invoke(context.Background())
		require.NoError(t, err)
		require.Equal(t, true, r)
	}
}

func TestNotEquals(t *testing.T) {
	desc, err := fnlib.Library().Descriptor("notEquals")
	require.NoError(t, err)

	cases := [][]interface{}{
		[]interface{}{"foobar", "barfoo"},
		[]interface{}{10, 50},
	}

	for _, c := range cases {
		invoker, err := desc.PositionalInvoker(c)
		require.NoError(t, err)
		r, err := invoker.Invoke(context.Background())
		require.NoError(t, err)
		require.Equal(t, true, r)
	}
}
