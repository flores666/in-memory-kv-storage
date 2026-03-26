package storage

import (
	"errors"
	"time"

	"github.com/cespare/xxhash/v2"
)

type StorageService interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
}

type storage struct {
	shards []*shard
}

const TTL = 10 * time.Second
const shardsCount = 256

var ErrItemNotFound = errors.New("item not found")

func NewStorageService() StorageService {
	s := make([]*shard, shardsCount)
	for i := range s {
		s[i] = NewShard()
	}

	storage := &storage{
		shards: s,
	}

	go storage.expirationHandler()

	return storage
}

func (s *storage) Get(key string) (string, error) {
	shard := s.shards[getShardIndex(key)]
	shard.mutex.Lock()
	defer shard.mutex.Unlock()

	item, ok := shard.m[key]
	if !ok {
		return "", ErrItemNotFound
	}

	if item.expirationDate.Before(time.Now()) {
		delete(shard.m, key)
		return "", ErrItemNotFound
	}

	item.expirationDate = time.Now().Add(TTL)

	return item.value, nil
}

func (s *storage) Set(key, value string) error {
	shard := s.shards[getShardIndex(key)]

	shard.mutex.Lock()
	shard.m[key] = &item{
		value:          value,
		expirationDate: time.Now().Add(TTL),
	}
	shard.mutex.Unlock()

	return nil
}

func (s *storage) Delete(key string) error {
	shard := s.shards[getShardIndex(key)]

	shard.mutex.Lock()
	delete(shard.m, key)
	shard.mutex.Unlock()

	return nil
}

func getShardIndex(value string) int {
	return int(xxhash.Sum64String(value) % uint64(shardsCount))
}

func (s *storage) expirationHandler() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()

		for i := range s.shards {
			shard := s.shards[i]
			shard.mutex.Lock()

			for key, value := range shard.m {
				if value.expirationDate.Before(now) {
					delete(shard.m, key)
				}
			}
			shard.mutex.Unlock()
		}
	}
}
