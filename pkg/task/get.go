package task

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/puppetlabs/horsehead/v2/encoding/transfer"
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
	a := transfer.JSONInterface{}
	_ = a.UnmarshalJSON(body)
	if path != "" {
		// If a path is specified, `ni get` returns the single value, not json
		ret = fmt.Sprintf("%s", a.Data)
	} else {
		// If no path is specified, `ni get` returns the full encoded json
		data, _ := a.MarshalJSON()
		ret = string(data)
	}

	return []byte(ret), nil
}
