package models

type Book struct {
	ID     int    `gorm:"primaryKey" example:"1" json:"id"`
	Title  string `gorm:"size:255;not null;uniqueIndex" json:"title" example:"LemonisTheBestFruit"`
	Author string `gorm:"not null" json:"author" example:"Lemon"`
	Stock  int    `gorm:"not null" json:"stock" example:"10"`
}
