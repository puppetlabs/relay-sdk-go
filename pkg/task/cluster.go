package task

import (
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/imdario/mergo"
	"github.com/puppetlabs/relay-sdk-go/pkg/model"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func (ti *TaskInterface) ProcessClusters(directory string) error {

	var spec model.ClusterSpec
	if err := taskutil.PopulateSpecFromDefaultPlan(&spec, ti.opts); err != nil {
		return err
	}

	cluster := retrieveClusterFromEnvironment()

	if cluster != nil && spec.Cluster != nil {
		if err := mergo.Merge(cluster, spec.Cluster); err != nil {
			return err
		}
	}
	if cluster == nil {
		return errors.New("unable to find cluster credentials")
	}

	if cluster.Name == "" {
		cluster.Name = DefaultName
	}

	// Prefer non-deprecated values over deprecated values
	if cluster.Connection == nil {
		cluster.Connection = &model.ClusterConnectionSpec{}
	}
	if cluster.Connection.Server == "" {
		cluster.Connection.Server = cluster.URL
		cluster.URL = ""
	}
	if cluster.Connection.CertificateAuthority == "" {
		cluster.Connection.CertificateAuthority = cluster.CAData
		cluster.CAData = ""
	}
	if cluster.Connection.Token == "" {
		cluster.Connection.Token = cluster.Token
		cluster.Token = ""
	}

	return CreateKubeconfigFile(directory, cluster)
}

func CreateKubeconfigFile(directory string, resource *model.ClusterDetails) error {
	if resource == nil {
		return nil
	}

	cluster := &clientcmdapi.Cluster{
		Server:                resource.Connection.Server,
		InsecureSkipTLSVerify: resource.Insecure,
	}

	if resource.Connection.CertificateAuthority != "" {
		ca, err := base64.StdEncoding.DecodeString(resource.Connection.CertificateAuthority)
		if err != nil {
			return err
		}
		cluster.CertificateAuthorityData = ca
	}

	//only one authentication technique per user is allowed in a kubeconfig, so clear out the password if a token is provided
	user := resource.Username
	pass := resource.Password
	if resource.Connection.Token != "" {
		user = ""
		pass = ""
	}
	auth := &clientcmdapi.AuthInfo{
		Token:    resource.Connection.Token,
		Username: user,
		Password: pass,
	}
	context := &clientcmdapi.Context{
		Cluster:  resource.Name,
		AuthInfo: resource.Username,
	}
	c := clientcmdapi.NewConfig()
	c.Clusters[resource.Name] = cluster
	c.AuthInfos[resource.Username] = auth
	c.Contexts[resource.Name] = context
	c.CurrentContext = resource.Name
	c.APIVersion = "v1"
	c.Kind = "Config"

	if directory == "" {
		directory = DefaultPath
	}

	destination := filepath.Join(directory, resource.Name, KubeConfigFile)
	return clientcmd.WriteToFile(*c, destination)
}

func retrieveClusterFromEnvironment() *model.ClusterDetails {

	cluster := &model.ClusterDetails{}

	if nameFromEnv := os.Getenv("CLUSTER_NAME"); nameFromEnv != "" {
		cluster.Name = nameFromEnv
	}
	if urlFromEnv := os.Getenv("CLUSTER_URL"); urlFromEnv != "" {
		cluster.Connection.Server = urlFromEnv
	}
	if caFromEnv := os.Getenv("CLUSTER_CADATA"); caFromEnv != "" {
		cluster.Connection.CertificateAuthority = caFromEnv
	}
	if tokenFromEnv := os.Getenv("CLUSTER_TOKEN"); tokenFromEnv != "" {
		cluster.Connection.Token = strings.TrimRight(tokenFromEnv, "\r\n")
	}
	if usernameFromEnv := os.Getenv("CLUSTER_USERNAME"); usernameFromEnv != "" {
		cluster.Username = usernameFromEnv
	}
	if passwordFromEnv := os.Getenv("CLUSTER_PASSWORD"); passwordFromEnv != "" {
		cluster.Password = passwordFromEnv
	}

	return cluster
}
