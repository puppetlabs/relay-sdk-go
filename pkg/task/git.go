package task

import (
	"path/filepath"

	"github.com/puppetlabs/relay-sdk-go/pkg/model"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
)

func (ti *TaskInterface) CloneRepository(revision, directory string) error {
	var spec model.GitSpec
	if err := taskutil.PopulateSpecFromDefaultPlan(&spec, ti.opts); err != nil {
		return err
	}

	resource := spec.GitRepository
	if resource == nil {
		return nil
	}

	if resource.Name == "" {
		resource.Name = DefaultName
	}

	if revision != "" {
		resource.Branch = revision
	}

	if resource.Branch == "" {
		resource.Branch = DefaultRevision
	}

	if directory == "" {
		directory = DefaultPath
	}

	if err := writeSSHConfig(resource); err != nil {
		return err
	}

	err := taskutil.Fetch(resource.Branch, filepath.Join(directory, resource.Name), resource.Repository)
	if err != nil {
		return err
	}

	return nil
}

func writeSSHConfig(resource *model.GitDetails) error {
	sshKey, found, err := resource.ConfiguredSSHKey()
	if err != nil || !found {
		return err
	}

	host, found, err := resource.ConfiguredRepository()
	if err != nil || !found {
		return err
	}

	knownHosts, found, err := resource.ConfiguredKnownHosts()
	if err != nil {
		return err
	}

	if !found {
		hostKeys, err := taskutil.SSHKeyScan(host)
		if err != nil {
			return err
		}

		knownHosts = string(hostKeys)
	}

	gitConfig := taskutil.SSHConfig{}

	gitConfig.Order = make([]string, 0)
	gitConfig.Order = append(gitConfig.Order, host)
	gitConfig.Entries = make(map[string]taskutil.SSHEntry)

	gitConfig.Entries[host] = taskutil.SSHEntry{
		Name:       resource.Name,
		PrivateKey: sshKey,
		KnownHosts: knownHosts,
	}

	return gitConfig.Write()
}
