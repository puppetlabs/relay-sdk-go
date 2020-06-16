package taskutil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/puppetlabs/horsehead/v2/encoding/transfer"
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

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("the spec was not found")
	} else if resp.StatusCode == http.StatusInternalServerError {
		return nil, errors.New("an unexpected server error was encountered when retrieving the spec")
	} else if resp.StatusCode == http.StatusUnprocessableEntity {
		// Spec evaluation failed
		// TODO: Is this the correct behavior?
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
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
	Client   *http.Client
	SpecURL  string
	SpecPath string
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
	tree, err := TreeFromDefaultPlan(opts)
	if err != nil {
		return err
	}

	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339Nano),
		),
		ZeroFields: true,
		Result:     target,
		TagName:    "spec",
	})

	return d.Decode(tree)
}

func TreeFromDefaultPlan(opts DefaultPlanOptions) (interface{}, error) {
	location := opts.SpecURL
	var err error

	if location == "" {
		location, err = MetadataSpecURL()
		if err != nil {
			return nil, err
		}
		if location == "" {
			return nil, fmt.Errorf("%s was empty", MetadataAPIURLEnvName)
		}
	}

	u, err := url.Parse(location)
	if err != nil {
		return nil, fmt.Errorf("parsing spec URL %s failed: %+v", location, err)
	}

	if opts.SpecPath != "" {
		qs := u.Query()
		qs.Set("q", opts.SpecPath)
		qs.Set("lang", "jsonpath-template")

		u.RawQuery = qs.Encode()
	}

	var loader SpecLoader

	switch u.Scheme {
	case "http", "https":
		loader = NewRemoteSpecLoader(u, opts.Client)
	default:
		return nil, fmt.Errorf("unknown scheme %s in spec URL", u.Scheme)
	}

	r, err := loader.LoadSpec()
	if err != nil || r == nil {
		return nil, err
	}

	var env struct {
		Value transfer.JSONInterface `json:"value"`
	}
	if err := json.NewDecoder(r).Decode(&env); err != nil {
		return nil, err
	}

	return env.Value.Data, nil
}
