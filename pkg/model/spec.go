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
	Name       string
	Connection *ClusterConnectionSpec

	// Deprecated
	URL string
	// Deprecated
	CAData string `spec:"cadata"`
	// Deprecated
	Token    string
	Insecure bool
	Username string
	Password string
}

type ClusterConnectionSpec struct {
	Server               string
	CertificateAuthority string
	Token                string
}

type AWSSpec struct {
	AWS *AWSDetails
}

type AWSConnection struct {
	AccessKeyID     string
	SecretAccessKey string
}

type AWSDetails struct {
	Connection AWSConnection

	// deprecated
	AccessKeyID string
	// deprecated
	SecretAccessKey string
	Region          string
}

func (ad *AWSDetails) GetAccessKeyID() string {
	if ad.Connection.AccessKeyID == "" {
		return ad.AccessKeyID

	}
	return ad.Connection.AccessKeyID
}

func (ad *AWSDetails) GetSecretAccessKey() string {
	if ad.Connection.SecretAccessKey == "" {
		return ad.SecretAccessKey

	}
	return ad.Connection.SecretAccessKey
}

type GCPSpec struct {
	Google *GCPDetails
}

type GCPDetails struct {
	Connection GCPConnection
	Project    string

	// deprecated
	ServiceAccountKey string `json:"serviceAccountKey"`
}

func (gd *GCPDetails) GetServiceAccountKey() string {
	if gd.Connection.ServiceAccountKey == "" {
		return gd.ServiceAccountKey
	}
	return gd.Connection.ServiceAccountKey
}

type GCPConnection struct {
	ServiceAccountKey string `json:"serviceAccountKey"`
}
