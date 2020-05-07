package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/evaluate"
)

func (ti *TaskInterface) ReadData(path string) ([]byte, error) {
	u, err := url.Parse(ti.opts.SpecURL)
	if err != nil {
		return nil, err
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}

	if path != "" {
		q.Add("q", path)
		q.Add("lang", "jsonpath-template")

		u.RawQuery = q.Encode()
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("the spec was not found")
	}
	if resp.StatusCode == http.StatusInternalServerError {
		return nil, errors.New("an unexpected server error was encountered when retrieving the spec")
	}
	if resp.StatusCode == http.StatusUnprocessableEntity {
		// Spec evaluation failed
		return nil, nil
	}
	if resp.StatusCode/100 != 2 {
		// This will be a non-200 code other than 404, 422, or 500
		return nil, errors.New(fmt.Sprintf("an unexpected response %d was encountered when retrieving the spec", resp.StatusCode))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var output []byte

	var result evaluate.JSONResultEnvelope
	err = json.Unmarshal(body, &result)
	if path != "" {
		// If a path is specified, `ni get` returns the single value, not json
		output = []byte(fmt.Sprintf("%s", result.Value.Data))
	} else {
		// If no path is specified, `ni get` returns the full encoded json
		output, _ = result.Value.MarshalJSON()
	}

	return output, nil
}
