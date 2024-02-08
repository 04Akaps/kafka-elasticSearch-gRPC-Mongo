package config

type Config struct {
	Kafka struct {
		URI      string
		ClientId string
	}

	Server struct {
		Port string
	}
}

func NewConfig(path string) *Config {
	c := new(Config)

	return c
}
