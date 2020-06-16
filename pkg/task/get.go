package task

import (
	"encoding/json"
	"fmt"

	"github.com/puppetlabs/horsehead/v2/encoding/transfer"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
)

func (ti *TaskInterface) ReadData(path string) ([]byte, error) {
	opts := ti.opts
	opts.SpecPath = path

	tree, err := taskutil.TreeFromDefaultPlan(opts)
	if err != nil || tree == nil {
		return nil, err
	}

	if opts.SpecPath != "" {
		// If a path is specified, `ni get` returns the single value, not json
		return []byte(fmt.Sprintf("%v", tree)), nil
	}

	return json.Marshal(transfer.JSONInterface{Data: tree})
}
