package config

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

type Config struct {
	Kafka struct {
		URI      string
		ClientId string
		Topics   []string
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

	Auth struct {
		ClientURL string
		PasetoKey string
	}

	Server struct {
		Port string
	}
}

func NewConfig(path string) *Config {
	c := new(Config)

	if file, err := os.Open(path); err != nil {
		panic(err)
	} else {
		defer file.Close()
		if err = toml.NewDecoder(file).Decode(c); err != nil {
			panic(err)
		} else {
			return c
		}
	}
}
