package network

import (
	"github.com/04Akaps/kafka-go/server/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type login struct {
	network *Network
}

func newLogin(n *Network) {
	l := &login{n}

	basePath := "/login"

	n.register(basePath, GET, l.login)
}

func (l *login) login(c *gin.Context) {
	var req types.LoginReq

	if err := c.ShouldBindQuery(&req); err != nil {
		response(c, http.StatusUnprocessableEntity, nil, err.Error())
	} else if res, err := l.network.auth.CreateAuth(req.Name); err != nil {
		response(c, http.StatusInternalServerError, nil, err.Error())
	} else {
		response(c, http.StatusOK, res, "")
	}
}
