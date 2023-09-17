package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"pets/internal/model"
	"pets/internal/server/handlers/requests"
	"pets/internal/server/handlers/responses"
	mock_service "pets/mocks/service"
)

func TestHandlers_GetPets(t *testing.T) {
	ctrl := gomock.NewController(t)
	srvMock := mock_service.NewMockIService(ctrl)
	defer ctrl.Finish()

	tests := []struct {
		name   string
		url    string
		limit  string
		offset string
		order  string

		goToSev bool
		srvErr  error
		pets    []*model.Pet
		total   int

		wantBody   *responses.GetPetsResp
		wantStatus int
		wantErr    string
	}{
		{
			name:       "check 200 query",
			url:        "/pets?limit=1&offset=2&order=asc",
			limit:      "1",
			offset:     "2",
			order:      "asc",
			goToSev:    true,
			pets:       []*model.Pet{{ID: 1, Name: "Velho"}, {ID: 2, Name: "Melho"}},
			total:      2,
			wantBody:   &responses.GetPetsResp{Pets: []*model.Pet{{ID: 1, Name: "Velho"}, {ID: 2, Name: "Melho"}}, Total: 2},
			wantStatus: http.StatusOK,
		},
		{
			name:       "check 200 no query",
			url:        "/pets",
			limit:      "",
			offset:     "",
			order:      "",
			goToSev:    true,
			pets:       []*model.Pet{{ID: 1, Name: "Velho"}, {ID: 2, Name: "Melho"}},
			total:      2,
			wantBody:   &responses.GetPetsResp{Pets: []*model.Pet{{ID: 1, Name: "Velho"}, {ID: 2, Name: "Melho"}}, Total: 2},
			wantStatus: http.StatusOK,
		},
		{
			name:       "check 404 not found",
			url:        "/pets",
			limit:      "",
			offset:     "",
			order:      "",
			goToSev:    true,
			total:      0,
			wantStatus: http.StatusNotFound,
			wantErr:    "pets not found",
		},
		{
			name:       "check 500 db error",
			url:        "/pets",
			limit:      "",
			offset:     "",
			order:      "",
			goToSev:    true,
			srvErr:     fmt.Errorf("db error occurred"),
			wantStatus: http.StatusInternalServerError,
			wantErr:    "db error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandlers(srvMock)
			getPets := h.GetPets()

			res := httptest.NewRecorder()
			var b []byte

			body := bytes.NewReader(b)
			req, _ := http.NewRequest("GET", tt.url, body)

			if tt.goToSev {
				srvMock.EXPECT().GetPets(tt.limit, tt.offset, tt.order).Return(tt.pets, tt.total, tt.srvErr)
			}

			getPets.ServeHTTP(res, req)

			want := httptest.NewRecorder()
			if tt.wantErr == "" {
				json.NewEncoder(want).Encode(tt.wantBody)
			} else {
				http.Error(want, tt.wantErr, tt.wantStatus)
			}

			require.Equal(t, tt.wantStatus, res.Code)
			require.Equal(t, want.Body.String(), res.Body.String())
		})
	}
}

func TestHandlers_CreatePet(t *testing.T) {
	ctrl := gomock.NewController(t)
	srvMock := mock_service.NewMockIService(ctrl)
	defer ctrl.Finish()

	tests := []struct {
		name string
		req  *requests.AddPetReq

		goToSev bool
		id      int
		srvErr  error
		pet     *model.Pet

		wantBody   *responses.AddPetResp
		wantStatus int
		wantErr    string
	}{
		{
			name:       "check 201",
			req:        &requests.AddPetReq{Name: "Velho"},
			goToSev:    true,
			pet:        &model.Pet{Name: "Velho"},
			id:         1,
			wantBody:   &responses.AddPetResp{ID: 1},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "check 400 no body",
			wantStatus: http.StatusBadRequest,
			wantErr:    `provide body params {"name":string}`,
		},
		{
			name:       "check 400 blank name",
			req:        &requests.AddPetReq{Name: ""},
			wantStatus: http.StatusBadRequest,
			wantErr:    "name cannot be blank",
		},
		{
			name:       "check 500 db error",
			req:        &requests.AddPetReq{Name: "Velho"},
			goToSev:    true,
			pet:        &model.Pet{Name: "Velho"},
			srvErr:     fmt.Errorf("db error occurred"),
			wantStatus: http.StatusInternalServerError,
			wantErr:    "db error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandlers(srvMock)
			createPet := h.CreatePet()

			res := httptest.NewRecorder()
			var b []byte

			if tt.req != nil {
				b, _ = json.Marshal(tt.req)
			}

			body := bytes.NewReader(b)
			req, _ := http.NewRequest("POST", "/pet", body)

			if tt.goToSev {
				srvMock.EXPECT().AddPet(tt.pet).Return(tt.id, tt.srvErr)
			}

			createPet.ServeHTTP(res, req)

			want := httptest.NewRecorder()
			if tt.wantErr == "" {
				json.NewEncoder(want).Encode(tt.wantBody)
			} else {
				http.Error(want, tt.wantErr, tt.wantStatus)
			}

			require.Equal(t, tt.wantStatus, res.Code)
			require.Equal(t, want.Body.String(), res.Body.String())
		})
	}
}

func TestHandlers_UpdatePet(t *testing.T) {
	ctrl := gomock.NewController(t)
	srvMock := mock_service.NewMockIService(ctrl)
	defer ctrl.Finish()

	tests := []struct {
		name string
		req  *requests.UpdateReq

		goToExist bool
		exist     bool

		goToSev bool
		srvErr  error
		pet     *model.Pet

		wantStatus int
		wantErr    string
	}{
		{
			name:       "check 200",
			req:        &requests.UpdateReq{Name: "Velho", ID: 1},
			goToExist:  true,
			exist:      true,
			goToSev:    true,
			pet:        &model.Pet{Name: "Velho", ID: 1},
			wantStatus: http.StatusOK,
		},
		{
			name:       "check 400 no body",
			wantStatus: http.StatusBadRequest,
			wantErr:    `provide body params {"name":string, "id": number}`,
		},
		{
			name:       "check 400 blank name",
			req:        &requests.UpdateReq{Name: "", ID: 1},
			wantStatus: http.StatusBadRequest,
			wantErr:    "name cannot be blank",
		},
		{
			name:       "check 400 0 id",
			req:        &requests.UpdateReq{Name: "Velho", ID: 0},
			wantStatus: http.StatusBadRequest,
			wantErr:    "id should be more than 0",
		},
		{
			name:       "check 400 not exist",
			req:        &requests.UpdateReq{Name: "Velho", ID: 1},
			goToExist:  true,
			exist:      false,
			wantStatus: http.StatusBadRequest,
			wantErr:    "pet does not exist",
		},
		{
			name:       "check 500 db error",
			req:        &requests.UpdateReq{Name: "Velho", ID: 1},
			goToExist:  true,
			exist:      true,
			goToSev:    true,
			pet:        &model.Pet{Name: "Velho", ID: 1},
			srvErr:     fmt.Errorf("db error occured"),
			wantStatus: http.StatusInternalServerError,
			wantErr:    "db error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandlers(srvMock)
			updatePet := h.UpdatePet()

			res := httptest.NewRecorder()
			var b []byte

			if tt.req != nil {
				b, _ = json.Marshal(tt.req)
			}

			body := bytes.NewReader(b)
			req, _ := http.NewRequest("PUT", "/pet", body)

			if tt.goToExist {
				srvMock.EXPECT().IsExist(tt.req.ID).Return(tt.exist)
			}

			if tt.goToSev {
				srvMock.EXPECT().UpdatePet(tt.pet).Return(tt.srvErr)
			}

			updatePet.ServeHTTP(res, req)

			want := httptest.NewRecorder()
			if tt.wantErr != "" {
				http.Error(want, tt.wantErr, tt.wantStatus)
			}

			require.Equal(t, tt.wantStatus, res.Code)
			require.Equal(t, want.Body.String(), res.Body.String())
		})
	}
}

func TestHandlers_DeletePet(t *testing.T) {
	ctrl := gomock.NewController(t)
	srvMock := mock_service.NewMockIService(ctrl)
	defer ctrl.Finish()

	tests := []struct {
		name string
		req  *requests.DeleteReq

		goToExist bool
		exist     bool

		goToSev bool
		srvErr  error
		pet     *model.Pet

		wantStatus int
		wantErr    string
	}{
		{
			name:       "check 200",
			req:        &requests.DeleteReq{ID: 1},
			goToExist:  true,
			exist:      true,
			goToSev:    true,
			pet:        &model.Pet{ID: 1},
			wantStatus: http.StatusOK,
		},
		{
			name:       "check 400 no body",
			wantStatus: http.StatusBadRequest,
			wantErr:    `provide body params {"id": number}`,
		},
		{
			name:       "check 400 0 id",
			req:        &requests.DeleteReq{ID: 0},
			wantStatus: http.StatusBadRequest,
			wantErr:    "id should be more than 0",
		},
		{
			name:       "check 400 not exist",
			req:        &requests.DeleteReq{ID: 1},
			goToExist:  true,
			exist:      false,
			wantStatus: http.StatusBadRequest,
			wantErr:    "pet does not exist",
		},
		{
			name:       "check 500 db error",
			req:        &requests.DeleteReq{ID: 1},
			goToExist:  true,
			exist:      true,
			goToSev:    true,
			pet:        &model.Pet{ID: 1},
			srvErr:     fmt.Errorf("db error occured"),
			wantStatus: http.StatusInternalServerError,
			wantErr:    "db error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandlers(srvMock)
			deletePet := h.DeletePet()

			res := httptest.NewRecorder()
			var b []byte

			if tt.req != nil {
				b, _ = json.Marshal(tt.req)
			}

			body := bytes.NewReader(b)
			req, _ := http.NewRequest("DELETE", "/pet", body)

			if tt.goToExist {
				srvMock.EXPECT().IsExist(tt.req.ID).Return(tt.exist)
			}

			if tt.goToSev {
				srvMock.EXPECT().DeletePet(tt.pet).Return(tt.srvErr)
			}

			deletePet.ServeHTTP(res, req)

			want := httptest.NewRecorder()
			if tt.wantErr != "" {
				http.Error(want, tt.wantErr, tt.wantStatus)
			}

			require.Equal(t, tt.wantStatus, res.Code)
			require.Equal(t, want.Body.String(), res.Body.String())
		})
	}
}
