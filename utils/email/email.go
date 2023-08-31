package email

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"time"

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

func GenericSendMail(subject string, content string, toEmail string, userName string) error {
	url := "https://send.api.mailtrap.io/api/send"
	method := "POST"

	data := RegistrationEmail{
		Subject: subject,
		From: EmailAddress{
			Email: "bookstore@anrdhmshr.tech",
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
		logger.Errorf("Email Error: %v", err)
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))

	if err != nil {
		logger.Errorf("Email Error: %v", err)
		return err
	}

	bearer := fmt.Sprintf("Bearer %s", viper.GetString("MAILTRAP_API_TOKEN"))
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		logger.Errorf("Email Error: %v", err)
		return err
	}
	defer res.Body.Close()
	return nil
}

func SendRegistrationMail(subject string, content string, toEmail string, userID uint, userName string, newUser bool) error {
	otp := ""
	if newUser {
		otp = GenerateOTP(6)
		content += "http://bookstore.anrdhmshr.tech/verify?email=" + toEmail + "&otp=" + otp
	}

	err := GenericSendMail(subject, content, toEmail, userName)
	if err != nil {
		return err
	}

	if newUser {
		entry := models.VerificationEntry{
			Email: toEmail,
			OTP:   otp,
		}
		database.DB.Create(&entry)
	}
	return nil
}

func SendDeletionMail(toEmail string, userID uint, userName string) {
	otp := ""
	otp = GenerateOTP(6)
	confirmationURL := ""
	confirmationURL += "http://bookstore.anrdhmshr.tech/delete-account?email=" + toEmail + "&otp=" + otp
	content := "A request for the deletion of the bookstore account associated with your user has been made. If this was not you, please change your password. Otherwise, click on this link to confirm account deletion: " + confirmationURL + " . This link will be active for 3 minutes."
	subject := "Request for account deletion."

	GenericSendMail(subject, content, toEmail, userName)

	entry := models.DeletionConfirmation{
		Email:     toEmail,
		OTP:       otp,
		ValidTill: time.Now().Add(3 * time.Minute),
	}
	database.DB.Create(&entry)
}

func SendForgotPasswordMail(toEmail string, userID uint, userName string) {
	otp := ""
	otp = GenerateOTP(6)
	verificationURL := ""
	verificationURL += "http://bookstore.anrdhmshr.tech/set-forgotten-password?email=" + toEmail + "&otp=" + otp
	content := "A forgot password request was made for the email associated with your account. If this was not you, feel free to ignore this email. Otherwise, click on this link to post your new password: " + verificationURL + " . This link will be active for 3 minutes."
	subject := "Forgot Password."

	GenericSendMail(subject, content, toEmail, userName)

	entry := models.ForgotPassword{
		Email:     toEmail,
		OTP:       otp,
		ValidTill: time.Now().Add(3 * time.Minute),
	}
	database.DB.Create(&entry)
}
