package network

import (
	"github.com/04Akaps/kafka-go/config"
	"github.com/04Akaps/kafka-go/server/repository/auth"
	"github.com/04Akaps/kafka-go/server/service"
	"github.com/gin-gonic/gin"
	"log"
)

type Network struct {
	config  *config.Config
	engine  *gin.Engine
	service service.ServiceImpl
	auth    *auth.RPCAuth

	port string
}

func NewNetwork(config *config.Config, service service.ServiceImpl, auth *auth.RPCAuth) *Network {
	n := &Network{
		config:  config,
		engine:  gin.New(),
		service: service,
		auth:    auth,
		port:    config.Server.Port,
	}

	n.setCors()

	newLogin(n)
	newLike(n, n.service)

	return n
}

func (n *Network) StartServer() error {
	log.Println("Start Server")
	return n.engine.Run(n.port)
}
