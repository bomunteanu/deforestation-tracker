package models

type JobSchedule struct {
	ID     uint `gorm:"primary_key"`
	AreaID uint `gorm:"not null"`
}
