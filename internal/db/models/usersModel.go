package models

type UsersModel struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique;not null"`
	Name     string `gorm:"not null"`
	Password string `gorm:"not null"`
}
