package services

import (
	// "crypto/md5" // MD5 is not secure for password hashing. It's still used for the ID here, but should be avoided for passwords.

	"encoding/json"
	"net/http"

	models "lotto-backend-api/models"

	// <-- IMPORT BCRYPT
	"gorm.io/gorm"
)

// JSON Request
type BetHistoryRequest struct {
	ReqAmount      string `json:"bet_amount"`
	ReqDescription string `json:"bet_description"`
	ReqNumber      string `json:"bet_number"`
	ReqPrize       string `json:"bet_prize"`
	ReqDate        string `json:"bet_date"`
	ReqTime        string `json:"bet_time"`
	ReqType        string `json:"bet_type"`
	ReqCurrent     string `json:"current_amount"`
	ReqMember      string `json:"member"`
}

func BetHistory(DB *gorm.DB, r *http.Request) map[string]string {
	var payload map[string]interface{}
	errPayload := json.NewDecoder(r.Body).Decode(&payload)
	if errPayload != nil {
		res := map[string]string{
			"code":    "-2",
			"message": "bet-history-service invalid Json #1",
		}
		return res
	}

	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		res := map[string]string{
			"code":    "-3",
			"message": "bet-history-service invalid Json #2",
		}
		return res
	}
	var betHistoryRequest BetHistoryRequest
	err = json.Unmarshal(jsonData, &betHistoryRequest)
	if err != nil {
		res := map[string]string{
			"code":    "-4",
			"message": "bet-history-service invalid Json #3",
		}
		return res
	}

	//UPDATE CREDIT
	Member := models.Members{}
	if err := DB.Model(&Member).
		Where("id = ?", betHistoryRequest.ReqMember).
		Update("credit_balance", betHistoryRequest.ReqCurrent).Error; err != nil {

		res := map[string]string{
			"code":    "-1",
			"message": err.Error(),
		}
		return res
	}

	//INSERT BET HISTORY

	History := models.Bets{
		BetDate:        betHistoryRequest.ReqDate,
		BetTime:        betHistoryRequest.ReqTime,
		BetNumber:      betHistoryRequest.ReqNumber,
		BetType:        betHistoryRequest.ReqType,
		BetDescription: betHistoryRequest.ReqDescription,
		BetPrize:       betHistoryRequest.ReqPrize,
		BetAmount:      betHistoryRequest.ReqAmount,
		MemberId:       betHistoryRequest.ReqMember,
	}

	result := DB.Create(&History)
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
