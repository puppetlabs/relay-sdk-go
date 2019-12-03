package task

import (
	"net/http/httptest"
	"testing"

	"github.com/puppetlabs/nebula-sdk/pkg/testutil"
	"github.com/stretchr/testify/require"
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

func TestGitURLMatching(t *testing.T) {
	cases := []struct {
		gitURL      string
		shouldMatch bool
	}{
		{gitURL: "git@github.com:example/example-repo", shouldMatch: true},
		{gitURL: "git@github.com:example/example-repo.git", shouldMatch: true},
		{gitURL: "git@github.com/example/example-repo.git", shouldMatch: false},
		{gitURL: "foobar@github.com:example/example-repo", shouldMatch: true},
		{gitURL: "foobar@github.com:example/example-repo.git", shouldMatch: true},
		{gitURL: "foobar@github.com/example/example-repo.git", shouldMatch: false},
		{gitURL: "https://example.com", shouldMatch: false},
		{gitURL: "git@example.com:example-org/example-repo", shouldMatch: true},
	}

	for _, c := range cases {
		t.Run(c.gitURL, func(t *testing.T) {
			matches, err := gitURLComponents(c.gitURL)

			if !c.shouldMatch {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.GreaterOrEqual(t, len(matches), 4)
			}
		})
	}
}
