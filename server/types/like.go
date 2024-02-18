package types

type LikeHistory struct {
	FromUser string `json:"fromUser"`
	ToUser   string `json:"toUser"`
	Point    int64  `json:"point"`
	Action   string `json:"action"`
}
