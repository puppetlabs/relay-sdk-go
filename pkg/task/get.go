package task

import (
	"context"

	"github.com/puppetlabs/nebula-sdk/pkg/taskutil"
	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/evaluate"
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
		eval := eval.Copy(evaluate.WithResultMapper(evaluate.NewJSONResultMapper()))
		output, err := eval.EvaluateAll(context.Background())
		if err == nil {
			return output.Value.([]byte), nil
		}
	}

	return nil, nil
}
