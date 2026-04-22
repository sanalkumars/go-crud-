package notes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)





func RegisterRoutes( r *gin.Engine , db *mongo.Database){
	// we are creating a new instance of the Repo struct and passing the database connection to it. We are then creating a new instance of the Handler struct and passing the Repo instance to it. Finally, we are registering the route for creating a note and associating it with the CreateNoteHandler function.

	repo := NewRepo(db)
	handler := NewHandler(repo)

	noteGroups := r.Group("/notes")
	{
		noteGroups.POST("/", handler.CreateNoteHandler)
		noteGroups.GET("/", handler.GetAllNotesHandler)
		noteGroups.GET("/:id",handler.GetNoteByIDHandler)
		noteGroups.PUT("/:id",handler.updateNoteHandler)
		noteGroups.DELETE("/:id", handler.deleteNoteHnadler)
		// we will add more routes here for reading, updating and deleting notes.

	}
}