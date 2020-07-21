package cmd

import "strings"

func quoteShell(data string) string {
	return `'` + strings.Replace(data, `'`, `'"'"'`, -1) + `'`
}
