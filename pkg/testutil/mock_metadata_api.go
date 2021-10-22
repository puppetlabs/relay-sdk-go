package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"

	"github.com/puppetlabs/relay-sdk-go/pkg/envelope"
	"github.com/stretchr/testify/require"
)

type MockMetadataAPIOptions struct {
	SpecResponse         interface{}
	SpecQueryResponses   map[string]interface{}
	WorkflowRunResponses map[string]interface{}
	ExpectedDecorators   map[string]map[string]string
}

func WithMockMetadataAPI(t *testing.T, fn func(ts *httptest.Server), opts MockMetadataAPIOptions) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler, subpath := shiftPath(r.URL.Path)
		switch handler {
		case "environment":
			resp := opts.SpecResponse
			if subpath != "" {
				name := path.Base(subpath)
				resp = opts.SpecQueryResponses[name]
			}

			if err := json.NewEncoder(w).Encode(resp); err != nil {
				panic(err)
			}

			return
		case "spec":
			resp := opts.SpecResponse
			if q := r.URL.Query().Get("q"); q != "" {
				sq, ok := opts.SpecQueryResponses[q]
				if !ok {
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}

				resp = sq
			}

			if err := json.NewEncoder(w).Encode(resp); err != nil {
				panic(err)
			}

			return
		case "workflows":
			if subpath == "" {
				t.Log("mock metadata api: missing subpath")
				break
			}

			var name string
			name, subpath = shiftPath(subpath)

			if path.Base(subpath) != "run" {
				t.Logf("mock metadata api: missing run subpath: %s", subpath)
				break
			}

			req := envelope.PostWorkflowRunRequestEnvelope{}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Logf("mock metadata api: failed to decode body: %v", err)
				w.WriteHeader(http.StatusNotAcceptable)
				return
			}

			resp, ok := opts.WorkflowRunResponses[name]
			if !ok {
				t.Logf("mock metadata api: missing response config: %s", name)
				break
			}

			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				panic(err)
			}

			return
		case "decorators":
			if subpath == "" {
				t.Log("mock metadata api: missing subpath")
				break
			}

			var name string
			name, _ = shiftPath(subpath)

			var env = make(map[string]string)
			if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
				t.Logf("mock metadata api: failed to decode body: %v", err)
				w.WriteHeader(http.StatusNotAcceptable)
				return
			}

			resp, ok := opts.ExpectedDecorators[name]
			if !ok {
				t.Logf("mock metadata api: missing response config: %s", name)
				break
			}

			require.Equal(t, resp, env)

			w.WriteHeader(http.StatusCreated)

			return
		}

		t.Logf("mock metadata api: request path: %s", r.URL.Path)
		http.NotFound(w, r)
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	fn(ts)
}

// shiftPath takes a URI path and returns the first segment as the head
// and the remaining segments as the tail.
func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], ""
	}

	return p[1:i], p[i:]
}
