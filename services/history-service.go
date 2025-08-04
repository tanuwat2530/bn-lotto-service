package services

import (
	// "crypto/md5" // MD5 is not secure for password hashing. It's still used for the ID here, but should be avoided for passwords.

	"encoding/json"
	"fmt"
	models "lotto-backend-api/models"
	"net/http"

	// <-- IMPORT BCRYPT
	"gorm.io/gorm"
)

// JSON Request
type HistoryRequest struct {
	Member string `json:"member_id"`
}

type HistoryResponse struct {
	BetDate        string `json:"bet_date"`
	BetTime        string `json:"bet_time"`
	BetNumber      string `json:"bet_number"`
	BetType        string `json:"bet_type"`
	BetDescription string `json:"bet_description"`
	BetPrize       string `json:"bet_prize"`
	BetAmount      string `json:"bet_amount"`
	MemberId       string `json:"member_id"`
}

func History(DB *gorm.DB, r *http.Request) map[string]string {

	var payload map[string]interface{}
	errPayload := json.NewDecoder(r.Body).Decode(&payload)
	if errPayload != nil {
		res := map[string]string{
			"code":    "-2",
			"message": "history-service invalid Json #1",
		}
		return res
	}

	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		res := map[string]string{
			"code":    "-3",
			"message": "history-service invalid Json #2",
		}
		return res
	}

	var historyReq HistoryRequest
	err = json.Unmarshal(jsonData, &historyReq)
	if err != nil {
		res := map[string]string{
			"code":    "-4",
			"message": "history-service invalid Json #3",
		}
		return res
	}

	// ✅ Use a slice to hold multiple records
	var betHistory []models.Bets
	// ✅ Use .Where before .Find
	result := DB.Where("member_id = ?", historyReq.Member).Find(&betHistory)

	if result.Error != nil {
		res := map[string]string{
			"code":    "-5",
			"message": result.Error.Error(),
		}
		return res
	}

	var responseDataList []HistoryResponse

	for _, bet := range betHistory {
		responseData := HistoryResponse{
			BetDate:        bet.BetDate,
			BetTime:        bet.BetTime,
			BetNumber:      bet.BetNumber,
			BetDescription: bet.BetDescription,
			BetType:        bet.BetType,
			BetPrize:       bet.BetPrize,
			BetAmount:      bet.BetAmount,
			MemberId:       bet.MemberId,
		}
		responseDataList = append(responseDataList, responseData)
	}

	// ✅ Now marshal to JSON
	jsonBytes, err := json.Marshal(responseDataList)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		res := map[string]string{
			"code":    "0",
			"message": "Error marshaling to JSON : " + err.Error(),
		}
		return res
	}

	jsonString := string(jsonBytes)
	defer r.Body.Close()
	res := map[string]string{
		"code":    "0",
		"message": jsonString,
	}
	return res

}
