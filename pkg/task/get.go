package task

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, nil
	}

	var ret string
	if path != "" {
		// If no path is specified, ni get returns the full json
		err = json.Unmarshal(body, &ret)
	} else {
		// If a path is specified, ni get returns the single value
		ret = string(body)
	}

	return []byte(ret), nil
}
