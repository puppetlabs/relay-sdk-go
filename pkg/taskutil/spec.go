package taskutil

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/evaluate"
	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/parse"
)

const MetadataAPIURLEnvName = "METADATA_API_URL"

// SpecLoader returns an io.Reader containing the bytes
// of a task spec. This is used as input to a spec unmarshaler.
// An error is returned if the operation fails.
type SpecLoader interface {
	LoadSpec() (io.Reader, error)
}

type RemoteSpecLoader struct {
	u      *url.URL
	client *http.Client
}

func (r RemoteSpecLoader) LoadSpec() (io.Reader, error) {
	var client = http.DefaultClient

	if r.client != nil {
		client = http.DefaultClient
	}

	resp, err := client.Get(r.u.String())
	if err != nil {
		return nil, fmt.Errorf("network request failed: %+v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	buf := &bytes.Buffer{}

	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, fmt.Errorf("reading response from remote service failed: %+v", err)
	}

	return buf, nil
}

func NewRemoteSpecLoader(u *url.URL, client *http.Client) RemoteSpecLoader {
	return RemoteSpecLoader{
		u:      u,
		client: client,
	}
}

type LocalSpecLoader struct {
	path string
}

func (l LocalSpecLoader) LoadSpec() (io.Reader, error) {
	b, err := ioutil.ReadFile(l.path)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %+v", l.path, err)
	}

	return bytes.NewBuffer(b), nil
}

func NewLocalSpecLoader(path string) LocalSpecLoader {
	return LocalSpecLoader{path: path}
}

type SpecDecoder interface {
	DecodeSpec(io.Reader) (interface{}, error)
}

type DefaultPlanOptions struct {
	Client  *http.Client
	SpecURL string
}

func MetadataSpecURL() (string, error) {
	if e := os.Getenv(MetadataAPIURLEnvName); e != "" {
		u, err := url.Parse(e)
		if err != nil {
			return "", fmt.Errorf("could not parse %s value as url: %s", MetadataAPIURLEnvName, e)
		}
		u = u.ResolveReference(&url.URL{Path: "spec"})
		return u.String(), nil
	}
	return "", nil
}

func PopulateSpecFromDefaultPlan(target interface{}, opts DefaultPlanOptions) error {
	eval, err := EvaluatorFromDefaultPlan(opts)
	if err != nil {
		return err
	}

	resultTarget := evaluate.Result{
		Value: target,
	}
	// The existing behavior is to substitute empty strings where secrets,
	// outputs, etc. are currently missing. We should revisit this.
	_, err = eval.EvaluateInto(context.Background(), &resultTarget)
	if err != nil {
		return err
	}

	return nil
}

func EvaluatorFromDefaultPlan(opts DefaultPlanOptions) (*evaluate.ScopedEvaluator, error) {
	location := opts.SpecURL
	var err error

	if location == "" {
		location, err = MetadataSpecURL()
		if err != nil {
			return nil, err
		}
		if location == "" {
			return nil, errors.New(fmt.Sprintf("%s was empty", MetadataAPIURLEnvName))
		}
	}

	u, err := url.Parse(location)
	if err != nil {
		return nil, fmt.Errorf("parsing spec URL %s failed: %+v", location, err)
	}

	var loader SpecLoader

	switch u.Scheme {
	case "http", "https":
		loader = NewRemoteSpecLoader(u, opts.Client)
	default:
		return nil, fmt.Errorf("unknown scheme %s in spec URL", u.Scheme)
	}

	r, err := loader.LoadSpec()
	if err != nil {
		return nil, err
	}

	tree, err := parse.ParseJSON(r)
	if err != nil {
		return nil, err
	}

	ev := evaluate.NewEvaluator(
		evaluate.WithLanguage(evaluate.LanguageJSONPathTemplate),
	).ScopeTo(tree)
	return ev, nil
}
