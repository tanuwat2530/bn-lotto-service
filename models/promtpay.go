package models

type Promtpays struct {
	QrId          string `json:"qr_id"`
	QrName        string `json:"qr_name"`
	TotalDeposit  int64  `json:"total_deposit"`
	TotalWithdraw int64  `json:"total_withdraw"`
	Counter       int64  `json:"counter"`
}
