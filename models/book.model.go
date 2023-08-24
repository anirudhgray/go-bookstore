package models

import (
	"time"

	"github.com/anirudhgray/balkan-assignment/infra/logger"
	"gorm.io/gorm"
)

type BookCategory string

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
	Author        string    `gorm:"not null"`
	Publisher     string    `gorm:"not null"`
	Date          time.Time `gorm:"not null"`
	Year          time.Time
	Price         int64 `gorm:"not null"`
	FilePath      string
	ISBN          string          `gorm:"size:255;not null"`
	Category      BookCategory    `gorm:"not null"`
	ShoppingCarts []*ShoppingCart `gorm:"many2many:cart_books;"`
	UserLibraries []*UserLibrary  `gorm:"many2many:user_library_books;"`
	Reviews       []*Review       `gorm:"foreignKey:BookID;not null"`
	CatalogDelete bool
}

func (book *Book) CalculateAvgRating() float64 {
	total := float64(0)
	count := float64(0)
	for _, review := range book.Reviews {
		total += float64(review.Rating)
		count += 1.0
	}
	if count == 0 {
		logger.Infof("Rip\n")
		return 0
	}
	logger.Infof("Rating: %v\n", total/count)
	return (total / count)
}
