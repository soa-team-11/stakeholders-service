package services

import (
	"fmt"

	"stakeholder-service/internal/repos"
	"stakeholder-service/models"

	"github.com/google/uuid"
)

type ProfileService struct {
	profileRepo repos.ProfileRepo
}

func NewProfileService() *ProfileService {
	return &ProfileService{
		profileRepo: repos.NewProfileRepo(),
	}
}

func (s *ProfileService) Create(profile models.Profile) (*models.Profile, error) {
	if profile.UserID == uuid.Nil {
		return nil, fmt.Errorf("UserId cannot be nil")
	}

	created_profile, err := s.profileRepo.Create(profile)

	if err != nil {
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
