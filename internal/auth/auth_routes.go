package auth

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(r *gin.Engine, db *mongo.Database, jwtSecret string) {
	repo := NewRepo(db)
	handler := NewHandler(repo, jwtSecret)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", handler.SignupHandler)
		authGroup.POST("/signin", handler.SigninHandler)
	}
}
