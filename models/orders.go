package models

type Orders struct {
	Id         int8   `json:"id"`
	Date       string `json:"date"`
	Time       string `json:"time"`
	MemberId   string `json:"member_id"`
	OrderData  string `json:"order_data"`
	ExpireTime int64  `json:"expire_time"`
}
