package entity

import "gorm.io/gorm"

type Store struct {
	gorm.Model
	UserID 		uint 	`gorm:"not null"`
	User 		User 	`gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	StoreName 	string 	`gorm:"constraint:OnUpdate:CASCADE,OnDelete:Cascade;"`
	Description string 	`gorm:"type:text"`
}