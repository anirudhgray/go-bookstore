package models

import "gorm.io/gorm"

type ShoppingCart struct {
	gorm.Model
	UserID uint    `gorm:"not null;unique;"` // FK to users
	Books  []*Book `gorm:"many2many:cart_books;"`
}

func (cart *ShoppingCart) CalculateTotalCartPrice() int64 {
	total := int64(0)
	for _, book := range cart.Books {
		total += book.Price
	}
	return total
}
