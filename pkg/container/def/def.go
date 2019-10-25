package def

const (
	DefaultFileName = "container.yaml"
)

type Options struct {
	resolver *Resolver
}

type Option func(o *Options)

func WithResolver(resolver *Resolver) Option {
	return func(o *Options) {
		o.resolver = resolver
	}
}

func WithRelativeToFileRef(fr *FileRef) Option {
	return WithResolver(fr.ResolverHere())
}
