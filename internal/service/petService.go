package service

import (
	"database/sql"
	"errors"
	"strconv"

	"pets/internal/model"
)

// GetPets is implementing IService.GetPets function
func (s *Service) GetPets(limit string, offset string, order string) ([]*model.Pet, int, error) {
	res, err := s.repository.GetPets(convertString(limit), convertString(offset), order)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, nil
	}

	if err != nil {
		return nil, 0, err
	}

	if res == nil {
		return nil, 0, nil
	}

	setLocalTimePets(res)

	return res, len(res), nil
}

// AddPet is implementing IService.AddPet function
func (s *Service) AddPet(pet *model.Pet) (int, error) {
	if err := s.repository.AddPet(pet); err != nil {
		return 0, err
	}

	return pet.ID, nil
}

// UpdatePet is implementing IService.UpdatePet function
func (s *Service) UpdatePet(pet *model.Pet) error {
	return s.repository.UpdatePet(pet)
}

// DeletePet is implementing IService.DeletePet function
func (s *Service) DeletePet(pet *model.Pet) error {
	return s.repository.DeletePet(pet)
}

// IsExist is implementing IService.IsExist function
func (s *Service) IsExist(id int) bool {
	res, _ := s.repository.GetPet(id)

	return res != nil && res.ID != 0
}

// convertString is used to convert string to int format. If string is not convertable will return 0.
func convertString(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	} else {
		return i
	}
}

// setLocalTimePets is used to set local time in all given model.Pet objects
func setLocalTimePets(pets []*model.Pet) {
	for _, p := range pets {
		p.SetLocal()
	}
}
