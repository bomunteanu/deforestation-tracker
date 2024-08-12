package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type History struct {
	gorm.Model
	Date            time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	ImagePath       string    `gorm:"type:varchar(256);not null"`
	MaskedImagePath string    `gorm:"type:varchar(256);not null"`
	DeforestedArea  float64   `gorm:"not null"`
	AreaID          uint      `gorm:"not null"`
	Area            Area      `gorm:"foreignkey:AreaID"`
}
