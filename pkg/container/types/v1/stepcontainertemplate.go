package v1

import (
	"io"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
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

	result, err := StepContainerTemplateSchema.Validate(gojsonschema.NewGoLoader(sc))
	if err != nil {
		return nil, err
	} else if !result.Valid() {
		// XXX: FIXME: Aggregate and return errors.
		spew.Dump(result.Errors())
		panic("no")
	}

	return sc, nil
}
