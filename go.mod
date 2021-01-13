module github.com/puppetlabs/relay-sdk-go

go 1.14

require (
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/magiconair/properties v1.8.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/puppetlabs/leg/encoding v0.1.0
	github.com/puppetlabs/leg/timeutil v0.2.0
	github.com/smartystreets/goconvey v0.0.0-20190731233626-505e41936337 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.6.1
	golang.org/x/oauth2 v0.0.0-20201208152858-08078c50e5b5 // indirect
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324 // indirect
	gopkg.in/ini.v1 v1.48.0
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/client-go v0.17.12
)

replace (
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.12
	k8s.io/client-go => k8s.io/client-go v0.17.12
)
