package config

import "os"

const (
	FlagOutputPath = "path"
	JsonConfigFileName = "servers.json"
	FlagJsonConfigPath = "config_file"
)

type Sentry struct {}

// Env for sentry env
func (Sentry) Env() string {
	return os.Getenv("SENTRY_ENVIRONMENT")
}

// Dsn for sentry
func (Sentry) Dsn() string {
	return os.Getenv("SENTRY_DSN")
}
