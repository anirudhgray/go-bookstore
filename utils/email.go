package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net/http"

	"github.com/anirudhgray/balkan-assignment/infra/database"
	"github.com/anirudhgray/balkan-assignment/infra/logger"
	"github.com/anirudhgray/balkan-assignment/models"
	"github.com/spf13/viper"
)

func GenerateOTP(maxDigits uint32) string {
	bi, err := rand.Int(
		rand.Reader,
		big.NewInt(int64(math.Pow(10, float64(maxDigits)))),
	)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%0*d", maxDigits, bi)
}

type EmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type RegistrationEmail struct {
	Subject  string         `json:"subject"`
	From     EmailAddress   `json:"from"`
	To       []EmailAddress `json:"to"`
	Category string         `json:"category"`
	Text     string         `json:"text"`
}

func SendRegistrationMail(subject string, content string, toEmail string, userID uint, userName string, newUser bool) {
	otp := ""
	if newUser {
		otp = GenerateOTP(6)
		content += "http://0.0.0.0:8000/v1/auth/verify?email=" + toEmail + "&otp=" + otp

	}
	url := "https://send.api.mailtrap.io/api/send"
	method := "POST"

	data := RegistrationEmail{
		Subject: subject,
		From: EmailAddress{
			Email: "mailtrap@anrdhmshr.tech",
			Name:  "BOOKSTORE ADMIN",
		},
		To: []EmailAddress{
			{
				Email: toEmail,
				Name:  userName,
			},
		},
		Category: "BookStore",
		Text:     content,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("Error: %v", err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))

	if err != nil {
		logger.Errorf("Error: %v", err)
		return
	}

	bearer := fmt.Sprintf("Bearer %s", viper.GetString("MAILTRAP_API_TOKEN"))
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		logger.Errorf("Error: %v", err)
		return
	}
	defer res.Body.Close()

	if newUser {
		entry := models.VerificationEntry{
			Email: toEmail,
			OTP:   otp,
		}
		database.DB.Create(&entry)
	}
}
