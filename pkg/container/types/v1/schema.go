package v1

import (
	"github.com/puppetlabs/nebula-sdk/pkg/container/asset"
	"github.com/xeipuuv/gojsonschema"
)

var (
	StepContainerTemplateSchema *gojsonschema.Schema
	StepContainerSchema         *gojsonschema.Schema
)

func init() {
	stepContainerTemplateSchemaLoader := gojsonschema.NewStringLoader(asset.MustAssetString("schemas/v1/StepContainerTemplate.json"))
	stepContainerSchemaLoader := gojsonschema.NewStringLoader(asset.MustAssetString("schemas/v1/StepContainer.json"))

	stepContainerTemplateSchema, err := gojsonschema.NewSchema(stepContainerTemplateSchemaLoader)
	if err != nil {
		panic(err)
	}

	stepContainerSchemaExternalLoader := gojsonschema.NewSchemaLoader()
	stepContainerSchemaExternalLoader.AddSchemas(stepContainerTemplateSchemaLoader)

	stepContainerSchema, err := stepContainerSchemaExternalLoader.Compile(stepContainerSchemaLoader)
	if err != nil {
		panic(err)
	}

	StepContainerTemplateSchema = stepContainerTemplateSchema
	StepContainerSchema = stepContainerSchema
}
