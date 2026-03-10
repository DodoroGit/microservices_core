package handlers

import (
	"net/http"
	"strings"

	"note-service/models"
	"note-service/services"

	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	service services.NoteServiceInterface
}

func NewNoteHandler(service services.NoteServiceInterface) *NoteHandler {
	return &NoteHandler{service: service}
}

func (h *NoteHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "note-service"})
}

func (h *NoteHandler) GetNotes(c *gin.Context) {
	category := c.Query("category")
	notes, err := h.service.GetAll(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]models.NoteResponse, len(notes))
	for i, note := range notes {
		responses[i] = note.ToResponse()
	}
	c.JSON(http.StatusOK, responses)
}

func (h *NoteHandler) GetNote(c *gin.Context) {
	id := c.Param("id")
	note, err := h.service.GetByID(id)
	if err != nil {
		if err.Error() == "note not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, note.ToResponse())
}

func (h *NoteHandler) CreateNote(c *gin.Context) {
	var req models.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := h.service.Create(req)
	if err != nil {
		msg := err.Error()
		if strings.HasPrefix(msg, "invalid") {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusCreated, note.ToResponse())
}

func (h *NoteHandler) UpdateNote(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := h.service.Update(id, req)
	if err != nil {
		switch err.Error() {
		case "note not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "invalid category":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, note.ToResponse())
}

func (h *NoteHandler) DeleteNote(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		if err.Error() == "note not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
