package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserUUID  string    `gorm:"type:char(36);uniqueIndex;not null"`
	Username  string 	`gorm:"size:100;not null;uniqueIndex"`
	Name      string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:100;uniqueIndex;not null"`
	Role      string    `gorm:"type:enum('seller','buyer');not null"`
	Password  string	`gorm:"size:255;not null"`
	ConfirmPassword string `gorm:"-"`
	
	AccessToken  string `gorm:"-"`
	RefreshToken string `gorm:"-"`
	
	Store 	  *Store 	`gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
