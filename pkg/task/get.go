package task

import (
	"context"
	"encoding/json"

	"github.com/puppetlabs/nebula-sdk/pkg/taskutil"
)

func (ti *TaskInterface) ReadData(path string) ([]byte, error) {
	eval, err := taskutil.EvaluatorFromDefaultPlan(ti.opts)
	if err != nil {
		return nil, err
	}

	if path != "" {
		output, _ := eval.EvaluateQuery(context.Background(), path)
		switch vt := output.Value.(type) {
		case string:
			return []byte(vt), nil
		}
	} else {
		output, _ := eval.EvaluateAll(context.Background())
		return json.Marshal(output.Value)
	}

	return nil, nil
}
