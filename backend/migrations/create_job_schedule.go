package migrations

import (
	"deforestation/models"

	"github.com/jinzhu/gorm"
)

func CreateJobSchedules(db *gorm.DB) {
	db.AutoMigrate(&models.JobSchedule{})
}
