package services

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"lotto-backend-api/models"
	"net/http"

	"gorm.io/gorm"
)

// LoginUserRequest defines the structure for the login JSON payload
type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserResponse defines the structure of the user data returned to the client.
// Notice it does NOT include the Password field, ensuring it's never exposed.
type UserResponse struct {
	ID                string `json:"id"`
	Username          string `json:"username"`
	Identity          string `json:"identity"`
	BankProviderId    string `json:"bank_provider_id"`
	BankAccountNumber string `json:"bank_account_number"`
	BankAccountOwner  string `json:"bank_account_owner"`
	CreditBalance     int8   `json:"credit_balance"`
}

// hashLoginPasswordSHA256 creates a consistent SHA-256 hash for a given password.
func hashLoginPasswordSHA256(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashedBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashedBytes)
}

// LoginUser handles user authentication and returns sanitized user data on success.
func LoginUser(DB *gorm.DB, r *http.Request) map[string]interface{} {
	defer r.Body.Close()
	var loginUserRequest LoginUserRequest

	if err := json.NewDecoder(r.Body).Decode(&loginUserRequest); err != nil {
		return map[string]interface{}{
			"code":    "-2",
			"message": "Invalid JSON format",
		}
	}

	hashedPassword := hashLoginPasswordSHA256(loginUserRequest.Password)
	var member models.Members // This is your full database model, including the password

	result := DB.Where("username = ? AND password = ?", loginUserRequest.Username, hashedPassword).First(&member)

	if result.Error != nil {
		return map[string]interface{}{
			"code":    "-1",
			"message": "Invalid username or password",
		}
	}

	// --- CREATE THE SECURE RESPONSE OBJECT ---
	// Instead of returning the 'member' object directly, we map its values
	// to our new UserResponse struct, which safely omits the password.
	responseData := UserResponse{
		ID:                member.Id,
		Username:          member.Username,
		Identity:          member.Identity,
		BankProviderId:    member.BankProviderId,
		BankAccountNumber: member.BankAccountNumber,
		BankAccountOwner:  member.BankAccountOwner,
		CreditBalance:     int8(member.CreditBalance),
	}

	// --- SUCCESS RESPONSE ---
	// Now, the 'data' field contains only the safe, non-sensitive information.
	return map[string]interface{}{
		"code":    "200",
		"message": responseData, // <-- Return the sanitized responseData object
	}
}
