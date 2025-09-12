package services

import (
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
}
type UserCredit struct {
	MemberId     string `json:"id"`
	MemberCredit string `json:"credit_balance"`
}

func PromtpayNoti(DB *gorm.DB, r *http.Request) map[string]string {

	fmt.Println("Param : ", r.URL.Query().Get("member_id"))
	fmt.Println("Param : ", r.URL.Query().Get("credit"))
	fmt.Println("Param : ", r.URL.Query().Get("secret_sign"))

	var member models.Members
	member_result := DB.Where("id = ?", r.URL.Query().Get("member_id")).First(&member)
	if member_result.Error != nil {
		res := map[string]string{
			"code":    "-1",
			"message": "Member Not Found",
		}
		return res
	}

	credit := UserCredit{
		MemberId:     member.Id,
		MemberCredit: member.CreditBalance,
	}
	fmt.Println("Member ID : " + credit.MemberId)
	fmt.Println("Member Credit : " + credit.MemberCredit)
	num1, _ := strconv.Atoi(member.CreditBalance)
	num2, _ := strconv.Atoi(r.URL.Query().Get("credit"))
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

	currentTime := time.Now()
	// Format as "YYYY-MM-DD"
	formattedDate := currentTime.Format("2006-01-02")
	formattedTime := currentTime.Format("15:04:05")
	orders := models.Orders{
		Date:      formattedDate,
		Time:      formattedTime,
		MemberId:  credit.MemberId,
		OrderData: r.URL.Query().Get("member_id") + "|" + r.URL.Query().Get("credit"),
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
