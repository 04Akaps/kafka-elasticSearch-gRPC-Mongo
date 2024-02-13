package config

type Config struct {
	Kafka struct {
		URI      string
		ClientId string
	}

	Elastic struct {
		URI      string
		User     string
		Password string
	}

	DB struct {
		URI string
		DB  string
	}

	Server struct {
		Port string
	}
}

func NewConfig(path string) *Config {
	c := new(Config)

	return c
}
