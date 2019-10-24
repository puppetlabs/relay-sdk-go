package taskutil

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/puppetlabs/horsehead/v2/encoding/transfer"
	"github.com/puppetlabs/nebula-sdk/pkg/testutil"
	"github.com/stretchr/testify/require"
)

type TestSpec struct {
	Name      string
	SomeNum   int
	SomeValue string `spec:"someEncodedValue"`
}

func TestDefaultSpecPlan(t *testing.T) {
	encodedValue, _ := transfer.EncodeJSON([]byte("Hello, \x90!"))

	// Make sure that this will actually get encoded.
	require.NotEqual(t, transfer.NoEncodingType, encodedValue.JSON.EncodingType)

	opts := testutil.SingleSpecMockMetadataAPIOptions("test1", testutil.MockSpec{
		ResponseObject: map[string]interface{}{
			"name":             "test1",
			"someNum":          12,
			"someEncodedValue": encodedValue,
		},
	})

	testutil.WithMockMetadataAPI(t, func(ts *httptest.Server) {
		testSpec := TestSpec{}

		opts := DefaultPlanOptions{
			Client:  ts.Client(),
			SpecURL: fmt.Sprintf("%s/specs/test1", ts.URL),
		}

		require.NoError(t, PopulateSpecFromDefaultPlan(&testSpec, opts))
		require.Equal(t, "test1", testSpec.Name)
		require.Equal(t, 12, testSpec.SomeNum)
		require.Equal(t, "Hello, \x90!", testSpec.SomeValue)
	}, opts)
}
