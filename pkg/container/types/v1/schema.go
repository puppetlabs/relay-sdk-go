package v1

import (
	"fmt"

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

	stepContainerTemplateSchemaExternalLoader := gojsonschema.NewSchemaLoader()
	stepContainerTemplateSchemaExternalLoader.Validate = true
	if err := stepContainerTemplateSchemaExternalLoader.AddSchemas(gojsonschema.NewStringLoader(asset.MustAssetString("schemas/v1/StepContainer-common.json"))); err != nil {
		panic(fmt.Errorf("failed to load common schema for StepContainerTemplate: %+v", err))
	}

	stepContainerTemplateSchema, err := stepContainerTemplateSchemaExternalLoader.Compile(stepContainerTemplateSchemaLoader)
	if err != nil {
		panic(fmt.Errorf("failed to load schema for StepContainerTemplate: %+v", err))
	}

	stepContainerSchemaExternalLoader := gojsonschema.NewSchemaLoader()
	stepContainerSchemaExternalLoader.Validate = true
	if err := stepContainerSchemaExternalLoader.AddSchemas(gojsonschema.NewStringLoader(asset.MustAssetString("schemas/v1/StepContainer-common.json"))); err != nil {
		panic(fmt.Errorf("failed to load common schema for StepContainer: %+v", err))
	}

	stepContainerSchema, err := stepContainerSchemaExternalLoader.Compile(stepContainerSchemaLoader)
	if err != nil {
		panic(fmt.Errorf("failed to load schema for StepContainer: %+v", err))
	}

	StepContainerTemplateSchema = stepContainerTemplateSchema
	StepContainerSchema = stepContainerSchema
}

func schemaError(errs []gojsonschema.ResultError) error {
	fes := make([]*SchemaValidationFieldError, len(errs))
	for i, err := range errs {
		fes[i] = &SchemaValidationFieldError{
			Context:     err.Context().String(),
			Field:       err.Field(),
			Description: err.Description(),
		}
	}

	return &SchemaValidationError{
		FieldErrors: fes,
	}
}
