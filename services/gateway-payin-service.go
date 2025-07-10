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
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// ApiPayload represents the structure of the JSON payload to be sent.
type PayInPayload struct {
	PaymentType string `json:"paymentType"`
	Merchant    string `json:"merchant"`
	Gold        int    `json:"gold"`
	Channel     int    `json:"channel"`
	OrderID     string `json:"orderId"`
	NotifyURL   string `json:"notify_url"`
	//OrderAttach string `json:"order_attach"`
	FeeType int    `json:"feeType"`
	Sign    string `json:"sign"`
}

// ApiPayin constructs the request, generates a signature, and sends the request.
func ApiPayin() string {

	paymentType := "1001"
	merchant := "tonybet168"
	amount := 50
	channel := 0
	orderId := uuid.New().String()
	notiUrl := "https://tonybet168.com/payin-noti"

	secretKey := "Secret Key"
	params := map[string]string{
		"paymentType": paymentType,
		"merchant":    merchant,
		"gold":        strconv.FormatInt(int64(amount), 10),
		"channel":     "0",
		"orderId":     orderId,
		"notify_url":  notiUrl,
		"feeType":     "0",
	}
	sortedKeys := make([]string, 0, len(params))
	for k := range params {
		sortedKeys = append(sortedKeys, k)
	}
	// Sort the keys alphabetically.
	sort.Strings(sortedKeys)

	// 2. Create Query String
	// Create a slice to hold the "key=value" parts.
	var queryStringParts []string
	for _, key := range sortedKeys {
		// Build the "key=value" string and add it to the slice.
		queryStringParts = append(queryStringParts, fmt.Sprintf("%s=%s", key, params[key]))
	}
	// Join the parts with an ampersand (&).
	queryString := strings.Join(queryStringParts, "&")

	// 3. Concatenate Key and Secret Key
	// Append the secret key to the query string.
	stringToSign := queryString + "&key=" + secretKey

	// 4. Calculate the MD5 Hash
	hasher := md5.New()
	hasher.Write([]byte(stringToSign))
	hashBytes := hasher.Sum(nil)

	// Convert the hash to a lowercase hexadecimal string.
	signature := hex.EncodeToString(hashBytes)

	// --- Step 3: Construct the final payload struct ---
	finalPayload := PayInPayload{
		PaymentType: paymentType,
		Merchant:    merchant,
		Gold:        amount,
		Channel:     channel,
		OrderID:     orderId,
		NotifyURL:   notiUrl,
		FeeType:     0,
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
