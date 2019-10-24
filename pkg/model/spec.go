package model

type CredentialSpec struct {
	Credentials map[string]string
}

type GitSpec struct {
	GitRepository *GitDetails `spec:"git"`
}

type GitDetails struct {
	Name       string
	Repository string
	Branch     string
	SSHKey     string `spec:"ssh_key"`
	KnownHosts string `spec:"known_hosts"`
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
