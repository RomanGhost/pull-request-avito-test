package domain

type User struct {
	UserID   string `gorm:"primaryKey;size:255"`
	Username string `gorm:"size:255;not null"`
	TeamName string `gorm:"size:255;index"`
	IsActive bool   `gorm:"default:true"`
}
