package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"
)

type MockMetadataAPIOptions struct {
	SpecResponse       interface{}
	SpecQueryResponses map[string]interface{}
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
		default:
			http.NotFound(w, r)
		}

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
