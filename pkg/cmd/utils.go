package cmd

import "strings"

func quoteShell(data string) string {
	return `'` + strings.ReplaceAll(data, `'`, `'"'"'`) + `'`
}
