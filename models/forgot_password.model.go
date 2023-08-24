package models

import (
	"time"

	"gorm.io/gorm"
)

type ForgotPassword struct {
	gorm.Model
	Email     string `gorm:"unique"`
	OTP       string
	ValidTill time.Time
}
