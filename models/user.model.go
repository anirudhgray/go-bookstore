package models

import (
	"html"
	"strings"

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
	Password     string `gorm:"size:255;not null;"` // TODO password strength checks
	Name         string `gorm:"size:255;not null;"`
	Role         UserRole
	ShoppingCart ShoppingCart `gorm:"foreignKey:UserID;"` // one to one
	UserLibrary  UserLibrary  `gorm:"foreignKey:UserID;"`
	UserReviews  []Review     `gorm:"foreignKey:UserID;"` // one to many
	// TODO wishlists, one to many
}

func (user *User) Associate() error {
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	user.Name = html.EscapeString(strings.TrimSpace(user.Name))
	user.Role = BaseUser
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
