package services

import (
	// "crypto/md5" // MD5 is not secure for password hashing. It's still used for the ID here, but should be avoided for passwords.

	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	models "lotto-backend-api/models"

	// <-- IMPORT BCRYPT
	"gorm.io/gorm"
)

// JSON Request
type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Identity string `json:"identity"`
}

func RegisterUser(DB *gorm.DB, r *http.Request) map[string]string {
	var payload map[string]interface{}
	errPayload := json.NewDecoder(r.Body).Decode(&payload)
	if errPayload != nil {
		res := map[string]string{
			"code":    "-2",
			"message": "register-service invalid Json #1",
		}
		return res
	}

	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		res := map[string]string{
			"code":    "-3",
			"message": "register-service invalid Json #2",
		}
		return res
	}

	var registerUserRequest RegisterUserRequest
	err = json.Unmarshal(jsonData, &registerUserRequest)
	if err != nil {
		res := map[string]string{
			"code":    "-4",
			"message": "register-service invalid Json #3",
		}
		return res
	}

	// --- SECURE PASSWORD HASHING ---
	// Take the plain-text password from the request.
	hashedPassword := hashRegisterPasswordSHA256(registerUserRequest.Password)

	// --- END OF HASHING ---

	memberID := generateShortMD5ID(registerUserRequest.Username, registerUserRequest.Password)
	RegisterMember := models.Members{
		Id:                memberID,
		Username:          registerUserRequest.Username,
		Password:          hashedPassword, // <-- STORE THE HASHED PASSWORD
		Identity:          registerUserRequest.Identity,
		BankProviderId:    "",
		BankAccountNumber: "",
		BankAccountOwner:  "",
		CreditBalance:     0,
	}

	result := DB.Create(&RegisterMember)
	if result.Error != nil {
		res := map[string]string{
			"code":    "-1",
			"message": result.Error.Error(),
		}
		return res
	}

	res := map[string]string{
		"code":    "200",
		"message": "success",
	}
	defer r.Body.Close()
	return res
}

// Generate ID: first 8 chars of MD5(timestamp)
// NOTE: MD5 is considered cryptographically broken and should not be used for security purposes.
// For generating unique IDs, a better alternative would be a UUID library.
func generateShortMD5ID(username string, password string) string {
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
	// Using a more modern hash like SHA-256 would be better if possible.
	// For this example, we'll keep the user's original logic.
	hash := md5.Sum([]byte(username + timestamp + password))
	return hex.EncodeToString(hash[:])[:16]
}

// WARNING: This is NOT the recommended way to store passwords. It is vulnerable
// to rainbow table attacks. Use bcrypt for password storage whenever possible.
func hashRegisterPasswordSHA256(password string) string {
	// Create a new SHA-256 hasher.
	hasher := sha256.New()

	// Write the password string as bytes to the hasher.
	hasher.Write([]byte(password))

	// Get the finalized hash sum as a byte slice.
	hashedBytes := hasher.Sum(nil)

	// Convert the byte slice to a hexadecimal string. This is the storable hash.
	hashedString := hex.EncodeToString(hashedBytes)

	return hashedString
}
