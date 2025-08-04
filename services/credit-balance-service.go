package services

import (
	"encoding/json"
	"lotto-backend-api/models"
	"net/http"

	"gorm.io/gorm"
)

// LoginUserRequest defines the structure for the login JSON payload
type CreditBalanceRequest struct {
	Id string `json:"member_id"`
}

// UserResponse defines the structure of the user data returned to the client.
// Notice it does NOT include the Password field, ensuring it's never exposed.
type CreditBalanceResponse struct {
	ID                string `json:"id"`
	Username          string `json:"username"`
	Identity          string `json:"identity"`
	BankProviderId    string `json:"bank_provider_id"`
	BankAccountNumber string `json:"bank_account_number"`
	BankAccountOwner  string `json:"bank_account_owner"`
	CreditBalance     string `json:"credit_balance"`
}

// LoginUser handles user authentication and returns sanitized user data on success.
func CreditBalanceService(DB *gorm.DB, r *http.Request) map[string]interface{} {
	defer r.Body.Close()
	var creditBalanceRequest CreditBalanceRequest

	if err := json.NewDecoder(r.Body).Decode(&creditBalanceRequest); err != nil {
		return map[string]interface{}{
			"code":    "-2",
			"message": "Invalid JSON format",
		}
	}

	var member models.Members // This is your full database model, including the password

	result := DB.Where("id = ?", creditBalanceRequest.Id).First(&member)

	if result.Error != nil {
		return map[string]interface{}{
			"code":    "-1",
			"message": "Invalid member_id",
		}
	}

	// --- CREATE THE SECURE RESPONSE OBJECT ---
	// Instead of returning the 'member' object directly, we map its values
	// to our new UserResponse struct, which safely omits the password.
	creditBalanceResponse := CreditBalanceResponse{
		ID:                member.Id,
		Username:          member.Username,
		Identity:          member.Identity,
		BankProviderId:    member.BankProviderId,
		BankAccountNumber: member.BankAccountNumber,
		BankAccountOwner:  member.BankAccountOwner,
		CreditBalance:     member.CreditBalance,
	}

	// --- SUCCESS RESPONSE ---
	// Now, the 'data' field contains only the safe, non-sensitive information.
	return map[string]interface{}{
		"code":    "200",
		"message": creditBalanceResponse, // <-- Return the sanitized responseData object
	}
}
