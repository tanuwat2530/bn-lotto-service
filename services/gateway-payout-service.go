package services

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"gorm.io/gorm"
)

type PayoutRequest struct {
	Amout           string `json:"amout"`
	Channel         string `json:"chanel"`
	TransferAccount string `"transferAccount`
	TransferName    string `transferName`
	TransferPhone   string `transferPhone`
}

// ApiPayload represents the structure of the JSON payload to be sent.
type PayOutPayload struct {
	PaymentType string `json:"paymentType"`
	Merchant    string `json:"merchant"`
	Gold        string `json:"gold"`
	Channel     string `json:"channel"` //0 = QR Code, 1 = Bank Transfer
	OrderID     string `json:"orderId"`
	NotifyURL   string `json:"notify_url"`
	//OrderAttach     string `json:"order_attach"`
	FeeType         string `json:"feeType"`
	TransferAccount string `json:"transferAccount"`
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	//BankCode        string `json:"bankCode"`
	//IdCard          string `json:"idCard"`
	Sign string `json:"sign"`
}

// ApiPayin constructs the request, generates a signature, and sends the request.
func ApiPayout(DB *gorm.DB, r *http.Request) string {
	secretKey := os.Getenv("SECRET_KEY")
	gatewayAccount := os.Getenv("GATEWAY_ACCOUNT")

	var payoutRequest PayoutRequest

	if err := json.NewDecoder(r.Body).Decode(&payoutRequest); err != nil {
		fmt.Println("Invalid JSON format")
	}

	params := map[string]string{
		"paymentType":     "2001",
		"merchant":        gatewayAccount,
		"gold":            payoutRequest.Amout,
		"channel":         payoutRequest.Channel,
		"notify_url":      "http://www.baidu.cim",
		"feeType":         "0",
		"transferAccount": payoutRequest.TransferAccount,
		"name":            payoutRequest.TransferName,
		"phone":           payoutRequest.TransferPhone,
	}

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys) // Sorts strings in ascending ASCII order

	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	// Join all parts with an ampersand.
	queryString := strings.Join(parts, "&")
	str := queryString + "&key=" + secretKey

	// 3. Perform an MD5 hash on the resulting string.
	hasher := md5.New()
	hasher.Write([]byte(str))
	hashInBytes := hasher.Sum(nil)
	signature := hex.EncodeToString(hashInBytes)
	fmt.Println(signature)

	// --- Step 3: Construct the final payload struct ---
	finalPayload := PayOutPayload{
		PaymentType:     "2001",
		Merchant:        gatewayAccount,
		Gold:            payoutRequest.Amout,
		Channel:         payoutRequest.Channel,
		NotifyURL:       "http://www.baidu.cim",
		FeeType:         "0",
		TransferAccount: payoutRequest.TransferAccount,
		Name:            payoutRequest.TransferName,
		Phone:           payoutRequest.TransferPhone,
		Sign:            signature,
	}

	// --- Step 4: Marshal the struct into a JSON byte slice ---
	payloadBytes, err := json.Marshal(finalPayload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)

	}

	// --- Step 5: Make the HTTP request ---
	url := "https://api.ghpay.vip/api/payOut"
	method := "POST"

	client := &http.Client{}
	// Create a reader from the JSON byte slice for the request body
	req, err := http.NewRequest(method, url, bytes.NewReader(payloadBytes))
	if err != nil {
		log.Printf("Error creating request: %v", err)

	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)

	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body) // Use io.ReadAll instead of deprecated ioutil
	if err != nil {
		log.Printf("Error reading response body: %v", err)

	}

	return string(body)
}
