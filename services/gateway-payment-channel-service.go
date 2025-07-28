package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ApiMerchantAisle() string {

	//url := "https://api.ghpay.vip/api/MerchantAisle"
	//method := "POST"
	//secKey := "696b98401cca4195b4e76d80ab58ecca"
	currentTime := time.Now()
	timeStamp := strconv.FormatInt(currentTime.Unix(), 10)
	params := map[string]string{
		"merchant":  "tonybet168",
		"noncestr":  "tonybet168tonybet168",
		"timestamp": timeStamp,
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
	str := queryString + "&key=696b98401cca4195b4e76d80ab58ecca"
	// 3. Perform an MD5 hash on the resulting string.
	hasher := md5.New()
	hasher.Write([]byte(str))
	hashInBytes := hasher.Sum(nil)
	signature := hex.EncodeToString(hashInBytes)
	fmt.Println(signature)

	// Define the API endpoint URL
	apiURL := "https://api.ghpay.vip/api/MerchantAisle"

	// Prepare the form data as url.Values
	// This automatically handles URL encoding for the key-value pairs.
	formData := url.Values{}
	formData.Set("merchant", "tonybet168")
	formData.Set("noncestr", "tonybet168tonybet168")
	formData.Set("timestamp", timeStamp) // Note: In a real scenario, this would be generated dynamically
	formData.Set("sign", signature)      // Note: In a real scenario, this would be generated dynamically

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new POST request
	// The body is a string reader of the URL-encoded form data
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return "ERROR #1"
	}

	// Set the Content-Type header to application/x-www-form-urlencoded
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return "ERROR #2"
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return "ERROR #3"
	}

	// Print the response status and body
	//fmt.Printf("Response Status: %s\n", resp.Status)
	//fmt.Printf("Response Body: %s\n", string(body))

	return string(body)
}
