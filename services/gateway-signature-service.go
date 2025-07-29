package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"strings"
)

// GenerateSignature creates a request signature based on sorted parameters,
// a secret key, and the MD5 hashing algorithm. This function is a Go
// implementation of the provided PHP logic.
func GatewaySignature() string {
	//func GatewaySignature(params map[string]string, secretKey string) string {
	// 1. Retrieve and Sort Query Parameter Keys
	// Create a slice to hold the keys from the map.
	secretKey := os.Getenv("SECRET_KEY")
	params := map[string]string{
		"merchant":    "Merchant ID",
		"paymentType": "Channel",
		"gold":        "Amount", // Converted Amount to a string
		"channel":     "0",      // Converted 0 to a string
		"notify_url":  "Callback URL",
		"feeType":     "0", // Converted 0 to a string
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

	// 5. Return the Signature in Lowercase (Go's hex encoding is already lowercase)
	return signature
}

// func GatewaySignature() {
// 	// Sample Data (matching the PHP example)
// 	// Note: In Go, it's common to handle all param values as strings.
// 	params := map[string]string{
// 		"merchant":    "Merchant ID",
// 		"paymentType": "Channel",
// 		"gold":        "Amount", // Converted Amount to a string
// 		"channel":     "0",      // Converted 0 to a string
// 		"notify_url":  "Callback URL",
// 		"feeType":     "0", // Converted 0 to a string
// 	}

// 	secretKey := "Secret Key"

// 	// Generate Signature
// 	signature := GenerateSignature(params, secretKey)

// 	fmt.Println("Generated Signature:", signature)
// }
