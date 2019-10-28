package def

import (
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	v1 "github.com/puppetlabs/nebula-sdk/pkg/container/types/v1"
)

type Container struct {
	*Common

	Name        string
	Title       string
	Description string
}

type ResolvedContainer struct {
	FileRef   *FileRef
	Container *Container
}

func NewFromTyped(sct *v1.StepContainer, opts ...CommonOption) (*Container, error) {
	co, err := NewCommonFromTyped(sct.StepContainerCommon, opts...)
	if err != nil {
		return nil, err
	}

	// At this point, all settings must have values.
	for name, setting := range co.Settings {
		if setting.Value == nil {
			return nil, &MissingSettingValueError{Name: name}
		}
	}

	name := sct.Name
	if name == "" {
		dir := slug(path.Base(co.Resolver.WorkingDirectory))
		if dir == "" {
			return nil, ErrMissingName
		}

		name = dir
	}

	c := &Container{
		Common:      co,
		Name:        name,
		Title:       sct.Title,
		Description: sct.Description,
	}
	return c, nil
}

func NewFromReader(r io.Reader, opts ...CommonOption) (*Container, error) {
	sct, err := v1.NewStepContainerFromReader(r)
	if err != nil {
		return nil, err
	}

	return NewFromTyped(sct, opts...)
}

func NewFromFileRef(ref *FileRef) (rc *ResolvedContainer, err error) {
	err = ref.WithFile(func(f http.File) (err error) {
		fi, err := f.Stat()
		if err != nil {
			return err
		} else if fi.IsDir() {
			rc, err = NewFromFileRef(ref.Join(DefaultFilename))
		} else {
			c, err := NewFromReader(f, WithResolver(ref.ResolverHere()))
			if err != nil {
				return err
			}

			rc = &ResolvedContainer{
				FileRef:   ref,
				Container: c,
			}
		}
		return
	})
	return
}

func NewFromFilePath(name string) (*ResolvedContainer, error) {
	if !filepath.IsAbs(name) {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		name = filepath.Clean(filepath.Join(wd, name))
	}

	return NewFromFileRef(NewFileRef(name))
}

var slugReplacer = regexp.MustCompile(`[^A-Za-z0-9_-]+`)

func slug(in string) string {
	return strings.Trim(slugReplacer.ReplaceAllLiteralString(in, "-"), "-")
}
