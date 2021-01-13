package task

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/puppetlabs/leg/encoding/transfer"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
	"github.com/puppetlabs/relay-sdk-go/pkg/testutil"
	"github.com/stretchr/testify/require"
)

func mustTransferEncodeJSON(t *testing.T, data []byte) transfer.JSONOrStr {
	j, err := transfer.EncodeJSON(data)
	require.NoError(t, err)

	return j
}

func TestGetOutput(t *testing.T) {
	testSpec := map[string]interface{}{
		"name":    "test1",
		"someNum": float64(12),
		"data":    []interface{}{"something", "else", "Hello, \x90!"},
	}

	testSpecData, err := json.Marshal(transfer.JSONInterface{Data: testSpec})
	require.NoError(t, err)

	opts := testutil.MockMetadataAPIOptions{
		SpecResponse: map[string]interface{}{
			"value": transfer.JSONInterface{Data: testSpec},
		},
		SpecQueryResponses: map[string]interface{}{
			"{.name}":    map[string]interface{}{"value": testSpec["name"]},
			"{.someNum}": map[string]interface{}{"value": testSpec["someNum"]},
			"{.data[0]}": map[string]interface{}{"value": mustTransferEncodeJSON(t, []byte(testSpec["data"].([]interface{})[0].(string)))},
			"{.data[1]}": map[string]interface{}{"value": mustTransferEncodeJSON(t, []byte(testSpec["data"].([]interface{})[1].(string)))},
			"{.data[2]}": map[string]interface{}{"value": mustTransferEncodeJSON(t, []byte(testSpec["data"].([]interface{})[2].(string)))},
		},
	}

	testutil.WithMockMetadataAPI(t, func(ts *httptest.Server) {
		opts := taskutil.DefaultPlanOptions{
			Client:      ts.Client(),
			SpecURL:     fmt.Sprintf("%s/spec", ts.URL),
			SpecTimeout: 5 * time.Second,
		}

		task := NewTaskInterface(opts)

		tests := []struct {
			Path     string
			Expected string
		}{
			{
				Path:     "{.name}",
				Expected: testSpec["name"].(string),
			},
			{
				Path:     "{.someNum}",
				Expected: strconv.Itoa(int(testSpec["someNum"].(float64))),
			},
			{
				Path:     "{.data[0]}",
				Expected: testSpec["data"].([]interface{})[0].(string),
			},
			{
				Path:     "{.data[1]}",
				Expected: testSpec["data"].([]interface{})[1].(string),
			},
			{
				Path:     "{.data[2]}",
				Expected: testSpec["data"].([]interface{})[2].(string),
			},
			{
				Path:     "{.nothing}",
				Expected: "",
			},
			{
				Path:     "",
				Expected: string(testSpecData),
			},
		}
		for _, test := range tests {
			t.Run(test.Path, func(t *testing.T) {
				output, err := task.ReadData(test.Path)
				require.NoError(t, err)
				require.Equal(t, test.Expected, string(output))
			})
		}
	}, opts)
}
