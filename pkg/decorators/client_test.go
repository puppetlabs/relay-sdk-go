package decorators_test

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/puppetlabs/relay-sdk-go/pkg/decorators"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
	"github.com/puppetlabs/relay-sdk-go/pkg/testutil"
	"github.com/stretchr/testify/require"
)

func TestSetDecorator(t *testing.T) {
	decValues := map[string]string{
		"type":        "link",
		"description": "some link",
		"uri":         "https://unit-testing.relay.sh/link-location",
	}

	opts := testutil.MockMetadataAPIOptions{
		ExpectedDecorators: map[string]map[string]string{
			"some-link": decValues,
		},
	}
	testutil.WithMockMetadataAPI(t, func(ts *httptest.Server) {
		ctx := context.Background()

		os.Setenv(taskutil.MetadataAPIURLEnvName, ts.URL)
		c, err := decorators.NewDefaultClientFromEnv()
		require.NoError(t, err)
		require.NotNil(t, c)

		for name, exp := range opts.ExpectedDecorators {
			require.NoError(t, c.Set(ctx, name, exp))
		}
	}, opts)
}
