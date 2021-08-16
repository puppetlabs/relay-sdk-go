package taskutil

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/puppetlabs/leg/encoding/transfer"
	"github.com/puppetlabs/leg/timeutil/pkg/retry"
)

func LoadData(c *http.Client, u *url.URL, timeout time.Duration) (interface{}, error) {
	loader := NewRemoteLoader(u, c)

	var env struct {
		Value transfer.JSONInterface `json:"value"`
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	waitOptions := []retry.WaitOption{}

	err := retry.Wait(ctx, func(ctx context.Context) (bool, error) {
		r, rerr := loader.Load()
		if rerr != nil {
			return false, rerr
		}

		if r == nil {
			return true, nil
		}

		if derr := json.NewDecoder(r).Decode(&env); derr != nil {
			return false, derr
		}

		return true, nil
	}, waitOptions...)

	if err != nil {
		return nil, err
	}

	return env.Value.Data, nil
}
