package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func GetImageByPath(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Param("path")

		// Validate the path parameter
		if path == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image path is required"})
			return
		}

		// Construct the full file path
		filePath := "/app/images/" + path

		// Check if the file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
			return
		}

		// Serve the file
		c.File(filePath)
	}
}
