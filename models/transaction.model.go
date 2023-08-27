package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	UserID         uint // FK to the user who initiated the transaction
	Amount         int64
	Status         string
	PaymentMethod  string
	CreditPurchase bool
	Books          []*Book `gorm:"many2many:transaction_books;"` // Books in the transaction
}

// could have gone with an alternative approach — storing a full cart instead of this books array, and subsequently creating a new cart for the user for future transactions. Decided against because of code complexity, and because this way I can reuse Transaction easily for single book direct purchases.
