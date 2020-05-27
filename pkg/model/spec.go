package model

import "encoding/base64"

type CredentialSpec struct {
	Credentials map[string]string
}

type GitSpec struct {
	GitRepository *GitDetails `spec:"git"`
}

type GitDetails struct {
	// Newer connection support for SSH keys.
	Connection *GitConnection

	// Older explicit Base64-encoded SSH keys.
	SSHKey string `spec:"ssh_key"`

	Name       string
	Repository string
	Branch     string
	KnownHosts string `spec:"known_hosts"`
}

func (gd *GitDetails) ConfiguredSSHKey() (string, bool, error) {
	if gd.Connection != nil {
		return gd.Connection.SSHKey, gd.Connection.SSHKey != "", nil
	}

	if gd.SSHKey == "" {
		return "", false, nil
	}

	sshKey, err := base64.StdEncoding.DecodeString(gd.SSHKey)
	if err != nil {
		return "", false, err
	}

	return string(sshKey), true, nil
}

type GitConnection struct {
	SSHKey string `spec:"sshKey"`
}

type ClusterSpec struct {
	Cluster *ClusterDetails
}

type ClusterDetails struct {
	Name     string
	URL      string
	CAData   string `spec:"cadata"`
	Token    string
	Insecure bool
	Username string
	Password string
}

type AWSSpec struct {
	AWS *AWSDetails
}

type AWSDetails struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}
