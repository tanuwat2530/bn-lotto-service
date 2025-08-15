package services

import (
	"encoding/json"
	"fmt"
	"io"
	models "lotto-backend-api/models"
	"net/http"
	"net/url"
	"time"

	"gorm.io/gorm"
)

type Member struct {
	MemberId string `json:"orderId"`
}

func OrderNoti(DB *gorm.DB, r *http.Request) map[string]string {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// Handle error
		res := map[string]string{
			"code":    "-2",
			"message": "order-noti-service can't read request #2",
		}
		return res
	}

	// Parse the URL-encoded data from the body
	parsedData, err := url.ParseQuery(string(body))
	if err != nil {
		res := map[string]string{
			"code":    "-3",
			"message": "order-noti-service can't read request #3",
		}
		return res
	}

	// Create a map to hold the key-value pairs
	// The map will be used to create the JSON object
	jsonData := make(map[string]string)
	for key, value := range parsedData {
		// url.ParseQuery returns a map[string][]string, so we take the first element
		jsonData[key] = value[0]
	}

	// Marshal the map into a JSON byte slice
	jsonOutput, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		// Handle error
		res := map[string]string{
			"code":    "-4",
			"message": "order-noti-service can't read request #4",
		}
		return res
	}
	fmt.Println(string(jsonOutput))

	// Create a variable of the struct type
	var member Member

	// Unmarshal the JSON string into the struct
	err = json.Unmarshal([]byte(string(jsonOutput)), &member)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
	}
	// Access the orderId field directly
	fmt.Println("Member ID :", member.MemberId)

	currentTime := time.Now()
	// Format as "YYYY-MM-DD"
	formattedDate := currentTime.Format("2006-01-02")
	formattedTime := currentTime.Format("15:04:05")
	orders := models.Orders{
		Date:      formattedDate,
		Time:      formattedTime,
		MemberId:  member.MemberId,
		OrderData: string(jsonOutput),
	}

	result := DB.Create(&orders)
	if result.Error != nil {
		res := map[string]string{
			"code":    "-1",
			"message": result.Error.Error(),
		}
		return res
	}

	res := map[string]string{
		"code":    "0",
		"message": "Success",
	}
	return res
}
