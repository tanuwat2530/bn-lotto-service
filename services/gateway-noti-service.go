package services

import (
	"encoding/json"
	"fmt"
	"io"
	models "lotto-backend-api/models"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type OrderNotiJson struct {
	OrderNo string `json:"order_no"`
	Gold    string `json:"gold"`
}
type Credit struct {
	MemberId     string `json:"id"`
	MemberCredit string `json:"credit_balance"`
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
	fmt.Println("Json Format : ", string(jsonOutput))

	// Create a variable of the struct type
	var orderNotiJson OrderNotiJson
	// Unmarshal the JSON string into the struct
	err = json.Unmarshal([]byte(string(jsonOutput)), &orderNotiJson)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
	}
	fmt.Println("order_no:", orderNotiJson.OrderNo)
	fmt.Println("gold:", orderNotiJson.Gold)

	var payIn models.Payins
	payin_result := DB.Where("order_id = ?", orderNotiJson.OrderNo).First(&payIn)
	if payin_result.Error != nil {
		res := map[string]string{
			"code":    "-1",
			"message": "Order Not Found",
		}
		return res
	}
	fmt.Println("Member ID :", payIn.MemberId)
	fmt.Println("Order ID : ", payIn.OrderId)
	fmt.Println("Request Data : ", payIn.RequestData)

	var member models.Members
	member_result := DB.Where("id = ?", payIn.MemberId).First(&member)
	if member_result.Error != nil {
		res := map[string]string{
			"code":    "-1",
			"message": "Member Not Found",
		}
		return res
	}

	credit := Credit{
		MemberId:     member.Id,
		MemberCredit: member.CreditBalance,
	}
	fmt.Println("Member ID : " + credit.MemberId)
	fmt.Println("Member Credit : " + credit.MemberCredit)
	num1, _ := strconv.Atoi(member.CreditBalance)
	num2, _ := strconv.Atoi(orderNotiJson.Gold)
	sum := num1 + num2
	fmt.Println("Current Credit : " + strconv.Itoa(sum))

	//UPDATE CREDIT
	if err := DB.Model(&member).
		Where("id = ?", credit.MemberId).
		Update("credit_balance", strconv.Itoa(sum)).Error; err != nil {
		res := map[string]string{
			"code":    "-1",
			"message": "Update member credit_balance : " + err.Error(),
		}
		return res
	}

	//UPDATE PAYIN
	if err := DB.Model(&payIn).
		Where("order_id = ?", payIn.OrderId).
		Update("order_id", payIn.OrderId+":DONE").Error; err != nil {
		res := map[string]string{
			"code":    "-1",
			"message": "Update payin order_id : " + err.Error(),
		}
		return res
	}

	currentTime := time.Now()
	// Format as "YYYY-MM-DD"
	formattedDate := currentTime.Format("2006-01-02")
	formattedTime := currentTime.Format("15:04:05")
	orders := models.Orders{
		Date:      formattedDate,
		Time:      formattedTime,
		MemberId:  payIn.MemberId,
		OrderData: string(jsonOutput),
	}

	orders_result := DB.Create(&orders)
	if orders_result.Error != nil {
		res := map[string]string{
			"code":    "-1",
			"message": "Insert orders error : " + orders_result.Error.Error(),
		}
		return res
	}

	res := map[string]string{
		"code":    "0",
		"message": "Success",
	}
	return res
}
