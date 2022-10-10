package model

import "gorm.io/gorm"

type Light struct {
	gorm.Model
	Question string `gorm:"not null;unique"`
	Answer   string `gotm:"not null"`
}
