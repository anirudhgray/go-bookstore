package models

import "gorm.io/gorm"

type VerificationEntry struct {
	gorm.Model
	Email string `gorm:"unique"`
	OTP   string
}
