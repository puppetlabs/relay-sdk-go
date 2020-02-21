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
		[]interface{}{10.5, 10.5},
		[]interface{}{[]string{"foo", "bar"}, []string{"foo", "bar"}},
		[]interface{}{[]int{1, 2}, []int{1, 2}},
		[]interface{}{[]float32{1.1, 2.0}, []float32{1.1, 2.0}},
		[]interface{}{[]float64{1.1, 2.0}, []float64{1.1, 2.0}},
		[]interface{}{map[string]string{"foo": "bar"}, map[string]string{"foo": "bar"}},
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
		[]interface{}{10.0, 50.5},
		[]interface{}{[]string{"foo", "bar", "baz"}, []string{"foo", "bar"}},
		[]interface{}{[]int{1, 2, 3}, []int{1, 2}},
		[]interface{}{[]float32{1.1, 2.0, 3.2}, []float32{1.1, 2.0}},
		[]interface{}{[]float64{1.1, 2.0, 3.2}, []float64{1.1, 2.0}},
		[]interface{}{map[string]string{"foo": "bar", "baz": "biz"}, map[string]string{"foo": "bar"}},
	}

	for _, c := range cases {
		invoker, err := desc.PositionalInvoker(c)
		require.NoError(t, err)

		r, err := invoker.Invoke(context.Background())
		require.NoError(t, err)

		require.Equal(t, true, r)
	}
}
