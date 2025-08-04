package models

type Bets struct {
	Id             int8   `json:"id"`
	BetDate        string `json:"bet_date"`
	BetTime        string `json:"bet_time"`
	BetNumber      string `json:"bet_number"`
	BetType        string `json:"bet_type"`
	BetDescription string `json:"bet_description"`
	BetPrize       string `json:"bet_prize"`
	MemberId       string `json:"member_id"`
	BetAmount      string `json:"bet_amount"`
}
