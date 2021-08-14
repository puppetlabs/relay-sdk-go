package task

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/puppetlabs/leg/encoding/transfer"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
)

const (
	METADATA_ENVIRONMENT_SUBPATH = "environment"
)

func (ti *TaskInterface) ReadEnvironmentVariables() ([]byte, error) {
	tree, err := ti.readEnvironmentFromMetadata(METADATA_ENVIRONMENT_SUBPATH)
	if err != nil {
		return nil, err
	}

	return json.Marshal(transfer.JSONInterface{Data: tree})
}

func (ti *TaskInterface) ReadEnvironmentVariable(name string) ([]byte, error) {
	tree, err := ti.readEnvironmentFromMetadata(path.Join(METADATA_ENVIRONMENT_SUBPATH, name))
	if err != nil {
		return nil, err
	}

	// For individual variables, return the single value, not json
	return []byte(fmt.Sprintf("%v", tree)), nil
}

func (ti *TaskInterface) readEnvironmentFromMetadata(subpath string) (interface{}, error) {
	opts := ti.opts

	if ti.opts.SpecURL == "" {
		url, err := taskutil.MetadataURL(subpath)
		if err != nil {
			return nil, err
		}

		opts.SpecURL = url
	}

	tree, err := taskutil.LoadEnvironment(opts)
	if err != nil || tree == nil {
		return nil, err
	}

	return tree, nil
}
