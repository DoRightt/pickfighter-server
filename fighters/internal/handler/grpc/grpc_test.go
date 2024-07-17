package grpc

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"fightbettr.com/fighters/gen/mocks"
	"fightbettr.com/fighters/internal/controller/fighters"
	"fightbettr.com/fighters/pkg/model"
	"fightbettr.com/gen"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/go-playground/assert.v1"
)

func TestNew(t *testing.T) {
	mockCtrl := &fighters.Controller{}
	handler := New(mockCtrl)

	if handler.ctrl != mockCtrl {
		t.Errorf("expected controller to be %v, got %v", mockCtrl, handler.ctrl)
	}
}

func TestSearchFightersCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCtrl := mocks.NewMockFightersController(ctrl)
	handler := &Handler{ctrl: mockCtrl}
	ctx := context.Background()

	tests := []struct {
		name            string
		req             *gen.FightersRequest
		mockBehavior    func()
		expectedResp    *gen.FightersCountResponse
		expectedErr     error
		expectedErrCode codes.Code
	}{
		{
			name:            "Nil request",
			req:             nil,
			mockBehavior:    func() {},
			expectedResp:    nil,
			expectedErr:     status.Errorf(codes.InvalidArgument, "nil request"),
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "Controller error",
			req:  &gen.FightersRequest{Status: "active"},
			mockBehavior: func() {
				mockCtrl.EXPECT().SearchFightersCount(ctx, gomock.Any()).Return(int32(0), fmt.Errorf("some error"))
			},
			expectedResp:    nil,
			expectedErr:     status.Errorf(codes.Internal, "some error"),
			expectedErrCode: codes.Internal,
		},
		{
			name: "Success",
			req:  &gen.FightersRequest{Status: "active"},
			mockBehavior: func() {
				mockCtrl.EXPECT().SearchFightersCount(ctx, gomock.Any()).Return(int32(10), nil)
			},
			expectedResp: &gen.FightersCountResponse{Count: 10},
			expectedErr:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			resp, err := handler.SearchFightersCount(ctx, tc.req)
			if tc.expectedErr != nil {
				if err == nil || status.Code(err) != tc.expectedErrCode {
					t.Errorf("expected error code %v, got %v", tc.expectedErrCode, status.Code(err))
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}

			if resp != nil && tc.expectedResp != nil {
				if resp.Count != tc.expectedResp.Count {
					t.Errorf("expected count %v, got %v", tc.expectedResp.Count, resp.Count)
				}
			} else if resp != tc.expectedResp {
				t.Errorf("expected response %v, got %v", tc.expectedResp, resp)
			}
		})
	}
}

func TestSearchFighters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCtrl := mocks.NewMockFightersController(ctrl)
	handler := &Handler{ctrl: mockCtrl}
	ctx := context.Background()

	tests := []struct {
		name          string
		req           *gen.FightersRequest
		mockResp      []*model.Fighter
		mockErr       error
		expectedResp  *gen.FightersResponse
		expectedError error
	}{
		{
			name:          "Nil request",
			req:           nil,
			mockResp:      nil,
			mockErr:       nil,
			expectedResp:  nil,
			expectedError: status.Errorf(codes.InvalidArgument, "nil request"),
		},
		{
			name:          "Controller error not found",
			req:           &gen.FightersRequest{Status: "inactive", FightersIds: []int32{-5}},
			mockResp:      nil,
			mockErr:       fighters.ErrNotFound,
			expectedResp:  nil,
			expectedError: status.Errorf(codes.NotFound, "not found"),
		},
		{
			name:          "Controller error",
			req:           &gen.FightersRequest{Status: "inactive"},
			mockResp:      nil,
			mockErr:       errors.New("internal error"),
			expectedResp:  nil,
			expectedError: status.Errorf(codes.Internal, "internal error"),
		},
		{
			name:          "Success",
			req:           &gen.FightersRequest{Status: "active", FightersIds: []int32{1, 2}},
			mockResp:      []*model.Fighter{{FighterId: 1}, {FighterId: 2}},
			mockErr:       nil,
			expectedResp:  &gen.FightersResponse{Fighters: model.FightersToProto([]*model.Fighter{{FighterId: 1}, {FighterId: 2}})},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.req != nil {
				fReq := &model.FightersRequest{Status: tc.req.Status, FightersIds: tc.req.FightersIds}
				mockCtrl.EXPECT().SearchFighters(gomock.Any(), fReq).Return(tc.mockResp, tc.mockErr)
			}

			resp, err := handler.SearchFighters(ctx, tc.req)

			assert.Equal(t, tc.expectedResp, resp)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
