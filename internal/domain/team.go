package domain

type Team struct {
	TeamName string `gorm:"primaryKey;size:255"`
	Members  []User `gorm:"foreignKey:TeamName;references:TeamName"`
}
