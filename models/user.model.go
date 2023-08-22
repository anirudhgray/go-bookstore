package models

import "gorm.io/gorm"

type UserRole string

const (
	Admin    UserRole = "admin"
	BaseUser UserRole = "base_user"
)

type User struct {
	gorm.Model
	Email        string       `gorm:"size:255;not null;unique;"`
	Name         string       `gorm:"size:255;not null;"`
	Role         UserRole     `gorm:"not null;"`
	ShoppingCart ShoppingCart `gorm:"foreignKey:UserID;not null;"` // one to one
	UserLibrary  UserLibrary  `gorm:"foreignKey:UserID;not null;"`
	UserReviews  []Review     `gorm:"foreignKey:UserID;not null"` // one to many
	// TODO wishlists, one to many
}
