package def

import (
	"io"
	"net/http"

	v1 "github.com/puppetlabs/nebula-sdk/pkg/container/types/v1"
)

type Container struct {
	Template    *Template
	Title       string
	Description string
}

func NewFromTyped(sct *v1.StepContainer, opts ...Option) (*Container, error) {
	t, err := NewTemplateFromTyped(sct.StepContainerCommon, opts...)
	if err != nil {
		return nil, err
	}

	// At this point, all settings must have values.
	for name, setting := range t.Settings {
		if setting.Value == nil {
			return nil, &MissingSettingValueError{Name: name}
		}
	}

	c := &Container{
		Template:    t,
		Title:       sct.Title,
		Description: sct.Description,
	}
	return c, nil
}

func NewFromReader(r io.Reader, opts ...Option) (*Container, error) {
	sct, err := v1.NewStepContainerFromReader(r)
	if err != nil {
		return nil, err
	}

	return NewFromTyped(sct, opts...)
}

func NewFromFileRef(ref *FileRef) (c *Container, err error) {
	err = ref.WithFile(func(f http.File) (err error) {
		fi, err := f.Stat()
		if err != nil {
			return err
		} else if fi.IsDir() {
			c, err = NewFromFileRef(ref.Join(DefaultFileName))
		} else {
			c, err = NewFromReader(f, WithRelativeToFileRef(ref))
		}
		return
	})
	return
}

func NewFromFilePath(name string) (*Container, error) {
	return NewFromFileRef(NewFileRef(name))
}
