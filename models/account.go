package models

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	UserID  uint
	Balance float64
}
