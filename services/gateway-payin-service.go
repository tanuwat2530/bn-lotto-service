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

type PayinRequest struct {
	Amout   string `json:"amout"`
	Channel string `json:"chanel"`
}

// ApiPayload represents the structure of the JSON payload to be sent.
type PayInPayload struct {
	PaymentType string `json:"paymentType"`
	Merchant    string `json:"merchant"`
	Gold        string `json:"gold"`
	Channel     string `json:"channel"`
	OrderID     string `json:"orderId"`
	NotifyURL   string `json:"notify_url"`
	OrderAttach string `json:"order_attach"`
	FeeType     string `json:"feeType"`
	Sign        string `json:"sign"`
}

// ApiPayin constructs the request, generates a signature, and sends the request.
func ApiPayin(DB *gorm.DB, r *http.Request) string {

	secretKey := os.Getenv("SECRET_KEY")
	gatewayAccount := os.Getenv("GATEWAY_ACCOUNT")

	var payinRequest PayinRequest

	if err := json.NewDecoder(r.Body).Decode(&payinRequest); err != nil {
		fmt.Println("Invalid JSON format")
	}

	//orderId := uuid.New().String()
	notiUrl := "https://tonybet168.com/payin-noti"

	params := map[string]string{
		"merchant":    gatewayAccount,
		"paymentType": "1059",
		"gold":        payinRequest.Amout,
		"channel":     payinRequest.Channel,
		"notify_url":  notiUrl,
		"feeType":     "0",
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
	finalPayload := PayInPayload{
		Merchant:    gatewayAccount,
		PaymentType: "1059",
		Gold:        payinRequest.Amout,
		Channel:     payinRequest.Channel,
		NotifyURL:   notiUrl,
		FeeType:     "0",
		Sign:        signature,
	}

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

	return string(body)
}
