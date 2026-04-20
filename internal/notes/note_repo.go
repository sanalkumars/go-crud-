package notes

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

//  this file will contain the repository functions for the notes collection in the database. It will contain functions to create, read, update and delete notes from the database. It will also contain functions to get all notes for a user and to get a note by its ID.

type Repo struct {
	coll *mongo.Collection
}

func NewRepo(db *mongo.Database) *Repo {
	return &Repo{coll: db.Collection("notes")}
}

func (r *Repo) CreateNote(ctx context.Context, note *Note) (Note, error) {

	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.coll.InsertOne(opCtx, note)
	if err != nil {
		return Note{}, fmt.Errorf("Insertion Failed")
	}

	return *note, nil
}

func (r *Repo) updateNote(ctx context.Context, id string, note *Note) (Note, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Note{}, fmt.Errorf("Invalid ID format")
	}

	update := bson.M{"$set": bson.M{
		"title":      note.Title,
		"content":    note.Content,
		"pinned":     note.Pinned,
		"updated_at": note.UpdatedAt,
	}}

	_, err = r.coll.UpdateOne(opCtx, bson.M{"_id": objId}, update)
	if err != nil {
		return Note{}, fmt.Errorf("Update Failed")
	}

	return r.getNoteById(ctx, id)
}

func (r *Repo) GetAllNotes(ctx context.Context) ([]Note, error) {

	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.coll.Find(opCtx, bson.M{})

	if err != nil {
		return nil, err
	}
	defer cursor.Close(opCtx)

	var Notes []Note
	if err = cursor.All(opCtx, &Notes); err != nil {
		return nil, err
	}
	return Notes, nil
}

func (r *Repo) getNoteById(ctx context.Context, id string) (Note, error) {
	fmt.Println("id got 2", id)
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Note{}, fmt.Errorf("Invalid ID format")
	}

	var note Note
	err = r.coll.FindOne(opCtx, bson.M{"_id": objID}).Decode(&note)
	fmt.Println("note got is ", note)
	if err != nil {
		return Note{}, fmt.Errorf("Note not found")
	}
	return note, nil

}
