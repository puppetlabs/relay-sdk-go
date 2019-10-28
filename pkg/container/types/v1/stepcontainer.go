package v1

import (
	"io"
	"strings"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

const (
	StepContainerKind = "StepContainer"
)

// StepContainer represents an object with kind "StepContainer".
type StepContainer struct {
	*StepContainerCommon `yaml:",inline"`

	Version     string `yaml:"version" json:"version,omitempty"`
	Kind        string `yaml:"kind" json:"kind,omitempty"`
	Title       string `yaml:"title" json:"title"`
	Description string `yaml:"description" json:"description"`
}

func NewStepContainerFromString(data string) (*StepContainer, error) {
	return NewStepContainerFromReader(strings.NewReader(data))
}

func NewStepContainerFromReader(r io.Reader) (*StepContainer, error) {
	sc := &StepContainer{}
	if err := yaml.NewDecoder(r).Decode(&sc); err != nil {
		return nil, err
	}

	if sc.Version != Version || sc.Kind != StepContainerKind {
		return nil, &InvalidVersionKindError{
			ExpectedVersion: Version,
			ExpectedKind:    StepContainerKind,
			GotVersion:      sc.Version,
			GotKind:         sc.Kind,
		}
	}

	result, err := StepContainerSchema.Validate(gojsonschema.NewGoLoader(sc))
	if err != nil {
		return nil, err
	} else if !result.Valid() {
		return nil, schemaError(result.Errors())
	}

	return sc, nil
}
