package services

import (
	"encoding/json"
	models "lotto-backend-api/models"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type WithdrawCreditRequest struct {
	MemberId   string `json:"member_id"`
	Credit     string `json:"amount"`
	SecretSign string `json:"secret_sign"`
	OrderId    string `json:"order_id"`
}
type MemberCredit struct {
	Id           string `json:"id"`
	MemberCredit string `json:"credit_balance"`
}

func WithdrawCredit(DB *gorm.DB, r *http.Request) map[string]string {

	var withdrawCreditRequest WithdrawCreditRequest
	if err := json.NewDecoder(r.Body).Decode(&withdrawCreditRequest); err != nil {
		return map[string]string{
			"code":    "-1",
			"message": "withdraw-credit-service Invalid JSON format",
		}
	}

	// fmt.Println("req : " + withdrawCreditRequest.MemberId)
	// fmt.Println("req : " + withdrawCreditRequest.Credit)

	// var orders models.Orders
	// orders_result := DB.Where("member_id = ?", withdrawCreditRequest.MemberId).First(&orders)
	// if orders_result.Error != nil {
	// 	res := map[string]string{
	// 		"code":    "-1",
	// 		"message": "Not found data : " + orders_result.Error.Error(),
	// 	}
	// 	return res
	// }

	//currentTime := time.Now()
	//timestampSeconds := currentTime.Unix()

	//fmt.Println("DECREASE CREDIT")
	var member models.Members
	member_result := DB.Where("id = ?", withdrawCreditRequest.MemberId).First(&member)
	if member_result.Error != nil {
		res := map[string]string{
			"code":    "-1",
			"message": "Member Not Found",
		}
		return res
	}

	credit := MemberCredit{
		Id:           member.Id,
		MemberCredit: member.CreditBalance,
	}

	num1, _ := strconv.Atoi(member.CreditBalance)
	num2, _ := strconv.Atoi(withdrawCreditRequest.Credit)
	sum := num1 - num2

	// fmt.Println("Member ID : " + credit.Id)
	// fmt.Println("Member Credit : " + credit.MemberCredit)
	// fmt.Println("Withdraw Credit : " + withdrawCreditRequest.Credit)
	// fmt.Println("Total Credit : " + strconv.Itoa(sum))

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

	// //UPDATE ORDER
	// if err := DB.Model(&orders).
	// 	Where("id = ?", orders.Id).
	// 	Update("expire_time", timestampSeconds).Error; err != nil {
	// 	res := map[string]string{
	// 		"code":    "-1",
	// 		"message": "Update order expire time : " + err.Error(),
	// 	}
	// 	return res
	// }

	res := map[string]string{
		"code":    "0",
		"message": "success",
	}
	return res

}
