package repository

import (
	"fmt"
	"strings"
	"time"

	"pets/internal/model"
	"pets/pkg/logger"
)

// GetPet is used to get pet from DB by given ID
func (r *Repository) GetPet(id int) (pet *model.Pet, err error) {
	pet = &model.Pet{}

	q := `SELECT * FROM pets WHERE id = $1 LIMIT 1`

	err = r.db.QueryRow(q, id).Scan(&pet.ID, &pet.Name, &pet.CreatedAt, &pet.UpdatedAt)
	if err != nil {
		logger.Log().WithField("layer", "Repository-GetPet").Errorf("err query: %v", err.Error())
		return nil, err
	}

	return pet, nil
}

// GetPets is used to get pet from DB. Pagination can be used by setting limit and offset values. Order should be
// "asc" or "desc" in any register, all other values will be ignored. 0 limit will be ignored.
func (r *Repository) GetPets(limit int, offset int, order string) (pets []*model.Pet, err error) {
	q := `SELECT id, name, created_at, updated_at	 FROM pets`

	if strings.ToLower(order) == "asc" {
		q = fmt.Sprintf("%v ORDER BY id ASC", q)
	}

	if strings.ToLower(order) == "desc" {
		q = fmt.Sprintf("%v ORDER BY id DESC", q)
	}

	if limit != 0 {
		q = fmt.Sprintf("%v LIMIT %v", q, limit)
	}

	q = fmt.Sprintf("%v OFFSET %v", q, offset)

	err = r.db.Select(&pets, q)
	if err != nil {
		logger.Log().WithField("layer", "Repository-GetPets").Errorf("err query: %v", err.Error())
		return nil, err
	}

	return pets, nil
}

// AddPet is used to add new pet to the DB. Only "name" field will be used. Fields id and created_at will be set automatically
func (r *Repository) AddPet(pet *model.Pet) error {
	q := `INSERT INTO pets (name, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id`

	pet.CreatedAt = time.Now()

	err := r.db.QueryRow(q, pet.Name, pet.CreatedAt, pet.UpdatedAt).Scan(&pet.ID)
	if err != nil {
		logger.Log().WithField("layer", "Repository-AddPet").Errorf("err query: %v", err.Error())
		return err
	}

	return nil
}

// UpdatePet is used to update existing pet to the DB by given id filed. Only "name" field will be used. Fields id and
// updated_at will be set automatically
func (r *Repository) UpdatePet(pet *model.Pet) error {
	q := `UPDATE pets SET name = $1, updated_at = $2 WHERE id=$3`

	now := time.Now()
	pet.UpdatedAt = &now

	_, err := r.db.Exec(q, pet.Name, pet.UpdatedAt, pet.ID)
	if err != nil {
		logger.Log().WithField("layer", "Repository-UpdatePet").Errorf("err query: %v", err.Error())
		return err
	}

	return nil
}

// DeletePet is used to delete pet from the DB by given id
func (r *Repository) DeletePet(pet *model.Pet) error {
	q := `DELETE FROM pets WHERE id=$1`

	pet.CreatedAt = time.Now()

	_, err := r.db.Exec(q, pet.ID)
	if err != nil {
		logger.Log().WithField("layer", "Repository-DeletePet").Errorf("err query: %v", err.Error())
		return err
	}

	return nil
}
