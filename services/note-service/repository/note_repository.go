package repository

import (
	"context"
	"errors"
	"time"

	"note-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NoteRepositoryInterface interface {
	FindAll(category string) ([]models.Note, error)
	FindByID(id string) (*models.Note, error)
	Create(note *models.Note) (*models.Note, error)
	Update(id string, updates map[string]interface{}) (*models.Note, error)
	Delete(id string) error
}

type NoteRepository struct {
	collection *mongo.Collection
}

func NewNoteRepository(db *mongo.Database) *NoteRepository {
	return &NoteRepository{
		collection: db.Collection("notes"),
	}
}

func (r *NoteRepository) FindAll(category string) ([]models.Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if category != "" {
		filter["category"] = category
	}

	opts := options.Find().SetSort(bson.D{bson.E{Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notes []models.Note
	if err := cursor.All(ctx, &notes); err != nil {
		return nil, err
	}
	if notes == nil {
		notes = []models.Note{}
	}
	return notes, nil
}

func (r *NoteRepository) FindByID(id string) (*models.Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, nil
	}

	var note models.Note
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&note)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *NoteRepository) Create(note *models.Note) (*models.Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now().UTC()
	note.ID = primitive.NewObjectID()
	note.CreatedAt = now
	note.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, note)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (r *NoteRepository) Update(id string, updates map[string]interface{}) (*models.Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, nil
	}

	updates["updated_at"] = time.Now().UTC()

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": updates},
	)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, nil
	}
	return r.FindByID(id)
}

func (r *NoteRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("note not found")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("note not found")
	}
	return nil
}
