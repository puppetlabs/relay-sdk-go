package fnlib

import (
	"context"
	"fmt"
	"reflect"

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
				return isEqual(args[0], args[1])
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
				result, err := isEqual(args[0], args[1])

				return !result, err
			})

			return fn, nil
		},
	}
)

func isEqual(left, right interface{}) (bool, error) {
	leftt := reflect.TypeOf(left)
	rightt := reflect.TypeOf(right)

	if !leftt.Comparable() {
		return false, fmt.Errorf("%w: %v", fn.ErrUncomparableType, left)
	}

	if !rightt.Comparable() {
		return false, fmt.Errorf("%w: %v", fn.ErrUncomparableType, right)
	}

	return left == right, nil
}
