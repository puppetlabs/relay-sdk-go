package task

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/puppetlabs/horsehead/v2/encoding/transfer"
	"github.com/puppetlabs/nebula-sdk/pkg/taskutil"
	"github.com/puppetlabs/nebula-sdk/pkg/testutil"
	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/evaluate"
	"github.com/stretchr/testify/require"
)

type TestGetSpec struct {
	Name    string   `json:"name"`
	SomeNum int      `json:"someNum"`
	Data    []string `json:"data"`
}

func TestGetOutput(t *testing.T) {
	testSpec := &TestGetSpec{
		Name:    "test1",
		SomeNum: 12,
		Data:    []string{"something", "else", "Hello, \x90!"},
	}

	es := func(v ...string) (r []transfer.JSONOrStr) {
		for i := range v {
			a, _ := transfer.EncodeJSON([]byte(v[i]))
			r = append(r, a)
		}
		return
	}
	opts := testutil.SingleSpecMockMetadataAPIOptions("test1", testutil.MockSpec{
		ResponseObject: map[string]interface{}{
			"name":    "test1",
			"someNum": 12,
			"data":    es("something", "else", "Hello, \x90!"),
		},
	})

	testutil.WithMockMetadataAPI(t, func(ts *httptest.Server) {
		opts := taskutil.DefaultPlanOptions{
			Client:  ts.Client(),
			SpecURL: fmt.Sprintf("%s/specs/test1", ts.URL),
		}

		task := NewTaskInterface(opts)

		output, _ := task.ReadData("{.name}")
		require.Equal(t, testSpec.Name, string(output))

		output, _ = task.ReadData("{.someNum}")
		require.Equal(t, strconv.Itoa(testSpec.SomeNum), string(output))

		output, _ = task.ReadData("{.data[0]}")
		require.Equal(t, testSpec.Data[0], string(output))

		output, _ = task.ReadData("{.data[1]}")
		require.Equal(t, testSpec.Data[1], string(output))

		output, _ = task.ReadData("{.data[2]}")
		require.Equal(t, testSpec.Data[2], string(output))

		output, _ = task.ReadData("")
		var outputSpec TestGetSpec
		a := transfer.JSONInterface{}
		_ = a.UnmarshalJSON(output)
		e := evaluate.NewEvaluator()
		_, err := e.EvaluateInto(context.Background(), a.Data, &outputSpec)
		require.NoError(t, err)
		require.Equal(t, testSpec.Data[2], outputSpec.Data[2])

	}, opts)
}
