package services

import (
	"encoding/json"
	"fmt"
	models "lotto-backend-api/models"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type AddCreditRequest struct {
	MemberId   string `json:"member_id"`
	Credit     string `json:"credit"`
	SecretSign string `json:"secret_sign"`
	OrderId    string `json:"order_id"`
}
type UserCredit struct {
	Id           string `json:"id"`
	MemberCredit string `json:"credit_balance"`
}

func PromtpayNoti(DB *gorm.DB, r *http.Request) map[string]string {

	var addCreditRequest AddCreditRequest
	if err := json.NewDecoder(r.Body).Decode(&addCreditRequest); err != nil {
		return map[string]string{
			"code":    "-1",
			"message": "promtpay-noti-service Invalid JSON format",
		}
	}

	var orders models.Orders

	orders_result := DB.Where("member_id = ? AND order_data = ?", addCreditRequest.MemberId, addCreditRequest.OrderId).First(&orders)
	if orders_result.Error != nil {
		res := map[string]string{
			"code":    "-1",
			"message": "Not found data : " + orders_result.Error.Error(),
		}
		return res
	}

	currentTime := time.Now()
	timestampSeconds := currentTime.Unix()
	if timestampSeconds <= orders.ExpireTime {
		fmt.Println("ADD CREDIT")

		var member models.Members
		member_result := DB.Where("id = ?", addCreditRequest.MemberId).First(&member)
		if member_result.Error != nil {
			res := map[string]string{
				"code":    "-1",
				"message": "Member Not Found",
			}
			return res
		}

		credit := UserCredit{
			Id:           member.Id,
			MemberCredit: member.CreditBalance,
		}
		fmt.Println("Member ID : " + credit.Id)
		fmt.Println("Member Credit : " + credit.MemberCredit)
		num1, _ := strconv.Atoi(member.CreditBalance)
		num2, _ := strconv.Atoi(addCreditRequest.Credit)
		sum := num1 + num2
		fmt.Println("Current Credit : " + strconv.Itoa(sum))

		//UPDATE CREDIT
		if err := DB.Model(&member).
			Where("id = ?", credit.Id).
			Update("credit_balance", strconv.Itoa(sum)).Error; err != nil {
			res := map[string]string{
				"code":    "-1",
				"message": "Update member credit_balance : " + err.Error(),
			}
			return res
		}

		//UPDATE ORDER
		if err := DB.Model(&orders).
			Where("id = ?", orders.Id).
			Update("expire_time", timestampSeconds).Error; err != nil {
			res := map[string]string{
				"code":    "-1",
				"message": "Update order expire time : " + err.Error(),
			}
			return res
		}
	}

	res := map[string]string{
		"code":    "0",
		"message": "Success",
	}
	return res
}
