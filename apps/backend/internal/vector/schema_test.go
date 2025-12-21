package vector_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/weaviate/weaviate/entities/models"
	"qurio/apps/backend/internal/vector"
)

// MockSchemaClient simulates the vector database schema operations
type MockSchemaClient struct {
	mock.Mock
}

func (m *MockSchemaClient) ClassExists(ctx context.Context, className string) (bool, error) {
	args := m.Called(ctx, className)
	return args.Bool(0), args.Error(1)
}

func (m *MockSchemaClient) CreateClass(ctx context.Context, class *models.Class) error {
	args := m.Called(ctx, class)
	return args.Error(0)
}

func TestEnsureSchema_CreatesClass_WhenNotExists(t *testing.T) {
	mockClient := new(MockSchemaClient)
	// Expect check for "DocumentChunk" -> returns false (not exists)
	mockClient.On("ClassExists", mock.Anything, "DocumentChunk").Return(false, nil)
	// Expect create class -> returns nil (success)
	mockClient.On("CreateClass", mock.Anything, mock.MatchedBy(func(c *models.Class) bool {
		return c.Class == "DocumentChunk"
	})).Return(nil)

	err := vector.EnsureSchema(context.Background(), mockClient)
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestEnsureSchema_DoNothing_WhenExists(t *testing.T) {
	mockClient := new(MockSchemaClient)
	mockClient.On("ClassExists", mock.Anything, "DocumentChunk").Return(true, nil)
	// CreateClass should NOT be called

	err := vector.EnsureSchema(context.Background(), mockClient)
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}
