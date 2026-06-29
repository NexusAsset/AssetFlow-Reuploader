package accounts

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Account struct {
	Name      string `json:"name"`
	APIKey    string `json:"apiKey"`
	CreatorID string `json:"creatorId"`
	IsGroup   bool   `json:"isGroup"`
	Cookie    string `json:"cookie,omitempty"`
}

type Store struct {
	path string
	mu   sync.Mutex
	list []Account
}

func Load(path string) *Store {
	s := &Store{path: path}
	if b, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(b, &s.list)
	}
	return s
}

func (s *Store) All() []Account {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Account, len(s.list))
	copy(out, s.list)
	return out
}

func (s *Store) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.list)
}

func (s *Store) save() error {
	b, _ := json.MarshalIndent(s.list, "", "  ")
	return os.WriteFile(s.path, b, 0o600)
}

func (s *Store) Add(a Account) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range s.list {
		if a.APIKey != "" && e.APIKey == a.APIKey {
			return fmt.Errorf("that API key is already added")
		}
		if a.APIKey == "" && a.Cookie != "" && e.Cookie == a.Cookie {
			return fmt.Errorf("that downloader cookie is already added")
		}
	}
	s.list = append(s.list, a)
	return s.save()
}

// Uploaders are accounts that can upload (have an Open Cloud API key).
func (s *Store) Uploaders() []Account {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []Account
	for _, a := range s.list {
		if a.APIKey != "" {
			out = append(out, a)
		}
	}
	return out
}

// DownloaderCookies are the cookies of accounts that can fetch assets.
func (s *Store) DownloaderCookies() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []string
	for _, a := range s.list {
		if a.Cookie != "" {
			out = append(out, a.Cookie)
		}
	}
	return out
}

func (s *Store) Remove(idx int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if idx < 0 || idx >= len(s.list) {
		return fmt.Errorf("index out of range")
	}
	s.list = append(s.list[:idx], s.list[idx+1:]...)
	return s.save()
}
