package auth

import (
	"github.com/04Akaps/go-util/auth"
	"github.com/04Akaps/kafka-go/config"
)

type RPCAuth struct {
	*auth.Auth
}

func NewRpcAuth(config *config.Config) (*RPCAuth, error) {
	a := new(RPCAuth)
	var err error

	if a.Auth, err = auth.NewAuth(config.Auth.ServerURL, config.Auth.ClientURL, config.Auth.PasetoKey); err != nil {
		return nil, err
	} else {
		return a, nil
	}
}
