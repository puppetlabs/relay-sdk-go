package workflows

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/puppetlabs/relay-client-go/client/pkg/client/openapi"
	"github.com/puppetlabs/relay-sdk-go/pkg/envelope"
	"github.com/puppetlabs/relay-sdk-go/pkg/model"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
)

var (
	ErrWorkflowsClientNameEmpty = errors.New("name is required but was empty")
)

type WorkflowsClient interface {
	Run(ctx context.Context, name string, parameters map[string]string) (*model.WorkflowRun, error)
}

type DefaultWorkflowsClient struct {
	apiURL *url.URL
}

func (d *DefaultWorkflowsClient) Run(ctx context.Context, name string, parameters map[string]string) (*model.WorkflowRun, error) {
	if name == "" {
		return nil, ErrWorkflowsClientNameEmpty
	}

	loc := *d.apiURL
	loc.Path = path.Join(loc.Path, name, "run")

	reqEnv := envelope.PostWorkflowRunRequestEnvelope{
		Parameters: make(map[string]openapi.WorkflowRunParameter),
	}

	for k, v := range parameters {
		reqEnv.Parameters[k] = openapi.WorkflowRunParameter{
			Value: v,
		}
	}

	encoded, err := json.Marshal(reqEnv)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", loc.String(), bytes.NewBuffer(encoded))
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", "application/json")
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	env := envelope.PostWorkflowRunResponseEnvelope{}

	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return nil, err
	}

	return env.WorkflowRun, nil
}

func NewDefaultWorkflowsClient(location *url.URL) *DefaultWorkflowsClient {
	return &DefaultWorkflowsClient{apiURL: location}
}

func NewDefaultWorkflowsClientFromEnv() (*DefaultWorkflowsClient, error) {
	locstr, err := taskutil.MetadataURL("workflows")
	if err != nil {
		return nil, err
	}

	loc, err := url.Parse(locstr)
	if err != nil {
		return nil, err
	}

	return NewDefaultWorkflowsClient(loc), nil
}
