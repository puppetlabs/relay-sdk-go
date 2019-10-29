package v1

import (
	"github.com/puppetlabs/nebula-sdk/pkg/workflow/asset"
	"github.com/xeipuuv/gojsonschema"
)

var WorkflowSchema *gojsonschema.Schema

func init() {
	workflowSchemaLoader := gojsonschema.NewStringLoader(asset.MustAssetString("schemas/v1/Workflow.json"))

	workflowSchema, err := gojsonschema.NewSchema(workflowSchemaLoader)
	if err != nil {
		panic(err)
	}

	WorkflowSchema = workflowSchema
}
