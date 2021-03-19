package taskutil

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/puppetlabs/leg/encoding/transfer"
	"github.com/puppetlabs/relay-sdk-go/pkg/testutil"
	"github.com/stretchr/testify/require"
)

type TestSpec struct {
	Name       string
	SomeNum    int
	SomeValue  string `spec:"someEncodedValue"`
	SomeStruct OtherStruct
}

type OtherStruct struct {
	Value string
}

func (o *OtherStruct) UnmarshalText(data []byte) (err error) {
	o.Value = string(data)
	return
}

func TestDefaultSpecPlan(t *testing.T) {
	encodedValue, _ := transfer.EncodeJSON([]byte("Hello, \x90!"))

	// Make sure that this will actually get encoded.
	require.NotEqual(t, transfer.NoEncodingType, encodedValue.JSON.EncodingType)

	opts := testutil.MockMetadataAPIOptions{
		SpecResponse: map[string]interface{}{
			"value": map[string]interface{}{
				"name":             "test1",
				"someNum":          12,
				"someEncodedValue": encodedValue,
				"someStruct":       "test2",
			},
		},
	}

	testutil.WithMockMetadataAPI(t, func(ts *httptest.Server) {
		testSpec := TestSpec{}

		opts := DefaultPlanOptions{
			Client:      ts.Client(),
			SpecURL:     fmt.Sprintf("%s/spec", ts.URL),
			SpecTimeout: 5 * time.Second,
		}

		require.NoError(t, PopulateSpecFromDefaultPlan(&testSpec, opts))
		require.Equal(t, "test1", testSpec.Name)
		require.Equal(t, 12, testSpec.SomeNum)
		require.Equal(t, "Hello, \x90!", testSpec.SomeValue)
		require.Equal(t, "test2", testSpec.SomeStruct.Value)
	}, opts)
}

func TestDefaultSpecPlanWithInvalidMetadataAPI(t *testing.T) {
	testSpec := TestSpec{}

	opts := DefaultPlanOptions{
		Client:      nil,
		SpecURL:     "http://ip/spec",
		SpecTimeout: 5 * time.Second,
	}

	require.Error(t, PopulateSpecFromDefaultPlan(&testSpec, opts))
}

func TestValidMetadataURL(t *testing.T) {
	os.Setenv(MetadataAPIURLEnvName, "http://10.20.30.40")
	u, err := MetadataSpecURL()
	require.NoError(t, err)
	require.Equal(t, "http://10.20.30.40/spec", u)
}

func TestInvalidMetadataURL(t *testing.T) {
	os.Setenv(MetadataAPIURLEnvName, " http://ip")
	_, err := MetadataSpecURL()
	require.Error(t, err)
}

func TestUnsetMetadataURL(t *testing.T) {
	os.Unsetenv(MetadataAPIURLEnvName)
	u, err := MetadataSpecURL()
	require.NoError(t, err)
	require.Equal(t, "", u)
}
