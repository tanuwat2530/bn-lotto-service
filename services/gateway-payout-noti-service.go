package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func PayOutNotification() string {

	url := "https://api.ghpay.vip/api/payOut"
	method := "POST"

	payload := strings.NewReader("merchant=&order_no=&orderId=&gold=&pay_amount=&trade_amount=&orderStatus=&statusMsg=&paymentType=&fee=&%20completeTime=&sign=&order_attach=")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)

	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)

	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)

	}
	return string(body)
}
