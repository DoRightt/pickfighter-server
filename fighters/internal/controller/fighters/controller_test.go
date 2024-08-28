package fighters

import (
	"context"
	"errors"
	"testing"

	"pickfighter.com/fighters/gen/mocks"
	"pickfighter.com/fighters/internal/repository/psql"
	"pickfighter.com/fighters/pkg/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	mockRepo := &psql.Repository{}
	ctrl := New(mockRepo)

	if ctrl.repo != mockRepo {
		t.Errorf("expected controller to be %v, got %v", mockRepo, ctrl.repo)
	}
}

func TestSearchFightersCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockFightersRepository(ctrl)

	controller := &Controller{
		repo: mockRepo,
	}

	tests := []struct {
		name          string
		req           *model.FightersRequest
		mockReturnVal int32
		mockReturnErr error
		expectedCount int32
		expectedErr   error
	}{
		{
			name:          "Success",
			req:           &model.FightersRequest{Status: "active"},
			mockReturnVal: 5,
			mockReturnErr: nil,
			expectedCount: 5,
			expectedErr:   nil,
		},
		{
			name:          "Error",
			req:           &model.FightersRequest{Status: "active"},
			mockReturnVal: 0,
			mockReturnErr: errors.New("database error"),
			expectedCount: 0,
			expectedErr:   errors.New("database error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().
				SearchFightersCount(gomock.Any(), tc.req).
				Return(tc.mockReturnVal, tc.mockReturnErr).
				Times(1)

			count, err := controller.SearchFightersCount(context.Background(), tc.req)

			assert.Equal(t, tc.expectedCount, count)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestSearchFighters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockFightersRepository(ctrl)

	controller := &Controller{
		repo: mockRepo,
	}

	tests := []struct {
		name             string
		req              *model.FightersRequest
		mockCount        int32
		mockCountErr     error
		mockFighters     []*model.Fighter
		mockFightersErr  error
		expectedFighters []*model.Fighter
		expectedErr      error
	}{
		{
			name:             "Count Error",
			req:              &model.FightersRequest{Status: "active"},
			mockCount:        0,
			mockCountErr:     errors.New("count error"),
			mockFighters:     nil,
			mockFightersErr:  nil,
			expectedFighters: []*model.Fighter{},
			expectedErr:      errors.New("count error"),
		},
		{
			name:             "Success with No Fighters",
			req:              &model.FightersRequest{Status: "active"},
			mockCount:        0,
			mockCountErr:     nil,
			mockFighters:     nil,
			mockFightersErr:  nil,
			expectedFighters: []*model.Fighter{},
			expectedErr:      nil,
		},
		{
			name:             "Fighters Error",
			req:              &model.FightersRequest{Status: "active"},
			mockCount:        2,
			mockCountErr:     nil,
			mockFighters:     nil,
			mockFightersErr:  errors.New("fighters error"),
			expectedFighters: []*model.Fighter{},
			expectedErr:      errors.New("fighters error"),
		},
		{
			name:             "Success with Fighters",
			req:              &model.FightersRequest{Status: "active"},
			mockCount:        2,
			mockCountErr:     nil,
			mockFighters:     []*model.Fighter{{}},
			mockFightersErr:  nil,
			expectedFighters: []*model.Fighter{{}},
			expectedErr:      nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().
				SearchFightersCount(gomock.Any(), tc.req).
				Return(tc.mockCount, tc.mockCountErr).
				Times(1)

			if tc.mockCountErr == nil && tc.mockCount > 0 {
				mockRepo.EXPECT().
					SearchFighters(gomock.Any(), tc.req).
					Return(tc.mockFighters, tc.mockFightersErr).
					Times(1)
			}

			// Вызов тестируемого метода
			fighters, err := controller.SearchFighters(context.Background(), tc.req)

			// Проверка результатов
			assert.Equal(t, tc.expectedFighters, fighters)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
