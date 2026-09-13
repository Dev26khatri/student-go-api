package student

import "gorm.io/gorm"

// Student represents a student record stored in the database.
type Student struct {
	gorm.Model
	Name  string `gorm:"not null"`
	Email string `gorm:"unique;not null;index"`
	Age   int    `gorm:"not null"`
	Photo string `gorm:"column:photo" json:"photo"`
}

// TableName sets the database table used for student records.
func (Student) TableName() string {
	return "public.students"
}
