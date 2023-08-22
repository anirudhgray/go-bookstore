package models

import (
	"gorm.io/gorm"
)

type Example struct {
	gorm.Model
	Data string `binding:"required"`
}

// TableName is Database TableName of this model
func (e *Example) TableName() string {
	return "examples"
}
