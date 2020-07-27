package task

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
	"github.com/puppetlabs/relay-sdk-go/pkg/testutil"
	"github.com/stretchr/testify/require"
)

func TestClusterOutput(t *testing.T) {
	data := base64.StdEncoding.EncodeToString([]byte("cadata"))

	opts := testutil.MockMetadataAPIOptions{
		SpecResponse: map[string]interface{}{
			"value": map[string]interface{}{
				"cluster": map[string]interface{}{
					"name": "test1",
					"connection": map[string]interface{}{
						"token":                "tokendata",
						"certificateAuthority": data,
						"server":               "url",
					},
				},
			},
		},
	}

	testutil.WithMockMetadataAPI(t, func(ts *httptest.Server) {
		opts := taskutil.DefaultPlanOptions{
			Client:  ts.Client(),
			SpecURL: fmt.Sprintf("%s/spec", ts.URL),
		}

		task := NewTaskInterface(opts)

		testutil.WithTemporaryDirectory(t, "output-", func(dir string) {
			err := task.ProcessClusters(dir)

			require.Nil(t, err, "err is not nil")
			require.FileExists(t, filepath.Join(dir, "test1", "kubeconfig"))
			content, err := ioutil.ReadFile(filepath.Join(dir, "test1", "kubeconfig"))
			assert.Matches(t, string(content), "token: tokendata")
			assert.Matches(t, string(content), "server: url")
			assert.Matches(t, string(content), fmt.Sprintf("certificate-authority-data: %s", data))
		})
	}, opts)
}

func TestOldClusterOutput(t *testing.T) {
	data := base64.StdEncoding.EncodeToString([]byte("cadata"))

	opts := testutil.MockMetadataAPIOptions{
		SpecResponse: map[string]interface{}{
			"value": map[string]interface{}{
				"cluster": map[string]interface{}{
					"name":   "test1",
					"token":  "tokendata",
					"cadata": data,
					"url":    "url",
				},
			},
		},
	}

	testutil.WithMockMetadataAPI(t, func(ts *httptest.Server) {
		opts := taskutil.DefaultPlanOptions{
			Client:  ts.Client(),
			SpecURL: fmt.Sprintf("%s/spec", ts.URL),
		}

		task := NewTaskInterface(opts)

		testutil.WithTemporaryDirectory(t, "output-", func(dir string) {
			err := task.ProcessClusters(dir)
			require.Nil(t, err, "err is not nil")
			require.FileExists(t, filepath.Join(dir, "test1", "kubeconfig"))
			content, err := ioutil.ReadFile(filepath.Join(dir, "test1", "kubeconfig"))
			assert.Matches(t, string(content), "token: tokendata")
			assert.Matches(t, string(content), "server: url")
			assert.Matches(t, string(content), fmt.Sprintf("certificate-authority-data: %s", data))
		})
	}, opts)
}
