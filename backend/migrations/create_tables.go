package migrations

import (
	"deforestation/models"

	"github.com/jinzhu/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Area{}, &models.History{}, &models.User{})
}
