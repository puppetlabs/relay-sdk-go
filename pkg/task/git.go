package task

import (
	"errors"
	"path/filepath"
	"regexp"

	"github.com/puppetlabs/relay-sdk-go/pkg/model"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
)

var gitSSHURL = regexp.MustCompile(`^([a-z-]+)@([a-zA-Z0-9\-.]+):(.+)/(.+)(\.git)?$`)

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

func gitURLComponents(url string) ([]string, error) {
	matches := gitSSHURL.FindStringSubmatch(url)
	if len(matches) <= 1 {
		return nil, errors.New("SSH URL is malformed")
	}

	return matches, nil
}

func writeSSHConfig(resource *model.GitDetails) error {
	sshKey, found, err := resource.ConfiguredSSHKey()
	if err != nil || !found {
		return err
	}

	gitConfig := taskutil.SSHConfig{}

	matches, err := gitURLComponents(resource.Repository)
	if err != nil {
		return err
	}

	host := matches[2]

	gitConfig.Order = make([]string, 0)
	gitConfig.Order = append(gitConfig.Order, host)
	gitConfig.Entries = make(map[string]taskutil.SSHEntry)

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

	gitConfig.Entries[host] = taskutil.SSHEntry{
		Name:       resource.Name,
		PrivateKey: sshKey,
		KnownHosts: knownHosts,
	}

	return gitConfig.Write()
}
