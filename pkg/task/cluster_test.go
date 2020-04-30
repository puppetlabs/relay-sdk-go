package task

import (
	"encoding/base64"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/puppetlabs/nebula-sdk/pkg/taskutil"
	"github.com/puppetlabs/nebula-sdk/pkg/testutil"
	"github.com/stretchr/testify/require"
)

func TestClusterOutput(t *testing.T) {
	data := base64.StdEncoding.EncodeToString([]byte("cadata"))

	opts := testutil.SingleSpecMockMetadataAPIOptions("test1", testutil.MockSpec{
		ResponseObject: map[string]interface{}{
			"cluster": map[string]interface{}{
				"name":   "test1",
				"token":  "tokendata",
				"cadata": data,
				"url":    "url",
			},
		},
	})

	testutil.WithMockMetadataAPI(t, func(ts *httptest.Server) {
		opts := taskutil.DefaultPlanOptions{
			Client:  ts.Client(),
			SpecURL: fmt.Sprintf("%s/spec", ts.URL),
		}

		task := NewTaskInterface(opts)

		testutil.WithTemporaryDirectory(t, "output-", func(dir string) {
			err := task.ProcessClusters(dir)
			require.Nil(t, err, "err is not nil")
		})
	}, opts)
}
