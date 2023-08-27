package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	UserID  uint `gorm:"not null"` // FK to user
	BookID  uint `gorm:"not null"` // FK to book
	Comment string
	Rating  int `gorm:"not null"`
}
