package taskutil

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/mitchellh/mapstructure"
)

type DefaultPlanOptions struct {
	Client      *http.Client
	SpecURL     string
	SpecPath    string
	SpecTimeout time.Duration
}

func MetadataSpecURL() (string, error) {
	return MetadataURL("spec")
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
			mapstructure.TextUnmarshallerHookFunc(),
		),
		ZeroFields: true,
		Result:     target,
		TagName:    "spec",
	})
	if err != nil {
		return err
	}

	return d.Decode(tree)
}

func TreeFromDefaultPlan(opts DefaultPlanOptions) (interface{}, error) {
	location := opts.SpecURL

	timeout := DefaultMetadataTimeout
	if opts.SpecTimeout > 0 {
		timeout = opts.SpecTimeout
	}

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

	return LoadData(opts.Client, u, timeout)
}
