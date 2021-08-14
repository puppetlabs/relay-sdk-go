package taskutil

import (
	"fmt"
	"net/url"
)

func LoadEnvironment(opts DefaultPlanOptions) (interface{}, error) {
	location := opts.SpecURL

	timeout := DefaultMetadataTimeout
	if opts.SpecTimeout > 0 {
		timeout = opts.SpecTimeout
	}

	u, err := url.Parse(location)
	if err != nil {
		return nil, fmt.Errorf("parsing URL %s failed: %+v", location, err)
	}

	return LoadData(opts.Client, u, timeout)
}
