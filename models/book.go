package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title  string `json:"title" gorm:"not null;uniqueIndex:title_idx"`
	Author string `json:"author" gorm:"not null"`
	Year   int    `json:"year" gorm:"not null"`
}
