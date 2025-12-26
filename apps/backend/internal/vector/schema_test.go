package vector

import (
	"context"
	"testing"

	"github.com/weaviate/weaviate/entities/models"
)

type MockSchemaClient struct {
	CreatedClass *models.Class
}

func (m *MockSchemaClient) ClassExists(ctx context.Context, className string) (bool, error) {
	return false, nil
}

func (m *MockSchemaClient) CreateClass(ctx context.Context, class *models.Class) error {
	m.CreatedClass = class
	return nil
}

func TestEnsureSchema_Types(t *testing.T) {
	client := &MockSchemaClient{}
	if err := EnsureSchema(context.Background(), client); err != nil {
		t.Fatalf("EnsureSchema failed: %v", err)
	}

	if client.CreatedClass == nil {
		t.Fatal("Class not created")
	}

	for _, prop := range client.CreatedClass.Properties {
		if prop.Name == "sourceId" || prop.Name == "url" {
			if len(prop.DataType) == 0 || prop.DataType[0] != "string" {
				t.Errorf("Property %s has wrong DataType: %v", prop.Name, prop.DataType)
			}
		}
	}
}