package models

type Members struct {
	Id                string `json:"id"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	Identity          string `json:"identity"`
	BankProviderId    string `json:"bank_provider_id"`
	BankAccountNumber string `json:"bank_account_number"`
	BankAccountOwner  string `json:"bank_account_owner"`
	CreditBalance     string `json:"credit_balance"`
}
