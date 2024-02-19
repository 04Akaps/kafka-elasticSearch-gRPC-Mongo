package types

type LikeHistory struct {
	FromUser   string `json:"fromUser"`
	ToUser     string `json:"toUser"`
	Point      int64  `json:"point"`
	Action     string `json:"action"`
	SearchText string `json:"searchText"`
	Time       int64  `json:"time"`
}
