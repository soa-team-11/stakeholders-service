package repos

import (
	"context"

	mg "stakeholder-service/internal/providers/mongo"
	"stakeholder-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProfileRepo interface {
	Create(profile models.Profile) (*models.Profile, error)
	GetByUserId(userId string) (*models.Profile, error)
	// Add more if needed: Update, Delete, etc.
}

type ProfileRepoImpl struct {
	collection *mongo.Collection
}

func NewProfileRepo() *ProfileRepoImpl {
	return &ProfileRepoImpl{
		collection: mg.GetDatabase().Collection("profiles"),
	}
}

func (r *ProfileRepoImpl) Create(profile models.Profile) (*models.Profile, error) {
	_, err := r.collection.InsertOne(context.Background(), profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepoImpl) GetByUserId(userId string) (*models.Profile, error) {
	var profile models.Profile
	err := r.collection.FindOne(context.Background(), bson.M{"user_id": userId}).Decode(&profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}
