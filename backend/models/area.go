package models

import "github.com/jinzhu/gorm"

type Area struct {
	gorm.Model
	AreaName       string  `gorm:"type:varchar(128);not null"`
	TopRightLat    float64 `gorm:"not null"`
	TopRightLon    float64 `gorm:"not null"`
	BottomLeftLat  float64 `gorm:"not null"`
	BottomLeftLon  float64 `gorm:"not null"`
	DeforestedArea float64 `gorm:"default:0.0"`
	UserID         uint    `gorm:"not null"`
}
