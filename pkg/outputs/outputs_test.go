package outputs_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/puppetlabs/horsehead/v2/encoding/transfer"
	"github.com/puppetlabs/nebula-sdk/pkg/outputs"
	"github.com/puppetlabs/nebula-sdk/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutputs(t *testing.T) {
	opts := testutil.MockMetadataAPIOptions{
		Outputs: map[testutil.MockOutputKey]testutil.MockOutput{
			testutil.MockOutputKey{TaskName: "test1", Key: "raw"}: {
				ResponseObject: outputs.Output{
					TaskName: "test1",
					Key:      "output",
					Value:    transfer.JSONInterface{Data: "Hello, test!"},
				},
			},
			testutil.MockOutputKey{TaskName: "test1", Key: "encoded"}: {
				ResponseObject: outputs.Output{
					TaskName: "test1",
					Key:      "output",
					Value: transfer.JSONInterface{Data: map[string]interface{}{
						"a": "Hello, \x90!",
						"b": "Goodbye, world!",
					}},
				},
			},
		},
	}

	testutil.WithMockMetadataAPI(t, func(ts *httptest.Server) {
		ctx := context.Background()

		u, err := url.Parse(ts.URL + `/outputs`)
		require.NoError(t, err)

		client := outputs.NewDefaultOutputsClient(u)

		tt := []struct {
			TaskName, Key string
			ExpectedValue interface{}
			ExpectedError error
		}{
			{
				TaskName:      "test1",
				Key:           "raw",
				ExpectedValue: "Hello, test!",
			},
			{
				TaskName: "test1",
				Key:      "encoded",
				ExpectedValue: map[string]interface{}{
					"a": "Hello, \x90!",
					"b": "Goodbye, world!",
				},
			},
			{
				TaskName:      "test2",
				Key:           "key",
				ExpectedError: outputs.ErrOutputsClientNotFound,
			},
		}
		for _, test := range tt {
			t.Run(fmt.Sprintf("%s/%s", test.TaskName, test.Key), func(t *testing.T) {
				s, err := client.GetOutput(ctx, test.TaskName, test.Key)
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
