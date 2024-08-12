package handlers

import (
	"deforestation/jobs"
	"deforestation/models"
	"deforestation/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type CreateAreaInput struct {
	AreaName      string  `json:"area_name" binding:"required"`
	TopRightLat   float64 `json:"top_right_lat" binding:"required"`
	TopRightLon   float64 `json:"top_right_lon" binding:"required"`
	BottomLeftLat float64 `json:"bottom_left_lat" binding:"required"`
	BottomLeftLon float64 `json:"bottom_left_lon" binding:"required"`
}

func CreateArea(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("userID")

		var input CreateAreaInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
			return
		}

		area := models.Area{
			AreaName:      input.AreaName,
			TopRightLat:   input.TopRightLat,
			TopRightLon:   input.TopRightLon,
			BottomLeftLat: input.BottomLeftLat,
			BottomLeftLon: input.BottomLeftLon,
			UserID:        userID,
		}

		db.Create(&area)

		// Save job schedule
		jobSchedule := models.JobSchedule{AreaID: area.ID}
		db.Create(&jobSchedule)

		// Start job for the new area here
		jobs.StartWeeklyJob(area.ID, utils.GetSatelliteImage)
		err := utils.GetSatelliteImage(area.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		}

		c.JSON(http.StatusOK, gin.H{"data": area})
	}
}

// GetArea fetches an area by its ID
func GetArea(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		userID := c.GetUint("userID")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid area ID"})
			return
		}

		var area models.Area
		if err := db.First(&area, id).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Area not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		if area.UserID != userID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You do not have access to this area"})
		}

		c.JSON(http.StatusOK, gin.H{"data": area})
	}
}

// GetAllAreas fetches an area by its ID
func GetAllAreas(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("userID")

		var areas []models.Area
		if err := db.Where("user_id = ?", userID).Find(&areas).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Area not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": areas})
	}
}

// DeleteArea deletes an area by its ID
func DeleteArea(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid area ID"})
			return
		}

		// Check if the area exists
		var area models.Area
		if err := db.First(&area, id).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Area not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		// Delete the area
		if err := db.Delete(&area).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Area deleted successfully"})
	}
}
