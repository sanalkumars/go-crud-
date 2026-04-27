package auth

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repo struct {
	coll *mongo.Collection
}

func NewRepo(db *mongo.Database) *Repo {
	return &Repo{coll: db.Collection("users")}
}

func (r *Repo) CreateUser(ctx context.Context, user *User) error {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.coll.InsertOne(opCtx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *Repo) GetUserByEmail(ctx context.Context, email string) (User, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user User
	err := r.coll.FindOne(opCtx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
