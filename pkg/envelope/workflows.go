package envelope

import (
	"github.com/puppetlabs/relay-client-go/client/pkg/client/openapi"
	"github.com/puppetlabs/relay-sdk-go/pkg/model"
)

type PostWorkflowRunResponseEnvelope struct {
	WorkflowRun *model.WorkflowRun `json:"workflow_run"`
}

type PostWorkflowRunRequestEnvelope struct {
	Parameters map[string]openapi.WorkflowRunParameter `json:"parameters"`
}
