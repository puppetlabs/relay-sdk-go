package def

import (
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/puppetlabs/nebula-sdk/pkg/container/asset"
	v1 "github.com/puppetlabs/nebula-sdk/pkg/container/types/v1"
)

// Resolver allows the generator to load dependent resources.
type Resolver struct {
	// FileSystem is the filesystem implementation to use to load dependent
	// resources.
	FileSystem http.FileSystem

	// WorkingDirectory is the directory to use to resolve relative paths in the
	// step container data.
	WorkingDirectory string
}

func (r *Resolver) Open(name string) (http.File, error) {
	if !path.IsAbs(name) {
		name = path.Join(r.WorkingDirectory, name)
	}

	if r.FileSystem != nil {
		return r.FileSystem.Open(name)
	}

	return os.Open(filepath.FromSlash(name))
}

var (
	SDKResolver = &Resolver{
		FileSystem:       asset.FileSystem,
		WorkingDirectory: "templates",
	}
)

// FileRef points to a given file in a resolver.
type FileRef struct {
	resolver *Resolver
	name     string
}

func (fr *FileRef) Join(name string) *FileRef {
	return NewFileRef(path.Join(fr.name, name), WithFileRefResolver(fr.resolver))
}

// ResolverHere returns a resolver that can look up files relative to this
// FileRef.
func (fr *FileRef) ResolverHere() *Resolver {
	resolver := &Resolver{}
	if fr.resolver != nil {
		*resolver = *fr.resolver
	}

	resolver.WorkingDirectory = path.Join(resolver.WorkingDirectory, path.Dir(fr.name))

	return resolver
}

func (fr *FileRef) WithFile(fn func(f http.File) error) error {
	resolver := fr.resolver
	if resolver == nil {
		resolver = &Resolver{}
	}

	f, err := resolver.Open(fr.name)
	if err != nil {
		return err
	}
	defer f.Close()

	return fn(f)
}

type FileRefOption func(fr *FileRef)

func WithFileRefResolver(resolver *Resolver) FileRefOption {
	return func(fr *FileRef) {
		fr.resolver = resolver
	}
}

func NewFileRef(name string, opts ...FileRefOption) *FileRef {
	fr := &FileRef{name: name}
	for _, opt := range opts {
		opt(fr)
	}

	if fr.resolver == nil || fr.resolver.FileSystem == nil {
		fr.name = filepath.ToSlash(fr.name)
	}

	return fr
}

func NewFileRefFromTyped(ref v1.FileRef, opts ...FileRefOption) (*FileRef, error) {
	switch ref.From {
	case v1.FileSourceSystem:
		return NewFileRef(ref.Name, opts...), nil
	case v1.FileSourceSDK:
		// Clean name so users can't traverse outside of the template directory
		// and override the resolver with the SDK one.
		name := path.Clean(`.` + path.Clean(`/`+ref.Name))
		opts = append([]FileRefOption{}, opts...)
		opts = append(opts, WithFileRefResolver(SDKResolver))
		return NewFileRef(name, opts...), nil
	default:
		return nil, &UnknownFileSourceError{Got: ref.From}
	}
}
