package services

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"net/http"
	"os"
	"sort"
	"strings"

	models "lotto-backend-api/models"

	"gorm.io/gorm"
)

type PayinRequest struct {
	MemberId    string `json:"member_id"`
	Amount      string `json:"amount"`
	Channel     string `json:"channel"`
	NotiUrl     string `json:"noti_url"`
	PaymentType string `json:"payment_type"`
	FeeType     string `json:"fee_type"`
}

// ApiPayload represents the structure of the JSON payload to be sent.
type PayInPayload struct {
	PaymentType string `json:"paymentType"`
	Merchant    string `json:"merchant"`
	Gold        string `json:"gold"`
	Channel     string `json:"channel"`

	//OrderID     string `json:"orderId"`
	NotifyURL string `json:"notify_url"`
	//OrderAttach string `json:"order_attach"`
	FeeType string `json:"feeType"`
	Sign    string `json:"sign"`
}

type DepositResponse struct {
	Code int8        `json:"code"`
	Data DepositData `json:"data"`
}
type DepositData struct {
	OrderNo string `json:"order_no"`
	PayURL  string `json:"payUrl"`
}

// ApiPayin constructs the request, generates a signature, and sends the request.
func ApiPayin(DB *gorm.DB, r *http.Request) string {
	defer r.Body.Close()
	secretKey := os.Getenv("SECRET_KEY")
	gatewayAccount := os.Getenv("GATEWAY_ACCOUNT")

	var payinRequest PayinRequest

	if err := json.NewDecoder(r.Body).Decode(&payinRequest); err != nil {
		fmt.Println("Invalid JSON format")
	}
	params := map[string]string{
		"merchant": gatewayAccount,
		//"orderId":     payinRequest.MemberId,
		"paymentType": payinRequest.PaymentType,
		"gold":        payinRequest.Amount,
		"channel":     payinRequest.Channel,
		"notify_url":  payinRequest.NotiUrl,
		"feeType":     payinRequest.FeeType,
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
	//fmt.Println(signature)

	// --- Step 3: Construct the final payload struct ---
	finalPayload := PayInPayload{
		Merchant: gatewayAccount,
		//OrderID:     payinRequest.MemberId,
		PaymentType: payinRequest.PaymentType,
		Gold:        payinRequest.Amount,
		Channel:     payinRequest.Channel,
		NotifyURL:   payinRequest.NotiUrl,
		FeeType:     payinRequest.FeeType,
		Sign:        signature,
	}
	//fmt.Println(finalPayload)
	// --- Step 4: Marshal the struct into a JSON byte slice ---
	payloadBytes, err := json.Marshal(finalPayload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)

	}
	// --- Step 5: Make the HTTP request ---
	url := "https://api.ghpay.vip/api/payIn" // Use a full URL
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
	fmt.Println("API : " + url)
	fmt.Printf("%+v\n", finalPayload)
	fmt.Println("Response : " + string(body))

	currentTime := time.Now()
	year := currentTime.Year()
	month := int(currentTime.Month())
	day := currentTime.Day()
	dateString := fmt.Sprintf("%d-%02d-%02d", year, month, day)

	hour := currentTime.Hour()
	minute := int(currentTime.Minute())
	second := currentTime.Second()
	timeString := fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)

	// Create a variable of the struct type
	var depositResponse DepositResponse

	// Unmarshal the JSON string into the struct
	err = json.Unmarshal([]byte(string(body)), &depositResponse)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
	}
	fmt.Println("Code:", depositResponse.Code)
	fmt.Println("Order Number:", depositResponse.Data.OrderNo)
	fmt.Println("Payment URL:", depositResponse.Data.PayURL)
	PayinData := models.Payins{
		//Id:          "",
		OrderId:     depositResponse.Data.OrderNo,
		RequestDate: dateString,
		RequestTime: timeString,
		GatewayData: string(body), //Response from Gateway
		RequestData: string(payloadBytes),
		MemberId:    payinRequest.MemberId,
	}
	result := DB.Create(&PayinData)
	if result.Error != nil {
		fmt.Println("PAYIN INSERT ERROR")
	}

	return string(body)
}
