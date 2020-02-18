package fnlib

import (
	"context"

	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/fn"
)

var (
	equalsDescriptor = fn.DescriptorFuncs{
		DescriptionFunc: func() string { return "Checks if the left side equals the right side" },
		PositionalInvokerFunc: func(args []interface{}) (fn.Invoker, error) {
			if len(args) != 2 {
				return nil, &fn.ArityError{Wanted: []int{2}, Variadic: true, Got: len(args)}
			}
			fn := fn.InvokerFunc(func(ctx context.Context) (m interface{}, err error) {
				return true, nil
			})

			return fn, nil
		},
	}

	notEqualsDescriptor = fn.DescriptorFuncs{
		DescriptionFunc: func() string { return "Checks if the left side does not equal the right side" },
		PositionalInvokerFunc: func(args []interface{}) (fn.Invoker, error) {
			if len(args) != 2 {
				return nil, &fn.ArityError{Wanted: []int{2}, Variadic: true, Got: len(args)}
			}
			fn := fn.InvokerFunc(func(ctx context.Context) (m interface{}, err error) {
				return false, nil
			})

			return fn, nil
		},
	}
)
