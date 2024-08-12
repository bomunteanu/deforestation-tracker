package main

import (
	db "deforestation/database"
	"deforestation/handlers"
	"deforestation/middleware"
	"deforestation/migrations"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	db.Init()
	defer db.GetDB().Close()

	// Run migrations
	migrations.Migrate(db.GetDB())

	r := gin.Default()

	// CORS middleware setup
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.POST("/signup", handlers.Signup(db.GetDB()))
	r.POST("/login", handlers.Login(db.GetDB()))

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	protected.GET("/auth/check", handlers.Check(db.GetDB()))

	protected.POST("/areas", handlers.CreateArea(db.GetDB()))
	protected.GET("/areas", handlers.GetAllAreas(db.GetDB()))
	protected.GET("/areas/:id", handlers.GetArea(db.GetDB()))
	protected.DELETE("/areas/:id", handlers.DeleteArea(db.GetDB()))

	protected.GET("/images/:path", handlers.GetImageByPath(db.GetDB()))

	protected.GET("/histories", handlers.GetAllHistories)
	protected.GET("/histories/:id", handlers.GetHistoryByID)
	protected.GET("/histories/area/:id", handlers.GetHistoriesByAreaID)

	if err := r.Run(); err != nil {
		fmt.Printf("Gin server encountered an error: %v\n", err)
	}
}
