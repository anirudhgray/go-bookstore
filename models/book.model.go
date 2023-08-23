package models

import (
	"time"

	"gorm.io/gorm"
)

type BookCategory string // TODO book category model so admins can add categories

const (
	Fantasy  BookCategory = "fantasy"
	Mystery  BookCategory = "mystery"
	Drama    BookCategory = "drama"
	Historic BookCategory = "historic"
	SelfHelp BookCategory = "self_help"
	Travel   BookCategory = "travel"
)

type Book struct {
	gorm.Model
	Name          string    `gorm:"not null;"`
	Author        string    `gorm:"not null"` // TODO FK to Author
	Publisher     string    `gorm:"not null"` // TODO FK to Publisher
	Date          time.Time `gorm:"not null"`
	Year          time.Time
	Price         int64 `gorm:"not null"`
	FilePath      string
	ISBN          string          `gorm:"size:255;not null"`
	Category      BookCategory    `gorm:"not null"`
	ShoppingCarts []*ShoppingCart `gorm:"many2many:cart_books;"` // TODO for a "In X carts currently!" feature
	UserLibraries []*UserLibrary  `gorm:"many2many:user_library_books;"`
	Reviews       []Review        `gorm:"foreignKey:BookID;not null"`
}
