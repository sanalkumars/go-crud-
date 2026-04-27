package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	repo      *Repo
	jwtSecret string
}

func NewHandler(repo *Repo, jwtSecret string) *Handler {
	return &Handler{repo: repo, jwtSecret: jwtSecret}
}

func (h *Handler) SignupHandler(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	_, err := h.repo.GetUserByEmail(c.Request.Context(), email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check existing user"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process password"})
		return
	}

	user := User{
		ID:           primitive.NewObjectID(),
		Name:         strings.TrimSpace(req.Name),
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now().UTC(),
	}

	if err := h.repo.CreateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	token, err := GenerateToken(user.ID.Hex(), user.Email, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "signup successful",
		"token":   token,
		"user": gin.H{
			"id":    user.ID.Hex(),
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func (h *Handler) SigninHandler(c *gin.Context) {
	var req SigninRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	user, err := h.repo.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := GenerateToken(user.ID.Hex(), user.Email, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "signin successful",
		"token":   token,
		"user": gin.H{
			"id":    user.ID.Hex(),
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
