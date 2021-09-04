package task

import (
	"github.com/puppetlabs/relay-sdk-go/pkg/model"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"
)

func (ti *TaskInterface) GetAzureSpec() (*model.AzureSpec, error) {
	var spec model.AzureSpec

	if err := taskutil.PopulateSpecFromDefaultPlan(&spec, ti.opts); err != nil {
		return nil, err
	}

	if spec.Azure == nil {
		return nil, nil
	}

	return &spec, nil
}

func (ti *TaskInterface) GetAzureARMEnvironmentVariables() (map[string]string, error) {
	spec, err := ti.GetAzureSpec()
	if err != nil {
		return nil, err
	}

	if spec.Azure == nil {
		return nil, nil
	}

	return map[string]string{
		"ARM_CLIENT_ID":       spec.Azure.Connection.ClientID,
		"ARM_CLIENT_SECRET":   spec.Azure.Connection.Secret,
		"ARM_SUBSCRIPTION_ID": spec.Azure.Connection.SubscriptionID,
		"ARM_TENANT_ID":       spec.Azure.Connection.TenantID,
	}, nil
}
