package services

import (
	"crypto/md5" // Package for MD5 hashing
	"encoding/json"
	"fmt"
	"io"
	"log"
	models "lotto-backend-api/models"
	"net/http"
	"os"
	"strconv"
	"time"

	promtpay "github.com/Frontware/promptpay"
	qrcode "github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

type QrPromtpayRequest struct {
	MemberId string `json:"member_id"`
	Amount   int64  `json:"amount"`
}

// ApiPayin constructs the request, generates a signature, and sends the request.
func ApiPromtpay(DB *gorm.DB, r *http.Request) map[string]interface{} {
	defer r.Body.Close()

	savePath := os.Getenv("PATH_TO_SAVE_FILE")

	var qrPromtpayRequest QrPromtpayRequest

	if err := json.NewDecoder(r.Body).Decode(&qrPromtpayRequest); err != nil {
		return map[string]interface{}{
			"code":    "-1",
			"message": "qr-protpay-service Invalid JSON format",
		}
	}

	// Assuming you've parsed the JSON body into a QrPromtpayRequest struct,
	// let's use some placeholder values for this example.
	req := QrPromtpayRequest{
		MemberId: qrPromtpayRequest.MemberId,
		Amount:   qrPromtpayRequest.Amount,
	}

	// 1. Generate a unique ID using the uuid package.
	//id := uuid.New().String()

	// 2. Get the current timestamp.
	//timestamp := time.Now().Format("20060102150405") // YYYYMMDDhhmmss

	inProgress := []string{"085xxxxxxx", "086xxxxxxx"}
	var promtpayTB models.Promtpays
	result := DB.Where("counter < ? AND total_deposit < ? AND qr_id NOT IN (?)", 399, 1990000, inProgress).First(&promtpayTB)

	if result.Error != nil {
		return map[string]interface{}{
			"code":    "-2",
			"message": result.Error,
		}
	}
	fmt.Println("promtpays table : ", promtpayTB.QrId)
	fmt.Println("promtpays table : ", promtpayTB.QrName)
	fmt.Println("promtpays table : ", promtpayTB.TotalDeposit)
	fmt.Println("promtpays table : ", promtpayTB.TotalWithdraw)
	fmt.Println("promtpays table : ", promtpayTB.Counter)

	now := time.Now()
	// Unix timestamp (seconds since epoch)
	fmt.Println("Unix timestamp (seconds):", now.Unix())

	// 3. Construct the filename with the required information.
	fileName := fmt.Sprintf("%s%s|%s|%s|%v.png",
		savePath,
		req.MemberId,
		promtpayTB.QrId,
		strconv.FormatInt(req.Amount, 10),
		now.Unix(),
	)

	// Your existing QR code generation logic
	payment := promtpay.PromptPay{
		PromptPayID: promtpayTB.QrId,
		Amount:      float64(req.Amount),
	}

	qrcodeString, err := payment.Gen()
	if err != nil {
		log.Fatalf("Error generating PromptPay QR string: %v", err)
	}

	fmt.Printf("Generated QR Code String: %s\n", qrcodeString)

	pngBytes, err := qrcode.Encode(qrcodeString, qrcode.Medium, 256)
	if err != nil {
		log.Fatalf("Error encoding string to QR code image: %v", err)
	}

	// 4. Save the QR code image to the newly constructed filename.
	err = os.WriteFile(fileName, pngBytes, 0644)
	if err != nil {
		log.Fatalf("Error saving QR code image: %v", err)
	}

	fmt.Printf("Successfully generated PromptPay QR code and saved to %s\n", fileName)
	qrImageName := qrPromtpayRequest.MemberId + "|" + promtpayTB.QrId + "|" + strconv.FormatInt(req.Amount, 10) + "|" + strconv.FormatInt(now.Unix(), 10) + ".png"

	md5Helper := md5.New()
	io.WriteString(md5Helper, qrImageName)
	hashSum := md5Helper.Sum(nil)
	hashString := fmt.Sprintf("%x", hashSum)

	// Save request to Orders table
	currentTime := time.Now()
	// Format as "YYYY-MM-DD"
	formattedDate := currentTime.Format("2006-01-02")
	formattedTime := currentTime.Format("15:04:05")

	// Add 1 minute to the current time
	oneMinuteLater := now.Add(4 * time.Minute)

	// Convert the new time to a Unix timestamp (in seconds)
	unixTimestamp := oneMinuteLater.Unix()
	orders := models.Orders{
		Date:       formattedDate,
		Time:       formattedTime,
		MemberId:   req.MemberId,
		OrderData:  hashString,
		ExpireTime: unixTimestamp,
	}

	orders_result := DB.Create(&orders)
	if orders_result.Error != nil {
		res := map[string]interface{}{
			"code":    "-1",
			"message": "Insert orders error : " + orders_result.Error.Error(),
		}
		return res
	}

	return map[string]interface{}{
		"code":           "200",
		"qr_img_name":    qrImageName,
		"qr_img_path":    fileName,
		"qr_id":          promtpayTB.QrId,
		"qr_name":        promtpayTB.QrName,
		"bank_provider":  promtpayTB.BankProvider,
		"member_request": qrPromtpayRequest.MemberId,
		"timestamp":      now.Unix(),
	}
}
