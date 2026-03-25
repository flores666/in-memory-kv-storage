package storage

import (
	"sync"
	"time"
)

type StorageService interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
}

type storage struct {
	m     map[string]string
	mutex sync.Mutex
}

const TTL = 1 * time.Minute

func NewStorageService() StorageService {
	return &storage{
		m: make(map[string]string),
	}
}

func (s *storage) Get(key string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.m[key], nil
}

func (s *storage) Set(key, value string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.m[key] = value

	go func() {
		time.Sleep(TTL)
		s.Delete(key)
	}()

	return nil
}

func (s *storage) Delete(key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.m, key)

	return nil
}
