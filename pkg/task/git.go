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

	ssh, err := configuredSSH(resource)
	if err != nil {
		return err
	}

	if ssh != nil {
		if err := writeSSHConfig(resource.Name, ssh); err != nil {
			return err
		}
	}

	err = taskutil.Fetch(resource.Branch, filepath.Join(directory, resource.Name), resource.Repository)
	if err != nil {
		return err
	}

	return nil
}

func configuredSSH(gd *model.GitDetails) (*model.GitSSHDetails, error) {
	if gd.Repository == "" {
		return nil, nil
	}

	matches := model.GitSSHURL.FindStringSubmatch(gd.Repository)
	if len(matches) < 4 {
		return nil, nil
	}

	host := matches[2]

	sshKey, found, err := gd.ConfiguredSSHKey()
	if err != nil || !found {
		return nil, err
	}

	knownHosts, found, err := gd.ConfiguredKnownHosts()
	if err != nil {
		return nil, err
	}

	if !found {
		hostKeys, err := taskutil.SSHKeyScan(host)
		if err != nil {
			return nil, err
		}

		knownHosts = string(hostKeys)
	}

	return &model.GitSSHDetails{
		Host:       host,
		SSHKey:     sshKey,
		KnownHosts: knownHosts,
	}, nil
}

func writeSSHConfig(name string, ssh *model.GitSSHDetails) error {
	gitConfig := taskutil.SSHConfig{}

	gitConfig.Order = make([]string, 0)
	gitConfig.Order = append(gitConfig.Order, ssh.Host)
	gitConfig.Entries = make(map[string]taskutil.SSHEntry)

	gitConfig.Entries[ssh.Host] = taskutil.SSHEntry{
		Name:       name,
		PrivateKey: ssh.SSHKey,
		KnownHosts: ssh.KnownHosts,
	}

	return gitConfig.Write()
}
