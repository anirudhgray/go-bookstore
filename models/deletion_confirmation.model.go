package models

import (
	"time"

	"gorm.io/gorm"
)

type DeletionConfirmation struct {
	gorm.Model
	Email     string `gorm:"unique"`
	OTP       string
	ValidTill time.Time
}
