package failcache

import (
	"encoding/json"
	"os"
	"sync"
)

type Store struct {
	path string
	mu   sync.Mutex
	data map[string]map[string]bool
}

func Load(path string) *Store {
	s := &Store{path: path, data: map[string]map[string]bool{}}
	if b, err := os.ReadFile(path); err == nil {
		raw := map[string][]string{}
		if json.Unmarshal(b, &raw) == nil {
			for k, ids := range raw {
				m := map[string]bool{}
				for _, id := range ids {
					m[id] = true
				}
				s.data[k] = m
			}
		}
	}
	return s
}

func (s *Store) save() {
	raw := map[string][]string{}
	for k, m := range s.data {
		ids := make([]string, 0, len(m))
		for id := range m {
			ids = append(ids, id)
		}
		raw[k] = ids
	}
	b, _ := json.Marshal(raw)
	_ = os.WriteFile(s.path, b, 0o600)
}

func (s *Store) Failed(key string, ids []string) map[string]bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := map[string]bool{}
	m := s.data[key]
	if m == nil {
		return out
	}
	for _, id := range ids {
		if m[id] {
			out[id] = true
		}
	}
	return out
}

func (s *Store) Mark(key string, ids []string) {
	if len(ids) == 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data[key] == nil {
		s.data[key] = map[string]bool{}
	}
	for _, id := range ids {
		s.data[key][id] = true
	}
	s.save()
}

func (s *Store) Clear(key string, ids []string) {
	if len(ids) == 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	m := s.data[key]
	if m == nil {
		return
	}
	for _, id := range ids {
		delete(m, id)
	}
	s.save()
}
