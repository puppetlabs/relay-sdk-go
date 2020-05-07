package task

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/puppetlabs/nebula-sdk/pkg/taskutil"
	"github.com/puppetlabs/nebula-sdk/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

type TestValues struct {
	Name    string   `yaml:"name" json:"name"`
	SomeNum int      `yaml:"someNum" json:"someNum"`
	Data    []string `yaml:"data" json:"data"`
}

type TestSpec struct {
	Values *TestValues `yaml:"values" json:"values"`
}

func TestGetFileOutput(t *testing.T) {
	testSpec := &TestSpec{
		Values: &TestValues{
			Name:    "test1",
			SomeNum: 12,
			Data:    []string{"something", "else"},
		},
	}

	opts := testutil.SingleSpecMockMetadataAPIOptions("test1", testutil.MockSpec{
		ResponseObject: testSpec,
	})

	testutil.WithMockMetadataAPI(t, func(ts *httptest.Server) {
		opts := taskutil.DefaultPlanOptions{
			Client:  ts.Client(),
			SpecURL: fmt.Sprintf("%s/spec", ts.URL),
		}

		task := NewTaskInterface(opts)

		testutil.WithTemporaryDirectory(t, "output-", func(dir string) {
			f := filepath.Join(dir, "values-test.yml")

			err := task.WriteFile(f, "values", "yaml")
			require.Nil(t, err, "err is not nil")

			b, err := ioutil.ReadFile(f)
			require.NoError(t, err)

			got := &TestValues{}
			require.NoError(t, yaml.Unmarshal(b, got))
			assert.EqualValues(t, testSpec.Values, got)
		})
	}, opts)
}
