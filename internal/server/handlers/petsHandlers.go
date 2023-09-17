package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"pets/internal/model"
	"pets/internal/server/handlers/requests"
	"pets/internal/server/handlers/responses"
	"pets/pkg/logger"
)

// GetPets is a handler func for GET /pet route
// Will return pets in responses.GetPetsResp format if pets found
// Will return 404 status if pets not found
// Can return 500 if unexpected DB error or encoding error occurred
func (h *Handlers) GetPets() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		l := request.URL.Query().Get("limit")
		o := request.URL.Query().Get("offset")
		ord := request.URL.Query().Get("order")

		res, total, err := h.srv.GetPets(l, o, ord)
		if err != nil {
			http.Error(writer, fmt.Sprintf("db error"), http.StatusInternalServerError)
			return
		}

		if total == 0 {
			logger.Log().WithField("layer", "Handlers-GetPets").Warningf("pets not found")
			http.Error(writer, fmt.Sprintf("pets not found"), http.StatusNotFound)
			return
		}

		log.Println(res[0].CreatedAt.Local())

		resp := &responses.GetPetsResp{
			Pets:  res,
			Total: total,
		}

		writer.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(writer).Encode(resp); err != nil {
			logger.Log().WithField("layer", "Handlers-GetPets").Errorf("error encode resp %v", err.Error())
			http.Error(writer, fmt.Sprintf("decode error"), http.StatusInternalServerError)
			return
		}
	}
}

// CreatePet is a handler func for POST /pet route
// Will return created pet ID in responses.AddPetResp format
// Will return 400 status if no request.Body provided or name in body is blank
// Can return 500 if unexpected DB error or encoding error occurred
func (h *Handlers) CreatePet() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		req := &requests.AddPetReq{}

		if err := json.NewDecoder(request.Body).Decode(req); err != nil {
			logger.Log().WithField("layer", "Handlers-CreatePet").Errorf("err decode body: %v", err.Error())
			http.Error(writer, fmt.Sprintf(`provide body params {"name":string}`), http.StatusBadRequest)
			return
		}

		if req.Name == "" {
			logger.Log().WithField("layer", "Handlers-CreatePet").Errorf("received blank name")
			http.Error(writer, fmt.Sprintf("name cannot be blank"), http.StatusBadRequest)
			return
		}

		id, err := h.srv.AddPet(model.GetPetFromReq(req))
		if err != nil {
			http.Error(writer, fmt.Sprintf("db error"), http.StatusInternalServerError)
			return
		}

		resp := &responses.AddPetResp{
			ID: id,
		}

		writer.WriteHeader(http.StatusCreated)
		if err = json.NewEncoder(writer).Encode(resp); err != nil {
			logger.Log().WithField("layer", "Handlers-CreatePet").Errorf("error encode resp %v", err.Error())
			http.Error(writer, fmt.Sprintf("decode error"), http.StatusInternalServerError)
			return
		}
	}
}

// UpdatePet is a handler func for PUT /pet route
// Will return 200 if request is successful
// Will return 400 status if no request.Body provided or name in body is blank or ID in body is less than 0 or ID not found
// Can return 500 if unexpected DB error occurred
func (h *Handlers) UpdatePet() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		req := &requests.UpdateReq{}

		if err := json.NewDecoder(request.Body).Decode(req); err != nil {
			logger.Log().WithField("layer", "Handlers-UpdatePet").Warningf("err decode body: %v", err.Error())
			http.Error(writer, fmt.Sprintf(`provide body params {"name":string, "id": number}`), http.StatusBadRequest)
			return
		}

		if req.Name == "" {
			logger.Log().WithField("layer", "Handlers-UpdatePet").Warningf("received blank name")
			http.Error(writer, fmt.Sprintf("name cannot be blank"), http.StatusBadRequest)
			return
		}

		if req.ID <= 0 {
			logger.Log().WithField("layer", "Handlers-UpdatePet").Warningf("received id less than 0: %v", req.ID)
			http.Error(writer, fmt.Sprintf("id should be more than 0"), http.StatusBadRequest)
			return
		}

		if !h.srv.IsExist(req.ID) {
			logger.Log().WithField("layer", "Handlers-UpdatePet").Warningf("pet does not exist id %v", req.ID)
			http.Error(writer, fmt.Sprintf("pet does not exist"), http.StatusBadRequest)
			return
		}

		if err := h.srv.UpdatePet(&model.Pet{ID: req.ID, Name: req.Name}); err != nil {
			http.Error(writer, fmt.Sprintf("db error"), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

// DeletePet is a handler func for DELETE /pet route
// Will return 200 if request is successful
// Will return 400 status if no request.Body provided or ID in body is less than 0 or ID not found
// Can return 500 if unexpected DB error occurred
func (h *Handlers) DeletePet() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		req := &requests.DeleteReq{}

		if err := json.NewDecoder(request.Body).Decode(req); err != nil {
			logger.Log().WithField("layer", "Handlers-DeletePet").Warningf("err decode body: %v", err.Error())
			http.Error(writer, fmt.Sprintf(`provide body params {"id": number}`), http.StatusBadRequest)
			return
		}

		if req.ID <= 0 {
			logger.Log().WithField("layer", "Handlers-DeletePet").Warningf("received id less than 0: %v", req.ID)
			http.Error(writer, fmt.Sprintf("id should be more than 0"), http.StatusBadRequest)
			return
		}

		if !h.srv.IsExist(req.ID) {
			logger.Log().WithField("layer", "Handlers-DeletePet").Warningf("pet does not exist id %v", req.ID)
			http.Error(writer, fmt.Sprintf("pet does not exist"), http.StatusBadRequest)
			return
		}

		if err := h.srv.DeletePet(&model.Pet{ID: req.ID}); err != nil {
			http.Error(writer, fmt.Sprintf("db error"), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}
