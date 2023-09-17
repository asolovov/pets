package requests

// AddPetReq is a form of request accepted in POST /pet route
type AddPetReq struct {
	// Name is a pet name to add
	Name string `json:"name"`
}

// UpdateReq is a form of request accepted in PUT /pet route
type UpdateReq struct {
	// ID is a pet ID to update
	ID int `json:"id"`
	// Name is a new pet Name
	Name string `json:"name"`
}

// DeleteReq is a form of request accepted in DELETE /pet route
type DeleteReq struct {
	// ID is a pet ID to delete
	ID int `json:"id"`
}
