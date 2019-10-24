package task

import (
	"net/http/httptest"
	"testing"

	"github.com/puppetlabs/nebula-sdk/pkg/testutil"
)

func TestGitOutput(t *testing.T) {
	t.Skip("Functional testing harness. Needs to be completed.")

	opts := testutil.SingleSpecMockMetadataAPIOptions("test1", testutil.MockSpec{
		ResponseObject: map[string]interface{}{
			"git": map[string]interface{}{
				"ssh_key":     "<ssh_key>",
				"known_hosts": "<known_hosts>",
				"name":        "<name>",
				"repository":  "<repository>",
			},
		},
	})

	testutil.WithMockMetadataAPI(t, func(ts *httptest.Server) {
		// opts := taskutil.DefaultPlanOptions{
		// 	Client:  ts.Client(),
		// 	SpecURL: fmt.Sprintf("%s/specs/test1", ts.URL),
		// }

		// task := NewTaskInterface(opts)

		// err := task.CloneRepository("master", "output")
		// require.Nil(t, err, "err is not nil")
	}, opts)
}
