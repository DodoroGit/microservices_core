package main

import (
	"log"
	"note-service/config"
	"note-service/database"
	"note-service/handlers"
	"note-service/repository"
	"note-service/routes"
	"note-service/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	mongoClient, err := database.InitMongoDB(cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer database.CloseMongoDB(mongoClient)

	db := mongoClient.Database(cfg.MongoDB.DBName)
	noteRepo := repository.NewNoteRepository(db)
	noteService := services.NewNoteService(noteRepo)
	noteHandler := handlers.NewNoteHandler(noteService)

	router := gin.Default()
	routes.SetupRoutes(router, noteHandler)

	log.Printf("note-service running on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
