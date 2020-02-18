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
		[]interface{}{123, 123},
		[]interface{}{123.05, 123.05},
		[]interface{}{true, true},
		[]interface{}{false, false},
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
		[]interface{}{123, 321},
		[]interface{}{123.05, 321.05},
		[]interface{}{true, false},
	}

	for _, c := range cases {
		invoker, err := desc.PositionalInvoker(c)
		require.NoError(t, err)
		r, err := invoker.Invoke(context.Background())
		require.NoError(t, err)
		require.Equal(t, false, r)
	}
}
