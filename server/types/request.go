package types

type LoginReq struct {
	Name string `form:"name" binding:"required"`
}

type LikeRequest struct {
	FromUser string `form:"fromUser" binding:"required"`
	ToUser   string `form:"toUser" binding:"required"`
	Point    int64  `form:"point" binding:"required"`
}

type UnLikeRequest struct {
	FromUser string `form:"fromUser" binding:"required"`
	ToUser   string `form:"toUser" binding:"required"`
	Point    int64  `form:"point" binding:"required"`
}

type LikeHistoryRequest struct {
	Paging
	Sort
	Search string `form:"search" binding:"required"`
}
