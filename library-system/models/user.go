package models

type User struct {
	ID       int    `gorm:"primaryKey" json:"id" example:"123"`
	Name     string `gorm:"type:varchar(255);uniqueIndex;not null" json:"name" example:"lemon"`
	Password string `gorm:"not null" json:"-"`
	Role     string `gorm:"not null" json:"role" example:"admin"`
}
