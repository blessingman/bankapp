package models

import "gorm.io/gorm"

type Budget struct {
	gorm.Model
	UserID   uint
	Category string
	Amount   float64
}
