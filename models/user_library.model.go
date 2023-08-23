package models

import "gorm.io/gorm"

// bought books by a user
type UserLibrary struct {
	gorm.Model
	UserID uint    `gorm:"not null;unique;"` // FK to users
	Books  []*Book `gorm:"many2many:user_library_books;"`
}
