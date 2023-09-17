package model

import (
	"time"

	"pets/internal/server/handlers/requests"
)

// Pet is a pet model struct
type Pet struct {
	// ID is a pet id
	ID int `json:"id"`
	// Name is a pet name
	Name string `json:"name"`
	// CreatedAt is a date when pet was created
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	// UpdatedAt is a date when pet was updated. Can be nil
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// GetPetFromReq is used to get Pet model from given requests.AddPetReq model
func GetPetFromReq(req *requests.AddPetReq) *Pet {
	return &Pet{
		Name: req.Name,
	}
}

// SetLocal is used to set local time format
func (p *Pet) SetLocal() {
	p.CreatedAt = p.CreatedAt.Local()

	if p.UpdatedAt != nil {
		l := p.UpdatedAt.Local()
		p.UpdatedAt = &l
	}
}
