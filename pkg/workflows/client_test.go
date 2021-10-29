package workflows_test

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/puppetlabs/relay-sdk-go/pkg/envelope"
	"github.com/puppetlabs/relay-sdk-go/pkg/model"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
	"github.com/puppetlabs/relay-sdk-go/pkg/testutil"
	"github.com/puppetlabs/relay-sdk-go/pkg/workflows"
	"github.com/stretchr/testify/require"
)

func TestRunWorkflow(t *testing.T) {
	opts := testutil.MockMetadataAPIOptions{
		WorkflowRunResponses: map[string]interface{}{
			"test-workflow": envelope.PostWorkflowRunResponseEnvelope{
				WorkflowRun: &model.WorkflowRun{
					Name:      "test-workflow",
					RunNumber: 1,
					AppURL:    "https://unit-testing.relay.sh/workflows/test-workflow/runs/1",
				},
			},
		},
	}

	testutil.WithMockMetadataAPI(t, func(ts *httptest.Server) {
		ctx := context.Background()

		os.Setenv(taskutil.MetadataAPIURLEnvName, ts.URL)
		c, err := workflows.NewDefaultWorkflowsClientFromEnv()
		require.NoError(t, err)
		require.NotNil(t, c)

		wfr, err := c.Run(ctx, "test-workflow", map[string]string{
			"test-param": "test-param-value",
		})
		require.NoError(t, err)
		require.NotNil(t, wfr)

		require.Equal(t, "test-workflow", wfr.Name)
		require.Equal(t, int32(1), wfr.RunNumber)
		require.Equal(t, "https://unit-testing.relay.sh/workflows/test-workflow/runs/1", wfr.AppURL)

		wfr, err = c.Run(ctx, "test-workflow-not", nil)
		require.Error(t, err)
		require.Nil(t, wfr)
	}, opts)
}
