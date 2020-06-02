package fnlib_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/convert"
	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/fnlib"
	"github.com/stretchr/testify/require"
)

func TestConvertMarkdown(t *testing.T) {
	desc, err := fnlib.Library().Descriptor("convertMarkdown")
	require.NoError(t, err)

	tcs := []struct {
		Name     string
		Markdown string
		Expected string
	}{
		{
			Name:     "Sample monitor event",
			Markdown: "%%% @contact [![imageTitle](imageUrl)](imageRedirect) `data{context} > threshold` Detailed description. - - - [[linkTitle1](link1)] · [[linkTitle2](link2)] %%%",
			Expected: "@contact \n\n[!imageUrl!|imageRedirect] {code}data{context} > threshold{code} Detailed description.\n----\n[[linkTitle1|link1]] · [[linkTitle2|link2]]",
		},
	}

	for _, test := range tcs {
		t.Run(fmt.Sprintf("%s", test.Name), func(t *testing.T) {
			invoker, err := desc.PositionalInvoker([]interface{}{convert.ConvertTypeJira.String(), test.Markdown})
			require.NoError(t, err)

			r, err := invoker.Invoke(context.Background())
			require.NoError(t, err)

			require.Equal(t, test.Expected, r)
		})
	}
}
