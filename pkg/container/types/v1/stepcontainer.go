package v1

import (
	"io"
	"io/ioutil"

	"github.com/puppetlabs/nebula-sdk/pkg/util/typeutil"
	"gopkg.in/yaml.v3"
)

const (
	StepContainerKind = "StepContainer"
)

var StepContainerVersionKindExpectation = typeutil.NewVersionKindExpectation(Version, StepContainerKind)

// StepContainer represents an object with kind "StepContainer".
type StepContainer struct {
	*typeutil.VersionKind `yaml:",inline"`
	*StepContainerCommon  `yaml:",inline"`

	Title       string `yaml:"title" json:"title"`
	Description string `yaml:"description" json:"description"`
}

func NewStepContainerFromString(data string) (*StepContainer, error) {
	if _, err := StepContainerVersionKindExpectation.NewFromYAMLString(data); err != nil {
		return nil, err
	}

	if err := typeutil.ValidateYAMLString(StepContainerSchema, data); err != nil {
		return nil, err
	}

	sc := &StepContainer{}
	if err := yaml.Unmarshal([]byte(data), &sc); err != nil {
		return nil, err
	}

	return sc, nil
}

func NewStepContainerFromReader(r io.Reader) (*StepContainer, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return NewStepContainerFromString(string(b))
}
