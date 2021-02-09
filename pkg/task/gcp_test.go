package task

import (
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/puppetlabs/relay-sdk-go/pkg/model"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
	"github.com/puppetlabs/relay-sdk-go/pkg/testutil"
	"github.com/stretchr/testify/require"
)

func TestGCPConnectionBackwardCompatibility(t *testing.T) {
	tests := []struct {
		Spec          map[string]interface{}
		ExpectedKey   string
		ExpectedError bool
	}{
		{
			Spec: map[string]interface{}{
				"serviceAccountKey": "TEST KEY",
			},
			ExpectedKey: "TEST KEY",
		},
		{
			Spec: map[string]interface{}{},
		},
		{
			Spec: map[string]interface{}{
				"somethingElse": "something else",
			},
			ExpectedError: true,
		},
		{
			Spec: map[string]interface{}{
				"connection": map[string]interface{}{
					"serviceAccountKey": "TEST KEY",
				},
			},
			ExpectedKey: "TEST KEY",
		},
	}
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			opts := testutil.MockMetadataAPIOptions{
				SpecResponse: map[string]interface{}{
					"value": test.Spec,
				},
			}

			testutil.WithMockMetadataAPI(t, func(ts *httptest.Server) {
				opts := taskutil.DefaultPlanOptions{
					Client:  ts.Client(),
					SpecURL: fmt.Sprintf("%s/spec", ts.URL),
				}

				var spec model.GCPDetails
				require.NoError(t, taskutil.PopulateSpecFromDefaultPlan(&spec, opts))

				key := spec.GetServiceAccountKey()
				if test.ExpectedError {
					require.Empty(t, key)
				} else if test.ExpectedKey != "" {
					require.Equal(t, test.ExpectedKey, key)
				} else {
					require.Empty(t, key)
				}
			}, opts)
		})
	}
}
