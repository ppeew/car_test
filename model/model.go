package model

import "gorm.io/gorm"

type Light struct {
	gorm.Model
	Question string `gorm:"type:TEXT;not null"`
	Answer   string `gotm:"not null"`
}
