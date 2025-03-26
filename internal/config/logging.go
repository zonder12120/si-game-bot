package config

type Logging struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}
