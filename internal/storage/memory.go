package storage

import (
	"strings"
	"sync"

	"github.com/yourname/logflow/internal/models"
)

// Store defines the storage interface — swap with PostgreSQL easily
type Store interface {
	Save(entry *models.LogEntry) error
	Query(filter models.QueryFilter) ([]*models.LogEntry, error)
	Stats() (*models.Stats, error)
}

// MemoryStore is a thread-safe in-memory log store
type MemoryStore struct {
	mu   sync.RWMutex
	logs []*models.LogEntry
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		logs: make([]*models.LogEntry, 0),
	}
}

func (m *MemoryStore) Save(entry *models.LogEntry) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logs = append(m.logs, entry)
	return nil
}

func (m *MemoryStore) Query(filter models.QueryFilter) ([]*models.LogEntry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	limit := filter.Limit
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	result := make([]*models.LogEntry, 0)

	// Iterate in reverse (latest first)
	for i := len(m.logs) - 1; i >= 0; i-- {
		entry := m.logs[i]

		if filter.Level != "" && entry.Level != filter.Level {
			continue
		}
		if filter.Service != "" && !strings.EqualFold(entry.Service, filter.Service) {
			continue
		}
		if filter.Search != "" && !strings.Contains(
			strings.ToLower(entry.Message),
			strings.ToLower(filter.Search),
		) {
			continue
		}

		result = append(result, entry)
		if len(result) >= limit {
			break
		}
	}

	return result, nil
}

func (m *MemoryStore) Stats() (*models.Stats, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := &models.Stats{
		Total:     len(m.logs),
		ByLevel:   make(map[models.Level]int),
		ByService: make(map[string]int),
	}

	for _, entry := range m.logs {
		stats.ByLevel[entry.Level]++
		stats.ByService[entry.Service]++
	}

	return stats, nil
}
