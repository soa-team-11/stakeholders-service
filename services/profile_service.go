package services

import (
	"context"
	"fmt"

	"stakeholder-service/internal/repos"
	"stakeholder-service/models"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type ProfileService struct {
	profileRepo repos.ProfileRepo
}

func NewProfileService() *ProfileService {
	return &ProfileService{
		profileRepo: repos.NewProfileRepo(),
	}
}

func (s *ProfileService) Create(ctx context.Context, profile models.Profile) (*models.Profile, error) {
	tracer := otel.Tracer("stakeholders-service")
	_, span := tracer.Start(ctx, "ProfileService.Create")
	defer span.End()

	if profile.UserID == uuid.Nil {
		return nil, fmt.Errorf("UserId cannot be nil")
	}

	created_profile, err := s.profileRepo.Create(profile)

	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("%s", err.Error())
	}

	return created_profile, nil
}

func (s *ProfileService) GetByUserId(UserId uuid.UUID) (*models.Profile, error) {
	profile, err := s.profileRepo.GetByUserId(UserId)

	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	return profile, nil
}

func (s *ProfileService) Update(profile models.Profile) (*models.Profile, error) {
	if profile.UserID == uuid.Nil {
		return nil, fmt.Errorf("UserID cannot be nil")
	}

	existing, err := s.profileRepo.GetByUserId(profile.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to load existing profile: %w", err)
	}

	if profile.ImageURL == "" {
		profile.ImageURL = existing.ImageURL
	}

	updatedProfile, err := s.profileRepo.Update(profile)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return updatedProfile, nil
}

func (s *ProfileService) GetRecommendations(userId uuid.UUID) ([]models.Profile, error) {
	if userId == uuid.Nil {
		return nil, fmt.Errorf("UserID cannot be nil")
	}

	profiles, err := s.profileRepo.GetRecommendations(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get recommendations: %w", err)
	}

	return profiles, nil
}
