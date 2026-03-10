package routes

import (
	"note-service/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, handler *handlers.NoteHandler) {
	router.GET("/health", handler.Health)

	notes := router.Group("/notes")
	{
		notes.GET("", handler.GetNotes)
		notes.GET("/:id", handler.GetNote)
		notes.POST("", handler.CreateNote)
		notes.PUT("/:id", handler.UpdateNote)
		notes.DELETE("/:id", handler.DeleteNote)
	}
}
