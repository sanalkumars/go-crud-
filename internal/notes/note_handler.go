package notes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) CreateNoteHandler(c *gin.Context) {

	var req NoteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON Body",
		})
		return
	}

	now := time.Now().UTC()
	//  here we are creating a new note object and populating it with the data from the request body. We are also setting the createdAt and updatedAt fields to the current time. We are also generating a new ObjectID for the note using the primitive.NewObjectID() function from the MongoDB driver.
	note := Note{
		ID:        primitive.NewObjectID(),
		Title:     req.Title,
		Content:   req.Content,
		Pinned:    *req.Pinned,
		CreatedAt: now,
		UpdatedAt: now,
	}

	createdNote, err := h.repo.CreateNote(c.Request.Context(), &note)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create note",
		})
		return
	}
	c.JSON(http.StatusCreated, createdNote)
}

func (h *Handler) updateNoteHandler(c *gin.Context) {
	var req NoteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON Body",
		})
		return
	}
	updatedNote, err := h.repo.updateNote(c.Request.Context(), c.Param("id"), &Note{
		Title:     req.Title,
		Content:   req.Content,
		Pinned:    *req.Pinned,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update note",
		})
		return
	}
	c.JSON(http.StatusOK, updatedNote)

}

func (h *Handler) GetAllNotesHandler(c *gin.Context) {
	notes, err := h.repo.GetAllNotes(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get notes",
		})
		return
	}
	c.JSON(http.StatusOK, notes)
}

func (h *Handler) GetNoteByIDHandler(c *gin.Context) {
	params := c.Param("id")
	fmt.Println("id got is ", params)
	note, err := h.repo.getNoteById(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get note",
		})
		return
	}
	c.JSON(http.StatusOK, note)
}
