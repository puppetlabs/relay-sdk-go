package fnlib

import (
	"context"
	"reflect"

	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/convert"
	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/fn"
)

var convertMarkdownDescriptor = fn.DescriptorFuncs{
	DescriptionFunc: func() string { return "Converts a string in markdown format to another applicable syntax" },
	PositionalInvokerFunc: func(args []interface{}) (fn.Invoker, error) {
		fn := fn.InvokerFunc(func(ctx context.Context) (m interface{}, err error) {
			if len(args) != 2 {
				return nil, &fn.ArityError{Wanted: []int{2}, Variadic: true, Got: len(args)}
			}

			convertType, ok := args[0].(string)
			if !ok {
				return nil, &fn.PositionalArgError{
					Arg: 0,
					Cause: &fn.UnexpectedTypeError{
						Wanted: []reflect.Type{reflect.TypeOf("")},
						Got:    reflect.TypeOf(args[0]),
					},
				}
			}

			switch arg := args[1].(type) {
			case string:
				r, err := convert.ConvertMarkdown(convert.ConvertType(convertType), []byte(arg))
				if err != nil {
					return nil, err
				}
				return string(r), nil
			default:
				return nil, &fn.PositionalArgError{
					Arg: 1,
					Cause: &fn.UnexpectedTypeError{
						Wanted: []reflect.Type{reflect.TypeOf("")},
						Got:    reflect.TypeOf(args[1]),
					},
				}
			}
		})
		return fn, nil
	},
}
