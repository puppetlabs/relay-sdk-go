package secrets_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/puppetlabs/horsehead/v2/encoding/transfer"
	"github.com/puppetlabs/nebula-sdk/pkg/secrets"
	"github.com/puppetlabs/nebula-sdk/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecrets(t *testing.T) {
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
		ctx := context.Background()

		u, err := url.Parse(ts.URL + `/secrets`)
		require.NoError(t, err)

		client := secrets.NewDefaultClient(u)

		tt := []struct {
			Key           string
			ExpectedValue string
			ExpectedError error
		}{
			{
				Key:           "raw",
				ExpectedValue: "Hello, test!",
			},
			{
				Key:           "encoded",
				ExpectedValue: "Hello, \x90!",
			},
			{
				Key:           "nope",
				ExpectedError: secrets.ErrClientNotFound,
			},
		}
		for _, test := range tt {
			t.Run(fmt.Sprintf("%s", test.Key), func(t *testing.T) {
				s, err := client.GetSecret(ctx, test.Key)
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
