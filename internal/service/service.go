package service

import (
	"pets/internal/model"
	"pets/internal/repository"
	"pets/pkg/logger"
)

// IService is an app service layer interface
type IService interface {
	// GetPets is used to get pets. Limit and offset can be used for pagination. 0 limit and 0 offset will return
	// all existing pets. Not convertable values limit and offset will be ignored. For order arg can be used "asc" and "desc"
	// string value to order pets by ID.
	// Function will return slice of pets model, total found pets or error
	GetPets(limit string, offset string, order string) ([]*model.Pet, int, error)

	// AddPet is used to add new pet to the DB. Only "name" field will be used.
	AddPet(pet *model.Pet) (int, error)

	// UpdatePet is used to update existing pet. If pet with given ID not exist, will return error. Only "name" and "id"
	// fields will be used.
	UpdatePet(pet *model.Pet) error

	// DeletePet is used to delete existing pet. If pet with given ID not exist, will return error. Only "id"
	// field will be used.
	DeletePet(pet *model.Pet) error

	// IsExist is used to check if Pet exists by given ID. If DB returns error function will return false.
	IsExist(id int) bool
}

// Service is a service struct implementing IService interface
type Service struct {
	repository repository.IRepository
}

// NewService is used to get new Service instance
func NewService(rep repository.IRepository) IService {
	s := &Service{}

	s.repository = rep

	logger.Log().WithField("layer", "Service-Init").Infof("service created")

	return s
}
