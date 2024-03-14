package repo

import (
	mock_repo "projects/fb-server/internal/repo/auth/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockFbAuthRepo(ctrl)

	authRepo := New(mockRepo)

	assert.NotNil(t, authRepo, "authRepo should not be nil")

	switch v := authRepo.(type) {
	case FbAuthRepo:
		return
	default:
		t.Errorf("Expected Repo interface is FbAuthRepo, but got %v", v)
	}
}
