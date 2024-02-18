package types

type LoginReq struct {
	Name string `form:"name" binding:"required"`
}

type LikeRequest struct {
	ToUser string `form:"toUser" binding:"required"`
	Point  int64  `form:"point" binding:"required"`
}

type UnLikeRequest struct {
	ToUser string `form:"toUser" binding:"required"`
	Point  int64  `form:"point" binding:"required"`
}
