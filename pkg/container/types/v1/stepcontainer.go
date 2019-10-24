package v1

import (
	"io"
	"strings"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

// StepContainer represents an object with kind "StepContainer".
type StepContainer struct {
	*StepContainerTemplate `yaml:",inline"`

	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

func NewStepContainerFromString(data string) (*StepContainer, error) {
	return NewStepContainerFromReader(strings.NewReader(data))
}

func NewStepContainerFromReader(r io.Reader) (*StepContainer, error) {
	sc := &StepContainer{}
	if err := yaml.NewDecoder(r).Decode(&sc); err != nil {
		return nil, err
	}

	result, err := StepContainerSchema.Validate(gojsonschema.NewGoLoader(sc))
	if err != nil {
		return nil, err
	} else if !result.Valid() {
		// TODO: Aggregate and return errors.
	}

	return sc, nil
}
