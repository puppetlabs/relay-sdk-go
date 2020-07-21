package task

import (
	"os"
	"path/filepath"

	"github.com/puppetlabs/relay-sdk-go/pkg/model"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
)

func (ti *TaskInterface) ProcessGCP(directory string) error {
	var spec model.GCPSpec

	if err := taskutil.PopulateSpecFromDefaultPlan(&spec, ti.opts); err != nil {
		return err
	}

	if spec.Google == nil {
		return nil
	}

	if directory == "" {
		directory = DefaultPath
	}

	destination := filepath.Join(directory, "credentials.json")
	err := taskutil.WriteToFile(destination, spec.Google.ServiceAccountKey)

	if err != nil {
		return err
	}

	err = os.Chmod(destination, 0400)

	if err != nil {
		return err
	}

	return nil
}
