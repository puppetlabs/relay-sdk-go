package typeutil

import "gopkg.in/yaml.v2"

type VersionKind struct {
	Version string `yaml:"version" json:"version,omitempty"`
	Kind    string `yaml:"kind" json:"kind,omitempty"`
}

type VersionKindExpectation struct {
	Version string
	Kind    string
}

func (vke *VersionKindExpectation) NewFromYAMLString(data string) (*VersionKind, error) {
	vk := &VersionKind{}
	if err := yaml.Unmarshal([]byte(data), &vk); err != nil {
		if _, ok := err.(*yaml.TypeError); ok {
			return nil, &InvalidVersionKindError{
				ExpectedVersion: vke.Version,
				ExpectedKind:    vke.Kind,
			}
		}
		return nil, err
	}

	if vk.Version != vke.Version || vk.Kind != vke.Kind {
		return nil, &InvalidVersionKindError{
			ExpectedVersion: vke.Version,
			ExpectedKind:    vke.Kind,
			GotVersion:      vk.Version,
			GotKind:         vk.Kind,
		}
	}

	return vk, nil
}

func NewVersionKindExpectation(version, kind string) *VersionKindExpectation {
	return &VersionKindExpectation{
		Version: version,
		Kind:    kind,
	}
}
