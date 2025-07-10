package services

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
)

type PayOutStatusPayload struct {
	Merchant string `json:"merchant"`
	OrderNo  string `json:"order_no"`
	Sign     string `json:"sign"`
}

func PayOutStatus() string {

	merchant := "tonybet168"
	orderId := "from-payout"

	secretKey := "Secret Key"
	params := map[string]string{
		"merchant": merchant,
		"orderId":  orderId,
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
	finalPayload := PayOutStatusPayload{
		Merchant: merchant,
		OrderNo:  orderId,
		Sign:     signature,
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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)

	}

	return string(body)
}
