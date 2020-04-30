package testutil

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"

	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/evaluate"
	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/parse"
	"github.com/stretchr/testify/require"
)

type MockSecret struct {
	ResponseObject interface{}
}

type MockOutputKey struct {
	TaskName string
	Key      string
}

type MockOutput struct {
	ResponseObject interface{}
}

type MockSpec struct {
	ResponseObject interface{}
}

type MockMetadataAPIOptions struct {
	Secrets map[string]MockSecret
	Outputs map[MockOutputKey]MockOutput
	Spec    MockSpec
}

func SingleSpecMockMetadataAPIOptions(name string, spec MockSpec) MockMetadataAPIOptions {
	return MockMetadataAPIOptions{
		Spec: spec,
	}
}

func WithMockMetadataAPI(t *testing.T, fn func(ts *httptest.Server), opts MockMetadataAPIOptions) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler, rest := shiftPath(r.URL.Path)
		switch handler {
		case "spec":
			// That is no mock...
			spec, _ := json.Marshal(opts.Spec.ResponseObject)
			// This returns a Tree which is an interface{} we don't want that
			tree, _ := parse.ParseJSONString(string(spec))

			ev := evaluate.NewEvaluator(
				evaluate.WithLanguage(evaluate.LanguageJSONPathTemplate),
			).ScopeTo(tree)

			var rv *evaluate.Result
			var err error
			ctx := context.Background()
			query := r.URL.Query().Get("q")
			if query != "" {
				rv, err = ev.EvaluateQuery(ctx, query)
			} else {
				rv, err = ev.EvaluateAll(ctx)
			}
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			err = json.NewEncoder(w).Encode(evaluate.NewJSONResultEnvelope(rv))
			require.NoError(t, err)

			return
		case "secrets":
			name, _ := shiftPath(rest)

			s, found := opts.Secrets[name]
			if found {
				if err := json.NewEncoder(w).Encode(s.ResponseObject); err != nil {
					panic(err)
				}

				return
			}
		case "outputs":
			var k MockOutputKey
			k.TaskName, rest = shiftPath(rest)
			k.Key, _ = shiftPath(rest)

			o, found := opts.Outputs[k]
			if found {
				if err := json.NewEncoder(w).Encode(o.ResponseObject); err != nil {
					panic(err)
				}

				return
			}
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
