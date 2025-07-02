package mockcache

import (
	storageerror "api-gateway/internal/storage"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type MockCache struct {
	log      *slog.Logger
	elements map[uuid.UUID]any
}

func New(log *slog.Logger) *MockCache {
	return &MockCache{
		log:      log,
		elements: make(map[uuid.UUID]any),
	}
}

func (m *MockCache) Get(id uuid.UUID) (any, error) {
	if value, ok := m.elements[id]; !ok {
		return nil, fmt.Errorf("%s: %w", "Not found", storageerror.ErrNotFound)
	} else {
		return value, nil
	}
}

func (m *MockCache) Set(id uuid.UUID, obj any) {
	m.elements[id] = obj
}

func (m *MockCache) Delete(id uuid.UUID) {
	delete(m.elements, id)
}
