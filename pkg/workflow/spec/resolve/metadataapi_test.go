package resolve_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/puppetlabs/horsehead/v2/encoding/transfer"
	"github.com/puppetlabs/nebula-sdk/pkg/outputs"
	"github.com/puppetlabs/nebula-sdk/pkg/secrets"
	"github.com/puppetlabs/nebula-sdk/pkg/testutil"
	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/resolve"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetadataAPISecretResolver(t *testing.T) {
	ctx := context.Background()

	rawValue, _ := (transfer.NoEncoding{}).EncodeJSON([]byte("Hello, test!"))
	encodedValue, _ := (transfer.Base64Encoding{}).EncodeJSON([]byte("Hello, \x90!"))

	opts := testutil.MockMetadataAPIOptions{
		Secrets: map[string]testutil.MockSecret{
			"raw": {
				ResponseObject: secrets.Secret{
					Key:   "raw",
					Value: rawValue,
				},
			},
			"encoded": {
				ResponseObject: secrets.Secret{
					Key:   "encoded",
					Value: encodedValue,
				},
			},
		},
	}

	testutil.WithMockMetadataAPI(t, func(ts *httptest.Server) {
		u, err := url.Parse(ts.URL)
		require.NoError(t, err)

		resolver := resolve.NewMetadataAPISecretTypeResolver(u)

		tt := []struct {
			Name          string
			ExpectedValue string
			ExpectedError error
		}{
			{
				Name:          "raw",
				ExpectedValue: "Hello, test!",
			},
			{
				Name:          "encoded",
				ExpectedValue: "Hello, \x90!",
			},
			{
				Name:          "nope",
				ExpectedError: &resolve.SecretNotFoundError{Name: "nope"},
			},
		}
		for _, test := range tt {
			t.Run(fmt.Sprintf("%s", test.Name), func(t *testing.T) {
				s, err := resolver.ResolveSecret(ctx, test.Name)
				if test.ExpectedError != nil {
					assert.Equal(t, test.ExpectedError, err)
					assert.Empty(t, s)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, test.ExpectedValue, s)
				}
			})
		}
	}, opts)
}

func TestMetadataAPIOutputResolver(t *testing.T) {
	ctx := context.Background()

	rawValue, _ := (transfer.NoEncoding{}).EncodeJSON([]byte("Hello, test!"))
	encodedValue, _ := (transfer.Base64Encoding{}).EncodeJSON([]byte("Hello, \x90!"))

	opts := testutil.MockMetadataAPIOptions{
		Outputs: map[testutil.MockOutputKey]testutil.MockOutput{
			testutil.MockOutputKey{TaskName: "test1", Key: "raw"}: {
				ResponseObject: outputs.Output{
					TaskName: "test1",
					Key:      "output",
					Value:    rawValue,
				},
			},
			testutil.MockOutputKey{TaskName: "test1", Key: "encoded"}: {
				ResponseObject: outputs.Output{
					TaskName: "test1",
					Key:      "output",
					Value:    encodedValue,
				},
			},
		},
	}

	testutil.WithMockMetadataAPI(t, func(ts *httptest.Server) {
		u, err := url.Parse(ts.URL)
		require.NoError(t, err)

		resolver := resolve.NewMetadataAPIOutputTypeResolver(u)

		tt := []struct {
			From, Name    string
			ExpectedValue string
			ExpectedError error
		}{
			{
				From:          "test1",
				Name:          "raw",
				ExpectedValue: "Hello, test!",
			},
			{
				From:          "test1",
				Name:          "encoded",
				ExpectedValue: "Hello, \x90!",
			},
			{
				From:          "test2",
				Name:          "key",
				ExpectedError: &resolve.OutputNotFoundError{From: "test2", Name: "key"},
			},
		}
		for _, test := range tt {
			t.Run(fmt.Sprintf("%s/%s", test.From, test.Name), func(t *testing.T) {
				o, err := resolver.ResolveOutput(ctx, test.From, test.Name)
				if test.ExpectedError != nil {
					assert.Equal(t, test.ExpectedError, err)
					assert.Empty(t, o)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, test.ExpectedValue, o)
				}
			})
		}
	}, opts)
}
