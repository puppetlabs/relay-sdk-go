package def

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/imdario/mergo"
	v1 "github.com/puppetlabs/nebula-sdk/pkg/container/types/v1"
)

type TemplateContainer struct {
	TemplateName string
	TemplateData string
	DependsOn    []string
}

type TemplateSetting struct {
	Description string
	Value       interface{}
}

type Template struct {
	SDKVersion string
	Containers map[string]*TemplateContainer
	Settings   map[string]*TemplateSetting
}

func NewTemplateFromTyped(sctt *v1.StepContainerCommon, opts ...Option) (*Template, error) {
	var o Options
	for _, opt := range opts {
		opt(&o)
	}

	containers := make(map[string]*TemplateContainer, len(sctt.Containers))
	for name, container := range sctt.Containers {
		tc := &TemplateContainer{
			TemplateName: container.Template.Name,
			DependsOn:    container.DependsOn,
		}

		fr, err := NewFileRefFromTyped(container.Template, WithFileRefResolver(o.resolver))
		if err != nil {
			return nil, err
		}

		if err := fr.WithFile(func(f http.File) error {
			b, err := ioutil.ReadAll(f)
			if err != nil {
				return err
			}

			tc.TemplateData = string(b)
			return nil
		}); err != nil {
			return nil, err
		}

		containers[name] = tc
	}

	settings := make(map[string]*TemplateSetting, len(sctt.Settings))
	for name, setting := range sctt.Settings {
		settings[name] = &TemplateSetting{
			Description: setting.Description,
			Value:       setting.Value,
		}
	}

	t := &Template{
		SDKVersion: sctt.SDKVersion,
		Containers: containers,
		Settings:   settings,
	}

	// We need to merge this file with its parents.
	if sctt.Inherit != nil {
		fr, err := NewFileRefFromTyped(*sctt.Inherit, WithFileRefResolver(o.resolver))
		if err != nil {
			return nil, err
		}

		parent, err := NewTemplateFromFileRef(fr)
		if err != nil {
			return nil, err
		}

		if err := mergo.Merge(t, parent); err != nil {
			return nil, err
		}
	}

	// Check that the SDK version has been set to a sane value.
	switch t.SDKVersion {
	case "v1": // TODO: Should this be formalized somewhere?
	case "":
		t.SDKVersion = "v1"
	default:
		return nil, &UnknownSDKVersionError{Got: t.SDKVersion}
	}

	return t, nil
}

func NewTemplateFromReader(r io.Reader, opts ...Option) (*Template, error) {
	sctt, err := v1.NewStepContainerTemplateFromReader(r)
	if err != nil {
		return nil, err
	}

	return NewTemplateFromTyped(sctt.StepContainerCommon, opts...)
}

func NewTemplateFromFileRef(ref *FileRef) (t *Template, err error) {
	err = ref.WithFile(func(f http.File) (err error) {
		fi, err := f.Stat()
		if err != nil {
			return err
		} else if fi.IsDir() {
			t, err = NewTemplateFromFileRef(ref.Join(DefaultFileName))
		} else {
			t, err = NewTemplateFromReader(f, WithRelativeToFileRef(ref))
		}
		return
	})
	return
}
