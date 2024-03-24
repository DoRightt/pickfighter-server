package repo

import (
	mock_repo "projects/fb-server/internal/repo/fighters/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockFbFightersRepo(ctrl)

	fighterRepo := New(mockRepo)

	assert.NotNil(t, fighterRepo, "fighterRepo should not be nil")
}
