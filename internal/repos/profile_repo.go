package repos

import (
	"context"

	mg "stakeholder-service/internal/providers/mongo"
	"stakeholder-service/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProfileRepo interface {
	Create(profile models.Profile) (*models.Profile, error)
	GetByUserId(userId uuid.UUID) (*models.Profile, error)
	Update(profile models.Profile) (*models.Profile, error)
	GetRecommendations(userId uuid.UUID) ([]models.Profile, error)
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

func (r *ProfileRepoImpl) GetByUserId(userId uuid.UUID) (*models.Profile, error) {
	var profile models.Profile
	err := r.collection.FindOne(context.Background(), bson.M{"user_id": userId}).Decode(&profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepoImpl) Update(profile models.Profile) (*models.Profile, error) {
	filter := bson.M{"user_id": profile.UserID}

	updateDoc := bson.M{
		"name":      profile.Name,
		"surname":   profile.Surname,
		"bio":       profile.Bio,
		"motto":     profile.Motto,
		"image_url": profile.ImageURL,
	}

	update := bson.M{"$set": updateDoc}

	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	return r.GetByUserId(profile.UserID)
}

func (r *ProfileRepoImpl) GetRecommendations(userId uuid.UUID) ([]models.Profile, error) {
	ctx := context.Background()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"user_id": bson.M{"$ne": userId}}}},
		{{Key: "$sample", Value: bson.M{"size": 5}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var profiles []models.Profile
	if err := cursor.All(ctx, &profiles); err != nil {
		return nil, err
	}

	return profiles, nil
}
