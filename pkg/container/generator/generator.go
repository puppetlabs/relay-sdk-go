package generator

import (
	"fmt"
	"io"
	"net/http"

	v1 "github.com/puppetlabs/nebula-sdk/pkg/container/types/v1"
)

type StepContainer struct {
}

type StepContainerOptions struct {
	// FileSystem is the filesystem implementation to use to load relative
	// paths.
	FileSystem http.FileSystem

	// FilePathContext is the directory to use to resolve relative paths in the
	// step container data.
	FilePathContext string
}

func NewStepContainerFromReader(r io.Reader, opts StepContainerOptions) (*StepContainer, error) {
	sct, err := v1.NewStepContainerFromReader(r)
	if err != nil {
		return nil, err
	}

	// We need to merge this file with each of its parents.
	for sct.Inherit != nil {
		switch sct.Inherit.From {
		case v1.FileSourceSystem:

		case v1.FileSourceSDK:

		default:
			panic(fmt.Errorf("unexpected file source %q", sct.Inherit.From))
		}
	}

	return nil, nil
}
