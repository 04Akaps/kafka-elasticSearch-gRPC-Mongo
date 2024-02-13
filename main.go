package main

import (
	"flag"
	"github.com/04Akaps/kafka-go/config"
	"github.com/04Akaps/kafka-go/docker"
	"github.com/04Akaps/kafka-go/server/app"
)

var confFlag = flag.String("config", "./config.toml", "configuration toml file path")

func main() {
	flag.Parsed()
	cfg := config.NewConfig(*confFlag)
	docker.DockerInit()
	app.NewApp(cfg)
}
