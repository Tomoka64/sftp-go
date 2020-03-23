package config

import "github.com/kelseyhightower/envconfig"

type (
	configuration struct {
		Proxy        server `envconfig:"PROXY"`
		TargetServer server `envconfig:"TARGET_SERVER"`
	}

	server struct {
		PrivateKey string `envconfig:"PRIVATE_KEY" default:""`
		PassPhrase string `envconfig:"PASS_PHRASE"`
		User       string `envconfig:"USER" default:"test"`
		Addr       string `envconfig:"ADDR"`
	}
)

const (
	prefix = "SFTP"
)

var conf = &configuration{}

func init() {
	envconfig.MustProcess(prefix, conf)
}

func Reload() error {
	return envconfig.Process(prefix, conf)
}

func Proxy() server {
	return conf.Proxy
}

func TargetServer() server {
	return conf.TargetServer
}
