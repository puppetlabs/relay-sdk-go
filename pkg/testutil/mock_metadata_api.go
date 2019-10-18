package testutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockMetadataAPIOptions struct {
	Name           string
	ResponseObject interface{}
}

func WithMockMetadataAPI(t *testing.T, fn func(ts *httptest.Server), opts MockMetadataAPIOptions) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, fmt.Sprintf("/specs/%s", opts.Name)) {
			if err := json.NewEncoder(w).Encode(opts.ResponseObject); err != nil {
				panic(err)
			}

			return
		}

		http.NotFound(w, r)
		return
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	fn(ts)
}
