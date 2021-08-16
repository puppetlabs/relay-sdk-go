package taskutil

import (
	"fmt"
	"net/url"
	"os"
	"time"
)

const (
	DefaultMetadataTimeout = 5 * time.Minute

	MetadataAPIURLEnvName = "METADATA_API_URL"
)

func MetadataURL(subpath string) (string, error) {
	if e := os.Getenv(MetadataAPIURLEnvName); e != "" {
		u, err := url.Parse(e)
		if err != nil {
			return "", fmt.Errorf("could not parse %s value as url: %s", MetadataAPIURLEnvName, e)
		}
		u = u.ResolveReference(&url.URL{Path: subpath})
		return u.String(), nil
	}

	return "", nil
}
