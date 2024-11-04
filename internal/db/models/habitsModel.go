package models

type Habit struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Name   string `json:"name"`
	Days   int    `json:"days"`
	Day    int    `json:"day" gorm:"not null;default:0"`
	UserID int    `json:"user_id"`
}
