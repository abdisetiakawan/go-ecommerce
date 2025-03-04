package entity

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	UserID 		uint 	   `gorm:"not null"`
	User 		User 	   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name        string     `gorm:"-"`
	Gender 		string	   `gorm:"type:enum('male', 'female');not null"`
	PhoneNumber string     `gorm:"size:15;not null"`
	Address     string     `gorm:"size:255;not null"`
	Avatar      string     `gorm:"size:255"`
	Bio         string     `gorm:"type:text"`
}