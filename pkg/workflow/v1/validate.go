package v1

import "github.com/puppetlabs/nebula-sdk/pkg/util/typeutil"

// Validates a yaml document according to the schema specification
func Validate(y string) error {

	return typeutil.ValidateYAMLString(WorkflowSchema, y)
}
