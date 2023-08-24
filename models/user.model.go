package models

import (
	"html"
	"strings"

	"github.com/anirudhgray/balkan-assignment/infra/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRole string

const (
	Admin    UserRole = "admin"
	BaseUser UserRole = "base_user"
)

type User struct {
	gorm.Model
	Email        string `gorm:"size:255;not null;unique;"`
	Password     string `gorm:"size:255;not null;"`
	Name         string `gorm:"size:255;not null;"`
	Role         UserRole
	Verified     bool
	ShoppingCart ShoppingCart   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"` // one to one
	UserLibrary  UserLibrary    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"` // one to one
	UserReviews  []*Review      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"` // one to many
	Transactions []*Transaction `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"` // one to many
	// TODO wishlists, one to many
}

func (user *User) BeforeDelete(tx *gorm.DB) error {
	// Set UserID to null for related transactions
	if err := tx.Model(&Transaction{}).Where("user_id = ?", user.ID).Update("user_id", nil).Error; err != nil {
		return err
	}
	return nil
}

func (user *User) Associate() error {
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	user.Name = html.EscapeString(strings.TrimSpace(user.Name))
	user.Role = BaseUser
	user.Verified = false

	return nil
}

func (user *User) AttachCartAndLibrary() error {
	shoppingCart := ShoppingCart{UserID: user.ID}
	userLibrary := UserLibrary{UserID: user.ID}
	if err := database.DB.Create(&shoppingCart).Error; err != nil {
		return err
	}
	if err := database.DB.Create(&userLibrary).Error; err != nil {
		return err
	}
	user.ShoppingCart = ShoppingCart{UserID: user.ID}
	user.UserLibrary = UserLibrary{UserID: user.ID}
	return nil
}

func (user *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

func (u *User) HasBookInLibrary(bookID uint) bool {
	if err := database.DB.Preload("Books").First(&u.UserLibrary, "user_id = ?", u.ID).Error; err != nil {
		return false
	}

	for _, book := range u.UserLibrary.Books {
		if book.ID == bookID {
			return true
		}
	}
	return false
}
