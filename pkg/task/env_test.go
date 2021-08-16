package task_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/puppetlabs/leg/encoding/transfer"
	"github.com/puppetlabs/relay-sdk-go/pkg/task"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
	"github.com/puppetlabs/relay-sdk-go/pkg/testutil"
	"github.com/stretchr/testify/require"
)

func TestGetEnvironmentVariables(t *testing.T) {
	testEnvironment := map[string]interface{}{
		"ENV1": "value1",
		"ENV2": float64(12),
		"ENV3": true,
	}

	testEnvironmentData, err := json.Marshal(transfer.JSONInterface{Data: testEnvironment})
	require.NoError(t, err)

	opts := testutil.MockMetadataAPIOptions{
		SpecResponse: map[string]interface{}{
			"value": transfer.JSONInterface{Data: testEnvironment},
		},
	}

	testutil.WithMockMetadataAPI(t, func(ts *httptest.Server) {
		opts := taskutil.DefaultPlanOptions{
			Client:      ts.Client(),
			SpecURL:     fmt.Sprintf("%s/environment", ts.URL),
			SpecTimeout: 5 * time.Second,
		}

		ti := task.NewTaskInterface(opts)

		output, err := ti.ReadEnvironmentVariables()
		require.NoError(t, err)
		require.Equal(t, string(testEnvironmentData), string(output))
	}, opts)
}

func TestGetEnvironmentVariable(t *testing.T) {
	testEnvironment := map[string]interface{}{
		"ENV1": "value1",
		"ENV2": float64(12),
		"ENV3": true,
	}

	opts := testutil.MockMetadataAPIOptions{
		SpecQueryResponses: map[string]interface{}{
			"ENV1": map[string]interface{}{"value": testEnvironment["ENV1"]},
			"ENV2": map[string]interface{}{"value": testEnvironment["ENV2"]},
			"ENV3": map[string]interface{}{"value": testEnvironment["ENV3"]},
		},
	}
	tests := []struct {
		Name     string
		Expected string
	}{
		{
			Name:     "ENV1",
			Expected: testEnvironment["ENV1"].(string),
		},
		{
			Name:     "ENV2",
			Expected: strconv.Itoa(int(testEnvironment["ENV2"].(float64))),
		},
		{
			Name:     "ENV3",
			Expected: strconv.FormatBool(testEnvironment["ENV3"].(bool)),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			testutil.WithMockMetadataAPI(t, func(ts *httptest.Server) {
				opts := taskutil.DefaultPlanOptions{
					Client:      ts.Client(),
					SpecURL:     fmt.Sprintf("%s/environment/%s", ts.URL, test.Name),
					SpecTimeout: 5 * time.Second,
				}

				ti := task.NewTaskInterface(opts)

				output, err := ti.ReadEnvironmentVariable(test.Name)
				require.NoError(t, err)
				require.Equal(t, test.Expected, string(output))
			}, opts)
		})
	}
}
