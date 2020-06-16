package task

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
	"github.com/puppetlabs/relay-sdk-go/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCredentialOutput(t *testing.T) {
	credentialSpec := make(map[string]string)

	data := base64.StdEncoding.EncodeToString([]byte("testdata"))
	credentialSpec["ca.pem"] = data
	credentialSpec["key.pem"] = data

	opts := testutil.MockMetadataAPIOptions{
		SpecResponse: map[string]interface{}{
			"value": map[string]interface{}{
				"credentials": credentialSpec,
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
			err := task.ProcessCredentials(dir)
			require.Nil(t, err, "err is not nil")

			b, err := ioutil.ReadFile(filepath.Join(dir, "ca.pem"))
			require.NoError(t, err)

			assert.Equal(t, "testdata", string(b))
		})
	}, opts)
}
