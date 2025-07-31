package models

type Payouts struct {
	Id          int64  `json:"id"`
	OrderId     string `json:"order_id"`
	RequestDate string `json:"request_date"`
	RequestTime string `json:"request_time"`
	GatewayData string `json:"gateway_data"`
	RequestData string `json:"request_data"`
	MemberId    string `json:"member_id"`
}
