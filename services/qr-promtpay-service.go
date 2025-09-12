package services

import (
	"encoding/json"
	"fmt"
	"log"
	"lotto-backend-api/models"
	"os"
	"strconv"
	"time"

	"net/http"

	promtpay "github.com/Frontware/promptpay"
	qrcode "github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

type QrPromtpayRequest struct {
	MemberId string `json:"member_id"`
	Amount   int64  `json:"amount"`
}

type QrPromtpayResponse struct {
	Id        string `json:"id"`
	MemberId  string `json:"member_id"`
	Amount    string `json:"amount"`
	QrImg     string `json:"qr_img"`
	Timestamp string `json:"timestamp"`
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
	return map[string]interface{}{
		"code":           "200",
		"qr_img_name":    qrPromtpayRequest.MemberId + "|" + promtpayTB.QrId + "|" + strconv.FormatInt(req.Amount, 10) + "|" + strconv.FormatInt(now.Unix(), 10) + ".png",
		"qr_img_path":    fileName,
		"qr_id":          promtpayTB.QrId,
		"qr_name":        promtpayTB.QrName,
		"bank_provider":  promtpayTB.BankProvider,
		"member_request": qrPromtpayRequest.MemberId,
		"timestamp":      now.Unix(),
	}
}
