package handlers

import (
	"deforestation/database"
	"deforestation/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetAllHistories returns all history items
func GetAllHistories(c *gin.Context) {
	db := database.GetDB()

	var histories []models.History
	if err := db.Find(&histories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": histories})
}

// GetHistoryByID returns a single history item by ID
func GetHistoryByID(c *gin.Context) {
	db := database.GetDB()
	historyID := c.Param("id")

	var history models.History
	if err := db.First(&history, historyID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "History not found"})
		return
	}

	var area models.Area
	if err := db.First(&area, history.AreaID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Area not found"})
		return
	}

	userID := c.GetUint("userID")
	if userID != area.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this entry"})
		return
	}

	// Serve the image file
	c.JSON(http.StatusOK, history)
}

// GetAllHistories returns all history items
func GetHistoriesByAreaID(c *gin.Context) {
	db := database.GetDB()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid area ID"})
		return
	}

	userID := c.GetUint("userID")
	var area models.Area
	if err := db.Where("id = ?", id).First(&area).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Area not found"})
		return
	}

	if userID != area.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this entry"})
		return
	}

	var histories []models.History
	if err := db.Where("area_id = ?", id).Find(&histories).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No histories found for this area"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": histories})
}
