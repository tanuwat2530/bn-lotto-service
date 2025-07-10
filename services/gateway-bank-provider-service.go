package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func ApiQueryDictionaryByCode() string {

	url := "https://api.ghpay.vip/api/QueryDictionaryByCode"
	method := "POST"

	// --- Generate a unique Order ID ---
	orderID := uuid.New().String()
	currentTime := time.Now()
	unixTimestampSeconds := strconv.FormatInt(currentTime.Unix(), 10)
	bankCode := "Kbank"

	secretKey := "Secret Key"
	params := map[string]string{
		"merchant":  "tonybet168",
		"noncestr":  orderID,
		"timestamp": unixTimestampSeconds,
		"bankCode":  bankCode,
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

	payload := strings.NewReader("merchant=tonybet168&noncestr=" + orderID + "&timestamp=" + unixTimestampSeconds + "&sign=" + signature + "&bankCode=" + bankCode)

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
