package main

import (
	"flag"
	"github.com/04Akaps/kafka-go/config"
	"github.com/04Akaps/kafka-go/docker"
	"github.com/04Akaps/kafka-go/server/app"
)

var confFlag = flag.String("config", "./config.toml", "toml file not found")

func main() {
	flag.Parsed()
	docker.DockerInit()
	cfg := config.NewConfig(*confFlag)
	app.NewApp(cfg)
}
