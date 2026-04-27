package server

import (
	"net/http"
	"notes-api/internal/auth"
	"notes-api/internal/notes"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(db *mongo.Database, jwtSecret string) *gin.Engine {
	router := gin.Default()

	// Define your routes here
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	//  we need to register the routes for
	// the notes module. We have created a new
	// file called note_routes.go in the
	// internal/notes package and we will d
	// efine a function called RegisterRoutes
	// that will take the router and
	// the database connection as parameters
	//  and it will register the routes for
	// the notes module.

	auth.RegisterRoutes(router, db, jwtSecret)
	notes.RegisterRoutes(router, db, auth.AuthMiddleware(jwtSecret))

	return router
}
