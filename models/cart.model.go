package models

import "gorm.io/gorm"

type ShoppingCart struct {
	gorm.Model
	UserID uint   `gorm:"not null;unique;"` // FK to users
	Books  []Book `gorm:"many2many:cart_books;"`
}
