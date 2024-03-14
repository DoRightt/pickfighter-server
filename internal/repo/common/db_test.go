package repo

import (
	mock_repo "projects/fb-server/internal/repo/common/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockFbCommonRepo(ctrl)

	cmRepo := New(mockRepo)

	assert.NotNil(t, cmRepo, "cmRepo should not be nil")
}
