package taskutil

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Loader interface {
	Load() (io.Reader, error)
}

type RemoteLoader struct {
	u      *url.URL
	client *http.Client
}

func (r RemoteLoader) Load() (io.Reader, error) {
	var client = http.DefaultClient

	if r.client != nil {
		client = http.DefaultClient
	}

	//nolint:noctx // FIXME Consider replacing this later
	resp, err := client.Get(r.u.String())
	if err != nil {
		return nil, fmt.Errorf("network request failed: %+v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("the endpoint was not found")
	} else if resp.StatusCode == http.StatusInternalServerError {
		return nil, errors.New("an unexpected server error was encountered")
	} else if resp.StatusCode == http.StatusUnprocessableEntity {
		// TODO: Is this the correct behavior?
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	buf := &bytes.Buffer{}

	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, fmt.Errorf("reading response from remote service failed: %+v", err)
	}

	return buf, nil
}

func NewRemoteLoader(u *url.URL, client *http.Client) RemoteLoader {
	return RemoteLoader{
		u:      u,
		client: client,
	}
}

type LocalLoader struct {
	path string
}

func (l LocalLoader) Load() (io.Reader, error) {
	b, err := ioutil.ReadFile(l.path)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %+v", l.path, err)
	}

	return bytes.NewBuffer(b), nil
}

func NewLocalLoader(path string) LocalLoader {
	return LocalLoader{path: path}
}

type Decoder interface {
	Decode(io.Reader) (interface{}, error)
}
