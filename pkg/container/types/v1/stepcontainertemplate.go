package v1

import (
	"io"
	"strings"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

const (
	StepContainerTemplateKind = "StepContainerTemplate"
)

// StepContainerTemplate represents an object with kind "StepContainerTemplate".
type StepContainerTemplate struct {
	*StepContainerCommon `yaml:",inline"`

	Version string `yaml:"version" json:"version,omitempty"`
	Kind    string `yaml:"kind" json:"kind,omitempty"`
}

func NewStepContainerTemplateFromString(data string) (*StepContainerTemplate, error) {
	return NewStepContainerTemplateFromReader(strings.NewReader(data))
}

func NewStepContainerTemplateFromReader(r io.Reader) (*StepContainerTemplate, error) {
	sc := &StepContainerTemplate{}
	if err := yaml.NewDecoder(r).Decode(&sc); err != nil {
		return nil, err
	}

	if sc.Version != Version || sc.Kind != StepContainerTemplateKind {
		return nil, &InvalidVersionKindError{
			ExpectedVersion: Version,
			ExpectedKind:    StepContainerTemplateKind,
			GotVersion:      sc.Version,
			GotKind:         sc.Kind,
		}
	}

	result, err := StepContainerTemplateSchema.Validate(gojsonschema.NewGoLoader(sc))
	if err != nil {
		return nil, err
	} else if !result.Valid() {
		return nil, schemaError(result.Errors())
	}

	return sc, nil
}
