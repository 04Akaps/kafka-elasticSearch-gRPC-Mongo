package network

import (
	"github.com/04Akaps/kafka-go/server/service"
	"github.com/04Akaps/kafka-go/server/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
 간단하게 좋아요, 싫어요 요청을 하면, DB를 업데이트 하고,
 해당 이력을 ElasticSearch로 저장
 이뗴 업데이트 하는 이력을 이벤트 토픽으로 전송을 하여, kafka측에서 히스토리성 데이트럴
 ElasticSearch에 업데이트 하는 방식으로 동작시켜보는 API
*/

type like struct {
	network *Network
	service service.ServiceImpl
}

func newLike(n *Network, service service.ServiceImpl) {
	l := &like{n, service}

	basePath := "/like"

	n.register(basePath+"/like-request", GET, n.verifyLogin(), l.like)
	n.register(basePath+"/unLike-request", GET, n.verifyLogin(), l.unLike)
}

func (l *like) like(c *gin.Context) {
	var req types.LikeRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		response(c, http.StatusUnprocessableEntity, nil, err.Error())
	} else if err = l.service.Like(req.FromUser, req.ToUser, req.Point); err != nil {
		response(c, http.StatusInternalServerError, nil, err.Error())
	} else {
		response(c, http.StatusOK, nil, "")
	}

}

func (l *like) unLike(c *gin.Context) {
	var req types.UnLikeRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		response(c, http.StatusUnprocessableEntity, nil, err.Error())
	} else if err = l.service.UnLike(req.FromUser, req.ToUser, req.Point); err != nil {
		response(c, http.StatusInternalServerError, nil, err.Error())
	} else {
		response(c, http.StatusOK, nil, "")
	}

}
