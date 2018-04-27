package config

import (
	"time"

	"github.com/spf13/viper"
)

const (
	defaultHTTPClientTimeout = 15 * time.Second
	defaultFlagVerbose       = false

)

const (
	envHTTPClientTimeout = "DCOS_TIMEOUT"
	envSecret            = "DCOS_SECRET"
	envVerbose           = "DCOS_DEBUG"
)

// Config represents the runtime configuration of the system
type Config struct {
	viper *viper.Viper
}


func (c *Config) Secret() string {
	envReturned := c.viper.GetString(envSecret)
	if envReturned == "" {
		return "empty"
	}
	return c.viper.GetString(envSecret)
}

func (c *Config) EnvVerbose() bool {
	envReturned := c.viper.GetString(envVerbose)
	if envReturned == "" {
		return defaultFlagVerbose
	}
	return c.viper.GetBool(envVerbose)
}

// HTTPClientTimeout returns the configured http timeout
func (c *Config) HTTPClientTimeout() time.Duration {
	envReturned := c.viper.GetString(envHTTPClientTimeout)
	if envReturned == "" {
		return defaultHTTPClientTimeout
	}
	return c.viper.GetDuration(envHTTPClientTimeout)
}


func envs(viper *viper.Viper) {
	viper.BindEnv(envHTTPClientTimeout)
	viper.BindEnv(envSecret)
	viper.BindEnv(envVerbose)
}

// Parse returns an initialized Config after parsing program arguments
func Parse(args []string) (*Config, error) {
	viper := viper.New()
	envs(viper)

	return &Config{viper}, nil
}
