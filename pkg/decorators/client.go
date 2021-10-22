package decorators

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
)

type DecoratorClient interface {
	Set(ctx context.Context, name string, values map[string]string) error
}

type DefaultClient struct {
	apiURL *url.URL
}

func (d *DefaultClient) Set(ctx context.Context, name string, values map[string]string) error {
	p := &url.URL{Path: path.Join(d.apiURL.Path, name)}
	loc := d.apiURL.ResolveReference(p)

	encoded, err := json.Marshal(values)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", loc.String(), bytes.NewBuffer(encoded))
	if err != nil {
		return err
	}

	req.Header.Set("content-type", "application/json")
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	return nil
}

func NewDefaultClient(location *url.URL) *DefaultClient {
	return &DefaultClient{apiURL: location}
}

func NewDefaultClientFromEnv() (*DefaultClient, error) {
	locstr, err := taskutil.MetadataURL("decorators")
	if err != nil {
		return nil, err
	}

	loc, err := url.Parse(locstr)
	if err != nil {
		return nil, err
	}

	return NewDefaultClient(loc), nil
}
