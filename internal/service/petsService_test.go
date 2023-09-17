package service

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"pets/internal/model"
	mock_repository "pets/mocks/repository"
)

func TestService_GetPets(t *testing.T) {
	ctrl := gomock.NewController(t)
	repMock := mock_repository.NewMockIRepository(ctrl)
	defer ctrl.Finish()

	type req struct {
		limit  string
		offset string
		order  string
	}
	type repReq struct {
		limit  int
		offset int
		order  string
	}
	tests := []struct {
		name string

		req    *req
		repReq *repReq

		repPets []*model.Pet
		repErr  error

		wantTotal int
		wantErr   bool
	}{
		{
			name:      "check total 3",
			req:       &req{},
			repReq:    &repReq{},
			repPets:   []*model.Pet{{ID: 1, Name: "Pet1"}, {ID: 2, Name: "Pet2"}, {ID: 3, Name: "Pet3"}},
			wantTotal: 3,
		},
		{
			name:      "check total 2",
			req:       &req{},
			repReq:    &repReq{},
			repPets:   []*model.Pet{{ID: 1, Name: "Pet1"}, {ID: 2, Name: "Pet2"}},
			wantTotal: 2,
		},
		{
			name:      "check total 1",
			req:       &req{},
			repReq:    &repReq{},
			repPets:   []*model.Pet{{ID: 1, Name: "Pet1"}},
			wantTotal: 1,
		},
		{
			name:      "check total 1",
			req:       &req{},
			repReq:    &repReq{},
			repPets:   []*model.Pet{{ID: 1, Name: "Pet1"}},
			wantTotal: 1,
		},
		{
			name:      "check total 0",
			req:       &req{},
			repReq:    &repReq{},
			repPets:   []*model.Pet{},
			wantTotal: 0,
		},
		{
			name:      "check total nil-0",
			req:       &req{},
			repReq:    &repReq{},
			wantTotal: 0,
		},
		{
			name:      "check convert limit",
			req:       &req{limit: "1"},
			repReq:    &repReq{limit: 1},
			repPets:   []*model.Pet{{ID: 1, Name: "Pet1"}},
			wantTotal: 1,
		},
		{
			name:      "check convert offset",
			req:       &req{limit: "2"},
			repReq:    &repReq{limit: 2},
			repPets:   []*model.Pet{{ID: 1, Name: "Pet1"}},
			wantTotal: 1,
		},
		{
			name:      "check convert invalid",
			req:       &req{limit: "notInt"},
			repReq:    &repReq{},
			repPets:   []*model.Pet{{ID: 1, Name: "Pet1"}},
			wantTotal: 1,
		},
		{
			name:      "check forward order",
			req:       &req{order: "desc"},
			repReq:    &repReq{order: "desc"},
			repPets:   []*model.Pet{{ID: 1, Name: "Pet1"}},
			wantTotal: 1,
		},
		{
			name:    "check rep error",
			req:     &req{},
			repReq:  &repReq{},
			repErr:  fmt.Errorf("rep error"),
			wantErr: true,
		},
		{
			name:      "check rep no rows",
			req:       &req{},
			repReq:    &repReq{},
			repErr:    sql.ErrNoRows,
			wantTotal: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(repMock)

			repMock.EXPECT().GetPets(tt.repReq.limit, tt.repReq.offset, tt.repReq.order).Return(tt.repPets, tt.repErr)

			res, resTotal, err := s.GetPets(tt.req.limit, tt.req.offset, tt.req.order)

			if !tt.wantErr {
				require.NoError(t, err)
				require.Equal(t, tt.wantTotal, resTotal)
				require.Equal(t, tt.repPets, res)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestService_AddPet(t *testing.T) {
	ctrl := gomock.NewController(t)
	repMock := mock_repository.NewMockIRepository(ctrl)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		pet     *model.Pet
		wantId  int
		repErr  error
		wantErr bool
	}{
		{
			name:   "check id return 1",
			pet:    &model.Pet{Name: "Velho"},
			wantId: 1,
		},
		{
			name:   "check id return 2",
			pet:    &model.Pet{Name: "Velho"},
			wantId: 2,
		},
		{
			name:   "check id return 3",
			pet:    &model.Pet{Name: "Velho"},
			wantId: 3,
		},
		{
			name:    "check id return 0 with err",
			pet:     &model.Pet{Name: "Velho"},
			repErr:  fmt.Errorf("rep error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(repMock)

			repMock.EXPECT().AddPet(tt.pet).DoAndReturn(func(p *model.Pet) {
				if tt.repErr == nil {
					p.ID = tt.wantId
				}
			}).Return(tt.repErr)

			id, err := s.AddPet(tt.pet)

			if !tt.wantErr {
				require.NoError(t, err)
				require.Equal(t, tt.wantId, id)
			} else {
				require.Equal(t, 0, id)
				require.Error(t, err)
			}
		})
	}
}

func TestService_UpdatePet(t *testing.T) {
	ctrl := gomock.NewController(t)
	repMock := mock_repository.NewMockIRepository(ctrl)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		pet     *model.Pet
		repErr  error
		wantErr bool
	}{
		{
			name: "no error",
			pet:  &model.Pet{Name: "Velho", ID: 1},
		},
		{
			name:    "error",
			pet:     &model.Pet{Name: "Velho", ID: 1},
			repErr:  fmt.Errorf("rep error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(repMock)

			repMock.EXPECT().UpdatePet(tt.pet).Return(tt.repErr)

			err := s.UpdatePet(tt.pet)

			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestService_DeletePet(t *testing.T) {
	ctrl := gomock.NewController(t)
	repMock := mock_repository.NewMockIRepository(ctrl)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		pet     *model.Pet
		repErr  error
		wantErr bool
	}{
		{
			name: "no error",
			pet:  &model.Pet{Name: "Velho", ID: 1},
		},
		{
			name:    "error",
			pet:     &model.Pet{Name: "Velho", ID: 1},
			repErr:  fmt.Errorf("rep error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(repMock)

			repMock.EXPECT().DeletePet(tt.pet).Return(tt.repErr)

			err := s.DeletePet(tt.pet)

			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestService_IsExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	repMock := mock_repository.NewMockIRepository(ctrl)
	defer ctrl.Finish()

	tests := []struct {
		name   string
		id     int
		pet    *model.Pet
		res    bool
		repErr error
	}{
		{
			name: "check good id 1",
			id:   1,
			pet:  &model.Pet{ID: 1},
			res:  true,
		},
		{
			name: "check good id 2",
			id:   2,
			pet:  &model.Pet{ID: 2},
			res:  true,
		},
		{
			name: "check good id 3",
			id:   3,
			pet:  &model.Pet{ID: 3},
			res:  true,
		},
		{
			name: "check return nil",
			id:   3,
			res:  false,
		},
		{
			name: "check return 0 id",
			id:   3,
			pet:  &model.Pet{ID: 0},
			res:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(repMock)

			repMock.EXPECT().GetPet(tt.id).Return(tt.pet, tt.repErr)

			res := s.IsExist(tt.id)

			require.Equal(t, tt.res, res)
		})
	}
}
