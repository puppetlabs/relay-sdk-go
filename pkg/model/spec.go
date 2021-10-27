package model

import (
	"encoding/base64"
	"regexp"
)

var (
	GitSSHURL = regexp.MustCompile(`^([a-z-]+)@([a-zA-Z0-9\-.]+):(.+)/(.+)(\.git)?$`)
)

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

type GitSSHDetails struct {
	Host       string
	KnownHosts string
	SSHKey     string
}

func (gd *GitDetails) ConfiguredSSHKey() (string, bool, error) {
	if gd.Connection == nil {
		gd.Connection = &GitConnection{
			SSHKey: gd.SSHKey,
		}
	}

	if gd.Connection.SSHKey == "" {
		return "", false, nil
	}

	if sshKey, err := base64.StdEncoding.DecodeString(gd.Connection.SSHKey); err == nil {
		return string(sshKey), true, nil
	}

	return gd.Connection.SSHKey, true, nil
}

func (gd *GitDetails) ConfiguredKnownHosts() (string, bool, error) {
	if gd.KnownHosts == "" {
		return "", false, nil
	}

	if knownHosts, err := base64.StdEncoding.DecodeString(gd.KnownHosts); err == nil {
		return string(knownHosts), true, nil
	}

	return gd.KnownHosts, true, nil
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

type AzureSpec struct {
	Azure *AzureDetails
}

type AzureDetails struct {
	Connection AzureConnection
}

type AzureConnection struct {
	SubscriptionID string
	ClientID       string
	TenantID       string
	Secret         string
}

type AWSSpec struct {
	AWS *AWSDetails
}

type AWSConnection struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
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

func (ad *AWSDetails) GetSessionToken() string {
	return ad.Connection.SessionToken
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
