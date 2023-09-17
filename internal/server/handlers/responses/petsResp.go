package responses

import "pets/internal/model"

// AddPetResp is a form of response for POST /pet route
type AddPetResp struct {
	// ID is an added pet ID
	ID int `json:"id"`
}

// GetPetsResp is a form of response for GET /pet route
type GetPetsResp struct {
	// Pets is a slice of model.Pet found
	Pets []*model.Pet `json:"pets"`
	// Total is a Pets length value
	Total int `json:"total"`
}
