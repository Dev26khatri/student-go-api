package user

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Name     string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"not null;default:user"`
}

func (User) TableName() string {
	return "public.users"
}
