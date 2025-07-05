package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
)

func Signature(arr map[string]string, keys string) string {
	// Sort keys
	var sortedKeys []string
	for key := range arr {
		if key != "" {
			sortedKeys = append(sortedKeys, key)
		}
	}
	sort.Strings(sortedKeys)

	// Build the sign string
	signStr := ""
	for _, key := range sortedKeys {
		val := arr[key]
		signStr += fmt.Sprintf("%s=%s&", key, val)
	}

	// Print the sign string (like echo in PHP)
	//fmt.Println("Sort Param : " + signStr)
	signStr += "key=" + keys
	//fmt.Println("Sort Param append Key : " + signStr)

	// Calculate MD5 hash
	hash := md5.Sum([]byte(signStr))
	return hex.EncodeToString(hash[:])
}
